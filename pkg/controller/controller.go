package controller

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/golang/glog"

	batch "k8s.io/api/batch/v1"
	batchv2 "k8s.io/api/batch/v2alpha1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"

	kubeinformers "k8s.io/client-go/informers"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	namespacelisters "k8s.io/client-go/listers/batch/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
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
	NamespaceSynced cache.InformerSynced // returns true if job has been synced. Added as member for testing

	NamespaceControl NamespaceControlInterface

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
	mwinjectorInformer := mwinjectorInformerFactory.MWInjector().V1().MWInjectors()

	wc := &MWInjectorController{
		MWInjectorClient: mwinjectorClient,
		KubeClient:       kubeClient,
		MWInjectorLister: mwinjectorInformer.Lister(),
		MWInjectorSynced: mwinjectorInformer.Informer().HasSynced,
		NamespaceLister:  namespaceInformer.Lister(),
		NamespaceSynced:  namespaceInformer.Informer().HasSynced,

		queueNamespaces workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "mwinjector"),
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
		cache.Namespaces{
			AddFunc:    wc.onAddJob,
			UpdateFunc: wc.onUpdateJob,
			DeleteFunc: wc.onDeleteJob,
		},
	)
	wc.NamespaceControl = &MWInjectorNamespaceControl{kubeClient, wc.Recorder}
	wc.updateHandler = wc.updateMWInjector

	return wc
}

TO CONTINUE

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
		glog.Errorf("MWInjectorController.sync Worfklow %s/%s not valid: %v", namespace, name, errs)
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

	// Init status.StartTime
	if mwinjector.Status.StartTime == nil {
		mwinjector.Status.Statuses = make([]mwiapi.MWInjectorStepStatus, 0)
		now := metav1.Now()
		mwinjector.Status.StartTime = &now
		if err := w.updateHandler(mwinjector); err != nil {
			glog.Errorf("mwinjector %s/%s: unable init startTime: %v", namespace, name, err)
			return err
		}
		glog.V(4).Infof("mwinjector %s/%s: startTime updated", namespace, name)
		return nil
	}

	if pastActiveDeadline(mwinjector, startTime) {
		now := metav1.Now()
		mwinjector.Status.Conditions = append(mwinjector.Status.Conditions, newDeadlineExceededCondition())
		mwinjector.Status.CompletionTime = &now
		if err := w.updateHandler(mwinjector); err != nil {
			glog.Errorf("mwinjector %s/%s unable to set DeadlineExceeded: %v", mwinjector.ObjectMeta.Namespace, mwinjector.ObjectMeta.Name, err)
			return fmt.Errorf("unable to set DeadlineExceeded for MWInjector %s/%s: %v", mwinjector.ObjectMeta.Namespace, mwinjector.ObjectMeta.Name, err)
		}
		if err := w.deleteMWInjectorJobs(mwinjector); err != nil {
			glog.Errorf("mwinjector %s/%s: unable to cleanup jobs: %v", mwinjector.ObjectMeta.Namespace, mwinjector.ObjectMeta.Name, err)
			return fmt.Errorf("mwinjector %s/%s: unable to cleanup jobs: %v", mwinjector.ObjectMeta.Namespace, mwinjector.ObjectMeta.Name, err)
		}
		return nil
	}

	return w.manageMWInjector(mwinjector)
}

func newDeadlineExceededCondition() mwiapi.MWInjectorCondition {
	return newCondition(mwiapi.MWInjectorFailed, "DeadlineExceeded", "MWInjector was active longer than specified deadline")
}

