helm install postgres bitnami/postgresql -f config/postgres-values.yaml

kubectl apply -f config/k8s-config.yaml
kubectl apply -f manifests