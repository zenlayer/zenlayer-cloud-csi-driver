# k8s proxy

## install
kubectl apply -f https://raw.githubusercontent.com/ketches/registry-proxy/master/deploy/manifests.yaml 

## yaml example
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: registry-proxy

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: registry-proxy
  namespace: registry-proxy

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: registry-proxy
rules:
  - apiGroups: [""]
    resources: ["namespaces", "configmaps", "secrets"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
  - apiGroups: ["admissionregistration.k8s.io"]
    resources: ["mutatingwebhookconfigurations"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: registry-proxy
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: registry-proxy
subjects:
  - kind: ServiceAccount
    name: registry-proxy
    namespace: registry-proxy

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: registry-proxy
  namespace: registry-proxy
spec:
  selector:
    matchLabels:
      app: registry-proxy
  template:
    metadata:
      labels:
        app: registry-proxy
    spec:
      serviceAccountName: registry-proxy
      containers:
        - name: registry-proxy
          image: registry.cn-hangzhou.aliyuncs.com/ketches/registry-proxy:v1.3.2
          imagePullPolicy: Always
          resources:
            requests:
              memory: "64Mi"
              cpu: "50m"
            limits:
              memory: "256Mi"
              cpu: "200m"
          ports:
            - containerPort: 443

---
apiVersion: v1
kind: Service
metadata:
  name: registry-proxy
  namespace: registry-proxy
spec:
  selector:
    app: registry-proxy
  ports:
    - port: 443
      targetPort: 443
  type: ClusterIP

```
