# Zenlayer cloud csi driver
[![GoReportCard Widget]][GoReportCardResult]

简体中文 | [English](./README.md)

## 插件描述
* Zenlayer CSI插件实现了容器存储接口（[CSI](https://github.com/container-storage-interface/)）,启用容器编排（CO）和Zenlayer的存储之间的接口,可以自动管理云盘
* 版本介绍: v 1(稳定发行版).0(补丁版本).0(BUG修正版本)

## Zeccsi 和 Kubernetes 版本兼容列表
| ZEC CSI Version | Container Orchestrator Name | Version Tested     |
| -----------------| --------------------------- | -------------------|
| v1.0.0          | Kubernetes                   |  v1.28.2   +         |

## Zeccsi 功能发布版本列表
| ZEC CSI Version | Feature                                          |
| -----------------| -------------------------------------------------|
| v1.0.0          | 创建删除云盘，挂载云盘到Pod，云盘在线扩容              |

## helm 兼容版本
| Helm Version   | 
| -----------------|
|  v3.18.1   +       |

## csi sidecar 兼容版本
| sidecar           |       Version |
| ------------------ | -----------------|
| csi-provisioner                | v5.2.0          |
| csi-attacher                   | v4.6.1          |
| csi-resizer                    | v1.8.0          |
| csi-node-driver-registrar      | v2.8.0          |
| livenessprobe                  | v2.8.0          |

## zenlayer openApi sdk 最低版本 [API Github](https://github.com/zenlayer/zenlayercloud-sdk-go)
| SDK Version |
| -----------------|
| v0.1.27    +      |


# 云盘CSI插件
提供磁盘CSI驱动程序，帮助简化存储管理。一旦用户使用磁盘存储类创建了PVC，磁盘和相应的PV对象就会动态创建，并准备好供K8s Pod工作负载使用

## 配置要求
* 具有Zenlayer控制平台的合法认证权限. [console](https://console.zenlayer.com)        
* k8s集群是部署到zec的虚拟机内.   

## 如何使用

### Step 1: 准备需要的环境
* 在ZEC的虚拟机内搭建的一套kubernetes环境.      
* Kubectl命令环境配置.      

### Step 2: 安装Zec CSI驱动
* 如果你只需要部署使用zeccsi插件，按照此文档部署 [部署指南](./doc/install-guide.md).          
* csi镜像和chart仓库位于docker hub,确保网络通畅. csi可以使用helm完整安装，不需要从github下载源码除非您有开发需求.                       

### Step 3: k8s集群内创建StorageClass
存储类是动态卷供应所必需的,可以按照此例子部署[Storage-class设置指南](./doc/storage-class.md).

### Step 4: 检查CSI Pod状态
检查所有的CSI Pod是否运行并准备就绪.
```shell
kubectl get pods -n kube-system -l app=csi-zecplugin
```
输出:
```
NAME                  READY   STATUS    RESTARTS   AGE
csi-zecplugin-2xxr9   3/3     Running   0          2m19s
```
```shell
kubectl get pods -n kube-system -l app=csi-zecplugin-provisioner
```
输出:
```
NAME                                         READY   STATUS    RESTARTS   AGE
csi-zecplugin-provisioner-5ccf9f74cc-h6qfj   5/5     Running   0          2m37s
```

### Step 5: 测试部署工作负载Pod
为了确保您的CSI插件正常工作，创建一个简单的工作负载来测试它:
```shell
kubectl apply -f deploy/example/nginx-statefulset.yaml
```
你可以在Zenlayer Console上看到由CSI自动创建的云盘：
```shell
kubectl get pvc
kubectl get pv
```

测试完成后删除测试Pod：
```shell
kubectl delete -f deploy/example/nginx-statefulset.yaml
```

### Step 6:测试扩容PVC
选择一个PVC修改 spec:resources:requests:storage
```shell
kubectl get pvc -o wide
kubectl edit pvc nginx-data-nginx-statefulset-0
```

## 注意事项

* 我们强烈建议您在kubernetes集群中部署您的日志收集服务，持久化csi pod的日志，以便问题追踪           

[GoReportCard Widget]: https://goreportcard.com/badge/github.com/zenlayer/zenlayer-cloud-csi-driver
[GoReportCardResult]: https://goreportcard.com/report/github.com/zenlayer/zenlayer-cloud-csi-driver
