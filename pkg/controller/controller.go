package controller

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/golang/glog"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"

	kubeinformers "k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	namespacelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	mwiapi "github.com/dbenque/mw-injector/pkg/api/mwinjector/v1"

	mwiinformers "github.com/dbenque/mw-injector/pkg/client/informers/externalversions"
	mwilisters "github.com/dbenque/mw-injector/pkg/client/listers/mwinjector/v1"
	mwiclientset "github.com/dbenque/mw-injector/pkg/client/versioned"
)

// MWInjectorControllerConfig contains info to customize MWInjector controller behaviour
type MWInjectorControllerConfig struct {
	RemoveIvalidMWInjector bool
	NumberOfThreads        int
}

// MWInjectorController represents the MWInjector controller
type MWInjectorController struct {
	MWInjectorClient mwiclientset.Interface
	KubeClient       clientset.Interface

	MWInjectorLister mwilisters.MWInjectorLister
	MWInjectorSynced cache.InformerSynced

	NamespaceLister namespacelisters.NamespaceLister
	NamespaceSynced cache.InformerSynced // returns true if Namespaces has been synced. Added as member for testing

	updateHandler func(*mwiapi.MWInjector) error // callback to upate MWInjector. Added as member for testing

	queue workqueue.RateLimitingInterface // MWInjectors to be synced

	Recorder record.EventRecorder

	config MWInjectorControllerConfig
}

// NewMWInjectorController creates and initializes the MWInjectorController instance
func NewMWInjectorController(
	mwinjectorClient mwiclientset.Interface,
	kubeClient clientset.Interface,
	kubeInformerFactory kubeinformers.SharedInformerFactory,
	mwinjectorInformerFactory mwiinformers.SharedInformerFactory) *MWInjectorController {

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: v1core.New(kubeClient.Core().RESTClient()).Events("")})

	namespaceInformer := kubeInformerFactory.Core().V1().Namespaces()
	mwinjectorInformer := mwinjectorInformerFactory.Mwinjector().V1().MWInjectors()

	wc := &MWInjectorController{
		MWInjectorClient: mwinjectorClient,
		KubeClient:       kubeClient,
		MWInjectorLister: mwinjectorInformer.Lister(),
		MWInjectorSynced: mwinjectorInformer.Informer().HasSynced,
		NamespaceLister:  namespaceInformer.Lister(),
		NamespaceSynced:  namespaceInformer.Informer().HasSynced,

		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "mwinjector"),
		config: MWInjectorControllerConfig{
			RemoveIvalidMWInjector: true,
			NumberOfThreads:        1,
		},

		Recorder: eventBroadcaster.NewRecorder(scheme.Scheme, apiv1.EventSource{Component: "mwinjector-controller"}),
	}
	mwinjectorInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    wc.onAddMWInjector,
			UpdateFunc: wc.onUpdateMWInjector,
			DeleteFunc: wc.onDeleteMWInjector,
		},
	)

	namespaceInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    wc.onAddNamespace,
			UpdateFunc: wc.onUpdateNamespace,
			DeleteFunc: wc.onDeleteNamespace,
		},
	)
	wc.updateHandler = wc.updateMWInjector

	return wc
}

// Run simply runs the controller. Assuming Informers are already running: via factories
func (w *MWInjectorController) Run(ctx context.Context) error {
	glog.Infof("Starting mwinjector controller")

	if !cache.WaitForCacheSync(ctx.Done(), w.NamespaceSynced, w.MWInjectorSynced) {
		return fmt.Errorf("Timed out waiting for caches to sync")
	}

	for i := 0; i < w.config.NumberOfThreads; i++ {
		go wait.Until(w.runWorker, time.Second, ctx.Done())
	}

	<-ctx.Done()
	return ctx.Err()
}

func (w *MWInjectorController) runWorker() {
	for w.processNextItem() {
	}
}

func (w *MWInjectorController) processNextItem() bool {
	key, quit := w.queue.Get()
	if quit {
		return false
	}
	defer w.queue.Done(key)
	err := w.sync(key.(string))
	if err == nil {
		w.queue.Forget(key)
		return true
	}

	utilruntime.HandleError(fmt.Errorf("Error syncing mwinjector: %v", err))
	w.queue.AddRateLimited(key)

	return true
}

