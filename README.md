# Zenlayer cloud csi driver
[![GoReportCard Widget]][GoReportCardResult]

English | [简体中文](./README-zh_CN.md)

## Description
* Zenlayer CSI plugin implements an interface between Container Storage Interface ([CSI](https://github.com/container-storage-interface/)) enabled Container Orchestrator (CO) and the storage of Zenlayer.
* Version Description: v 1(Stable Release).0(Patch).0(BugFix)

## zec-csi and k8s Matrix
| ZEC CSI Version | Container Orchestrator Name | Version Tested     |
| -----------------| --------------------------- | -------------------|
| v1.0.0          | Kubernetes                   |  v1.28.2 +            |

## zec-csi feature Matrix
| ZEC CSI Version | Feature                                          |
| -----------------| -------------------------------------------------|
| v1.0.0          | Create cloud disk/Attach to Vm/Online Resize     |

## helm version
| Helm Version     |
| -----------------|
|  v3.18.1 +        |

## csi sidecar version
| sidecar               |       Version         |
| --------------------- | --------------------- |
| csi-provisioner                | v5.2.0          |
| csi-attacher                   | v4.6.1          |
| csi-resizer                    | v1.8.0          |
| csi-node-driver-registrar      | v2.8.0          |
| livenessprobe                  | v2.8.0          |

## zenlayer openApi sdk version [API Github](https://github.com/zenlayer/zenlayercloud-sdk-go)
| SDK Version |
| -----------------|
| v0.1.27  +        |


# Disk CSI Driver
Disk CSI driver is available to help simplify storage management.             
Once user creates PVC with the reference to a Disk storage class, disk and corresponding PV object get dynamically created and become ready to be used by workloads.          

## Configuration Requirements
* Authorizations to access related cloud resources. [console](https://console.zenlayer.com)        
* k8s cluster deployed on zec vm. 

## How to Use

### Step 1: Prepare Requirements Environment
* A working Kubernetes cluster deployed on ZEC Vms.         
* Local `kubectl` configured to communicate with this cluster.          

### Step 2: Install the CSI driver
If you only want to deploy the csi plugin， Please refer to the [ZecCSI Installation guide](./doc/install-guide.md) for detailed instructions.                            
The csi image and chart repository are located on docker hub, make ensure network connection. zenlayer csi can be fully installed with helm and there is no need to download the source code from github unless you have development requirements.               

### Step 3: Create StorageClass
Storage class is necessary for dynamic volume provisioning.       
Please refer to the [Storage-class and Topology config guide](./doc/storage-class.md) for detailied intructions.            

### Step 4: Check the Status of CSI driver
Checks that all pods are running and ready.         
```shell
kubectl get pods -n kube-system -l app=csi-zecplugin
```
Expected output:
```
NAME                  READY   STATUS    RESTARTS   AGE
csi-zecplugin-2xxr9   3/3     Running   0          2m19s
```
```shell
kubectl get pods -n kube-system -l app=csi-zecplugin-provisioner
```
Expected output:
```
NAME                                         READY   STATUS    RESTARTS   AGE
csi-zecplugin-provisioner-5ccf9f74cc-h6qfj   5/5     Running   0          2m37s
```

### Step 5: Test WorkLoad Pod
To make sure your CSI plugin is working, create a simple workload to test it out:           
```shell
kubectl apply -f deploy/example/nginx-statefulset.yaml
```

Now check that pod `nginx` is running, a new disk is created, attached, formatted, and mounted into the new pod.            
You can check pvc cloud-disk on Zenlayer Console.           
```shell
kubectl get pvc
kubectl get pv
```

After you are done, remove the test workload:           
```shell
kubectl delete -f deploy/example/nginx-statefulset.yaml
```

### Step 6:Test Expand Disk
choose on pvc, modify spec:resources:requests:storage           
```shell
kubectl get pvc -o wide
kubectl edit pvc nginx-data-nginx-statefulset-0
```

[GoReportCard Widget]: https://goreportcard.com/badge/github.com/zenlayer/zenlayer-cloud-csi-driver
[GoReportCardResult]: https://goreportcard.com/report/github.com/zenlayer/zenlayer-cloud-csi-driver