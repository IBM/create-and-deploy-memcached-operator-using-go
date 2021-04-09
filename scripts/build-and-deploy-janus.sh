set -x
set -e

make generate
make manifests
make install

# namespace is just your openshift project
# for example...
# export namespace=horea-demo-project

export namespace=<add-namespace-here>

# img is where you plan to push your image 
# for example...
# export namespace=docker.io/horeaporutiu/janusgraph-operator:latest

export img=docker.io/<username-goes-here>/janusgraph-operator:latest

cd config/manager
kustomize edit set namespace $namespace
kustomize edit set image controller=$img
cd ../../
cd config/default
kustomize edit set namespace $namespace
cd ../../

make docker-build IMG=$img
make docker-push IMG=$img
make deploy IMG=$img

kubectl apply -f config/samples/graph_v1alpha1_janusgraph.yaml
