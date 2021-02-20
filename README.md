# Build a simple Golang-based operator

In this tutorial we will be creating a simple Go-based operator that walks you through an example of building a simple memcached-operator using operator-sdk.

Operators make it easy to manage complex stateful applications on top of Kubernetes or OpenShift.

## Expectations (What you have)
* You have some experience developing operators.
* You've finished the `SIMPLE_OPERATOR.md` tutorial
* You've read articles and blogs on the basic idea of a Kubernetes Operators, and you know the basic Kubernetes resource types.

## Expectations (What you want)
* You want deep technical knowledge of the code which enables operators to run.
* You want to understand how the reconcile loop works, and how you can use it to manage Kubernetes resources
* You want to learn more about the basic Get, Update, and Create functions used to save resources to your Kubernetes cluster.
* You want to learn more about Kubebuilder markers and how to use them to set role based access control.

**IMPORTANT**
This tutorial is inspired from the operator-sdk tutorial - https://sdk.operatorframework.io/docs/building-operators/golang/tutorial/. **All credit goes to the operator-sdk team** for 
a great tutorial.

Note that this tutorial goes into much more depth in terms of the controller code
and the logic needed to understand a Kubernetes Operator. 

🚧 If you want a simple, fast, and quick
way to deploy your first operator, without all of the deep technical explanations of what is happening
behind the scenes, see the  `SIMPLE_OPERATOR.md` file. 🚧

## Flow

![Flow](images/architecture.png)

1. Create a new operator project using the SDK Command Line Interface(CLI)
2. Define new resource APIs by adding Custom Resource Definitions(CRD)
3. Define Controllers to watch and reconcile resources
4. Write the reconciling logic for your Controller using the SDK and controller-runtime APIs
5. Use the SDK CLI to build and generate the operator deployment manifests
6. Use the SDK CLI to build operator docker image, push and deploy to OpenShift
7. Operator docker image is deployed to OpenShift cluster creating manager and application replicas.
8. Reconcile loop watches and heals the resources as needed.

## Environment Setup

**IMPORTANT**
If you haven't setup your environment for building Kubernetes operators, setup your environment from these [instructions](installation.md).

