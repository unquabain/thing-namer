#! /usr/bin/env bash

version=$(date -u +%F-%H-%M)

docker build -t "unquabain/projectnamer:$version" --platform linux/amd64 .
docker push "unquabain/projectnamer:$version"
kubectl config use-context do-sfo3-unquabain-k8s-01
kubectl get deployment/thing-namer -oyaml \
	| sed "s/image: unquabain\/projectnamer:.*/image: unquabain\/projectnamer:$version/g" \
	> thing-namer.yaml
kubectl apply -f thing-namer.yaml
sleep 10
curl --verbose --header "Referer: http://localhost:8080" https://wizard-bacon.unquabain.com/api.json
