kubectl delete -f manifests/service-monitor.yaml -n monitoring
helm uninstall prometheus -n monitoring

helm uninstall postgres -n architect

kubectl delete -f config/k8s-config.yaml -n architect
kubectl delete -f manifests/deployment.yaml -n architect
kubectl delete -f manifests/ingress.yaml -n architect
kubectl delete -f manifests/service.yaml -n architect