func newCompletedCondition() mwiapi.MWInjectorCondition {
	return newCondition(mwiapi.MWInjectorComplete, "", "")
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
	if IsMWInjectorFinished(mwinjector) {
		glog.Warningf("Update event received on complete MWInjector: %s/%s", mwinjector.Namespace, mwinjector.Name)
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
	if _, err := w.MWInjectorClient.MWInjectorV1().MWInjectors(wfl.Namespace).Update(wfl); err != nil {
		glog.V(6).Infof("MWInjector %s/%s updated", wfl.Namespace, wfl.Name)
	}
	return nil
}

func (w *MWInjectorController) deleteMWInjector(namespace, name string) error {
	if err := w.MWInjectorClient.MWInjectorV1().MWInjectors(namespace).Delete(name, nil); err != nil {
		return fmt.Errorf("unable to delete MWInjector %s/%s: %v", namespace, name, err)
	}

	glog.V(6).Infof("MWInjector %s/%s deleted", namespace, name)
	return nil
}

func (w *MWInjectorController) onAddJob(obj interface{}) {
	job := obj.(*batch.Job)
	mwinjectors, err := w.getMWInjectorsFromJob(job)
	if err != nil {
		glog.Errorf("unable to get mwinjectors from job %s/%s: %v", job.Namespace, job.Name, err)
		return
	}
	for i := range mwinjectors {
		w.enqueue(mwinjectors[i])
	}
}

func (w *MWInjectorController) onUpdateJob(oldObj, newObj interface{}) {
	oldJob := oldObj.(*batch.Job)
	newJob := newObj.(*batch.Job)
	if oldJob.ResourceVersion == newJob.ResourceVersion { // Since periodic resync will send update events for all known jobs.
		return
	}
	glog.V(6).Infof("onUpdateJob old=%v, cur=%v ", oldJob.Name, newJob.Name)
	mwinjectors, err := w.getMWInjectorsFromJob(newJob)
	if err != nil {
		glog.Errorf("MWInjectorController.onUpdateJob cannot get mwinjectors for job %s/%s: %v", newJob.Namespace, newJob.Name, err)
		return
	}
	for i := range mwinjectors {
		w.enqueue(mwinjectors[i])
	}

	// TODO: in case of relabelling ?
	// TODO: in case of labelSelector relabelling?
}

func (w *MWInjectorController) onDeleteJob(obj interface{}) {
	job, ok := obj.(*batch.Job)
	glog.V(6).Infof("onDeleteJob old=%v", job.Name)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			glog.Errorf("Couldn't get object from tombstone %+v", obj)
			return
		}
		job, ok = tombstone.Obj.(*batch.Job)
		if !ok {
			glog.Errorf("Tombstone contained object that is not a job %+v", obj)
			return
		}
	}
	mwinjectors, err := w.getMWInjectorsFromJob(job)
	if err != nil {
		glog.Errorf("MWInjectorController.onDeleteJob: %v", err)
		return
	}
	for i := range mwinjectors {
		//ww.expectations.DeleteiontObjserved(keyFunc(mwinjectors[i])
		w.enqueue(mwinjectors[i])
	}
}

func (w *MWInjectorController) getMWInjectorsFromJob(job *batch.Job) ([]*mwiapi.MWInjector, error) {
	mwinjectors := []*mwiapi.MWInjector{}
	if len(job.Labels) == 0 {
		return mwinjectors, fmt.Errorf("no mwinjectors found for job. Job %s/%s has no labels", job.Namespace, job.Name)
	}
	mwinjectorList, err := w.MWInjectorLister.List(labels.Everything())
	if err != nil {
		return mwinjectors, fmt.Errorf("no mwinjectors found for job. Job %s/%s. Cannot list mwinjectors: %v", job.Namespace, job.Name, err)
	}
	//mwinjectorStore.List()
	for i := range mwinjectorList {
		mwinjector := mwinjectorList[i]
		if mwinjector.Namespace != job.Namespace {
			continue
		}
		if labels.SelectorFromSet(mwinjector.Spec.Selector.MatchLabels).Matches(labels.Set(job.Labels)) {
			mwinjectors = append(mwinjectors, mwinjector)
		}
	}
	return mwinjectors, nil
}

func pastActiveDeadline(mwinjector *mwiapi.MWInjector, now time.Time) bool {
	if mwinjector.Spec.ActiveDeadlineSeconds == nil || mwinjector.Status.StartTime == nil {
		return false
	}
	start := mwinjector.Status.StartTime.Time
	duration := now.Sub(start)
	allowedDuration := time.Duration(*mwinjector.Spec.ActiveDeadlineSeconds) * time.Second
	return duration >= allowedDuration
}

func (w *MWInjectorController) manageMWInjector(mwinjector *mwiapi.MWInjector) error {
	mwinjectorToBeUpdated := false
	mwinjectorComplete := true
	for i := range mwinjector.Spec.Steps {
		stepName := mwinjector.Spec.Steps[i].Name
		stepStatus := mwiapi.GetStepStatusByName(mwinjector, stepName)
		if stepStatus != nil && stepStatus.Complete {
			continue
		}

		mwinjectorComplete = false
		switch {
		case mwinjector.Spec.Steps[i].JobTemplate != nil: // Job step
			if w.manageMWInjectorJobStep(mwinjector, stepName, &(mwinjector.Spec.Steps[i])) {
				mwinjectorToBeUpdated = true
				break
			}
		case mwinjector.Spec.Steps[i].ExternalRef != nil: // TODO handle: external object reference
			if w.manageMWInjectorReferenceStep(mwinjector, stepName, &(mwinjector.Spec.Steps[i])) {
				mwinjectorToBeUpdated = true
				break
			}
		}
	}
	if mwinjectorComplete {
		now := metav1.Now()
		mwinjector.Status.Conditions = append(mwinjector.Status.Conditions, newCompletedCondition())
		mwinjector.Status.CompletionTime = &now
		glog.Infof("MWInjector %s/%s complete.", mwinjector.Namespace, mwinjector.Name)
		mwinjectorToBeUpdated = true
	}

	if mwinjectorToBeUpdated {
		if err := w.updateHandler(mwinjector); err != nil {
			utilruntime.HandleError(err)
			w.enqueue(mwinjector)
			return err
		}
	}
	return nil
}

