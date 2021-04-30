# kubeproxy-testing
A repo for testing kube-proxy bits

[![iptables nodeport](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/iptables_kubeproxy_nodeport.yml/badge.svg)](https://github.com/thekubeworld/kubeproxy-testing/actions/workflows/iptables_kubeproxy_nodeport.yml)

**1)** Download the [k8devel Go library for Kubernetes](https://github.com/thekubeworld/k8devel)
```
$ go get github.com/thekubeworld/k8devel
```

**2)** Clone this repo and run the tests
```
$ git clone https://github.com/thekubeworld/kubeproxy-testing
$ cd kubeproxy-testing
$ go run tests/iptables/clusterip/main.go 
```

```
$ go run tests/clusterip/main.go 
INFO[2021-04-26 21:34:34] Finished logrus log format settings...       
INFO[2021-04-26 21:34:34]                                              
INFO[2021-04-26 21:34:34] kube-proxy tests has started...              
INFO[2021-04-26 21:34:34]                                              
INFO[2021-04-26 21:34:34] Test #1) Pod connect via kube-proxy to a service and Pod 
INFO[2021-04-26 21:34:34]                                              
INFO[2021-04-26 21:34:34]            POD                               
INFO[2021-04-26 21:34:34]             |                                
INFO[2021-04-26 21:34:34]          Traffic                             
INFO[2021-04-26 21:34:34]             |                                
INFO[2021-04-26 21:34:34]             |                                
INFO[2021-04-26 21:34:34]    +-------------------+                     
INFO[2021-04-26 21:34:34]    |     kube-proxy    |                     
INFO[2021-04-26 21:34:34]    +-------------------+                     
INFO[2021-04-26 21:34:34]        |           |                         
INFO[2021-04-26 21:34:34] +------------------------------+             
INFO[2021-04-26 21:34:34] |     Service (Cluster IP)     |             
INFO[2021-04-26 21:34:34] |+-----------------------------+             
INFO[2021-04-26 21:34:34] |       |           |          |             
INFO[2021-04-26 21:34:34] |   +-------+   +-------+      |             
INFO[2021-04-26 21:34:34] |   |  Pod  |   |  Pod  |      |             
INFO[2021-04-26 21:34:34] |   +-------+   +-------+      |             
INFO[2021-04-26 21:34:34] |                              |             
INFO[2021-04-26 21:34:34] | kubernetes cluster           |             
INFO[2021-04-26 21:34:34] +------------------------------+             
INFO[2021-04-26 21:34:34]                                              
INFO[2021-04-26 21:34:34] Executing command: [iptables -t nat -L -n -v] 
INFO[2021-04-26 21:34:35]                                              
INFO[2021-04-26 21:34:35] Creating deployment: kptestingnginxdeploymentdcneai 
INFO[2021-04-26 21:34:35] Created deployment: kptestingnginxdeploymentdcneai namespace: kptesting 
INFO[2021-04-26 21:34:35]                                              
INFO[2021-04-26 21:34:35] Creating service: kptestingservicedcneai namespace: kptesting 
INFO[2021-04-26 21:34:35] Created service: kptestingservicedcneai namespace: kptesting 
INFO[2021-04-26 21:34:35]                                              
INFO[2021-04-26 21:34:35] Executing command: [iptables -t nat -L -n -v] 
INFO[2021-04-26 21:34:35] Diffing /tmp/iptables.993384280 and /tmp/iptables.487252695 
INFO[2021-04-26 21:34:35] /usr/bin/diff -r -u -N /tmp/iptables.993384280 /tmp/iptables.487252695 
INFO[2021-04-26 21:34:35] --- /tmp/iptables.993384280	2021-04-26 21:34:35.104491490 -0400
+++ /tmp/iptables.487252695	2021-04-26 21:34:35.749519073 -0400
@@ -1,6 +1,6 @@
 Chain PREROUTING (policy ACCEPT 1 packets, 60 bytes)
  pkts bytes target     prot opt in     out     source               destination         
-  402 25060 KUBE-SERVICES  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* kubernetes service portals */
+  403 25120 KUBE-SERVICES  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* kubernetes service portals */
     0     0 DOCKER_OUTPUT  all  --  *      *       0.0.0.0/0            172.18.0.1          
 
 Chain INPUT (policy ACCEPT 1 packets, 60 bytes)
@@ -8,12 +8,12 @@
 
 Chain OUTPUT (policy ACCEPT 1 packets, 60 bytes)
  pkts bytes target     prot opt in     out     source               destination         
-  542 34577 KUBE-SERVICES  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* kubernetes service portals */
+  543 34637 KUBE-SERVICES  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* kubernetes service portals */
   150 11376 DOCKER_OUTPUT  all  --  *      *       0.0.0.0/0            172.18.0.1          
 
 Chain POSTROUTING (policy ACCEPT 1 packets, 60 bytes)
  pkts bytes target     prot opt in     out     source               destination         
-  577 37110 KUBE-POSTROUTING  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* kubernetes postrouting rules */
+  578 37170 KUBE-POSTROUTING  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* kubernetes postrouting rules */
     0     0 DOCKER_POSTROUTING  all  --  *      *       0.0.0.0/0            172.18.0.1          
   170 10544 KIND-MASQ-AGENT  all  --  *      *       0.0.0.0/0            0.0.0.0/0            ADDRTYPE match dst-type !LOCAL /* kind-masq-agent: ensure nat POSTROUTING directs all non-LOCAL destination traffic to our custom KIND-MASQ-AGENT chain */
 
@@ -812,12 +812,22 @@
 
 Chain KUBE-SERVICES (2 references)
  pkts bytes target     prot opt in     out     source               destination         
+    0     0 KUBE-MARK-MASQ  tcp  --  *      *      !10.244.0.0/16        10.96.72.104         /* kptesting/kptestingserviceegqoro cluster IP */ tcp dpt:80
+    0     0 KUBE-SVC-V4XUFF4MJCFIIXWZ  tcp  --  *      *       0.0.0.0/0            10.96.72.104         /* kptesting/kptestingserviceegqoro cluster IP */ tcp dpt:80
+    0     0 KUBE-MARK-MASQ  tcp  --  *      *      !10.244.0.0/16        10.96.230.2          /* kptesting/kptestingservicelrsfvu cluster IP */ tcp dpt:80
+    0     0 KUBE-SVC-QLITVQESK2JB2TVM  tcp  --  *      *       0.0.0.0/0            10.96.230.2          /* kptesting/kptestingservicelrsfvu cluster IP */ tcp dpt:80
     0     0 KUBE-MARK-MASQ  tcp  --  *      *      !10.244.0.0/16        10.96.72.83          /* kptesting/kptestingservicerehidd cluster IP */ tcp dpt:80
     0     0 KUBE-SVC-FIROCJWWWYHBVEOY  tcp  --  *      *       0.0.0.0/0            10.96.72.83          /* kptesting/kptestingservicerehidd cluster IP */ tcp dpt:80
     0     0 KUBE-MARK-MASQ  tcp  --  *      *      !10.244.0.0/16        10.96.72.174         /* kptesting/kptestingservicewysntb cluster IP */ tcp dpt:80
     0     0 KUBE-SVC-OX7VTBBXQY2DY7GQ  tcp  --  *      *       0.0.0.0/0            10.96.72.174         /* kptesting/kptestingservicewysntb cluster IP */ tcp dpt:80
     0     0 KUBE-MARK-MASQ  tcp  --  *      *      !10.244.0.0/16        10.96.44.213         /* kptesting/kptestingservicemdngts cluster IP */ tcp dpt:80
     0     0 KUBE-SVC-RGRW4VBJ5GWIUEZV  tcp  --  *      *       0.0.0.0/0            10.96.44.213         /* kptesting/kptestingservicemdngts cluster IP */ tcp dpt:80
+    0     0 KUBE-MARK-MASQ  tcp  --  *      *      !10.244.0.0/16        10.96.137.38         /* kptesting/kptestingserviceqpkiba cluster IP */ tcp dpt:80
+    0     0 KUBE-SVC-MPXFVIPVS4IVSHL7  tcp  --  *      *       0.0.0.0/0            10.96.137.38         /* kptesting/kptestingserviceqpkiba cluster IP */ tcp dpt:80
+    0     0 KUBE-MARK-MASQ  tcp  --  *      *      !10.244.0.0/16        10.96.190.176        /* kptesting/kptestingserviceycgfbc cluster IP */ tcp dpt:80
+    0     0 KUBE-SVC-ARJVIMYQ6ZXXZUD4  tcp  --  *      *       0.0.0.0/0            10.96.190.176        /* kptesting/kptestingserviceycgfbc cluster IP */ tcp dpt:80
+    0     0 KUBE-MARK-MASQ  tcp  --  *      *      !10.244.0.0/16        10.96.40.49          /* kptesting/kptestingservicejgakgl cluster IP */ tcp dpt:80
+    0     0 KUBE-SVC-JQBU5HKBFMVVIUWS  tcp  --  *      *       0.0.0.0/0            10.96.40.49          /* kptesting/kptestingservicejgakgl cluster IP */ tcp dpt:80
     0     0 KUBE-MARK-MASQ  tcp  --  *      *      !10.244.0.0/16        10.96.0.1            /* default/kubernetes:https cluster IP */ tcp dpt:443
     0     0 KUBE-SVC-NPX46M4PTMTKRN6Y  tcp  --  *      *       0.0.0.0/0            10.96.0.1            /* default/kubernetes:https cluster IP */ tcp dpt:443
     0     0 KUBE-MARK-MASQ  tcp  --  *      *      !10.244.0.0/16        10.96.0.10           /* kube-system/kube-dns:metrics cluster IP */ tcp dpt:9153
@@ -826,12 +836,6 @@
     0     0 KUBE-SVC-TCOU7JCQXEZGVUNU  udp  --  *      *       0.0.0.0/0            10.96.0.10           /* kube-system/kube-dns:dns cluster IP */ udp dpt:53
     0     0 KUBE-MARK-MASQ  tcp  --  *      *      !10.244.0.0/16        10.96.11.83          /* kptesting/kptestingservicekuzuwh cluster IP */ tcp dpt:80
     0     0 KUBE-SVC-M6WUAEDO7WUYC6D3  tcp  --  *      *       0.0.0.0/0            10.96.11.83          /* kptesting/kptestingservicekuzuwh cluster IP */ tcp dpt:80
-    0     0 KUBE-MARK-MASQ  tcp  --  *      *      !10.244.0.0/16        10.96.137.38         /* kptesting/kptestingserviceqpkiba cluster IP */ tcp dpt:80
-    0     0 KUBE-SVC-MPXFVIPVS4IVSHL7  tcp  --  *      *       0.0.0.0/0            10.96.137.38         /* kptesting/kptestingserviceqpkiba cluster IP */ tcp dpt:80
-    0     0 KUBE-MARK-MASQ  tcp  --  *      *      !10.244.0.0/16        10.96.190.176        /* kptesting/kptestingserviceycgfbc cluster IP */ tcp dpt:80
-    0     0 KUBE-SVC-ARJVIMYQ6ZXXZUD4  tcp  --  *      *       0.0.0.0/0            10.96.190.176        /* kptesting/kptestingserviceycgfbc cluster IP */ tcp dpt:80
-    0     0 KUBE-MARK-MASQ  tcp  --  *      *      !10.244.0.0/16        10.96.40.49          /* kptesting/kptestingservicejgakgl cluster IP */ tcp dpt:80
-    0     0 KUBE-SVC-JQBU5HKBFMVVIUWS  tcp  --  *      *       0.0.0.0/0            10.96.40.49          /* kptesting/kptestingservicejgakgl cluster IP */ tcp dpt:80
     0     0 KUBE-MARK-MASQ  tcp  --  *      *      !10.244.0.0/16        10.96.0.10           /* kube-system/kube-dns:dns-tcp cluster IP */ tcp dpt:53
     0     0 KUBE-SVC-ERIFXISQEP7F7OF4  tcp  --  *      *       0.0.0.0/0            10.96.0.10           /* kube-system/kube-dns:dns-tcp cluster IP */ tcp dpt:53
     0     0 KUBE-MARK-MASQ  tcp  --  *      *      !10.244.0.0/16        10.96.115.181        /* kptesting/kptestingservicegwmkiw cluster IP */ tcp dpt:80
@@ -840,10 +844,6 @@
     0     0 KUBE-SVC-VA6KRVCWMSY7S4MY  tcp  --  *      *       0.0.0.0/0            10.96.205.231        /* kptesting/kptestingservicemfkonm cluster IP */ tcp dpt:80
     0     0 KUBE-MARK-MASQ  tcp  --  *      *      !10.244.0.0/16        10.96.193.59         /* kptesting/kptestingserviceaygojf cluster IP */ tcp dpt:80
     0     0 KUBE-SVC-NIUWNFLJG2W3FLXY  tcp  --  *      *       0.0.0.0/0            10.96.193.59         /* kptesting/kptestingserviceaygojf cluster IP */ tcp dpt:80
-    0     0 KUBE-MARK-MASQ  tcp  --  *      *      !10.244.0.0/16        10.96.72.104         /* kptesting/kptestingserviceegqoro cluster IP */ tcp dpt:80
-    0     0 KUBE-SVC-V4XUFF4MJCFIIXWZ  tcp  --  *      *       0.0.0.0/0            10.96.72.104         /* kptesting/kptestingserviceegqoro cluster IP */ tcp dpt:80
-    0     0 KUBE-MARK-MASQ  tcp  --  *      *      !10.244.0.0/16        10.96.230.2          /* kptesting/kptestingservicelrsfvu cluster IP */ tcp dpt:80
-    0     0 KUBE-SVC-QLITVQESK2JB2TVM  tcp  --  *      *       0.0.0.0/0            10.96.230.2          /* kptesting/kptestingservicelrsfvu cluster IP */ tcp dpt:80
     2   120 KUBE-NODEPORTS  all  --  *      *       0.0.0.0/0            0.0.0.0/0            /* kubernetes service nodeports; NOTE: this must be the last rule in this chain */ ADDRTYPE match dst-type LOCAL
 
 Chain KUBE-SVC-ARJVIMYQ6ZXXZUD4 (1 references)
INFO[2021-04-26 21:34:35]                                              
INFO[2021-04-26 21:34:35] Executing command: [curl http://10.96.197.161] 
INFO[2021-04-26 21:34:36] <!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
 /html>
[Test #1 PASSED]
```
