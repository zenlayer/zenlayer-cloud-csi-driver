# Installation


You can manage your Zec-Vm and Zec-CloudDisk and secret on the [Zenlayer console](https://console.zenlayer.com/).

## Prerequisites

* Kubernetes version >= 1.28.2
* `kubectl` configured to communicate with the cluster
* Helm 3.18.1
* Controller type csi pod need Connect to External public network 

Please review the [chart values file](../chart/values.yaml) before installing.

## Setup Permissions

The driver requires All permissions to invoke Zenlayer Cloud OpenAPIs to manage the volume on user's behalf.

* use a secret for access key.
  1. Create a Console user, enable OpenAPI access. Once the user is created, record the Access key ID and Access key password.
  2. Store the AccessKey to Cluster as a secret.
  ```shell
  kubectl create secret -n kube-system generic csi-access-key --from-literal=AccessKeyID='***********'  --from-literal=AccessKeyPassword='***********'
  ```

## Deploy the Drivers

* You can deploy the drivers using Helm. zeccsi have two roles, controller and node(csi-zecplugin-provisioner and csi-zecplugin)    
  1. controller is Deployment mode, pod need connected to External public network.     
  2. node is DaemonSet mode, will be deployed on each node.       

* You can choose one of the following two installation methods. Automatic installation is recommended. You can use helm for a complete installation. The image will be automatically pulled from docker hub.          
* If you want to manually package the image, please refer to [ZecCSI Build guide](../build-guide.md)        
### Automatic installation csi pod && Automatic pull image

```shell
## view chart/values.yaml View the configurable parameters during installation
helm install zeccsi oci://registry-1.docker.io/zenlayer297/zenlayer-cloud-csi-driver --version 1.0.0 --set defaultResourceGroup="" --set defaultZone="" --set maxVolume=6
```
### Manual installation csi pod && Automatic pull image

```shell
helm package ./chart
helm install zeccsi ./zenlayer-cloud-csi-driver-1.0.0.tgz
```
* You can also install csi with args, for example set controller-csi replica or specify controller-csi which node will be deployed.          
* If only node02 nodes in your k8s cluster can be connected to the External public network, you need deploy controller-csi on node02        
  1. first add lable: kubectl label nodes node02 zeccsiType=Controller       
  2. helm install zeccsi ./zenlayer-cloud-csi-driver-1.0.0.tgz --set controllerSelectorkey=zeccsiType --set controllerSelectorval=Controller --set replicaCount=1         
  3. Then Controller csi provisioner will deployed on node02, only have one replica        

## Verify

Check the csi pods are running and ready:      

```shell
kubectl get pods -n kube-system -l app=csi-zecplugin
kubectl get pods -n kube-system -l app=csi-zecplugin-provisioner
```