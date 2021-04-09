set -x
set -e

# read -p 'Enter your namespace: ' namespace
# read -p 'Enter your IMG: ' img

make generate
make manifests
make install

# namespace is just your openshift project
# for example...
# export namespace=horea-demo-project
export namespace=janusgraph-demo-project
# img is where you plan to push your image 
# for example...
# export namespace=docker.io/horeaporutiu/memcached-operator:latest

export img=docker.io/horeaporutiu/janusgraph-operator:latest

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

kubectl apply -f config/samples/graph_v1alpha1_memcached.yaml