// enqueue adds key in the controller queue
func (w *MWInjectorController) enqueue(mwinjector *mwiapi.MWInjector) {
	key, err := cache.MetaNamespaceKeyFunc(mwinjector)
	if err != nil {
		glog.Errorf("MWInjectorController:enqueue: couldn't get key for MWInjector   %s/%s", mwinjector.Namespace, mwinjector.Name)
		return
	}
	w.queue.Add(key)
}

// mwinjector sync method
func (w *MWInjectorController) sync(key string) error {
	startTime := time.Now()
	defer func() {
		glog.V(6).Infof("Finished syncing MWInjector %q (%v", key, time.Now().Sub(startTime))
	}()

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	glog.V(6).Infof("Syncing %s/%s", namespace, name)
	sharedMWInjector, err := w.MWInjectorLister.MWInjectors(namespace).Get(name)
	if err != nil {
		glog.Errorf("unable to get MWInjector %s/%s: %v. Maybe deleted", namespace, name, err)
		return nil
	}

	if !mwiapi.IsMWInjectorDefaulted(sharedMWInjector) {
		defaultedMWInjector := mwiapi.DefaultMWInjector(sharedMWInjector)
		if err := w.updateHandler(defaultedMWInjector); err != nil {
			glog.Errorf("MWInjectorController.sync unable to default MWInjector %s/%s: %v", namespace, name, err)
			return fmt.Errorf("unable to default MWInjector %s/%s: %v", namespace, name, err)
		}
		glog.V(6).Infof("MWInjectorController.sync Defaulted %s/%s", namespace, name)
		return nil
	}

	// Validation. We always revalidate MWInjector since any update must be re-checked.
	if errs := mwiapi.ValidateMWInjector(sharedMWInjector); errs != nil && len(errs) > 0 {
		glog.Errorf("MWInjectorController.sync MWInjector %s/%s not valid: %v", namespace, name, errs)
		if w.config.RemoveIvalidMWInjector {
			glog.Errorf("Invalid mwinjector %s/%s is going to be removed", namespace, name)
			if err := w.deleteMWInjector(namespace, name); err != nil {
				glog.Errorf("unable to delete invalid mwinjector %s/%s: %v", namespace, name, err)
				return fmt.Errorf("unable to delete invalid mwinjector %s/%s: %v", namespace, name, err)
			}
		}
		return nil
	}

	// TODO: add test the case of graceful deletion
	if sharedMWInjector.DeletionTimestamp != nil {
		return nil
	}

	mwinjector := sharedMWInjector.DeepCopy()

	return w.manageMWInjector(mwinjector)
}

func newCondition(conditionType mwiapi.MWInjectorConditionType, reason, message string) mwiapi.MWInjectorCondition {
	return mwiapi.MWInjectorCondition{
		Type:               conditionType,
		Status:             apiv1.ConditionTrue,
		LastProbeTime:      metav1.Now(),
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	}
}

func (w *MWInjectorController) onAddMWInjector(obj interface{}) {
	mwinjector, ok := obj.(*mwiapi.MWInjector)
	if !ok {
		glog.Errorf("adding mwinjector, expected mwinjector object. Got: %+v", obj)
		return
	}
	if !reflect.DeepEqual(mwinjector.Status, mwiapi.MWInjectorStatus{}) {
		glog.Errorf("mwinjector %s/%s created with non empty status. Going to be removed", mwinjector.Namespace, mwinjector.Name)
		if _, err := cache.MetaNamespaceKeyFunc(mwinjector); err != nil {
			glog.Errorf("couldn't get key for MWInjector (to be deleted) %s/%s: %v", mwinjector.Namespace, mwinjector.Name, err)
			return
		}
		// TODO: how to remove a mwinjector created with an invalid or even with a valid status. What in case of error for this delete?
		if err := w.deleteMWInjector(mwinjector.Namespace, mwinjector.Name); err != nil {
			glog.Errorf("unable to delete non empty status MWInjector %s/%s: %v. No retry will be performed.", mwinjector.Namespace, mwinjector.Name, err)
		}
		return
	}
	w.enqueue(mwinjector)
}

