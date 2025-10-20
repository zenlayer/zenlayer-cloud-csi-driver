# Zenlayer cloud csi driver
[![GoReportCard Widget]][GoReportCardResult]


## Description
* Zenlayer CSI plugin implements an interface between Container Storage Interface ([CSI](https://github.com/container-storage-interface/)) enabled Container Orchestrator (CO) and the storage of Zenlayer.
* Version Description: v 1(Stable Release).0(Patch).0(BugFix).

## Zec-csi and Kubernetes Version Matrix
| ZEC CSI Version | Container Orchestrator Name | Version Tested      |
| -----------------| --------------------------- | -------------------|
| v1.0.0          | Kubernetes                   |  v1.28.2 +         |

## Zec-csi Feature Matrix
| ZEC CSI Version  | Feature                                                          |
| -----------------| -----------------------------------------------------------------|
| v1.0.0           | Create/Delete/Attach/Detach/Resize/Snapshot/Topology Volume      |

## External-csi-sidecar Version Description
| sidecar                             |    Current Version    |     Min CSI Spec Version  |       Container Image                                                        |       Min K8s Version    |   Recommended K8s Version     |
| ----------------------------------- | --------------------- | ------------------------- | ---------------------------------------------------------------------------- | ------------------------ | ----------------------------- |
| external-provisioner                | v5.2.0                | v1.0.0                    | registry.k8s.io/sig-storage/csi-provisioner:v5.2.0                           | v1.20                    | v1.29                         |
| external-attacher                   | v4.8.0                | v1.0.0                    | registry.k8s.io/sig-storage/csi-attacher:v4.8.0                              | v1.17                    | v1.29                         |
| external-resizer                    | v1.11.0               | v1.5.0                    | registry.k8s.io/sig-storage/csi-resizer:v1.11.0                              | v1.16                    | v1.29                         |
| node-driver-registrar               | v2.13.0               | v1.0.0                    | registry.k8s.io/sig-storage/csi-node-driver-registrar:v2.13.0                | v1.13                    | v1.25                         |
| livenessprobe                       | v2.15.0               | v1.0.0                    | registry.k8s.io/sig-storage/livenessprobe:v2.15.0                            | v1.13                    | -                             |
| external-snapshotter                | v8.2.0                | v1.11.0                   | registry.k8s.io/sig-storage/snapshot-controller&csi-snapshotter:v8.2.0       | v1.25                    | v1.25                         |

## Zenlayer-openApi-sdk Version [API Github](https://github.com/zenlayer/zenlayercloud-sdk-go)
| ZEC CSI Version  | SDK Version                 |
| -----------------| --------------------------- |
| v1.0.0           | v0.2.0+                    |

## Helm Version [Help Doc](./tools/helm.md)
| Helm Version     |
| -----------------|
|  v3.18.1+        |

# Disk CSI Driver
Disk CSI driver is available to help simplify storage management.Once user creates PVC with the reference to a Disk storage class, disk and corresponding PV object get dynamically created and become ready to be used by workloads.          

## How to Use

### Step 1: Prepare Requirements Environment
* Authorizations to access related cloud resources. [console](https://console.zenlayer.com)        
* A working Kubernetes cluster deployed on ZEC Vms.         
* Local kubectl configured to communicate with this cluster.          

### Step 2: Install the CSI driver
* If you only want to deploy the csi plugin， Please refer to the [ZecCSI Installation guide](./doc/install-guide.md) for detailed instructions.  (The csi image and chart repository are located on docker hub, make ensure network connection. zenlayer csi can be fully installed with helm and there is no need to download the source code from github unless you have development requirements)           

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
csi-zecplugin-provisioner-678df8c5f4-9dhcx   7/7     Running   0          16s
```
```shell
kubectl get pods -n kube-system -l app=csi-fluent-bit
```
Expected output:
```
NAME                  READY    STATUS    RESTARTS   AGE
csi-fluent-bit-5gccr   1/1     Running   0          2m28s
```

### Step 5: Test WorkLoad Pod Use PVC
To make sure your CSI plugin is working, create a simple workload to test it out:           
```shell
kubectl apply -f deploy/simple-example/sc.yaml
kubectl apply -f deploy/simple-example/nginx-statefulset.yaml
kubectl get pvc
kubectl get pv
kubectl delete -f deploy/simple-example/nginx-statefulset.yaml
```

### Step 6: Test Expand Disk
choose on pvc, modify spec:resources:requests:storage           
```shell
kubectl get pvc -o wide
kubectl edit pvc nginx-data-nginx-statefulset-0
```

## Step 7: Test SnapShot
```shell
kubectl apply -f deploy/snapshot-example/sc.yaml
kubectl get vsclass
kubectl apply -f deploy/snapshot-example/OriginPvc.yaml
kubectl apply -f deploy/snapshot-example/CreateSnapFromExistPvc.yaml
kubectl get vs
kubectl get vsc
```

## Notice
* The logs of zeccsi driver are persisted to /var/log/zenlayer_csi_logsbackups_fluent.log. This log file will not be rotate and will not be automatically deleted. It will be continuously appended.

## Not supported feature Now
* v1.0.0 do not support clone volume, only support create pvc from snapshot (dataSource:kind:VolumeSnapshot).
* v1.0.0 do not support volumegroupsnapshots.
* v1.0.0 Snapshots rely on pv， If pv is deleted, the snapshot created with this pv will be automatically deleted in the storage system, so will leaving vs and vsc resource in the k8s cluster. The vs resource is no longer available, you need to clean it up manually.


[GoReportCard Widget]: https://goreportcard.com/badge/github.com/zenlayer/zenlayer-cloud-csi-driver
[GoReportCardResult]: https://goreportcard.com/report/github.com/zenlayer/zenlayer-cloud-csi-driver