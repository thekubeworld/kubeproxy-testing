[![iptables clusterip - CNI: kindnetd](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/iptables_kubeproxy_clusterip.yml/badge.svg)](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/iptables_kubeproxy_clusterip.yml)
[![iptables nodeport - CNI: kindnetd](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/iptables_kubeproxy_nodeport.yml/badge.svg)](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/iptables_kubeproxy_nodeport.yml)
[![ipvs nodeport - CNI: kindnetd](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/ipvs_kubeproxy_nodeport.yml/badge.svg)](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/ipvs_kubeproxy_nodeport.yml)
[![ipvs clusterip - CNI: kindnetd](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/ipvs_kubeproxy_clusterip.yml/badge.svg)](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/ipvs_kubeproxy_clusterip.yml)

# kubeproxy-testing
A repo for testing kube-proxy bits

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
