---
name: kubernetes-patterns
description: Kubernetes manifests, Helm charts, operators, and container orchestration patterns.
allowed-tools: Read, Write, Edit, Glob, Grep
---

# Kubernetes Patterns

> Principles for Kubernetes deployment and container orchestration.
> **Declarative infrastructure, immutable deployments.**

---

## 1. Resource Selection

### Workload Types

| Resource | Use Case |
|----------|----------|
| **Deployment** | Stateless apps, most common |
| **StatefulSet** | Databases, ordered deployment |
| **DaemonSet** | Node-level agents, logging |
| **Job** | One-time tasks |
| **CronJob** | Scheduled tasks |

### Decision Tree

```
Does it need persistent identity? → StatefulSet
Does it need to run on every node? → DaemonSet
Is it a one-time task? → Job/CronJob
Otherwise → Deployment
```

---

## 2. Pod Configuration

### Essential Fields

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
  labels:
    app: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: app
        image: myapp:1.0.0
        ports:
        - containerPort: 8080
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 512Mi
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
```

### Resource Guidelines

| Resource | Request | Limit |
|----------|---------|-------|
| **CPU** | Typical usage | 2-3x request |
| **Memory** | Typical usage | = request (avoid OOM) |

---

## 3. Service Networking

### Service Types

| Type | Use Case |
|------|----------|
| **ClusterIP** | Internal communication (default) |
| **NodePort** | External access (development) |
| **LoadBalancer** | External access (cloud) |
| **Headless** | StatefulSet discovery |

### Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app-ingress
  annotations:
    kubernetes.io/ingress.class: nginx
spec:
  rules:
  - host: app.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: app
            port:
              number: 80
```

---

## 4. Configuration

### ConfigMaps

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  APP_ENV: production
  LOG_LEVEL: info
```

### Secrets

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
type: Opaque
data:
  DATABASE_URL: base64encoded...
```

### Using in Pods

```yaml
env:
  - name: APP_ENV
    valueFrom:
      configMapKeyRef:
        name: app-config
        key: APP_ENV
  - name: DATABASE_URL
    valueFrom:
      secretKeyRef:
        name: app-secrets
        key: DATABASE_URL
```

---

## 5. Health Checks

### Probe Types

| Probe | Purpose | Failure Action |
|-------|---------|----------------|
| **Liveness** | Is container alive? | Restart container |
| **Readiness** | Can it serve traffic? | Remove from service |
| **Startup** | Has it started? | Delay other probes |

### Configuration

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 10
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  periodSeconds: 5
```

---

## 6. Scaling

### Horizontal Pod Autoscaler

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: app-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: app
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

---

## 7. Helm Charts

### Chart Structure

```
mychart/
├── Chart.yaml
├── values.yaml
├── templates/
│   ├── _helpers.tpl
│   ├── deployment.yaml
│   ├── service.yaml
│   └── ingress.yaml
└── charts/
```

### values.yaml Pattern

```yaml
replicaCount: 3

image:
  repository: myapp
  tag: "1.0.0"
  pullPolicy: IfNotPresent

resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 500m
    memory: 512Mi
```

---

## ✅ Checklist

- [ ] **Resource limits set?** (CPU, memory)
- [ ] **Health probes configured?** (liveness, readiness)
- [ ] **Replicas > 1 for production?** (high availability)
- [ ] **Secrets not in git?** (external secrets)
- [ ] **Labels consistent?** (selector matching)
- [ ] **Resource requests reasonable?** (not over-provisioned)
- [ ] **PodDisruptionBudget set?** (for availability)

---

> **Remember:** Kubernetes is declarative. Define the desired state and let the system reconcile.
