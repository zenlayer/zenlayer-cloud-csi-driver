# Zenlayer cloud csi driver
[![CI Status]][CI Result]


## Description
* The Zenlayer CSI plugin implements the interface between a Container Storage Interface ([CSI](https://github.com/container-storage-interface/))-enabled Container Orchestrator (CO) and Zenlayer storage.
* Version scheme: `v<major>.<minor>.<patch>`, for example `v1.2.0`.

## ZEC CSI and Kubernetes Version Matrix
| ZEC CSI Version | Container Orchestrator Name | Version Tested      |
| -----------------| --------------------------- | -------------------|
| v1.2.0          | Kubernetes                   |  v1.28.2 +         |

## ZEC CSI Feature Matrix
| ZEC CSI Version  | Feature                                                          |
| -----------------| -----------------------------------------------------------------|
| v1.2.0           | Create/Delete/Attach/Detach/Resize/Snapshot/Topology Volume      |

## External-csi-sidecar Version Description
| sidecar                             |    Current Version    |     Min CSI Spec Version  |       Container Image                                                        |       Min K8s Version    |   Recommended K8s Version     |
| ----------------------------------- | --------------------- | ------------------------- | ---------------------------------------------------------------------------- | ------------------------ | ----------------------------- |
| external-provisioner                | v5.2.0                | v1.0.0                    | registry.k8s.io/sig-storage/csi-provisioner:v5.2.0                           | v1.20                    | v1.29                         |
| external-attacher                   | v4.8.0                | v1.0.0                    | registry.k8s.io/sig-storage/csi-attacher:v4.8.0                              | v1.17                    | v1.29                         |
| external-resizer                    | v1.11.0               | v1.5.0                    | registry.k8s.io/sig-storage/csi-resizer:v1.11.0                              | v1.16                    | v1.29                         |
| node-driver-registrar               | v2.13.0               | v1.0.0                    | registry.k8s.io/sig-storage/csi-node-driver-registrar:v2.13.0                | v1.13                    | v1.25                         |
| livenessprobe                       | v2.15.0               | v1.0.0                    | registry.k8s.io/sig-storage/livenessprobe:v2.15.0                            | v1.13                    | -                             |
| external-snapshotter                | v8.2.0                | v1.11.0                   | registry.k8s.io/sig-storage/snapshot-controller&csi-snapshotter:v8.2.0       | v1.25                    | v1.25                         |

## Zenlayer OpenAPI SDK Version [API GitHub](https://github.com/zenlayer/zenlayercloud-sdk-go)
| ZEC CSI Version  | SDK Version                 |
| -----------------| --------------------------- |
| v1.2.0           | v0.2.49+                    |

## Helm Version [Helm Doc](./tools/helm.md)
| Helm Version     |
| -----------------|
|  v3.18.1+        |

# Disk CSI Driver
Disk CSI driver is available to help simplify storage management. Once a user creates a PVC with a reference to a Disk storage class, the disk and its corresponding PV object are dynamically created and become ready to be used by workloads.          

## How to Use

### Step 1: Prepare the Required Environment
* Authorization to access the related cloud resources. [console](https://console.zenlayer.com)        
* A working Kubernetes cluster deployed on ZEC VMs.         
* A local kubectl configured to communicate with this cluster.          

### Step 2: Install the CSI Driver
* If you only want to deploy the CSI plugin, refer to the [ZecCSI installation guide](./doc/install-guide.md) for detailed instructions. The CSI image and the chart repository are hosted on Docker Hub, so make sure the cluster has network connectivity to it. Zenlayer CSI can be installed entirely with Helm; there is no need to download the source code from GitHub unless you intend to develop against it.           

### Step 3: Create a StorageClass
A storage class is required for dynamic volume provisioning.       
Refer to the [StorageClass and topology configuration guide](./doc/storage-class.md) for detailed instructions.            

### Step 4: Check the Status of the CSI Driver
Check that all pods are running and ready.         
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

### Step 5: Test Workload Pod Using PVC
To make sure your CSI plugin is working, create a simple workload to test it out:           
```shell
kubectl apply -f deploy/simple-example/sc.yaml
kubectl apply -f deploy/simple-example/nginx-statefulset.yaml
kubectl get pvc
kubectl get pv
kubectl delete -f deploy/simple-example/nginx-statefulset.yaml
```

### Step 6: Test Disk Expansion
Choose a PVC and increase `spec.resources.requests.storage`:           
```shell
kubectl get pvc -o wide
kubectl edit pvc nginx-data-nginx-statefulset-0
```

### Step 7: Test Snapshots
```shell
kubectl apply -f deploy/snapshot-example/sc.yaml
kubectl get vsclass
kubectl apply -f deploy/snapshot-example/OriginPvc.yaml
kubectl apply -f deploy/snapshot-example/CreateSnapFromExistPvc.yaml
kubectl get vs
kubectl get vsc
```

## Notice
* The logs of the zeccsi driver are persisted to `/var/log/zenlayer_csi_logsbackups_fluent.log`. This log file is neither rotated nor deleted automatically; it is appended to continuously.        
* By default, each Elastic Compute instance can mount only two cloud disks: one boot disk and one data disk. Because `chart/values.yaml` sets `maxVolume: 9`, you must first raise the `Disks_per_instance` quota in the console (Products -> Service Quotas -> Elastic Compute -> Disks_per_instance) to the value you need, up to a maximum of 10. Otherwise only one PV can be attached to a virtual machine.            
* An RWO (ReadWriteOnce) volume can be attached to only one node at a time. When the original node shuts down or becomes unreachable, any new node that wants to attach the volume must first make sure the volume has been safely detached from the original node. The zec-csi-driver takes a conservative approach and never migrates an RWO volume automatically, in order to protect your data. If a node shuts down and cannot be recovered, manual intervention is required: confirm that the original pod has terminated, detach the disk from the original node manually, and clean up the corresponding VolumeAttachment. The CSI driver must guarantee that a new node cannot access the volume at the same time as the old one, but it cannot always determine whether the original pod has really terminated. In practice, do not rely on automatic migration of RWO volumes when a node fails.            
* When a volume is staged on a node for the first time, the driver formats it with discard turned off: `-E nodiscard` for ext3 and ext4, `-K` for xfs. A newly created ZEC cloud disk is already empty, so discarding every block before the filesystem is written gains nothing, and it makes formatting take longer the larger the disk is. This applies only to the one-time `mkfs`; the discard/TRIM behaviour of the mounted filesystem is unaffected, and you can still request it through the StorageClass `mountOptions` (which needs `--set featureGates=enable_mount_opt` at install time).            

## Currently Unsupported Features
* v1.2.0 does not support volume cloning; it only supports creating a PVC from a snapshot (`dataSource.kind: VolumeSnapshot`).
* v1.2.0 does not support volume group snapshots.
* In v1.2.0, snapshots depend on their source PV. If a PV is deleted, the snapshots created from it are also deleted in the storage system, but the corresponding VolumeSnapshot (`vs`) and VolumeSnapshotContent (`vsc`) objects remain in the Kubernetes cluster. Those VolumeSnapshot objects are no longer usable, so you have to clean them up manually.


[CI Status]: https://github.com/zenlayer/zenlayer-cloud-csi-driver/actions/workflows/ci.yml/badge.svg?branch=main
[CI Result]: https://github.com/zenlayer/zenlayer-cloud-csi-driver/actions/workflows/ci.yml
