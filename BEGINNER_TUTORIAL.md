# Develop and Deploy a Memcached Operator on OpenShift Container Platform

In this tutorial, we will provide all of the information needed to get quickly up and running building and deploying
your first operator to OpenShift. 

**Note:** this tutorial can apply to other Kubernetes clusters as well. The commands may differ slightly however. 

## Expectations (What you have)
* You have little or no experience developing operators
* You have some knowledge of Kubernetes Operators concepts
* You've read the [Intro to Operators](https://github.ibm.com/TT-ISV-org/operator/blob/main/INTRO_TO_OPERATORS.md)
* You've setup your environment as shown in the [Setup your Environment](https://github.ibm.com/TT-ISV-org/operator/blob/main/installation.md) tutorial

## Expectations (What you want)
* You want hands-on experience developing and deploying your first operator to OpenShift
* You want to learn the basic concepts and steps needed to develop a Golang based operator to manage Kubernetes resources

If you already know all of the basic concepts of operators and have developed and deployed an operator before you can move on to the [Deep dive into Memcached ](https://github.ibm.com/TT-ISV-org/operator/blob/main/INTERMEDIATE_TUTORIAL.md), which will explain the low-level functions within the Operator reconcile function in more detail. It will also explain the KubeBuilder markers, creating the CRDs from the API, and other important operator-specific details.

**IMPORTANT**
This tutorial is inspired from the Operator SDK tutorial - https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/. **All credit goes to the Operator SDK team** for 
a great tutorial.

## Flow

![Flow](images/architecture.png)

1. Create a new operator project using the Operator SDK Command Line Interface(CLI)
2. Define new resource APIs by adding Custom Resource Definitions(CRD)
3. Define Controllers to watch and reconcile resources
4. Write the reconciling logic for your Controller using the SDK and controller-runtime APIs
5. Use the SDK CLI to build and generate the operator deployment manifests
6. Use the SDK CLI to build the operator image, push to image registry and then deploy to OpenShift
7. Operator docker image is deployed to OpenShift cluster creating manager and application replicas.
8. Reconcile loop watches and heals the resources as needed.

## Environment Setup

**IMPORTANT**
If you haven't setup your environment for building Kubernetes operators, setup your environment from these [instructions](installation.md).

## Steps
1. [Create a new project using Operator SDK](#1-create-a-new-project-using-operator-sdk)
1. [Create API and Custom Controller](#2-Create-API-and-Custom-Controller)
1. [Update API](#3-Update-API)
1. [Implement Controller Logic](#4-implement-controller-logic)
1. [Compile, build and push](#5-compile-build-and-push)
1. [Deploy the operator](#6-deploy-the-operator)
1. [Create the Custom Resource](#7-Create-the-Custom-Resource)
1. [Test and Verify](#8-test-and-verify)

## 1. Create a new project using Operator SDK

First check your Go version. This tutorial is tested with the following Go version:

```bash
$ go version
$ go version go1.15.6 darwin/amd64
```
Next, create a directory for where you will hold your project files. 

```bash
$ mkdir $HOME/projects/memcached-operator
$ cd $HOME/projects/memcached-operator
```
<!-- 
Since we are not in our $GOPATH, we can activate module support by running the 
`export GO111MODULE=on` command before using the operator-sdk. -->

Note: before your run the `operator-sdk init` command your `memcached-operator` directory must be completely empty, otherwise KubeBuilder will complain with an error. This means no `.git` folder, etc.

Next, run the `operator-sdk init` command to create a new memcached-operator project:

```bash
$ operator-sdk init --domain=example.com --repo=github.com/example/memcached-operator
```

* The `--domain` flag is used to uniquely identify the operator resources that are created by
this project. The `example.com` domain will be used as part of the Kubernetes [API group](https://kubernetes.io/docs/reference/using-api/#api-groups).

When we use the command `oc api-resources` later, the `example.com` domain 
will be listed there by our `memcached` in the `APIGROUP` category. 

* Let's discuss [Go Modules](https://blog.golang.org/using-go-modules). This is very important since if this is not setup properly, you will not be able to develop and run your operator. By using the `--repo` flag you are setting the name to use for your Go module, which is specified at the top of your `go.mod` file:

```go
module github.com/example/memcached-operator
``` 

* Setting up your Go Module will enable us to work outside of our [GOPATH](https://golang.org/doc/gopath_code#GOPATH), as long as the working directory of the project is the same as the name of the module in the top of the `go.mod` file. Again,
make sure that your directory is called `memcached-operator` and that your `go.mod` file shows the 
following go module:

```go
module github.com/example/memcached-operator
```

* For Go Modules to work properly, make sure you activate GO module support by running the following command:

```bash
$ export GO111MODULE=on
```

To verify that GO module support is turned on, issue the following command and ensure you get the same output: 

```bash
$ echo $GO111MODULE
on
```

This will create the basic scaffold for your operator, such as the `bin`, `config` and `hack` directories, and will create the `main.go` file which initializes the manager.

## 2. Create API and Custom Controller

Next, we will use the `operator-sdk create api` command to create an API which will be in your `api` directory and a blank custom controller file, which will be in your 
`controllers` directory.

We will use the --group, --version, and --kind flags to pass in the resource 
group and version. The <b>--group, --version, and --kind</b> flags together form the fully qualified name of a Kubernetes resource type. This name must be unique across a cluster.

```bash
$ operator-sdk create api --group=cache --version=v1alpha1 --kind=Memcached --controller --resource
Writing scaffold for you to edit...
api/v1alpha1/memcached_types.go
controllers/memcached_controller.go
```

* The `--group` flag defines an `API Group` in Kubernetes. It is a collection of related functionality.
* Each group has one or more `versions`, which allows us to change how an API works over time. This is what the `--version` flag represents.
* Each API group-version contains one or more API types, called `Kinds`. This is the name of the API type that we are creating as part of this operator. 
  * There are more nuances when it comes to versioning which we will not cover. Read more about Groups, Versions, Kinds, and Resources from this [KubeBuilder reference](https://book.kubebuilder.io/cronjob-tutorial/gvks.html).
* The `--controller` flag signifies that we want the sdk to scaffold a controller file.
* The `--resource` flag signifies that we want the sdk to scaffold the schema for a resource.


Once you deploy this operator, you can use the `kubectl api-resources` to see the name
`cache.example.com` as the api-group, and `Memcached` as the `Kind`. We can try this command 
later after we've deployed the operator.


### (Optional) Troubleshooting the create api command

If you get an error during the create api command, that means you will likely have to install the modules manually. 
Here is an example error: 

```bash
Error: go exit status 1: go: github.com/example/memcached-operator/controllers: package github.com/go-logr/logr imported from implicitly required module; to add missing requirements, run:
        go get github.com/go-logr/logr@v0.3.0
```

you will have to install the modules manually by running the following commands:

```bash
$ go get github.com/go-logr/logr@v0.3.0
$ go get github.com/onsi/ginkgo@v1.14.1
$ go get github.com/onsi/gomega@v1.10.2
```

## 3. Update API

One of the two main parts of the operator pattern is defining an API, which will be used to create our Custom Resource Definition(CRD).
We will do that in the `api/v1alpha1/memcached_types.go` file.

First, we need to understand the struct which defines our schema. Note that it 
implements the [Object interface](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/client#Object) (which means it is a kubernetes object), and also,
it has the `Spec` and `Status` fields. More on these soon.

```go 
type Memcached struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MemcachedSpec   `json:"spec,omitempty"`
	Status MemcachedStatus `json:"status,omitempty"`
}
```

### What is the Spec?

The `MemcachedSpec` struct, or the `Spec` defines the desired state of the resource. 
A good way to think about `Spec` is that any inputs (values tweaked by the user) to our controller go in the spec section.
You'll see the `Spec` section being referenced in the [controller code](https://github.ibm.com/TT-ISV-org/operator/blob/main/artifacts/memcached_controller.go#L104) to determine how many replicas to deploy.

### What is the Status? 

The `MemcachedStatus` struct, or the `Status` defines the current, observed state of the resource.
The status contains information that we want users or other controllers to be able to easily obtain. You'll 
see the status being updated in the [controller code](https://github.ibm.com/TT-ISV-org/operator/blob/main/artifacts/memcached_controller.go#L132), and that is for updating the current state of the cluster. 

Each of those structs, the `MemcachedStatus struct` and the `MemcachedSpec struct` will each
have their own fields to describe the observed state and the desired state respectively.

First, add a `Size int32` field to your `MemcachedSpec` struct, along with their JSON 
encoded string representation of the field name, in lowercase. See [Golangs JSON encoding page](https://golang.org/pkg/encoding/json/) for more details.

In our example, since `Size` is the field name, and the JSON encoding must be lowercase, it 
would just look like `json:"size"`. 

Add the following to your struct: 


```go
type MemcachedSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Memcached. Edit Memcached_types.go to remove/update
	Size int32 `json:"size"`
}
```

When we create a custom resource later, we will need to fill out the size, which is the number of `Memcached` replicas we want as the `desired state` of my system. 

Next, add a `Nodes []string` field to your `MemcachedStatus` struct, as shown below:

```go
// MemcachedStatus defines the observed state of Memcached
type MemcachedStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Nodes []string `json:"nodes"`
}
```


The `MemcachedStatus` struct will use a string array to list the name of the Memcached pods in the current state.

Lastly, the `Memcached` struct will have the fields `Spec` and `Status` to denote the desired state (spec) and the observed state (status). At a high-level, when the system recognizes there is a difference in the spec and the status, the operator will use custom controller logic defined in our 
`controllers/memcached_controller.go` file to update the 
system to be in the desired state.

Now that we've modified the file `api/v1alpha1/memcached_types.go`, it should look like the [file in the artifacts directory](https://github.ibm.com/TT-ISV-org/operator/blob/main/artifacts/memcached_types.go):

```go
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MemcachedSpec defines the desired state of Memcached
type MemcachedSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Memcached. Edit Memcached_types.go to remove/update
	Size int32 `json:"size"`
}

// MemcachedStatus defines the observed state of Memcached
type MemcachedStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Nodes []string `json:"nodes"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Memcached is the Schema for the memcacheds API
type Memcached struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MemcachedSpec   `json:"spec,omitempty"`
	Status MemcachedStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MemcachedList contains a list of Memcached
type MemcachedList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Memcached `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Memcached{}, &MemcachedList{})
}
```

Note that above our `type Memcached struct` you'll see two lines of code starting with `+kubebuilder`. Note that these are actually commented out. These are important, since they 
tell our controller-tools extra information. For example, this one 

```golang
// +kubebuilder:object:root=true
```

tells the `object` generator that this type represents a Kind. The generator will then 
implement the `runtime.Object` interface for us, which all Kinds must implement.

This one: 

```golang
// +kubebuilder:subresource:status
```

Will add the status subresource in the Custom Resource Definition. If you run `make manifests` it will generate YAML under `config/crds/<kind_types.yaml`. It will add a `subresources`
section like so: 

```yaml
subresources:
    status: {}
```

We will see how to get and update the status subresource in the controller code in the section below.

Just know that
each of these markers, starting with `// +kubebuilder` will generate utility code (such as role based access control) and Kubernetes YAML. When you run `make generate` and `make manifests` 
your KubeBuilder Markers will be read in order to create RBAC roles, CRDs, and code. Read more about KubeBuilder markers [here](https://book.kubebuilder.io/reference/markers.html?highlight=markers#marker-syntax).


## 4. Implement controller logic

<b>Note: If you want to learn more in depth about the controller logic that is written here,
please view our [Deep dive into Memcached Operator Code](https://github.ibm.com/TT-ISV-org/operator/blob/main/INTERMEDIATE_TUTORIAL.md) article.</b>

Now that we have our API updated, our next step is to implement our controller logic in `controllers/memcached_controller.go`. First, go ahead and copy the code from the 
[artifacts/memcached_controller.go](https://github.ibm.com/TT-ISV-org/operator/blob/main/artifacts/memcached_controller.go) file, and replace your current controller code.

Once this is complete, your controller should look like the following:

```go
/*
Copyright 2021.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cachev1alpha1 "github.com/example/memcached-operator/api/v1alpha1"
)

// MemcachedReconciler reconciles a Memcached object
type MemcachedReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// generate rbac to get, list, watch, create, update and patch the memcached status the nencached resource
// +kubebuilder:rbac:groups=cache.example.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete

// generate rbac to get, update and patch the memcached status the memcached/finalizers
// +kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/status,verbs=get;update;patch

// generate rbac to update the memcached/finalizers
// +kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/finalizers,verbs=update

// generate rbac to get, list, watch, create, update, patch, and delete deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// generate rbac to get,list, and watch pods
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Memcached object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *MemcachedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("memcached", req.NamespacedName)

	// Fetch the Memcached instance
	memcached := &cachev1alpha1.Memcached{}
	err := r.Get(ctx, req.NamespacedName, memcached)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Memcached resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Memcached")
		return ctrl.Result{}, err
	}

	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.Get(ctx, req.NamespacedName, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForMemcached(memcached)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// Ensure the deployment size is the same as the spec
	size := memcached.Spec.Size
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err = r.Update(ctx, found)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}

	// Update the Memcached status with the pod names
	// List the pods for this memcached's deployment
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(memcached.Namespace),
		client.MatchingLabels(labelsForMemcached(memcached.Name)),
	}
	if err = r.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "Memcached.Namespace", memcached.Namespace, "Memcached.Name", memcached.Name)
		return ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, memcached.Status.Nodes) {
		memcached.Status.Nodes = podNames
		err := r.Status().Update(ctx, memcached)
		if err != nil {
			log.Error(err, "Failed to update Memcached status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// deploymentForMemcached returns a memcached Deployment object
func (r *MemcachedReconciler) deploymentForMemcached(m *cachev1alpha1.Memcached) *appsv1.Deployment {
	ls := labelsForMemcached(m.Name)
	replicas := m.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   "memcached:1.4.36-alpine",
						Name:    "memcached",
						Command: []string{"memcached", "-m=64", "-o", "modern", "-v"},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 11211,
							Name:          "memcached",
						}},
					}},
				},
			},
		},
	}
	// Set Memcached instance as the owner and controller
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}

// labelsForMemcached returns the labels for selecting the resources
// belonging to the given memcached CR name.
func labelsForMemcached(name string) map[string]string {
	return map[string]string{"app": "memcached", "memcached_cr": name}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

// SetupWithManager sets up the controller with the Manager.
func (r *MemcachedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.Memcached{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
```

## 5. Compile, build and push

At this point, we are ready to compile, build the image of our operator, and push the image to an 
image repository. You can use the image registry of your choice, but here we will use Docker Hub. If we 
plan on deploying to an OpenShift cluster, we need to login to our cluster now.

From your provisioned cluster which we set up in the `installation.md` file, select the cluster and go to `OpenShift web console` by clicking the button from top right corner of the page.

![OpenShift](images/openshift-2.png)

From the OpenShift web console, copy the login command from the account drop down menu.

![OpenShift](images/openshift-3.png)

and from your terminal run the command to login to your cluster. Once you've logged in, you should see output like the following:

```bash
$ oc login --token=fFQ-HbFVBT4qHKl1n0b*****63U --server=https://c****-e.us-south.containers.cloud.ibm.com:31047
s-south.containers.cloud.ibm.com:31047

Logged into "https://c116-e.us-south.containers.cloud.ibm.com:31047" as "IAM#horea.porutiu@ibm.com" using the token provided.

You have access to 84 projects, the list has been suppressed. You can list all projects with 'oc projects'

Using project "horea-test-scc".
```

<b>This is extremely important.</b> By running the login command, we should now be able to run `oc project` to see which project we 
are currently in. The project we are in is our namespace as well, which is very important since our operator will only run in the namespace which we deploy it to. OpenShift is connecting to our cluster by using the login command, and if you do not do this step properly, you will not be able to deploy your operator.

Next create a new project using the following command:

```bash
$ oc new-project <new-project-name>
```

Once you've created a new project, you will be automatically switched to that project, as shown in the output below:

```bash
$ oc new-project memcache-demo-project
Now using project "memcache-demo-project" on server "https://c116-e.us-south.containers.cloud.ibm.com:31047".
```

Now, for the rest of the tutorial, you will use `memcache-demo-project` or whatever you named your project, as your namespace. More on this in the following steps. Just know that your project is the same as your namespace in terms of OpenShift. 

### Edit the manager.yaml file

The `manager.yaml` file defines a Deployment manifest used to deploy the operator. That manifest includes a security context that tells Kubernetes to run the pods as a specific user (uid=65532). OpenShift already manages the users employed to run pods which is behavior the manifest should not override, so we will remove that from the manifest.

To do this, we can modify the `config/manager/manager.yaml` file to remove the following line:

```
runAsUser: 65532
```

This will enable OpenShift to run its default security constraint. Once you've saved the file after you've removed the `runAsUser`
line, your file should look like the following, and the same as the one in the [artifacts directory](https://github.ibm.com/TT-ISV-org/operator/blob/main/artifacts/manager.yaml): 

```yaml
apiVersion: v1
kind: Namespace
metadata:
  labels:
    control-plane: controller-manager
  name: system
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: controller-manager
  namespace: system
  labels:
    control-plane: controller-manager
spec:
  selector:
    matchLabels:
      control-plane: controller-manager
  replicas: 1
  template:
    metadata:
      labels:
        control-plane: controller-manager
    spec:
      securityContext:
      containers:
      - command:
        - /manager
        args:
        - --leader-elect
        image: controller:latest
        name: manager
        securityContext:
          allowPrivilegeEscalation: false
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8081
          initialDelaySeconds: 15
          periodSeconds: 20
        readinessProbe:
          httpGet:
            path: /readyz
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 10
        resources:
          limits:
            cpu: 100m
            memory: 30Mi
          requests:
            cpu: 100m
            memory: 20Mi
      terminationGracePeriodSeconds: 10
```

### Create CRD and RBAC

The generated code from the `operator-sdk` creates a `Makefile` which allows you to use `make` command to compile your `go` operator code.

Now that we have our controller code and API implemented, run the following command to implement the required Go type interfaces:

```bash
$ make generate
```

The above command will update our `api/v1alpha1/zz_generated.deepcopy.go` file to implement the [metav1.Object](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Object) and [runtime.Object](https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime#Object) interfaces. This will enable our Custom Resource to be treated like a native Kubernetes resource.

Once we've generated the code for our custom resource, we can use the `make manifests` command to generate CRD manifests and RBAC from KubeBuilder Markers:

```bash
$ make manifests
```

This command will invoke controller-gen to generate the CRD manifests at `config/crd/bases/cache.example.com_memcacheds.yaml` - you can see the yaml representation of the object we specified in our `_types.go` file. It will also generate RBAC yaml files in the `config/rbac` directory based on 
your KubeBuilder markers.


Don't worry about [KubeBuilder Markers](https://book.kubebuilder.io/reference/markers.html) for now, we will cover them in the [deep-dive article](https://github.ibm.com/TT-ISV-org/operator/blob/main/INTERMEDIATE_TUTORIAL.md#6-understanding-kubebuilder-markers).

### Compile your Operator

To compile the code run the following command in the terminal from your project root:
```bash
$ make install
```

### Set the Operator Namespace

Next, we need to make sure to update our config to tell our operator to run in our own project namespace. Do this by issuing the following Kustomize 
commands:

```bash
$ export IMG=docker.io/<username>/memcached-operator:<version>
$ export NAMESPACE=<oc-project-name>

$ cd config/manager
$ kustomize edit set image controller=${IMG}
$ kustomize edit set namespace "${NAMESPACE}"
$ cd ../../

$ cd config/default
$ kustomize edit set namespace "${NAMESPACE}"
$ cd ../../
```

`<username>` is your Docker Hub (or Quay.io) username, and `<version>` is the 
version of the operator image you will deploy. Note that each time you 
make a change to your operator code, it is good practice to increment the 
version. `NAMESPACE` is your oc project name in which you plan to deploy your operator. For me, this would be `memcache-demo-project`.

For example, my export statements would look like the following:

```bash
$ export IMG=docker.io/horeaporutiu/memcached-operator:latest
$ export NAMESPACE=memcache-demo-project
```

### Build and Push your Image

**Note:** You will need to have an account to a image repository like Docker Hub to be able to push your 
operator image. Use `Docker login` to login.

To build the Docker image run the following command. Note that you can also 
use the regular `docker build -t` command to build as well. 

```bash

$ make docker-build IMG=$IMG
```
and push the docker image to your registry using following from your terminal:

```bash
$ make docker-push IMG=$IMG
```

## 6. Deploy the operator

#### Deploy the operator to OpenShift cluster

To Deploy the operator run the following command from your terminal:

```bash
$ make deploy IMG=$IMG
```

The output of the deployment should look like the following:
```bash
...go-workspace/src/memcached-operator/bin/controller-gen "crd:trivialVersions=true,preserveUnknownFields=false" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
cd config/manager && ...go-workspace/src/memcached-operator/bin/kustomize edit set image controller=sanjeevghimire/memcached-operator:v0.0.5
.../go-workspace/src/memcached-operator/bin/kustomize build config/default | kubectl apply -f -
Warning: kubectl apply should be used on resource created by either kubectl create --save-config or kubectl apply
namespace/sanjeev-operator-prj configured
customresourcedefinition.apiextensions.k8s.io/memcacheds.cache.example.com configured
role.rbac.authorization.k8s.io/memcached-operator-leader-election-role created
clusterrole.rbac.authorization.k8s.io/memcached-operator-manager-role configured
clusterrole.rbac.authorization.k8s.io/memcached-operator-metrics-reader unchanged
clusterrole.rbac.authorization.k8s.io/memcached-operator-proxy-role unchanged
rolebinding.rbac.authorization.k8s.io/memcached-operator-leader-election-rolebinding created
clusterrolebinding.rbac.authorization.k8s.io/memcached-operator-manager-rolebinding configured
clusterrolebinding.rbac.authorization.k8s.io/memcached-operator-proxy-rolebinding configured
configmap/memcached-operator-manager-config created
service/memcached-operator-controller-manager-metrics-service created
deployment.apps/memcached-operator-controller-manager created
```

To make sure everything is working correctly, use the `oc get pods` command.

```bash
$ oc get pods

NAME                                                     READY   STATUS    RESTARTS   AGE
memcached-operator-controller-manager-54c5864f7b-znwws   2/2     Running   0          14s
```

This means your operator is up and running. Great job!

## 7. Create the Custom Resource

Next, let's create the custom resource.

Update your custom resource, by modifying the `config/samples/cache_v1alpha1_memcached.yaml` file
to look like the following:

```yaml
apiVersion: cache.example.com/v1alpha1
kind: Memcached
metadata:
  name: memcached-sample
spec:
  # Add fields here
  size: 3
``` 
Note that all we did is set the size of the Memcached replicas to be 3.

And finally create the custom resources using the following command:

```bash
$ kubectl apply -f config/samples/cache_v1alpha1_memcached.yaml
```

#### Verify that resources are Running

From the terminal run `kubectl get all` or `oc get all` to make sure that controllers, managers and pods have been successfully created and is in `Running` state with the right number of pods as defined in the spec.

```bash
kubectl get all 
```

Output:
![kubectl get all](images/kubectl-get-all.png)

Also from your cluster you can see the logs by going to your project in `OpenShift web console`

![kubectl get all](images/os-logs.png)

You can also now run `oc api-resources` to view the memcache resource we have created:
```bash
oc api-resources
NAME                APIGROUP                  NAMESPACED   KIND
memcacheds         cache.example.com          true         Memcached
```

Lastly, if you want to make some code changes in the controller and build and deploy a new version of your operator, you can simply use the `build-and-deploy.sh` script. Just make sure
to set your namespace and img in that file.

Run the script by issuing the following command:

```bash
$ ./build-and-deploy.sh
```

## 8. Test and verify

Update `config/samples/<group>_<version>_memcached.yaml` to change the `spec.size` field in the Memcached CR. This will increase the application pods from 3 to 5. It can also be patched directly in the cluster as follows:

```bash
$ oc patch memcached memcached-sample -p '{"spec":{"size": 5}}' --type=merge
```

You can also update the spec.size from `OpenShift web console` by going to `Deployments` and selecting `memcached-sample` and increase/decrease using the up or down arrow:

![kubectl get all](images/inc-dec-size.png)

Next, let's verify that our pods have scaled up. Run the following command:

```bash
$ kubectl get pods
```

You should now see that there are 5 total `memcached-sample` pods.

**Congratulations!** You've successfully deployed an Memcached operator using the `operator-sdk`! To learn more, go ahead and read
the [Deep dive into Memcached Operator Code](https://github.ibm.com/TT-ISV-org/operator/blob/main/INTERMEDIATE_TUTORIAL.md) tutorial, 
which explains the controller logic from step 4 in more depth.


## Cleanup

The `Makefile` part of generated project has a target called `undeploy` which deletes all the resources associated with your project. It can be run as follows:

```bash
make undeploy
```

# License

This code pattern is licensed under the Apache Software License, Version 2.  Separate third party code objects invoked within this code pattern are licensed by their respective providers pursuant to their own separate licenses. Contributions are subject to the [Developer Certificate of Origin, Version 1.1 (DCO)](https://developercertificate.org/) and the [Apache Software License, Version 2](https://www.apache.org/licenses/LICENSE-2.0.txt).

[Apache Software License (ASL) FAQ](https://www.apache.org/foundation/license-faq.html#WhatDoesItMEAN)
