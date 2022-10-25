// Package greedyquad contains an out-of-tree plugin based on the Kubernetes
// scheduling framework.
package greedyquad

import (
	"context"
	"fmt"
	"time"
	//"flag"
	//"os"
	//"path/filepath"

	"github.com/dimitrisdol/testschedulergreedy/greedyquad/hardcoded"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	//"k8s.io/apiserver/pkg/server"
	"k8s.io/apimachinery/pkg/util/wait"
	//apierrs "k8s.io/apimachinery/pkg/api/errors"
	
	//"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	//"k8s.io/client-go/tools/record"
	//"k8s.io/client-go/kubernetes/scheme"
	//"k8s.io/client-go/rest"
	
	//v1alpha1 "github.com/ckatsak/acticrds-go/apis/acti.cslab.ece.ntua.gr/v1alpha1"
	//"k8s.io/apimachinery/pkg/selection"
	//"k8s.io/apimachinery/pkg/labels"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/api/errors"
	listers "github.com/ckatsak/acticrds-go/client/listers/acti.cslab.ece.ntua.gr/v1alpha1"
	informers "github.com/ckatsak/acticrds-go/client/informers/externalversions/acti.cslab.ece.ntua.gr/v1alpha1"
	//informersacti "github.com/ckatsak/acticrds-go/client/informers/externalversions"
	clientset "github.com/ckatsak/acticrds-go/client/clientset/versioned"
	//actischeme "github.com/ckatsak/acticrds-go/client/clientset/versioned/scheme"
	//acticlient "github.com/ckatsak/acticrds-go/client/clientset/versioned/typed/acti.cslab.ece.ntua.gr/v1alpha1"
	//utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

const (
	// Name is the "official" external name of this scheduling plugin.
	Name = "GreedyQuadPlugin"

	// sla is the maximum slowdown that is allowed for an application when
	// it is being scheduled along another one.
	sla = 60
	
	// SuccessSynced is used as part of the Event 'reason' when a Foo is synced
	//SuccessSynced = "Synced"
	
	// MessageResourceSynced is the message used for an Event fired when a Foo
	// is synced successfully
	//MessageResourceSynced = "Foo synced successfully"

	// greedyquadLabelKey is the key of the Kubernetes Label which every
	// application that needs to be tracked by GreedyQuadPlugin should have.
	greedyquadLabelKey = "category"
)

// GreedyQuadPlugin is an out-of-tree plugin for the kube-scheduler, which takes into
// account information about the slowdown of colocated applications when they
// are wrapped into Pods and scheduled on the Kubernetes cluster.
type GreedyQuadPlugin struct {
	handle framework.Handle
	model  InterferenceModel
	ActinodesLister  listers.ActiNodeLister
	//actinodeInformer informers.ActiNodeInformer
	Acticlientset clientset.Interface
	ActinodesListerSynced cache.InformerSynced
	ActiQueue workqueue.RateLimitingInterface
	//recorder record.EventRecorder
}

var (

	_ framework.Plugin          = &GreedyQuadPlugin{}
	_ framework.FilterPlugin    = &GreedyQuadPlugin{}
	_ framework.ScorePlugin     = &GreedyQuadPlugin{}
	_ framework.ScoreExtensions = &GreedyQuadPlugin{}
)

// New instantiates a GreedyQuadPlugin.
func New(configuration runtime.Object, f framework.Handle) (framework.Plugin, error) {

	return &GreedyQuadPlugin{
		handle: f,
		model:  hardcoded.New(greedyquadLabelKey,),
	}, nil
}

// Name returns the official name of the GreedyQuadPlugin.
func (_ *GreedyQuadPlugin) Name() string {
	return Name
}

func NewController(
	Acticlientset clientset.Interface,
	ActiInformer informers.ActiNodeInformer,) *GreedyQuadPlugin {
	
	//utilruntime.Must(actischeme.AddToScheme(scheme.Scheme))
	//klog.V(4).Info("Creating event broadcaster")
	//eventBroadcaster := record.NewBroadcaster()
	//eventBroadcaster.StartStructuredLogging(0)
	//eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: acticlientset.ActiV1alpha1()})
	//recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: Name})
	
	Greedyquadplugin := &GreedyQuadPlugin{
		Acticlientset: Acticlientset,
		ActinodesLister: ActiInformer.Lister(),
		ActiQueue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Controller acti scheduling-queue"),
		//recorder:          recorder,
		ActinodesListerSynced: ActiInformer.Informer().HasSynced,
		}
		
	klog.Info("Setting up event handlers")
		
	/*greedyquadplugin.actinodesLister = actiInformer.Lister()
	greedyquadplugin.actinodesListerSynced = actiInformer.Informer().HasSynced
	greedyquadplugin.acticlientset = acticlientset
	*/
	
	return Greedyquadplugin
	}
/*	
func startActi() {

	//ctx := context.Background()
	
	home, exists := os.LookupEnv("TEST")
	if !exists {
	   home = "/etc/kubernetes"
	}
	
	configPath := filepath.Join(home, "scheduler.conf")
	
	cfg, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		klog.Fatalf("Error building config: %s", err.Error())
	}
	
	actiClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err.Error())
	}
	
	actiInformerFactory := informersacti.NewSharedInformerFactory(actiClient, 0)
	stopCh := server.SetupSignalHandler()

	greedyquadplugin := NewController(actiClient, actiInformerFactory.Acti().V1alpha1().ActiNodes())
	
	actiInformerFactory.Start(stopCh)
	
	if err = greedyquadplugin.Run(1); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}

}
	*/
func (ap *GreedyQuadPlugin) Run(workers int, stopCh <-chan struct{}) error {
	//defer utilruntime.HandleCrash()
	defer ap.ActiQueue.ShutDown()
	
	klog.Info("Starting scheduling controller")
	
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, ap.ActinodesListerSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch two workers to process Acti resources
	for i := 0; i < workers; i++ {
		go wait.Until(ap.runWorker, 4 * time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
} 


// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (ap *GreedyQuadPlugin) runWorker() {
	for ap.processNextWorkItem() {
	}
}

func (ap *GreedyQuadPlugin) processNextWorkItem() bool {
	obj, shutdown := ap.ActiQueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer ap.ActiQueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			ap.ActiQueue.Forget(obj)
			klog.Errorf("expected string in workqueue but got %#v", obj)
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Foo resource to be synced.
		/*if err := ap.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			ap.actiQueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		} */
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		ap.ActiQueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		klog.Errorf("error", err)
		return true
	}

	return true
}

/*
// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Foo resource
// with the current status of the resource.
func (ap *GreedyQuadPlugin) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		klog.Errorf("invalid resource key: %s", key)
		return nil
	}

	// Get the Foo resource with this namespace/name
	acti, err := ap.actinodesLister.ActiNodes(namespace).Get(name)
	if err != nil {
		// The Foo resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			klog.Errorf("foo '%s' in work queue no longer exists", key)
			return nil
		}

		return err
	}
	
	klog.Infof("Here is the actinode", acti)
	//ap.recorder.Event(acti, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}
*/

// enqueueActiNode takes an ActiNode resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than ActiNode.
func (ap *GreedyQuadPlugin) EnqueueActiNode(obj interface{}) {
	
	key := fmt.Sprint(obj)
	klog.Info("here is the key to be added", key)
	/*var key string
	
	if key, ok = string(obj); !ok {
			klog.Errorf("expected string in workqueue but got %#v", obj)
			return
		} */
	//var err error
	//if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
	//	klog.Errorf("error", err)
	//	return
	//}
	ap.ActiQueue.Add(key)
}


/*func (ap *GreedyQuadPlugin) sync() {

	//gp := startActi

	keyObj, quit := ap.actiQueue.Get()
	if quit {
		return
	}
	defer ap.actiQueue.Done(keyObj)
	
	key := keyObj.(string)
	namespace, actiName, err := cache.SplitMetaNamespaceKey(key)
	klog.V(4).Infof("Started acti processing %q", actiName)
	
	//get acti to process
	acti, err := ap.actinodesLister.ActiNodes(namespace).Get(actiName)
	klog.V(2).Infof("HERE IS THE ACTINODE ", acti)
	ctx := context.TODO()
	if err != nil {
		if apierrs.IsNotFound(err) {
			acti, err = ap.acticlientset.ActiV1alpha1().ActiNodes(namespace).Get(ctx, actiName, metav1.GetOptions{})
			if err != nil && apierrs.IsNotFound(err) {
				//acti was deleted in the meantime, ignore.
				klog.V(3).Infof("Acti %q deleted", actiName)
				return
			}
		}
		klog.Errorf("Error getting ActiNode %q: %v", actiName, err)
		ap.actiQueue.AddRateLimited(keyObj)
		return
	}
} */

// findCurrentOccupants returns all Pods that are being tracked by GreedyQuadPlugin
// and are already scheduled on the Node represented by the given NodeInfo.
//
// NOTE: For now, the number of the returned Pods should *always* be at most 4;
// otherwise, there must be some error in our scheduling logic.
func (_ *GreedyQuadPlugin) findCurrentOccupants(nodeInfo *framework.NodeInfo) []*corev1.Pod {
	ret := make([]*corev1.Pod, 0, 4)
	for _, podInfo := range nodeInfo.Pods {
		for key := range podInfo.Pod.Labels {
			if greedyquadLabelKey == key {
				ret = append(ret, podInfo.Pod)
			}
		}
	}
	return ret
}

// Filter is called by the scheduling framework.
//
// All FilterPlugins should return "Success" to declare that
// the given node fits the pod. If Filter doesn't return "Success",
// it will return "Unschedulable", "UnschedulableAndUnresolvable" or "Error".
//
// For the node being evaluated, Filter plugins should look at the passed
// nodeInfo reference for this particular node's information (e.g., pods
// considered to be running on the node) instead of looking it up in the
// NodeInfoSnapshot because we don't guarantee that they will be the same.
//
// For example, during preemption, we may pass a copy of the original
// nodeInfo object that has some pods removed from it to evaluate the
// possibility of preempting them to schedule the target pod.
func (ap *GreedyQuadPlugin) Filter(
	ctx context.Context,
	state *framework.CycleState,
	pod *corev1.Pod,
	nodeInfo *framework.NodeInfo,
	//actinodeInformer informers.ActiNodeInformer,
) *framework.Status {
	if pod == nil {
		return framework.NewStatus(framework.Error, "pod cannot be nil")
	}
	if nodeInfo.Node() == nil {
		return framework.NewStatus(framework.Error, "node cannot be nil")
	}
	nodeName := nodeInfo.Node().Name

	// If the given Pod does not have the greedyquadLabelKey, approve it and let
	// the other plugins decide for its fate.
	if _, exists := pod.Labels[greedyquadLabelKey]; !exists {
		klog.V(2).Infof("blindly approving Pod '%s/%s' as it does not have GreedyQuadPlugin's label %q", pod.Namespace, pod.Name, greedyquadLabelKey)
		return framework.NewStatus(framework.Success, "pod is not tracked by GreedyQuadPlugin")
	}

	// For the Node at hand, find all occupant Pods tracked by GreedyQuadPlugin.
	// These should *always* be fewer than or equal to 4, but we take the
	// opportunity to assert this invariant later anyway.
	occupants := ap.findCurrentOccupants(nodeInfo)
		
	//actinodelister : listers.ActiNodeLister{}
	//sel := labels.NewSelector()
	
	//req, err := labels.NewRequirement("component", selection.Equals, []string{"actinodes"})
	//if err != nil {
	//	panic(err.Error())
	//}
	//sel = sel.Add(*req)
	
	//klog.V(2).Infof("HERE IS THE SELECTOR '%s'", sel)
	
	//actinode := v1alpha1.ActiNode{
	//	Spec: v1alpha1.ActiNodeSpec
	//	        }
	
	//sel, err := metav1.LabelSelectorAsSelector(actinode.Spec)
	//if err != nul {
	//	panic(errr.Error())
	//}
	
	/*home, exists := os.LookupEnv("TEST")
	if !exists {
	   home = "/etc/kubernetes"
	}
	
	configPath := filepath.Join(home, "scheduler.conf")
	
	cfg, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		klog.Fatalf("Error building config: %s", err.Error())
	}
	
	actiClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err.Error())
	}
	
	actiInformerFactory := informersacti.NewSharedInformerFactory(actiClient, 0)
	stopCh := server.SetupSignalHandler()

	greedyquadplugin := NewController(actiClient, actiInformerFactory.Acti().V1alpha1().ActiNodes())
	
	actiInformerFactory.Start(stopCh)
	run := func(ctx context.Context) {
		greedyquadplugin.Run(1, ctx.Done())
		}
	run(ctx)
	//name := "termi7"
	//klog.V(2).Infof("HERE IS THE NAME '%s'", name)
	//klog.V(2).Infof("HERE IS THE NAMESPACE '%s'", pod.Namespace)
	
	<-stopCh */
	
	/*name := "termi7"
	klog.Info("HERE IS THE NAME '%s'", name)
	klog.Info("HERE IS THE NAMESPACE '%s'", pod.Namespace)
	
	actiactilister := ap.actinodesLister
	actinamespacer := actiactilister.ActiNodes(pod.Namespace)
	
	actinodes, err := actinamespacer.Get(name)
	if err != nil {
		klog.Info("ignoring orphaned object of acti")
		}
	ap.enqueueActiNode(actinodes) */
	
	//_ : ap.actinodeInformer.Lister()
	
	/*sel := labels.NewSelector()
	
	req, err := labels.NewRequirement("namespace", selection.Equals, []string{"acti-ns"})
	if err != nil {
		panic(err.Error())
	}
	sel = sel.Add(*req)
	
	*/
	
	/*sel := labels.SelectorFromSet(map[string]string{"app.kubernetes.o/component": "actinodes"})
	
	klog.V(2).Infof("HERE IS THE SELECTOR '%s'", sel)
	*/
	/*actiactilister := greedyquadplugin.actinodesLister
	actinamespacer := actiactilister.ActiNodes(pod.Namespace)
	
	actinodes, err := actinamespacer.Get(name)
	*/
	
	// !!!!!//
	
	/*actinodes, err := greedyquadplugin.acticlientset.ActiV1alpha1().ActiNodes(pod.Namespace).Get(context.TODO(), name, metav1.GetOptions{})
	//actinodes, err := ap.actinodesLister.ActiNodes(pod.Namespace).Get(name)
	if err != nil {
		fmt.Print(err.Error())
	}
	klog.V(2).Infof("HERE IS THE ACTINODE ", actinodes) */
	

	// Decide on how to proceed based on the number of current occupants
	switch len(occupants) {
	// If the Node is full (i.e., 4 applications tracked by GreedyQuadPlugin are
	// already scheduled on it), filter it out.
	case 4:
			klog.V(2).Infof("filtering Node %q out because 4 GreedyQuadPlugin applications are already scheduled there", nodeName)
		return framework.NewStatus(framework.Unschedulable, fmt.Sprintf("Node '%s' already has 2 GreedyQuadPlugin occupants", nodeName))
		
		
	// If the existing occupant is slowed down prohibitively much by the
	// new Pod's attack, filter the Node out.
	// Now cheking the cases of 1, 2 and 3 pods already being in the Node.
	case 3:
		occ1 := occupants[0] // the first, currently scheduled Pod
		occ2 := occupants[1] // second Pod
		occ3 := occupants[2] // third one 
		score1, err1 := ap.model.Attack(pod, occ1)
		if err1 != nil {
			err1 = fmt.Errorf("new Pod '%s/%s' on Node '%s': %v", occ1.Namespace, occ1.Name, nodeName, err1)
			klog.Warning(err1)
			return framework.NewStatus(framework.Error, err1.Error())
		}
		score2, err2 := ap.model.Attack(pod, occ2)
		if err2 != nil {
			err2 = fmt.Errorf("new Pod '%s/%s' on Node '%s': %v", occ2.Namespace, occ2.Name, nodeName, err2)
			klog.Warning(err2)
			return framework.NewStatus(framework.Error, err2.Error())
		}
		score3, err3 := ap.model.Attack(pod, occ3)
		if err3 != nil {
			err3 = fmt.Errorf("new Pod '%s/%s' on Node '%s': %v", occ3.Namespace, occ3.Name, nodeName, err3)
			klog.Warning(err3)
			return framework.NewStatus(framework.Error, err3.Error())
		}
		score := score1 + score2 + score3
		if score > sla {
			msg := fmt.Sprintf("filtering Node '%s': new pod '%s/%s' ('%s') incurs huge slowdown on pod '%s/%s' ('%s')",
				nodeName, pod.Namespace, pod.Name, pod.Labels[greedyquadLabelKey],)
			klog.V(2).Infof(msg)
			return framework.NewStatus(framework.Unschedulable, msg)
		}
		fallthrough
	case 2:
		occ1 := occupants[0] // the first scheduled Pod
		occ2 := occupants[1] // second one
		score1, err1 := ap.model.Attack(pod, occ1)
		if err1 != nil {
			err1 = fmt.Errorf("new Pod '%s/%s' on Node '%s': %v", occ1.Namespace, occ1.Name, nodeName, err1)
			klog.Warning(err1)
			return framework.NewStatus(framework.Error, err1.Error())
		}
		score2, err2 := ap.model.Attack(pod, occ2)
		if err2 != nil {
			err2 = fmt.Errorf("new Pod '%s/%s' on Node '%s': %v", occ2.Namespace, occ2.Name, nodeName, err2)
			klog.Warning(err2)
			return framework.NewStatus(framework.Error, err2.Error())
		}
		score := score1 + score2
		if score > sla {
			msg := fmt.Sprintf("filtering Node '%s': new pod '%s/%s' ('%s') incurs huge slowdown on pod '%s/%s' ('%s')",
				nodeName, pod.Namespace, pod.Name, pod.Labels[greedyquadLabelKey],)
			klog.V(2).Infof(msg)
			return framework.NewStatus(framework.Unschedulable, msg)
		}
		fallthrough
	case 1:
		occ := occupants[0] // the single, currently scheduled Pod
		score, err := ap.model.Attack(pod, occ)
		if err != nil {
			err = fmt.Errorf("new Pod '%s/%s' on Node '%s': %v", occ.Namespace, occ.Name, nodeName, err)
			klog.Warning(err)
			return framework.NewStatus(framework.Error, err.Error())
		}
		if score > sla {
			msg := fmt.Sprintf("filtering Node '%s': new pod '%s/%s' ('%s') incurs huge slowdown on pod '%s/%s' ('%s')",
				nodeName, pod.Namespace, pod.Name, pod.Labels[greedyquadLabelKey], occ.Namespace, occ.Name, occ.Labels[greedyquadLabelKey])
			klog.V(2).Infof(msg)
			return framework.NewStatus(framework.Unschedulable, msg)
		}
		fallthrough
	// If the Node is empty, or fell through from above (i.e., SLA allows
	// the single current occupant to be attacked by the new Pod), approve.
	case 0:
		klog.V(2).Infof("approving Node %q for pod '%s/%s'", nodeName, pod.Namespace, pod.Name)
		return framework.NewStatus(framework.Success)
	// If more than 2 occupants are found to be already scheduled on the
	// Node at hand, we must have fucked up earlier; report the error.
	default:
		klog.Errorf("detected %d occupant Pods tracked by GreedyQuadPlugin on Node %q", len(occupants), nodeName)
		return framework.NewStatus(framework.Error, fmt.Sprintf("found %d occupants on '%s' already", len(occupants), nodeName))
	}
}

// Score is called on each filtered node. It must return success and an integer
// indicating the rank of the node. All scoring plugins must return success or
// the pod will be rejected.
//
// In the case of GreedyQuadPlugin, scoring is reversed; i.e., higher score indicates
// worse scheduling decision.
// This is taken into account and "fixed" later, during the normalization.
func (ap *GreedyQuadPlugin) Score(
	ctx context.Context,
	state *framework.CycleState,
	p *corev1.Pod,
	nodeName string,
) (int64, *framework.Status) {
	// Retrieve the Node at hand from the cycle's snapshot
	nodeInfo, err := ap.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		return -1, framework.NewStatus(framework.Error, fmt.Sprintf("failed to get Node '%s' from snapshot: %v", nodeName, err))
	}
	
	// If the given Pod does not have the greedyquadLabelKey, approve it and let
	// the other plugins decide for its fate.
	if _, exists := p.Labels[greedyquadLabelKey]; !exists {
		klog.V(2).Infof("blindly scoring Pod '%s/%s' as it does not have GreedyQuadPlugin's label %q", p.Namespace, p.Name, greedyquadLabelKey)
		return 0, framework.NewStatus(framework.Success, fmt.Sprintf("Node '%s' is empty: interim score = 0", nodeName))
	}

	occupants := ap.findCurrentOccupants(nodeInfo)

	// If the Node is empty, for now, assume it is a perfect candidate.
	// Therefore, the scheduled applications are expected to tend to spread
	// among the available Nodes as much as possible.
	if len(occupants) == 0 {
		return 0, framework.NewStatus(framework.Success, fmt.Sprintf("Node '%s' is empty: interim score = 0", nodeName))
	}

	// Otherwise, evaluate the slowdown
	if len(occupants) == 1{
	occ := occupants[0]
	scoreFp, err := ap.model.Attack(p, occ)
	if err != nil {
		err = fmt.Errorf("new Pod '%s/%s' on Node '%s': %v", occ.Namespace, occ.Name, nodeName, err)
		klog.Warning(err)
		return -1, framework.NewStatus(framework.Error, err.Error())
	}
	score := int64(ap.model.ToInt64Multiplier() * scoreFp)
	return score, framework.NewStatus(framework.Success, fmt.Sprintf("Node '%s': interim score = %d", nodeName, score))
	}
	if len(occupants) == 2{
	occ1 := occupants[0]
	occ2 := occupants[1]
	scoreFp1, err1 := ap.model.Attack(p, occ1)
	if err1 != nil {
		err1 = fmt.Errorf("new Pod '%s/%s' on Node '%s': %v", occ1.Namespace, occ1.Name, nodeName, err1)
		klog.Warning(err1)
		return -1, framework.NewStatus(framework.Error, err1.Error())
	}
	scoreFp2, err2 := ap.model.Attack(p, occ2)
	if err2 != nil {
		err2 = fmt.Errorf("new Pod '%s/%s' on Node '%s': %v", occ2.Namespace, occ2.Name, nodeName, err2)
		klog.Warning(err2)
		return -1, framework.NewStatus(framework.Error, err2.Error())
	}
	scoreFp := scoreFp1 + scoreFp2
	score := int64(ap.model.ToInt64Multiplier() * scoreFp)
	return score, framework.NewStatus(framework.Success, fmt.Sprintf("Node '%s': interim score = %d", nodeName, score))
	}
	//3 occupants now
	occ1 := occupants[0]
	occ2 := occupants[1]
	occ3 := occupants[2]
	scoreFp1, err1 := ap.model.Attack(p, occ1)
	if err1 != nil {
		err1 = fmt.Errorf("new Pod '%s/%s' on Node '%s': %v", occ1.Namespace, occ1.Name, nodeName, err1)
		klog.Warning(err1)
		return -1, framework.NewStatus(framework.Error, err1.Error())
	}
	scoreFp2, err2 := ap.model.Attack(p, occ2)
	if err2 != nil {
		err2 = fmt.Errorf("new Pod '%s/%s' on Node '%s': %v", occ2.Namespace, occ2.Name, nodeName, err2)
		klog.Warning(err2)
		return -1, framework.NewStatus(framework.Error, err2.Error())
	}
	scoreFp3, err3 := ap.model.Attack(p, occ3)
	if err3 != nil {
		err3 = fmt.Errorf("new Pod '%s/%s' on Node '%s': %v", occ3.Namespace, occ3.Name, nodeName, err3)
		klog.Warning(err3)
		return -1, framework.NewStatus(framework.Error, err3.Error())
	}
	scoreFp := scoreFp1 + scoreFp2 + scoreFp3
	score := int64(ap.model.ToInt64Multiplier() * scoreFp)
	return score, framework.NewStatus(framework.Success, fmt.Sprintf("Node '%s': interim score = %d", nodeName, score))
	
}

