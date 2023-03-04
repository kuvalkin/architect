kubectl create namespace architect
kubectl create namespace monitoring

helm install postgres bitnami/postgresql -f config/postgres-values.yaml -n architect --atomic

kubectl apply -f config/k8s-config.yaml -n architect
kubectl apply -f manifests/deployment.yaml -n architect
kubectl apply -f manifests/ingress.yaml -n architect
kubectl apply -f manifests/service.yaml -n architect

kubectl apply -f manifests/service-monitor.yaml -n monitoring
helm install prometheus prometheus-community/kube-prometheus-stack -f config/prometheus-values.yaml -n monitoring --atomic