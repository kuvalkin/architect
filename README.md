# Architect

Simple study-project to learn more about microservises.

## Helm repos
```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
```

## How to

### Run
```bash
bash scripts/run.sh
```

### Stop and delete everithing (but not PVC)
```bash
bash scripts/delete.sh
```

### Run tests
```
newman run --env-var baseUrl=http://{node ip or url}:{architect service port} app/tests/postman/*
```