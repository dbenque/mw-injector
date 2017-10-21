package controller

import (
	"context"
	"time"

	"github.com/golang/glog"
	"k8s.io/client-go/informers"
	v1_informer "k8s.io/client-go/informers/core/v1"
	clientset "k8s.io/client-go/kubernetes"
)

type NamespaceController struct {
	KubeClient clientset.Interface
	informer   v1_informer.NamespaceInformer
}

type NamespaceEventHandler struct {
}

func (h *NamespaceEventHandler) OnAdd(obj interface{}) {

}
func (h *NamespaceEventHandler) OnUpdate(oldObj, newObj interface{}) {

}
func (h *NamespaceEventHandler) OnDelete(obj interface{}) {

}

/// NewNamespaceController creates and initializes the NamespaceController instance
func NewNamespaceController(kubeClient clientset.Interface) *NamespaceController {

	nc := &NamespaceController{
		KubeClient: kubeClient,
		informer:   informers.NewSharedInformerFactory(kubeClient, time.Minute).Core().V1().Namespaces(),
	}

	nc.informer.Informer().AddEventHandlerWithResyncPeriod(&NamespaceEventHandler{}, time.Minute)

	return nc
}

// Run simply runs the controller
func (nc *NamespaceController) Run(ctx context.Context) error {
	glog.Infof("Starting Namespace controller")

	// go .JobInformer.Run(ctx.Done())
	// go w.workflowInformer.Run(ctx.Done())

	// if !cache.WaitForCacheSync(ctx.Done(), w.JobInformer.HasSynced, w.workflowInformer.HasSynced) {
	// 	return fmt.Errorf("Timed out waiting for caches to sync")
	// }

	// for i := 0; i < w.config.NumberOfThreads; i++ {
	// 	go wait.Until(w.runWorker, time.Second, ctx.Done())
	// }

	<-ctx.Done()
	return ctx.Err()
}
