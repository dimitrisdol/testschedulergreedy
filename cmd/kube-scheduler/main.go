package main

import (
	//"context"
	"github.com/dimitrisdol/testschedulergreedy/greedyquad"

	"k8s.io/klog/v2"
	sched "k8s.io/kubernetes/cmd/kube-scheduler/app"
	//"flag"
	//"io"
	//"time"
	//"os"
	//"path/filepath"
	//"k8s.io/apiserver/pkg/server"
	//"k8s.io/client-go/tools/clientcmd"
	//clientset "github.com/ckatsak/acticrds-go/client/clientset/versioned"
	//informersacti "github.com/ckatsak/acticrds-go/client/informers/externalversions"
	
)

func main() {
	
	cmd := sched.NewSchedulerCommand(
		sched.WithPlugin(greedyquad.Name, greedyquad.New),
	)
	if err := cmd.Execute(); err != nil {
		klog.Fatalf("failed to execute %q: %v", greedyquad.Name, err)
	}
}
