# Setup your environment for Golang-based operators

To setup your environment for developing Golang-based operators, you'll need the 
following prerequisites installed on your machine. Note that the homebrew 
version is the easiest, but is only available for macOS.

### Prerequisites for installing operator-sdk via homebrew (available for macOS)
* Access to a Kubernetes v1.11.3+ cluster (v1.16.0+ if using apiextensions.k8s.io/v1 CRDs). 
* User logged with admin permission. See how to grant yourself cluster-admin privileges or be logged in as admin
* Access to a container registry such as [Quay.io](https://quay.io) or [DockerHub](https://hub.docker.com/)

### Prerequisites for installing operator-sdk via GitHub (for Linux and Windows)
* [dep](https://golang.github.io/dep/docs/installation.html) v0.5.0+
* [Git](https://git-scm.com/downloads)
* [Go](https://golang.org/dl/) v1.10+
* [Docker](https://docs.docker.com/get-docker/) v17.03+
* OpenShift CLI (oc) v4.1+ installed
* Access to a Kubernetes v1.11.3+ cluster (v1.16.0+ if using apiextensions.k8s.io/v1 CRDs). See [minikube](https://minikube.sigs.k8s.io/docs/start/) or [CodeReady Containers](https://code-ready.github.io/crc/#installing-codeready-containers_gsg) to access a cluster for free.
* Access to a container registry such as [Quay.io](https://quay.io) or [DockerHub](https://hub.docker.com/)

## Install operator-sdk and its prerequisites 

* Use the homebrew command `brew install operator-sdk`
to install operator-sdk for macOS. If you don't have homebrew 
installed, install it from [here](https://docs.brew.sh/Installation).

* For Linux or Windows, install from GitHub release [here](https://sdk.operatorframework.io/docs/installation/#install-from-github-release).