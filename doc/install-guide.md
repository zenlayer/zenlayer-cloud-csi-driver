# Installation
You can manage your ZEC VMs, ZEC cloud disks, and secrets in the [Zenlayer console](https://console.zenlayer.com/).

## Prerequisites

* Kubernetes version >= 1.28.2
* `kubectl` configured to communicate with the cluster
* Helm 3.18.1+
* The controller-type CSI pod needs to reach the public internet

Review the [chart values file](../chart/values.yaml) before installing.

## Set Up Permissions

The driver needs permission to call the Zenlayer Cloud OpenAPIs that manage volumes on the user's behalf.

* Use a secret for the access key.
  1. Create a console user and enable OpenAPI access for it. Once the user is created, record the access key ID and the access key password.
  2. Store the access key in the cluster as a secret.
  ```shell
  kubectl create secret -n kube-system generic csi-access-key --from-literal=AccessKeyID='***********'  --from-literal=AccessKeyPassword='***********'
  ```

## Deploy the Drivers

* You can deploy the drivers using Helm. The zeccsi pod has two roles: controller and node (`csi-zecplugin-provisioner` and `csi-zecplugin`).    
  1. The controller runs as a Deployment; its pod needs to reach the public internet.     
  2. The node runs as a DaemonSet and is deployed on every node.       

* Choose one of the two installation methods below. Automatic installation is recommended, because Helm performs the whole installation and the image is pulled from Docker Hub automatically.          
* If you want to build and package the image yourself, refer to the [ZecCSI build guide](../deploy/build/build-guide.md).        
### Automatic Installation (image pulled automatically)

```shell
# See chart/values.yaml for the parameters you can configure at install time, for example --set maxVolume=5
helm install zeccsi oci://registry-1.docker.io/zenlayer297/zenlayer-cloud-csi-driver --version 1.2.0
```
### Manual Installation (image pulled automatically)

```shell
helm package ./chart
helm install zeccsi ./zenlayer-cloud-csi-driver-1.2.0.tgz
```
* You can also pass arguments at install time, for example to set the number of controller replicas or to pin the controller to a specific node.          
* If node02 is the only node in your Kubernetes cluster that can reach the public internet, deploy the controller there:        
  1. Label the node: `kubectl label nodes node02 zeccsiType=Controller`       
  2. Install with a matching selector: `helm install zeccsi ./zenlayer-cloud-csi-driver-1.2.0.tgz --set controllerSelectorkey=zeccsiType --set controllerSelectorval=Controller --set replicaCount=1`         
  3. The controller CSI provisioner is then deployed on node02 with a single replica.        

## Verify

Check that the CSI pods are running and ready:      

```shell
kubectl get pods -n kube-system -l app=csi-zecplugin
kubectl get pods -n kube-system -l app=csi-zecplugin-provisioner
kubectl get pods -n kube-system -l app=csi-fluent-bit
```