func (w *MWInjectorController) manageMWInjectorJobStep(mwinjector *mwiapi.MWInjector, stepName string, step *mwiapi.MWInjectorStep) bool {
	mwinjectorUpdated := false
	for _, dependencyName := range step.Dependencies {
		dependencyStatus := GetStepStatusByName(mwinjector, dependencyName)
		if dependencyStatus == nil || !dependencyStatus.Complete {
			glog.V(4).Infof("MWInjector %s/%s: dependecy %q not satisfied for %q", mwinjector.Namespace, mwinjector.Name, dependencyName, stepName)
			return mwinjectorUpdated
		}
	}
	glog.V(6).Infof("MWInjector %s/%s: All dependecy satisfied for %q", mwinjector.Namespace, mwinjector.Name, stepName)
	jobs, err := w.retrieveJobsStep(mwinjector, step.JobTemplate, stepName)
	if err != nil {
		glog.Errorf("unable to retrieve step jobs for MWInjector %s/%s, step:%q: %v", mwinjector.Namespace, mwinjector.Name, stepName, err)
		w.enqueue(mwinjector)
		return mwinjectorUpdated
	}
	switch len(jobs) {
	case 0: // create job
		_, err := w.NamespaceControl.CreateJob(mwinjector.Namespace, step.JobTemplate, mwinjector, stepName)
		if err != nil {
			glog.Errorf("Couldn't create job: %v : %v", err, step.JobTemplate)
			w.enqueue(mwinjector)
			defer utilruntime.HandleError(err)
			return mwinjectorUpdated
		}
		//w.expectations.CreationObserved(key)
		mwinjectorUpdated = true
		glog.V(4).Infof("Job created for step %q", stepName)
	case 1: // update status
		job := jobs[0]
		reference, err := reference.GetReference(scheme.Scheme, job)
		if err != nil || reference == nil {
			glog.Errorf("Unable to get reference from %v: %v", job.Name, err)
			return false
		}
		jobFinished := IsJobFinished(job)
		stepStatus := GetStepStatusByName(mwinjector, stepName)
		if stepStatus == nil {
			mwinjector.Status.Statuses = append(mwinjector.Status.Statuses, mwiapi.MWInjectorStepStatus{
				Name:      stepName,
				Complete:  jobFinished,
				Reference: *reference})
			mwinjectorUpdated = true
		}
		if jobFinished {
			stepStatus.Complete = jobFinished
			glog.V(4).Infof("MWInjector %s/%s Job finished for step %q", mwinjector.Namespace, mwinjector.Name, stepName)
			mwinjectorUpdated = true
		}
	default: // reconciliate
		glog.Errorf("manageMWInjectorJobStep %v too many jobs reported... Need reconciliation", mwinjector.Name)
		return false
	}
	return mwinjectorUpdated
}

func (w *MWInjectorController) manageMWInjectorReferenceStep(mwinjector *mwiapi.MWInjector, stepName string, step *mwiapi.MWInjectorStep) bool {
	return true
}

func (w *MWInjectorController) retrieveJobsStep(mwinjector *mwiapi.MWInjector, template *batchv2.JobTemplateSpec, stepName string) ([]*batch.Job, error) {
	jobSelector := createMWInjectorJobLabelSelector(mwinjector, template, stepName)
	jobs, err := w.NamespaceLister.Jobs(mwinjector.Namespace).List(jobSelector)
	if err != nil {
		return jobs, err
	}
	return jobs, nil
}

func (w *MWInjectorController) deleteMWInjectorJobs(mwinjector *mwiapi.MWInjector) error {
	glog.V(6).Infof("deleting all jobs for mwinjector %s/%s", mwinjector.Namespace, mwinjector.Name)
	jobsSelector := inferrMWInjectorLabelSelectorForJobs(mwinjector)
	jobs, err := w.NamespaceLister.Jobs(mwinjector.Namespace).List(jobsSelector)
	if err != nil {
		return fmt.Errorf("mwinjector %s/%s, unable to retrieve jobs to remove: %v", mwinjector.Namespace, mwinjector.Name, err)
	}
	errs := []error{}
	for i := range jobs {
		jobToBeRemoved := jobs[i].DeepCopy()
		if IsJobFinished(jobToBeRemoved) { // don't remove already finished job
			glog.V(4).Info("skipping job %s since finished", jobToBeRemoved.Name)
			continue
		}
		if err := w.NamespaceControl.DeleteJob(jobToBeRemoved.Namespace, jobToBeRemoved.Name, jobToBeRemoved); err != nil {
			errs = append(errs, err)
		}
	}
	return utilerrors.NewAggregate(errs)
}
