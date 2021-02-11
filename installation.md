# Environment setup

### Prerequisites for installing operator-sdk for macOS

To setup your environment for developing Golang-based operators, you'll need the 
following prerequisites installed on your machine. Note that the homebrew 
version is the easiest, but is only available for macOS. 

* [Homebrew](https://brew.sh/)
* [Go](https://golang.org/dl/) 1.10+
* Access to a Kubernetes v1.11.3+ cluster (v1.16.0+ if using apiextensions.k8s.io/v1 CRDs). See [minikube](https://minikube.sigs.k8s.io/docs/start/) or [CodeReady Containers](https://code-ready.github.io/crc/#installing-codeready-containers_gsg) to access a cluster for free.
* User logged with admin permission. See how to grant yourself cluster-admin privileges or be logged in as admin.
* Access to a container registry such as [Quay.io](https://quay.io) or [DockerHub](https://hub.docker.com/)
* [Kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/homebrew/)
* Either OpenShift CLI [oc](https://docs.openshift.com/container-platform/4.5/cli_reference/openshift_cli/getting-started-cli.html) or [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

### Prerequisites for installing for Linux and Windows
* [dep](https://golang.github.io/dep/docs/installation.html) v0.5.0+
* [Git](https://git-scm.com/downloads)
* [Go](https://golang.org/dl/) v1.10+
* [Docker](https://docs.docker.com/get-docker/) v17.03+
* OpenShift CLI (oc) v4.1+ installed
* Access to a Kubernetes v1.11.3+ cluster (v1.16.0+ if using apiextensions.k8s.io/v1 CRDs). See [minikube](https://minikube.sigs.k8s.io/docs/start/) or [CodeReady Containers](https://code-ready.github.io/crc/#installing-codeready-containers_gsg) to access a cluster for free.
* Access to a container registry such as [Quay.io](https://quay.io) or [DockerHub](https://hub.docker.com/)
* [Kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/)


## Steps

1. [Install Operator SDK](#1-install-operator-sdk)
1. [Install oc or kubectl cli](#2-install-oc-or-kubectl-cli)

## 1. Install Operator SDK

### Install operator-sdk (version 1.0+) and Kustomize for macOS

* Use the homebrew command `brew install operator-sdk`
to install operator-sdk for macOS. Note that this guide 
is tested for operator-sdk version 1.0+, since the commands have changed with the 1.0 release. 
 If you don't have homebrew 
installed, install it from [here](https://docs.brew.sh/Installation).

* Use the homebrew command `brew install kustomize` to install Kustomize.

### Install operator-sdk (version 1.0+) and Kustomize for Linux or Windows

* For Linux or Windows, install the operator-sdk (version 1.0+) from the GitHub release [here](https://sdk.operatorframework.io/docs/installation/#install-from-github-release). Note that
commands have changed with the 1.0 release.

* You can use the following script to install Kustomize for Windows or Linux but note that it doesn't work for ARM architecture. For ARM architecture download 
Kustomize from the [releases page](https://github.com/kubernetes-sigs/kustomize/releases).

```
curl -s "https://raw.githubusercontent.com/\
kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"  | bash
```

### Test your environment for operator-sdk

Run the following command in the terminal of your choice:

```
operator-sdk version
```

You should see output like this:

```
operator-sdk version: "v1.3.0", commit: "1abf57985b43bf6a59dcd18147b3c574fa57d3f6", kubernetes version: "v1.19.4", go version: "go1.15.5", GOOS: "darwin", GOARCH: "amd64"
```

Now, let's ensure kustomize is installed.

```
kustomize version
```

You should see output like this:

```
{Version:kustomize/v3.9.1 GitCommit:7439f1809e5ccd4677ed52be7f98f2ad75122a93 BuildDate:2020-12-30T01:08:17+00:00 GoOs:darwin GoArch:amd64}
```

## 2. Install oc or kubectl cli
If you plan to use an OpenShift cluster, then you can install the OpenShift CLI via
your web console. Otherwise you can install kubectl from [here](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

First, go to your OpenShift console and click on the question mark in the 
top-right corner. From there, click on `Command Line Tools` and then choose
the `oc` CLI binary for your operating system. Once you've downloaded it,
ensure that the command is in your `PATH`.

Test your cli by issuing the following command to see the version of your cli:

```
$ oc version
Client Version: openshift-clients-4.5.0-202006231303.p0-18-g6082e941e
Kubernetes Version: v1.19.2
```

## Make sure OpenShift Lifecycle Manager (OLM) is up to date

As a note, if you still need to provision an OpenShift cluster, it takes some time
so it is recommended to do that **now** if you don't have one already. Skip down to
[step 9](https://github.ibm.com/TT-ISV-org/operator#9-deploy-the-operator) to see 
how to create an OpenShift cluster on IBM Cloud.

First, we need to take care of some cluster admin tasks. We will need to make sure our OpenShift Lifecycle Manager is 
up to date and running properly before we develop our operator. To do this, run the `operator-sdk olm status` command:

```
operator-sdk olm status
INFO[0003] Fetching CRDs for version "0.16.1"           
INFO[0003] Using locally stored resource manifests      
INFO[0005] Successfully got OLM status for version "0.16.1" 
NAME                                            NAMESPACE    KIND                        STATUS
operators.operators.coreos.com                               CustomResourceDefinition    Installed
operatorgroups.operators.coreos.com                          CustomResourceDefinition    Installed
installplans.operators.coreos.com                            CustomResourceDefinition    Installed
clusterserviceversions.operators.coreos.com                  CustomResourceDefinition    Installed
olm-operator                                    olm          Deployment                  Installed
subscriptions.operators.coreos.com                           CustomResourceDefinition    Installed
olm-operator-binding-olm                                     ClusterRoleBinding          Installed
operatorhubio-catalog                           olm          CatalogSource               Installed
olm-operators                                   olm          OperatorGroup               Installed
aggregate-olm-view                                           ClusterRole                 Installed
catalog-operator                                olm          Deployment                  Installed
aggregate-olm-edit                                           ClusterRole                 Installed
olm                                                          Namespace                   Installed
global-operators                                operators    OperatorGroup               Installed
operators                                                    Namespace                   Installed
packageserver                                   olm          ClusterServiceVersion       Installed
olm-operator-serviceaccount                     olm          ServiceAccount              Installed
catalogsources.operators.coreos.com                          CustomResourceDefinition    Installed
system:controller:operator-lifecycle-manager                 ClusterRole                 Installed
```

Note: if you see an error, you should still be able to proceed. Now you are ready to start developing your first operator.