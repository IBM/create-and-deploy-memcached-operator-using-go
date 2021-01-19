#!/bin/bash
set -x
set -e

# set ENV variable so we can run outside of $GOPATH
export GO111MODULE=on

# set namespace and image vars
NAMESPACE="kkbankol"
IMG="docker.io/kkbankol/memcached-operator:v0.0.13"

# switch to namespace
oc project ${NAMESPACE}

# initialize project directory + api
operator-sdk init --domain="example.com" --repo="github.com/example/memcached-operator"
operator-sdk create api --group=cache --version=v1alpha1 --kind=Memcached --controller --resource

# update Dockerfile and manager.yaml, default user (65532) won't work on OS due to security constraints
# TODO, valid user range seems to change every time we run "make install"
cp ../artifacts/Dockerfile .
cp ../artifacts/manager.yaml config/manager/

# overwrite controller, generated version is outdated as referenced here
# https://github.com/operator-framework/operator-sdk/issues/4381#issuecomment-760473278
cp ../artifacts/memcached_controller.go controllers/

# update CRD with size spec
cp ../artifacts/memcached_types.go ./api/v1alpha1/

# Generate CRD definitions
make generate
make manifests

# Register CRDs
make install
oc get Memcached
# Optional, test run operator locally
# make run ENABLE_WEBHOOKS=false

# Update manager manifest to ensure it's pointing to correct image and namespace
cd config/manager
kustomize edit set image controller=${IMG}
cd ../../

cd config/default
kustomize edit set namespace "${NAMESPACE}"
cd ../../

# Build and push docker image
docker build -t ${IMG} .
docker push ${IMG}

# Deploy operator
kustomize build config/default | oc apply -f -
make deploy IMG=${IMG}
oc get deployment

# Add size key to sample
cp ../artifacts/cache_v1alpha1_memcached.yaml config/samples/

# Deploy memcached samples
oc apply -f config/samples/cache_v1alpha1_memcached.yaml
oc get pods
oc patch memcached memcached-sample -p '{"spec":{"size": 4}}' --type=merge
oc get pods
