# Greedy Quad Scheduler

This is a scheduler out-of-tree plugin in order to greedily asssign different types of applications to servers.

It implements the greedy algorithm of matching as many "good" quadruples between insensitive and passive applications before allowing "bad" fits to happen.

In order to run:

   1. cd to root directory and type *make*
   2. wait for the build to happen
   3. *optional* docker push the docker image to your hub if you want
   4. docker container exec and transfer the greedyquad-config.yml file to /etc/kubernetes and the kube-scheduler.yml file to /etc/kubernetes/manifests folder.
   5. validate the plugin with the proposed Pod template.
