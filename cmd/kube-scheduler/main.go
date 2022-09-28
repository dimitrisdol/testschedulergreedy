package main

import (
	"github.com/dimitrisdol/testschedulergreedy/greedyquad"

	"k8s.io/klog/v2"
	sched "k8s.io/kubernetes/cmd/kube-scheduler/app"
	//"flag"
	//"io"
	//"time"
	//"os"
	//"path/filepath"
	//"k8s.io/client-go/tools/clientcmd"
	//clientset "github.com/ckatsak/acticrds-go/client/clientset/versioned"
	//informersacti "github.com/ckatsak/acticrds-go/client/informers/externalversions"
	
)

func main() {

	//flag.CommandLine.SetOutput(io.Discard)
	
	//initParseCommandLineFlags()
	
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
	
	actiInformerFactory := informersacti.NewSharedInformerFactory(actiClient, time.Second*30)

	greedyquadplugin := greedyquad.NewController(actiClient, actiInformerFactory.Acti().V1alpha1().ActiNodes())
	
	_ = greedyquadplugin */
	
	cmd := sched.NewSchedulerCommand(
		sched.WithPlugin(greedyquad.Name, greedyquad.New),
	)
	if err := cmd.Execute(); err != nil {
		klog.Fatalf("failed to execute %q: %v", greedyquad.Name, err)
	}
}
