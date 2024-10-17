#!/bin/bash


# Apply all Kubernetes manifests
kubectl apply -f postgres/postgres.yaml
kubectl apply -f nginx.yaml
kubectl apply -f api/api-deployment.yaml
kubectl apply -f frontend/web-frontend.yaml
kubectl apply -f api/otel-collector-config.yaml
kubectl apply -f tempo.yaml
kubectl apply -f grafana.yaml
kubectl apply -f postgres/postgres-configmap.yaml

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
kubectl apply -f ingress.yaml

echo "Deployment complete! Access Grafana at http://localhost:3000"