// ScoreExtensions returns the GreedyQuadPlugin itself, since it implements the
// framework.ScoreExtensions interface.
func (ap *GreedyQuadPlugin) ScoreExtensions() framework.ScoreExtensions {
	return ap
}

// NormalizeScore is called for all node scores produced by the same plugin's
// "Score" method. A successful run of NormalizeScore will update the scores
// list and return a success status.
//
// In the case of the GreedyQuadPlugin, its "Score" method produces scores of reverse
// priority (i.e., the lower the score, the better the result). Therefore all
// scores have to be reversed during the normalization, so that higher score
// indicates a better scheduling result in terms of slowdowns.
func (_ *GreedyQuadPlugin) NormalizeScore(
	ctx context.Context,
	state *framework.CycleState,
	p *corev1.Pod,
	scores framework.NodeScoreList,
) *framework.Status {
	// Find the max score for the normalization
	var maxScore int64
	for i := range scores {
		if scores[i].Score > maxScore {
			maxScore = scores[i].Score
		}
	}
	// When no Pod (tracked by GreedyQuadPlugin) is scheduled on the Node,
	// maxScore will be 0.
	if maxScore == 0 {
		for i := range scores {
			scores[i].Score = framework.MaxNodeScore // reverse priority
		}
		return framework.NewStatus(framework.Success)
	}

	// Normalize them & reverse their priority
	for i := range scores {
		score := scores[i].Score                          // load
		score = framework.MaxNodeScore * score / maxScore // normalize
		score = framework.MaxNodeScore - score            // reverse priority
		scores[i].Score = score                           // store
	}
	return framework.NewStatus(framework.Success)
}
