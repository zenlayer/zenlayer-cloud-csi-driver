# The zeccsi plugin uses the Storage-class method to manage storage

## Configuration Description
### storage-class yaml Description          
```yaml
provisioner: disk.csi.zenlayer.com                       //csi driver name,  Cannot be modified        

parameters:                                              //None of the parameters are mandatory
  fsType: "ext4"                                         //mount filesystem type. support ext4 ext3 or xfs. If not setting, default is ext4
  type: "1"                                              //cloud disk type：1 Basic NVMe SSD, 2 Standard NVMe SSD. If not setting, default is Standard NVME                                    
  zoneID: "asia-north-1a"                                //cloud disk zone，It only takes effect when volumeBindingMode=Immediate. If not setting, You must specific this val when helm install    
  placeGroupID: "xxx"                                    //cloud disk zenlayer console resource group ID. If not setting, You must specific this val when helm install     

reclaimPolicy: Delete                                    //support "Delete" and "Retain". It is not recommended to use the Retain mode. Users need to manually delete the cloud disk specifically. It may cause data residue      

allowVolumeExpansion: true                               //support cloud disk online Expansion          

volumeBindingMode: Immediate                             //support "Immediate" and "WaitForFirstConsumer" mode.  Immediate mode, pv will be created and bound immediately. The pv will be created in the area specified by the zoneID, and the Pod workloads using the storageclass will also be allocated to this zone.  WaitForFirstConsumer mode, pv will not be created immediately. It will remain pending until a Pod workload uses this pvc. The pv and the used workload Pods will then be randomly assigned locations within the cluster by the k8s scheduler.           

mountOptions:                                            //This feature needs to be enabled with configuration parameters during deployment. helm install featureGates=enable_mount_opt
  - discard
```
### pvc yaml Description    
Not support "dataSource"      
```yaml 
spec:       
  accessModes:        
    - ReadWriteOnce                                    //only support this mod     
  volumeMode: Filesystem                               //Filesystem is default, Cannot be modified      
  resources:      
    requests:     
      storage: 80Gi                                    //cloud disk size,  Greater than 20G    
  storageClassName: csi-zec                            //storage-class name     
```
## The k8s cluster is deployed in virtual machine in One zenlayer area
* In the configuration of storage-class, the zoneID information is set to the zoneID of this area. The automatically created cloud disk belongs to this zone, and all virtual machine nodes in the k8s cluster also belong to this zone. When the cloud disk is mounted to the virtual machine, it complies with the pod scheduling policy of k8s.    
* You can omit every parameters such as zoneID and placeGroupID in storage-class, but you need specific these when use helm install csi-driver.  (Only Support if k8s cluster in One zenlayer area) 

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

## The k8s cluster is deployed in virtual machines across multiple zenlayer regions

### If you want specify the region of the virtual machine where the pod is located
* If the k8s cluster deployed by the user on the zenlayer platform spans multiple zone regions, for instance, in a k8s cluster with 18 nodes, 6 virtual machine nodes belong to the Shanghai region, 6 virtual machine nodes belong to the Singapore region, and 6 virtual machine nodes belong to Los Angeles.
* zec's cloud disk can only be mounted on virtual machines in the same area. So users can create three Storage-classes. The Pods created using the workload of the storage-class csi-zec-shanghai and the corresponding cloud disks will only be mounted on the six virtual machines in the Shanghai area.

* sc-shanghai.yaml for work-podA
``` yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: csi-zec-shanghai
provisioner: disk.csi.zenlayer.com                
parameters:
  type: "1"                                        
  zoneID: "asia-east-1a"                                ##shanghai
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
  zoneID: "asia-southwest-1a"                                ##Singapore
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
  zoneID: "na-west-1a"                                   ##Los Angeles
  placeGroupID: "xxx"   
reclaimPolicy: Delete                                    
allowVolumeExpansion: true                              
volumeBindingMode: Immediate                 
``` 
### If you don't care about the region where the pod is located

``` yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: csi-zec
provisioner: disk.csi.zenlayer.com                       
parameters:
  type: "2"                                              
  zoneID: "na-west-1a"                                   ## Whatever. It doesn't take effect
  placeGroupID: "xxx"   
reclaimPolicy: Delete                                    
allowVolumeExpansion: true                              
volumeBindingMode: WaitForFirstConsumer                 ## workload Pods will be randomly distributed all node
``` 