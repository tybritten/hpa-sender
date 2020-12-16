# This script will deploy a kind cluster and configure knative eventing and the metrics server for HPA

curl -sL https://raw.githubusercontent.com/csantanapr/knative-kind/master/01-kind.sh | bash

kubectl apply -f https://github.com/knative/eventing/releases/download/v0.19.0/eventing-crds.yaml
kubectl apply -f https://github.com/knative/eventing/releases/download/v0.19.0/eventing-core.yaml
kubectl apply -f https://github.com/knative/eventing/releases/download/v0.19.0/in-memory-channel.yaml
kubectl apply -f https://github.com/knative/eventing/releases/download/v0.19.0/mt-channel-broker.yaml



kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/download/v0.3.6/components.yaml
kubectl patch deployment metrics-server -n kube-system -p '{"spec":{"template":{"spec":{"containers":[{"name":"metrics-server","args":["--cert-dir=/tmp", "--secure-port=4443", "--kubelet-insecure-tls","--kubelet-preferred-address-types=InternalIP"]}]}}}}'
