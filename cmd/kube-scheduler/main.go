package main

import (
	"github.com/dimitrisdol/testschedulergreedy/greedyquad"

	"k8s.io/klog/v2"
	sched "k8s.io/kubernetes/cmd/kube-scheduler/app"
	//"flag"
	//"io"
	
)

/*func initParseCommandLineFlags() {
		authenticationkubeconfig := flag.String("authentication-kubeconfig", "~/etc/kubernetes/scheduler.conf", "path to authentication-kubeconfig file")
	
	flag.Parse()
	
	_ = authenticationkubeconfig
} */

func main() {

	//flag.CommandLine.SetOutput(io.Discard)
	
	//initParseCommandLineFlags()

	cmd := sched.NewSchedulerCommand(
		sched.WithPlugin(greedyquad.Name, greedyquad.New),
	)
	if err := cmd.Execute(); err != nil {
		klog.Fatalf("failed to execute %q: %v", greedyquad.Name, err)
	}
}
