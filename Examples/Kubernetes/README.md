# Using Kubernetes deployment manifest
Use the provided example [deployment manifests](./calendar-to-mqtt-deployment/) as a starting point:

# [Deployment](./calendar-to-mqtt-deployment/deployment.yaml)
```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: calendar-to-mqtt
  labels:
    app: calendar-to-mqtt
spec:
  replicas: 1
  strategy:
    type: RollingUpdate
  selector:
    matchLabels:
      app: calendar-to-mqtt
  template:
    metadata:
      labels:
        app: calendar-to-mqtt
    spec:
      containers:
        - image: ghcr.io/teknikens-hus/calendar-to-mqtt:latest
          name: calendar-to-mqtt
          env:
            - name: TZ
              value: Europe/Stockholm
          # Adjust the resource limits as needed
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
              ephemeral-storage: "200Mi"
            limits:
              memory: "256Mi"
              cpu: "200m"
              ephemeral-storage: "1Gi"
          volumeMounts:
            - name: config-volume
              mountPath: /config.yaml
              subPath: config.yaml
          # Make sure we run as non-root "app"
          securityContext:
            runAsUser: 1001
            runAsGroup: 2001
            fsGroup: 2001
      Volumes:
        - name: config-volume
          configMap:
            name: config-configmap
```
## [Kustomization](./calendar-to-mqtt-deployment/kustomization.yaml)
The kustomization has a configMapGenerator
```yaml
---
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: calendar-to-mqtt
metadata:
  name: kustomize-calendar-to-mqtt
resources:
- deployment.yaml
- namespace.yaml
configMapGenerator:
  - name: config-configmap
    files:
      - config.yaml=config.yaml
    options:
      disableNameSuffixHash: false
```

## [Namespace](./calendar-to-mqtt-deployment/namespace.yaml)
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: calendar-to-mqtt
```


## Config.yaml
Adjust the values in the config.yaml file

To see options for the config file, check the main [README.md](../../README.md)