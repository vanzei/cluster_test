# Create a kind cluster config file (cluster-config.yaml)
cat <<EOF > cluster-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
  - role: worker
  - role: worker
  - role: worker
  - role: worker
EOF

# Create the cluster
kind create cluster --config cluster-config.yaml
