[![Go Report Card](https://goreportcard.com/badge/github.com/thekubeworld/kubeproxy-testing)](https://goreportcard.com/report/github.com/thekubeworld/kubeproxy-testing)

# kubeproxy-testing
A repo for testing kube-proxy and related bits

**Current status**:  

[![iptables clusterip - CNI: kindnetd](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/iptables_kubeproxy_clusterip.yml/badge.svg)](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/iptables_kubeproxy_clusterip.yml)
[![ipvs clusterip - CNI: kindnetd](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/ipvs_kubeproxy_clusterip.yml/badge.svg)](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/ipvs_kubeproxy_clusterip.yml)

[![iptables nodeport - CNI: kindnetd](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/iptables_kubeproxy_nodeport.yml/badge.svg)](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/iptables_kubeproxy_nodeport.yml)
[![ipvs nodeport - CNI: kindnetd](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/ipvs_kubeproxy_nodeport.yml/badge.svg)](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/ipvs_kubeproxy_nodeport.yml)

[![iptables loadbalancer - CNI: kindnetd](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/iptables_kubeproxy_loadbalancer.yml/badge.svg)](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/iptables_kubeproxy_loadbalancer.yml)
[![ipvs loadbalancer - CNI: kindnetd](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/ipvs_kubeproxy_loadbalancer.yml/badge.svg)](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/ipvs_kubeproxy_loadbalancer.yml)

[![iptables kubeproxy externalname - CNI: kindnetd](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/iptables_kubeproxy_externalname.yml/badge.svg)](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/iptables_kubeproxy_externalname.yml)
[![ipvs kubeproxy externalname - CNI: kindnetd](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/ipvs_kubeproxy_externalname.yml/badge.svg)](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/ipvs_kubeproxy_externalname.yml)


To see the logs/output from each job, follow these steps:

- Select the artifact from the list above.
- Click on the workflow/job
- Go to Artifacts and download logs.zip. It will contain the output logs from the job.

**Looking to run the tests locally**?  

**1)** Download the [k8devel Go library for Kubernetes](https://github.com/thekubeworld/k8devel)
```
$ go get github.com/thekubeworld/k8devel
```

**2)** Clone this repo and run the tests
```
$ git clone https://github.com/thekubeworld/kubeproxy-testing
$ cd kubeproxy-testing
$ go run tests/clusterip/main.go 
```