## Steps
1. [Create a new project using Operator SDK](#1-create-a-new-project-using-operator-sdk)
1. [Create CRD and Custom Controller](#2-Create-CRD-and-Custom-Controller)
1. [Update CRD and generate CRD manifest](#3-Update-CRD-and-generate-CRD-manifest)
1. [Implement Controller Logic](#4-implement-controller-logic)
1. [Compile, build and push](#5-compile-build-and-push)
1. [Deploy the operator](#6-deploy-the-operator)
1. [Test and verify](#7-test-and-verify)

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

Next, run the `operator-sdk init` command to create a new memcached-operator project:

```bash
$ operator-sdk init --domain=example.com --repo=github.com/example/memcached-operator --owner="Memcache Operator authors"
```

* The `--domain` flag is used to uniquely identify the operator resources that are created by
this project. When we use the command `oc api-resources` later, the `example.com` domain 
will be listed there by our `memcached` in the `APIGROUP` category.

* The `--repo` flag enables us to create this project outside of the standard 
`$GOPATH/src` strucutre. 
  * To work properly, make sure you activate GO module support by running the following command:

```bash
$ export GO111MODULE=on
```

To verify that GO module support is turned on, issue the following command and ensure you get the same output: 

```bash
$ echo $GO111MODULE
on
```

This will create the basic scaffold for your operator, such as the `bin`, `config` and `hack` directories, and will create the `main.go` file which initializes the manager.

## 2. Create CRD and Custom Controller

Next, we will use the `operator-sdk create api` command to create a blank <b>custom resource definition,
or CRD</b> which will be in your `api` directory and a blank custom controller file, which will be in your 
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
* Each API group-verison contains one or more API types, called `Kinds`. This is the name of the API type that we are creating as part of this operator. 
  * There are more nuances when it comes to versioning which we will not cover. Read more about Groups, Versions, Kinds, and Resources from this [Kubebuilder reference](https://book.kubebuilder.io/cronjob-tutorial/gvks.html).
* The `--controller` flag signifies that we want the sdk to scaffold a controller file.
* The `--resource` flag signifies that we want the sdk to scaffold the schema for a resource.


Now, once you deploy this operator, you can use the `kubectl api-resources` to see the name
`cache.example.com` as the api-group, and `Memcached` as the `Kind`. We can try this command 
later after we've deployed the operator.

## 3. Update CRD and generate CRD manifest

One of the two main parts of the operator pattern is defining a Custom Resource Definition(CRD). We
will do that in the `api/v1alpha1/memcached_types.go` file.

Let's first understand the basic foundation of our custom resource. There are 
three main structures to understand:

First, we need to understand the struct which defines our schema. Note that it 
implements the [Object interface](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Object) (which means it is a kubernetes object), and also,
it has the `Spec` and `Status` fields. More on those soon.

```go 
type Memcached struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MemcachedSpec   `json:"spec,omitempty"`
	Status MemcachedStatus `json:"status,omitempty"`
}
```

The `MemcachedSpec` struct, or the `Spec` defines the desired state of the resource. 

### What is the Spec?

A good way to think about `Spec` is that any inputs (values tweaker by the user) to our controller go in the spec section. 

```go
type MemcachedSpec struct {}
```

The `MemcachedStatus` struct, or the `Status` defines the current, observed state of the resource.

### What is the Status? 

The status contains information that we want users or other controllers to be able to easily obtain.

```go
type MemcachedStatus struct {}
```

Each of those structs, the `MemcachedStatus struct` and the `MemcachedSpec struct` will each
have their own fields to describe the observed state and the desired state respectively.

First, add a `Size int32` field to your `MemcachedSpec` struct, along with their json 
encoded string representation of the field name, in lowercase. See [golangs json encoding page](https://golang.org/pkg/encoding/json/) for more details.

In our example, since `Size` is the field name, and the json encoding must be lowercase, it 
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

Modify the `api/v1alpha1/memcached_types.go` to look like the the [file in the artifacts directory](https://github.ibm.com/TT-ISV-org/operator/blob/main/artifacts/memcached_types.go).

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

Will enable the status subresource in the in the Custom Resource Definition. If you run `make manifests` it will generate YAML under `config/crds/<kind_types.yaml`. It will add a `subresources`
section like so: 

```yaml
subresources:
    status: {}
```

We will see how to get and update the status subresource in the controller code in the section below.

Just know that
each of these markers, starting with `// +kubebuilder` will generate utility code (such as role based access control) and Kubernetes YALM. When you run `make generate` and `make manifests` 
your KubeBuilder Markers will be read in order to create RBAC roles, CRDs, and code, such as runtime.Object/DeepCopy implementations. Read more about KubeBuilder markers [here](https://book.kubebuilder.io/reference/markers.html?highlight=markers#marker-syntax).


## 4. Implement controller logic

Now that we have our CRDs registered, our next step is to implement our controller logic in `controllers/memcached_controller.go`. First, go ahead and copy the code from the 
[artifacts/memcached_controller.go](https://github.ibm.com/TT-ISV-org/operator/blob/main/artifacts/memcached_controller.go) file, and replace your current controller code. The next
few paragraphs will explain the controller code in detail. This is the heart of the operator.
If you're already experienced with operators, you can skip down to [build manifests and go files](https://github.ibm.com/TT-ISV-org/operator#build-manifests-and-go-files).

The controller "Reconcile" method contains the logic responsible for monitoring and applying the requested state for specific deployments. It does so by sending client requests to Kubernetes APIs, and will run every time a Custom Resource is modified by a user or changes state (ex. pod fails). If the reconcile method fails, it can be re-queued to run again.

After scaffolding our controller via the operator-sdk, we'll have an empty Reconciler function.

In this example, we want our Reconciler to
1. Check for an existing memcached deployment, and create one if it does not exist.
2. Retrieve the current state of our memcached deployment, and compare it to our desired state. More specifically, we'll compare the memcached deployment ReplicaSet value to the "Size" parameter that we defined earlier in our `memcached_types.go` file.
3. If the number of pods in the deployment ReplicaSet does not match the provided `Size`, then our Reconciler will update the ReplicaSet value, and re-queue the Reconciler until the desired state is achieved.

So, we'll start out by adding logic to our empty Reconciler function. First, we'll reference the instance we'd like to observe, which is the `Memcached` object defined in our `api/v1alpha1/memcached_types.go` file. We'll do this by retrieving the Memcached CRD from the `cachev1alpha1` object, which is listed in our import statements. Note that the trailing endpoint of the url maps to the files in our `/api/v1alpha1/` directory.

```go
import (
  ...
  cachev1alpha1 "github.com/example/memcached-operator/api/v1alpha1"  
)
```

Here we'll simply use `cachev1alpha1.<Object>{}` to reference any of the defined objects within that `memcached_types.go` file.

```go
memcached := &cachev1alpha1.Memcached{}
```


### Understanding the Get(ctx context.Context, key types.NamespacedName, obj client.Object) function

Next, we'll need to confirm that the `Memcached` resource is defined within our namespace.

This can be done using the [`Get` function](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Reader.Get), which retrieves an object from a Kubernetes cluster based on the arguments passed in. 

**Important:** The `Get` function expects the Reconciler context, the object key (which is just the namespace, and the name of the object), and the object itself as arguments. The object has to implement the [Object interface](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Object) which means that it is serializable (runtime.Object) and identifiable (metav1.Object). This object should be able to be written via YAML and then be created 
via `Kubectl create`.

The Reconcile function gives you two things, the context i.e. `ctx` and request i.e. `req`. 
The request parameter has all of the information we need to reconcile a Kubernetes object i.e. a
`memcached` object in our case. More specifically, the `req` struct contains the `NamespacedName` field which is the name and the namespace of the object to reconcile. That 
is what we will pass in to the `Get` function.

 If the resource doesn't exist, we'll receive an error.
```go
err := r.Get(ctx, req.NamespacedName, memcached)
```

If the Memcached object does not exist in the namespace yet, the Reconciler will return an error and try again.
```go
return ctrl.Result{}, err
```

### Understanding the Reconcile(ctx context.Context, req ctrl.Request) (Result, error) return types

Now, let's talk a bit about what the [reconcile function](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/reconcile#Reconciler) returns. This can be a bit 
tricky since there are various return types. 

The reconcile function returns a (Result, err). Now, more specifically, 
the [Result struct](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/reconcile#Result) has two fields, the `Requeue` bool, which just tells the reconcile 
function to requeue again. This bool defaults to false. The other field is 
`RequeueAfter` which expects a `time.Duration`. This tell the reconciler to requeue after a specific amount of time. 

For example the following code would requeue after 30 seconds.
```go
return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
```

Furthermore, the controller will requeue the request to be processed again if an error
is non-nil or `Result.Requeue` is true.

Here are three of the most common return types:

1. `return ctrl.Result{Requeue: true}, nil` when you want to return and requeue the request. This is done usually when we have updated the state of the cluster, i.e. created a deployment, or updated the spec. 
2. `return ctrl.Result{}, err` when there is an error. This will requeue the request.
3. `return ctrl.Result{}, nil` when everything goes fine and you do not want to requeue. This is
the return at the bottom of the reconcile loop. This means the observed state of the 
cluster is the same as the desired state (i.e. the `MemcachedSpec` is the same as the `MemcachedStatus`).

So at this point, our Reconciler function should look like: 

```go
func (r *MemcachedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
  // reference Memcached object
  memcached := &cachev1alpha1.Memcached{}
  // check if Memcached object is within namespace
  err := r.Get(ctx, req.NamespacedName, memcached)
  if err != nil {
    // throw error if Memcached object hasn't been defined yet
    return ctrl.Result{}, err
  }
}
```

Assuming the resource is defined, we can continue on by observing the state of our Memcached Deployment.

First, we'll want to confirm that a Memcached deployment exists within the namespace. To do so, we'll need to use the [k8s.io/api/apps/v1](https://godoc.org/k8s.io/api/apps/v1#Deployment) package, which is defined in our import statement.
```go
import (
	appsv1 "k8s.io/api/apps/v1"
  ...
)
```

Use the `apps` package to reference a [Deployment object](https://pkg.go.dev/k8s.io/api/apps/v1#Deployment) (note that a deployment object is a Kubernetes object which implements the [Object interface](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Object) as noted above), and then use the reconciler `Get` function to check whether the Memcached deployment exists with the provided name within our namespace.

```go
found := &appsv1.Deployment{}
err = r.Get(ctx, req.NamespacedName, found)
```

### Create a new memcached deployment if one is not found

If a deployment is not found, then we can use `Deployment` definition within the the `apps` package to create a new one using the reconciler [`Create`](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Writer) method:

```go

found := &appsv1.Deployment{}
err = r.Get(ctx, req.NamespacedName, found)

if err != nil && errors.IsNotFound(err) {
  dep := r.deploymentForMemcached(memcached)
  log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
  err = r.Create(ctx, dep)
  ...
  // if successful, return and re-queue Reconciler method
  return ctrl.Result{Requeue: true}, nil
```

For improved readability, the deployment definition has been placed in a different function called [`deploymentForMemcached`](https://github.ibm.com/TT-ISV-org/operator/blob/main/artifacts/memcached_controller.go#L134). This function includes the pod runtime specs (ports, startup command, image name), and the `Memcached.Spec.Size` value to determine how many replicas should be deployed. This function returns the deployment resource i.e. a Kubernetes object.

```go
func (r *MemcachedReconciler) deploymentForMemcached(m *cachev1alpha1.Memcached) *appsv1.Deployment {
	ls := labelsForMemcached(m.Name)
	replicas := m.Spec.Size

  dep := &appsv1.Deployment{
    ...
    Spec: appsv1.DeploymentSpec{
      Replicas: &replicas,
      ...
      Template: corev1.PodTemplateSpec{
        ...
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
	return dep
```

### Use Create(ctx context.Context, obj client.Object) to save the object

Once we create that deployment, we use the `r.Create(ctx context.Context, obj client.Object)` function to actually save the 
object in the Kuberenetes cluster. The `r.Create(ctx context.Context, obj client.Object)` 
function takes in the context (which is passed into our reconcile function) and the Kubernetes object we want to save (which in our case is the deployment we just created) in the `deploymentForMemcached` function:

```go
dep := r.deploymentForMemcached(memcached)
log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
err = r.Create(ctx, dep)
```
Since we've made an update to our cluster, we will requeue:

```go
return ctrl.Result{Requeue: true}, nil
```

Next, we'll add logic to our method to adjust the number of replicas in our deployment whenever the `Size` parameter is adjusted. This is assuming the deployment already exists in our namespace.



### Use Update(ctx context.Context, obj Object) to update the replicas in the Spec
First, request the desired `Size` and then compare the desired size to the number of replicas running in the deployment. If the states don't match, we'll use the `Update` method to adjust the amount of replicas in the deployment.

```go
size := memcached.Spec.Size
if *found.Spec.Replicas != size {  
  found.Spec.Replicas = &size
  err = r.Update(ctx, found)
  ...
}
```
If all goes well, the spec is updated, and we requeue. Otherwise, we return an error.

```go
if err != nil {
  log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
  return ctrl.Result{}, err
}
// Spec updated - return and requeue
return ctrl.Result{Requeue: true}, nil
```

### Use Update(ctx context.Context, obj Object) to update the status

Lastly, we will retrieve the list of pods in a specific namespace by using the 
[r.List](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Reader.List) function. The r.List function will create a `.Items` field in the 
ObjectList we pass in which will be populated with the objects for a given namespace.

```go
podList := &corev1.PodList{}
listOpts := []client.ListOption{
  client.InNamespace(memcached.Namespace),
  client.MatchingLabels(labelsForMemcached(memcached.Name)),
}
if err = r.List(ctx, podList, listOpts...); err != nil {
  log.Error(err, "Failed to list pods", "Memcached.Namespace", memcached.Namespace, "Memcached.Name", memcached.Name)
  return ctrl.Result{}, err
}
```

Then, we have a function to convert the PodList into a string array, since that 
is what our `MemcachedStatus` struct is expecting, as we have defined it in our `memcached_types.go` file.

```go
podNames := getPodNames(podList.Items)
```

Lastly, we will check if the podnames that we've just listed from `r.List` are the same 
as the `memcached.Status.Nodes`. If they are not the same, we will use [`Update(ctx context.Context, obj Object)` function](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Writer) to update the `MemcachedStatus` struct:

```go
// Update status.Nodes if needed
if !reflect.DeepEqual(podNames, memcached.Status.Nodes) {
  memcached.Status.Nodes = podNames
  err := r.Status().Update(ctx, memcached)
  if err != nil {
    log.Error(err, "Failed to update Memcached status")
    return ctrl.Result{}, err
  }
}
```

If all goes well, we return without an error. 
```go
return ctrl.Result{}, nil
```

### Role Based Access Control and Kubebuilder markers

Now, one more thing to understand before we deploy our operator. The Kubebuilder 
markers which you can see at the top of the file:

```go
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
```

**Kubebuilder markers are extrmely imporant and tricky since they are written in comments.**

For example, the marker belows generates and updates the rbac yaml files in our `config/rbac` directory. Once we deploy these 
updated files, our operator will have the permission to get, list, watch, create, update, path, and delete the `memcacheds` resources, as shown below: 

```go
// generate rbac to get, list, watch, create, update and patch the memcached status the nencached resource
// +kubebuilder:rbac:groups=cache.example.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
```

For example, if our memcached resource didn't have the `List` verb listed in the kubebuilder marker, we would not be able to use r.List() on our memcached resource - we would get a permissions error such as `Failed to list *v1.Pod`. Once we change these markers and add the `list` command, we have to run `make generate` and `make manifests` and that will in turn apply the changes from our kubebuilder commands into our `config/rbac` yaml files. To 
learn more about kubebuilder markets, see the docs [here](https://book.kubebuilder.io/reference/markers/rbac.html).


Once this is complete, your controller should look like the file in [artifacts/memcached_controller.go](artifacts/memcached_controller.go):

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

### Build manifests and go files

Before we compile our code, we need to change a couple of things. 

1. Make sure to change 
your Dockerfile so it looks exactly as the [one in the Artifacts directory](https://github.ibm.com/TT-ISV-org/operator/blob/main/artifacts/Dockerfile). It should look like this:

```Dockerfile
# Build the manager binary
FROM golang:1.15 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .

ENTRYPOINT ["/manager"]
```

2. Make sure to change 
your `manager.yaml` file in your `config/manager` directory so it looks exactly as the [one in the Artifacts directory](https://github.ibm.com/TT-ISV-org/operator/blob/main/artifacts/manager.yaml). It 
should look like the following: 

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


Now that we have our controller code and memcached types implemented, run the following command to update the generated code for that resource type:

```bash
$ make generate
```

The above command will use the controller-gen utility in `bin/controller-gen` to update the api/v1alpha1/zz_generated.deepcopy.go file to ensure our API’s Go type definitions implement the `runtime.Object` interface that all Kind types must implement.

Once the API is defined with spec/status fields and CRD validation markers, the CRD manifests can be generated and updated with the following command:

```bash
$ make manifests
```

This command will invoke controller-gen to generate the CRD manifests at `config/crd/bases/cache.example.com_memcacheds.yaml` - you can see the yaml representation 
of the object we specified in our `_types.go` file. 

## 5. Compile, build and push

At this point, we are ready to compile, build the image of our operator, and push the image to an 
image repository. You can use the image registry of your choice, but here we will use Docker Hub. If we 
plan on deploying to an OpenShift cluster, we need to login to our cluster now.

### Prepare your OpenShift Cluster

(If you haven't already) provision an OpenShift cluster by going to `https://cloud.ibm.com/` and clicking `Red Hat OpenShift on IBM Cloud` and get into 

![OpenShift](images/openshift-1.png)

Once you provisioned the cluster, select the cluster and go to `OpenShift web console` by clicking the button from top right corner of the page.

![OpenShift](images/openshift-2.png)

From the OpenShift web console, copy the login command from the account drop down menu.

![OpenShift](images/openshift-3.png)

and from your terminal run the command to login to your cluster.

If you haven't created a project, create a project by going to projects and clicking `Create Project`. From the terminal after you logged in change the project by running following in your terminal.

Note: you can also use the `oc new-project <new-project-name>` command to create a new project.
The below command simply switches you to an existing project.

```bash
oc project <project name>

```

### Create Operator Image

The generated code from the `operator-sdk` creates a `Makefile` which allows you to use `make` command to compile your `go` operator code. The same make command also allows you to build and push the docker image.

To compile the code run the following command in the terminal from your project root:
```bash
make install
```

**Note:** You will need to have an account to a image repository like Docker Hub to be able to push your 
operator image. Use `Docker login` to login.

To build the Docker image run the following command. Note that you can also 
use the regular `docker build -t` command to build as well. 

`<username>` is your Docker Hub (or Quay.io) username, and `<version>` is the 
version of the operator image you will deploy. Note that each time you 
make a change to your operator code, it is good practice to increment the 
version.


```bash
make docker-build IMG=docker.io/<username>/memcached-operator:<version>
```
and push the docker image to your registry using following from your terminal:

 ```bash
make docker-push IMG=docker.io/<username>/memcached-operator:<version>

 ```

## 6. Deploy the operator

#### Deploy the operator to OpenShift cluster

Make sure that the controller manager manifest has the right namespace and docker image. Apply the same to the default manifest as well by running following command:

```bash
export IMG=docker.io/<username>/memcached-operator:<version>
export NAMESPACE=<oc-project-name>

cd config/manager
kustomize edit set image controller=${IMG}
kustomize edit set namespace "${NAMESPACE}"
cd ../../

cd config/default
kustomize edit set namespace "${NAMESPACE}"
cd ../../
```


To Deploy the operator run the following command from your terminal:

```bash
make deploy IMG=docker.io/<username>/memcached-operator:<version>
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
oc get pods

NAME                                                     READY   STATUS    RESTARTS   AGE
memcached-operator-controller-manager-54c5864f7b-znwws   2/2     Running   0          14s
```

This means your operator is up and running. Next, let's create some custom resources via our operator.

Next, update your custom resource, by modifying the `config/samples/cache_v1alpha1_memcached.yaml` file
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

## 7. Test and verify

Update `config/samples/<group>_<version>_memcached.yaml` to change the `spec.size` field in the Memcached CR. This will increase the application pods from 3 to 5.

```bash
oc patch memcached memcached-sample -p '{"spec":{"size": 5}}' --type=merge
```

You can also update the spec.size from `OpenShift web console` by going to `Deployments` and selecting `memcached-sample` and increase/decrease using the up or down arrow:

![kubectl get all](images/inc-dec-size.png)

**Congratulations!** You've successfully deployed an operator using the `operator-sdk`!


## Cleanup

The `Makefile` part of generated project has a target called `undeploy` which deletes all the resource. You can run following to cleanup all the resources:

```bash
make undeploy
```

# License

This code pattern is licensed under the Apache Software License, Version 2.  Separate third party code objects invoked within this code pattern are licensed by their respective providers pursuant to their own separate licenses. Contributions are subject to the [Developer Certificate of Origin, Version 1.1 (DCO)](https://developercertificate.org/) and the [Apache Software License, Version 2](https://www.apache.org/licenses/LICENSE-2.0.txt).

[Apache Software License (ASL) FAQ](https://www.apache.org/foundation/license-faq.html#WhatDoesItMEAN)
