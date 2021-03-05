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
* [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* [Docker](https://docs.docker.com/get-docker/) v17.03+
* OpenShift CLI (**If you plan to deploy to OpenShift Cluster**) [oc](https://docs.openshift.com/container-platform/4.5/cli_reference/openshift_cli/getting-started-cli.html)

### Prerequisites for installing for Linux and Windows
* [dep](https://golang.github.io/dep/docs/installation.html) v0.5.0+
* [Git](https://git-scm.com/downloads)
* [Go](https://golang.org/dl/) v1.10+
* [Docker](https://docs.docker.com/get-docker/) v17.03+
* OpenShift CLI (oc) v4.1+ installed
* Access to a Kubernetes v1.11.3+ cluster (v1.16.0+ if using apiextensions.k8s.io/v1 CRDs). See [minikube](https://minikube.sigs.k8s.io/docs/start/) or [CodeReady Containers](https://code-ready.github.io/crc/#installing-codeready-containers_gsg) to access a cluster for free.
* Access to a container registry such as [Quay.io](https://quay.io) or [DockerHub](https://hub.docker.com/)
* [Kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/)
* OpenShift CLI (**If you plan to deploy to OpenShift Cluster**) [oc](https://docs.openshift.com/container-platform/4.5/cli_reference/openshift_cli/getting-started-cli.html)


## Steps

1. [Install Operator SDK](#1-install-operator-sdk)
1. [Install oc or kubectl cli](#2-install-oc-or-kubectl-cli)

## 1. Install Operator SDK

### Install operator-sdk (version 1.0+) for macOS

* Use the homebrew command `brew install operator-sdk`
to install operator-sdk for macOS. Note that this guide 
is tested for operator-sdk version 1.0+, since the commands have changed with the 1.0 release. 
 If you don't have homebrew 
installed, install it from [here](https://docs.brew.sh/Installation).

### Install operator-sdk (version 1.0+) for Linux or Windows

* For Linux or Windows, install the operator-sdk (version 1.0+) from the GitHub release [here](https://sdk.operatorframework.io/docs/installation/#install-from-github-release). Note that
commands have changed with the 1.0 release.

### Test your environment for operator-sdk

Run the following command in the terminal of your choice:

```bash
$ operator-sdk version
```

You should see output like this:

```bash 
operator-sdk version: "v1.3.0", commit: "1abf57985b43bf6a59dcd18147b3c574fa57d3f6", kubernetes version: "v1.19.4", go version: "go1.15.5", GOOS: "darwin", GOARCH: "amd64"
```

Now, let's ensure kustomize is installed.

```bash 
$ kustomize version
```

You should see output like this:

```bash
{Version:kustomize/v3.9.1 GitCommit:7439f1809e5ccd4677ed52be7f98f2ad75122a93 BuildDate:2020-12-30T01:08:17+00:00 GoOs:darwin GoArch:amd64}
```

## 2. Install oc or kubectl cli
If you plan to use an OpenShift cluster, then you can install the OpenShift CLI using [these instructions](https://docs.openshift.com/container-platform/latest/cli_reference/openshift_cli/getting-started-cli.html).

Otherwise you can install kubectl from [here](https://kubernetes.io/docs/tasks/tools/install-kubectl/).

Alternatively, here is another way to install the `oc` cli, through the OpenShift web console, as shown in [this article](https://developers.redhat.com/openshift/command-line-tools):

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

If you plan to use `kubectl` instead of `oc`:
```bash
$ kubectl version
Client Version: version.Info{Major:"1", Minor:"20", GitVersion:"v1.20.2", GitCommit:"faecb196815e248d3ecfb03c680a4507229c2a56", GitTreeState:"clean", BuildDate:"2021-01-14T05:15:04Z", GoVersion:"go1.15.6", Compiler:"gc", Platform:"darwin/amd64"}
Server Version: version.Info{Major:"1", Minor:"18+", GitVersion:"v1.18.3+e574db2", GitCommit:"e574db2", GitTreeState:"clean", BuildDate:"2021-01-30T06:33:00Z", GoVersion:"go1.13.15", Compiler:"gc", Platform:"linux/amd64"}
```

## Make sure OpenShift Lifecycle Manager (OLM) is up to date

As a note, if you still need to provision an OpenShift cluster, it takes some time
so it is recommended to do that **now** if you don't have one already. Skip down to
[the prepare your OpenShift Cluster step](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md#prepare-your-openshift-cluster) to see 
how to create an OpenShift cluster on IBM Cloud.

Lastly, we will need to make sure our OpenShift Lifecycle Manager is 
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

As you can see from my output above, all of the components of the OLM are in the `Installed` status.
If your components are in the `Installed` status, that means your Operator Lifecycle Manager is working properly.
<b>Note: if you see an error, you can read [this guide](https://sdk.operatorframework.io/docs/olm-integration/quickstart-bundle/#enabling-olm) which will show you how to install OLM on your cluster.</b>

### (Optional) Troubleshooting OLM error
If you've ran into an error like this one:

```bash
$ operator-sdk olm status 
FATA[0002] Failed to get OLM status: error getting installed OLM version (set --version to override the default version): no existing installation found 
```

or something like this: 

```bash 
$ operator-sdk olm status 
FATA[0002] Failed to get OLM status: error getting installed OLM version (set --version to override the default version): no existing installation found 


$ operator-sdk olm install
FATA[0005] Failed to install OLM version "latest": detected existing OLM resources: OLM must be completely uninstalled before installation 


$ operator-sdk olm uninstall
FATA[0002] Failed to uninstall OLM: error getting installed OLM version (set --version to override the default version): no existing installation found 

Sometimes, you will have to uninstall a specific version. For my OpenShift Cluster which is version `4.5.31_1531`,
I had to uninstall version `0.16.1`

$ operator-sdk olm uninstall --version 0.16.1
INFO[0009] Successfully uninstalled OLM version "0.16.1" 
```

Then I just installed the same version, and checked the status:

```bash 
$ operator-sdk olm install --version 0.16.1 
INFO[0072] Successfully installed OLM version "0.16.1"  

$ operator-sdk olm status 
INFO[0004] Successfully got OLM status for version "0.16.1" 

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

Once everything says installed, <b>congratulations</b> you are ready to start developing!


## Conclusion
<b>Congratulations!!</b> You've now setup your environment to develop an operator 
and deploy it to an OpenShift (or Kubernetes) cluster. You are ready to move on to the [Develop and Deploy a Memcached Operator](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md) tutorial.