func (w *MWInjectorController) onUpdateMWInjector(oldObj, newObj interface{}) {
	mwinjector, ok := newObj.(*mwiapi.MWInjector)
	if !ok {
		glog.Errorf("Expected mwinjector object. Got: %+v", newObj)
		return
	}
	w.enqueue(mwinjector)
}

func (w *MWInjectorController) onDeleteMWInjector(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		glog.Errorf("Unable to get key for %#v: %v", obj, err)
		return
	}
	w.queue.Add(key)
}

func (w *MWInjectorController) updateMWInjector(wfl *mwiapi.MWInjector) error {
	if _, err := w.MWInjectorClient.MwinjectorV1().MWInjectors(wfl.Namespace).Update(wfl); err != nil {
		glog.V(6).Infof("MWInjector %s/%s updated", wfl.Namespace, wfl.Name)
	}
	return nil
}

func (w *MWInjectorController) deleteMWInjector(namespace, name string) error {
	if err := w.MWInjectorClient.MwinjectorV1().MWInjectors(namespace).Delete(name, nil); err != nil {
		return fmt.Errorf("unable to delete MWInjector %s/%s: %v", namespace, name, err)
	}

	glog.V(6).Infof("MWInjector %s/%s deleted", namespace, name)
	return nil
}

func (w *MWInjectorController) onAddNamespace(obj interface{}) {
	ns := obj.(*apiv1.Namespace)
	mwinjectors, err := w.getMWInjectorsFromNamespace(ns)
	if err != nil {
		glog.Errorf("unable to get mwinjectors from namespace %s: %v", ns.Name, err)
		return
	}
	for i := range mwinjectors {
		w.enqueue(mwinjectors[i])
	}
}

func (w *MWInjectorController) onUpdateNamespace(oldObj, newObj interface{}) {
	oldNs := oldObj.(*apiv1.Namespace)
	newNs := newObj.(*apiv1.Namespace)
	if oldNs.ResourceVersion == newNs.ResourceVersion { // Since periodic resync will send update events for all known namespaces.
		return
	}
	glog.V(6).Infof("onUpdateNamespace old=%v, cur=%v ", oldNs.Name, newNs.Name)
	mwinjectors, err := w.getMWInjectorsFromNamespace(newNs)
	if err != nil {
		glog.Errorf("MWInjectorController.onUpdateNamespace cannot get mwinjectors for namespace %s: %v", newNs.Name, newNs.Name, err)
		return
	}
	for i := range mwinjectors {
		w.enqueue(mwinjectors[i])
	}
}

func (w *MWInjectorController) onDeleteNamespace(obj interface{}) {
	ns, ok := obj.(*apiv1.Namespace)
	glog.V(6).Infof("onDeleteNamespace old=%v", ns.Name)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			glog.Errorf("Couldn't get object from tombstone %+v", obj)
			return
		}
		ns, ok = tombstone.Obj.(*apiv1.Namespace)
		if !ok {
			glog.Errorf("Tombstone contained object that is not a ns %+v", obj)
			return
		}
	}
	mwinjectors, err := w.getMWInjectorsFromNamespace(ns)
	if err != nil {
		glog.Errorf("MWInjectorController.onDeleteNamespace: %v", err)
		return
	}
	for i := range mwinjectors {
		w.enqueue(mwinjectors[i])
	}
}

func (w *MWInjectorController) getMWInjectorsFromNamespace(ns *apiv1.Namespace) ([]*mwiapi.MWInjector, error) {
	mwinjectors := []*mwiapi.MWInjector{}
	if len(ns.Labels) == 0 {
		return mwinjectors, fmt.Errorf("no mwinjectors found for namespace. Namespace %s has no labels", ns.Name)
	}
	mwinjectorList, err := w.MWInjectorLister.List(labels.Everything())
	if err != nil {
		return mwinjectors, fmt.Errorf("no mwinjectors found for namespace. Namespace %s. Cannot list mwinjectors: %v", ns.Name, err)
	}

	for i := range mwinjectorList {
		mwinjector := mwinjectorList[i]
		if labels.SelectorFromSet(mwinjector.Spec.NamespaceSelector.MatchLabels).Matches(labels.Set(ns.Labels)) {
			mwinjectors = append(mwinjectors, mwinjector)
		}
	}
	return mwinjectors, nil
}

func (w *MWInjectorController) manageMWInjector(mwi *mwiapi.MWInjector) error {
	return nil
}
