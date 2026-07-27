# The zeccsi plugin manages storage through StorageClasses

## StorageClass and PVC configuration
### StorageClass YAML description          
```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: csi-zec
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"

provisioner: disk.csi.zenlayer.com                       //Zenlayer CSI driver name. Cannot be changed.

parameters:                                              //None of the parameters is mandatory
  fsType: "ext4"                                         //Mount filesystem type. Supports ext4, ext3, or xfs. Defaults to ext4 if not set.
  type: "1"                                              //Cloud disk type: 1 = Basic NVMe SSD, 2 = Standard NVMe SSD. Defaults to Standard NVMe SSD if not set.
  zoneID: "asia-north-1a"                                //Cloud disk zone. Only takes effect when volumeBindingMode is Immediate. If it is not set here, you must specify it at install time with: helm install --set defaultZone=... --set defaultResourceGroup=...
  placeGroupID: "xxx"                                    //Resource group ID of the cloud disk in the Zenlayer console. If it is not set here, you must specify it at install time with: helm install --set defaultZone=... --set defaultResourceGroup=...
  burstEnable: "false"                                   //Whether to enable cloud disk QoS bursting: true or false.

reclaimPolicy: Delete                                    //Supports "Delete" and "Retain". Retain is not recommended: you then have to delete the cloud disk yourself, which may leave residual data behind.

allowVolumeExpansion: true                               //Supports online cloud disk expansion.

volumeBindingMode: Immediate                             //Supports "Immediate" and "WaitForFirstConsumer". In Immediate mode, the PV is created and bound right away; it is created in the zone given by zoneID, and pods that use this StorageClass are scheduled to that zone as well. In WaitForFirstConsumer mode, the PV is not created right away: the PVC stays Pending until a pod uses it, and the Kubernetes scheduler then decides where both the PV and the pod are placed.

mountOptions:                                            //This feature must be enabled at install time: helm install --set featureGates=enable_mount_opt
  - xxx
```
### PVC YAML description (without a snapshot)    
```yaml 
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pvc1
  namespace: default
spec:       
  accessModes:        
    - ReadWriteOncePod                                   //Supports ReadWriteOncePod and ReadWriteOnce.
  volumeMode: Filesystem                                 //Supports Filesystem and Block.
  resources:      
    requests:     
      storage: 80Gi                                      //Cloud disk size. Must be at least 20Gi.
  storageClassName: csi-zec                              //StorageClass name.
```

## VolumeSnapshotClass and VolumeSnapshot Configuration
### VolumeSnapshotClass YAML description          
```yaml
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshotClass
metadata:
  name: csi-zec-snap
driver: disk.csi.zenlayer.com                           //Zenlayer CSI driver name. Cannot be changed.
deletionPolicy: Delete
parameters:
  tags: value
```

### VolumeSnapshot YAML description
```yaml
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshot
metadata:
  name: pvc-origin-snap
spec:
  volumeSnapshotClassName: csi-zec-snap
  source:
    persistentVolumeClaimName: pvc-origin
```

### PVC YAML description (restoring from a snapshot)
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pvc-clone
spec:
  storageClassName: csi-zec                           //Target StorageClass of the new PVC.
  accessModes:
    - ReadWriteOncePod
  resources:
    requests:
      storage: 22Gi
  dataSource:
    name: pvc-origin-snap                             //Name of the VolumeSnapshot.
    kind: VolumeSnapshot                              //Only VolumeSnapshot is supported.
    apiGroup: snapshot.storage.k8s.io
```

## Topology support: Kubernetes cluster deployed on virtual machines in a single Zenlayer zone
* Set `zoneID` in the StorageClass to the ID of that zone. Every cloud disk that is created automatically then belongs to this zone, as do all virtual machine nodes in the cluster, so attaching a cloud disk to a virtual machine always follows the normal Kubernetes pod scheduling policy.    
* You may omit `zoneID` and `placeGroupID` from the StorageClass entirely, but in that case you must supply them when you install the CSI driver with Helm. This is supported only when the whole cluster lives in a single Zenlayer zone. 

``` yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: csi-zec
provisioner: disk.csi.zenlayer.com                      
parameters:
  type: "1"                                             
  zoneID: "asia-north-1a"                                 
  placeGroupID: "xxx"   
reclaimPolicy: Delete                                      
allowVolumeExpansion: true                                        
volumeBindingMode: Immediate                 
``` 

## Topology support: Kubernetes cluster deployed on virtual machines across multiple Zenlayer zones
### If you want to pin a pod to a specific zone
* Consider a cluster on the Zenlayer platform that spans several zones, for example 18 nodes of which 6 are in Shanghai, 6 are in Singapore, and 6 are in Los Angeles.
* A ZEC cloud disk can only be attached to a virtual machine in the same zone, so create one StorageClass per zone. Pods created from a workload that uses the `csi-zec-shanghai` StorageClass, and the cloud disks backing them, are then placed only on the six virtual machines in Shanghai.

* sc-shanghai.yaml for work-podA
``` yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: csi-zec-shanghai
provisioner: disk.csi.zenlayer.com                
parameters:
  type: "1"                                        
  zoneID: "asia-east-1a"                                     ## Shanghai
  placeGroupID: "xxx"    
reclaimPolicy: Delete                                    
allowVolumeExpansion: true                                  
volumeBindingMode: Immediate                 
``` 

* sc-sin.yaml for work-podB
``` yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: csi-zec-sin
provisioner: disk.csi.zenlayer.com                      
parameters:
  type: "1"                                             
  zoneID: "asia-southwest-1a"                                ## Singapore
  placeGroupID: "xxx"   
reclaimPolicy: Delete                                   
allowVolumeExpansion: true                              
volumeBindingMode: Immediate                 
``` 

* sc-los.yaml for work-podC
``` yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: csi-zec-los
provisioner: disk.csi.zenlayer.com                       
parameters:
  type: "2"                                              
  zoneID: "na-west-1a"                                     ##Los Angeles
  placeGroupID: "xxx"   
reclaimPolicy: Delete                                    
allowVolumeExpansion: true                              
volumeBindingMode: Immediate                 
``` 
### If you don't care which zone the pod lands in

``` yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: csi-zec
provisioner: disk.csi.zenlayer.com                       
parameters:
  type: "2"                                              
  zoneID: "na-west-1a"                                    ## Any value; it has no effect in this mode
  placeGroupID: "xxx"   
reclaimPolicy: Delete                                    
allowVolumeExpansion: true                              
volumeBindingMode: WaitForFirstConsumer                   ## Workload pods may be placed on any node
``` 