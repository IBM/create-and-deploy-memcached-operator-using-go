# Deep dive into Memcached Operator Code
In this article, we will discuss the low-level functions needed to write your own operator. This article builds off of the 
previous [Develop and Deploy a Memcached Operator on OpenShift Container Platform](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md) tutorial, so if you want the complete steps to develop and deploy the Memcached operator, view that tutorial. This 
article will discuss the Memcached custom controller code in depth.

## Expectations (What you have)
* You have some experience developing operators.
* You've finished the beginner tutorials in this learning path, including  [Develop and Deploy a Memcached Operator on OpenShift Container Platform](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md)
* You've read articles and blogs on the basic idea of a Kubernetes Operators, and you know the basic Kubernetes resource types.

## Expectations (What you want)
* You want deep technical knowledge of the code which enables operators to run.
* You want to understand how the reconcile loop works, and how you can use it to manage Kubernetes resources
* You want to learn more about the basic Get, Update, and Create functions used to save resources to your Kubernetes cluster.
* You want to learn more about KubeBuilder markers and how to use them to set role based access control.

## Outline
1. [Reconcile function overview](#1-reconcile-function-overivew)
1. [Understanding the Get function](#2-Understanding-the-get-function)
1. [Understanding the Reconcile function return types](#3-Understanding-the-reconcile-function-return-types)
1. [Create Deployment](#4-Create-deployment)
1. [Understanding the Update function](#5-Understanding-the-Update-function)
1. [Understanding KubeBuilder Markers](#6-Understanding-KubeBuilder-Markers)

In this article, we will cover the custom controller code for the Memcached Operator, found [here](https://github.ibm.com/TT-ISV-org/operator/blob/main/artifacts/memcached_controller.go). The complete code is shown below too, for convenience:

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

## 1. Reconcile function overview

The controller `[Reconcile]`(https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/reconcile#Reconciler) method contains the logic responsible for monitoring and applying the requested state for specific deployments. It does so by sending client requests to Kubernetes APIs, and will run every time a Custom Resource is modified by a user or changes state (ex. pod fails). If the reconcile method fails, it can be re-queued to run again.

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

## 2. Understanding the Get function

Next, we'll need to confirm that the `Memcached` resource is defined within our namespace.

This can be done using the [`Get` function](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Reader.Get), which retrieves an object from a Kubernetes cluster based on the arguments passed in. 

The function definition is the following: <b>Get(ctx context.Context, key types.NamespacedName, obj client.Object)</b>

### Understanding context in Go

The `Get` function expects the [context](https://pkg.go.dev/context#Context), the object key (which is the namespace, and the name of the object), and the object as arguments. Let's talk a bit about the context, since you will see it in many function calls in the controller code. The context carries a deadline, a cancellation signal, and other values across API boundaries. The context takes into account the
identity of the end user, auth tokens, and the request's deadline. 

To see your current context, run the following command: 

```bash
$ kubectl config view
```

You should see output like the following: 

```bash
apiVersion: v1
clusters:
- cluster:
    server: https://c116-e.us-south.containers.cloud.ibm.com:31047
  name: c116-e-us-south-containers-cloud-ibm-com:31047
contexts:
- context:
    cluster: c116-e-us-south-containers-cloud-ibm-com:31047
    namespace: test-tekton2-horea
    user: IAM#horea.porutiu@ibm.com
  name: test-tekton2-horea/c116-e-us-south-containers-cloud-ibm-com:31047/IAM#horea.porutiu@ibm.com
current-context: test-tekton2-horea/c116-e-us-south-containers-cloud-ibm-com:31047/IAM#horea.porutiu@ibm.com
kind: Config
preferences: {}
users:
- name: IAM#horea.porutiu@ibm.com
  user:
    token: REDACTED
```

Read more about context in Golang [here](https://blog.golang.org/context).

### Understanding Objects in Go

Now, let's understand the object that we pass into the `Get` function. The object has to implement the [Object interface](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Object) which means that it needs to embed both [runtime.Object](https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime#Object), and [metav1.Object](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Object). This object should be able to be written via YAML and then be created 
via `Kubectl create`. All of this means that the object will be treated like a Kubernetes native object. You will see that in the later 
parts of the code we will pass in a different type of resource (a Deployment) to the `Get` function. Since the `Get` function is accepting any Kubernetes 
object that implements the Object interface, it doesn't matter if our object is a custom resource we created (Memcached) or a native 
Kubernetes resource, like a [`Deployment`](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/).

The Reconcile function gives you two things, the context i.e. `ctx` and request i.e. `req`. 
The request parameter has all of the information we need to reconcile a Kubernetes object i.e. a
`memcached` object in our case. More specifically, the `req` struct contains the `NamespacedName` field which is the name and the namespace 
of the object to reconcile. That is what we will pass in to the `Get` function.

If the resource doesn't exist, we'll receive an error.
```go
err := r.Get(ctx, req.NamespacedName, memcached)
```

If the Memcached object does not exist in the namespace yet, the Reconciler will return an error and try again.
```go
return ctrl.Result{}, err
```

## 3. Understanding the Reconcile function return types

Now, let's talk a bit about what the [reconcile function](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/reconcile#Reconciler) returns. This can be a bit 
tricky since there are various return types. 

The function definition is the following: <b>Reconcile(ctx context.Context, req ctrl.Request) (Result, error)</b>

The reconcile function returns a (Result, err). Now, more specifically, 
the [Result struct](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/reconcile#Result) has two fields, the `Requeue` bool, which tells the reconcile 
function to requeue again. This bool defaults to false. The other field is 
`RequeueAfter` which expects a `time.Duration`. This tell the reconciler to requeue after a specific amount of time. 

For example the following code would requeue after 30 seconds.
```go
return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
```

Furthermore, the controller will requeue the request to be processed again if an error
is non-nil or `Result.Requeue` is true.

### Most common return types

Here are three of the most common return types:

1. `return ctrl.Result{Requeue: true}, nil` when you want to return and requeue the request. This is done usually when we have updated the state of the cluster, i.e. created a deployment, or updated the spec. 
2. `return ctrl.Result{}, err` when there is an error. This will requeue the request.
3. `return ctrl.Result{}, nil` when everything goes fine and you do not want to requeue. This is
the return at the bottom of the reconcile loop. This means the observed state of the 
cluster is the same as the desired state (i.e. the `MemcachedSpec` is the same as the `MemcachedStatus`).

Let's summarize. The main idea is that if we reach an error, or we update the state of the cluster, we will requeue. 
Otherwise, if the current state is the same as the desired state, we do not requeue. We wait for the desired state to
change instead.

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

## 4. Create Deployment

Assuming the resource is defined, we can continue on by observing the state of our Memcached Deployment. When we say "Memcached Deployment"
we are referring to the standard `Deployment` Kubernetes resource, but the difference is that that these deployments are created by the Custom Resource,
instead of a SRE or Kubernetes admin. 

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

If a deployment is not found, then we can use `Deployment` definition within the the `apps` package to create a new one using the reconciler [`Create`](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Writer) method:

```go

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

Creating a deployment, and more specifically creating a [`PodSpec`](https://pkg.go.dev/k8s.io/api/core/v1#PodSpec) is extremely important. Specifically the [`Image`](https://kubernetes.io/docs/concepts/containers/images/) and Ports field are important. 
In the code above, we are using the Docker Hub's Official [`Memcached Image`](https://hub.docker.com/_/memcached) and using version 1.4.36-alpine and we are exposing container port 11211 in our PodSpec.

### Use the Create function to save a new object to the cluster 

Once we create that deployment, we use the [`r.Create(ctx context.Context, obj client.Object)`](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Writer) function to save the 
object in the Kubernetes cluster. We use this function only if this object does not exist yet. If the object does exist and we want to save changes that we've made to it, we will use the Update() function. More on that soon.

The `r.Create(ctx context.Context, obj client.Object)` 
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

To summarize: using the [Create](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Writer) function is an important step in changing the current state of the 
cluster. The difference between Create and Update is that Create is used the first time when a user wants to create an object while [`Update`](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Writer) is used after the first time to update an object.


## 5. Understanding the Update Function

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
If all goes well, the spec is updated, and we requeue. Otherwise, we return an error. Again this is important. 
We always want to requeue after we update the state of the cluster. If the actual state is equal to the desired state,
then we do not have to requeue. 

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

This code is really important since it uses the [ListOption](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#ListOption) package which offers options for filtering results. In our case, we want to filter all the the pods which
are in our given namespace and have the same labels as our Memcached custom resource. Matching labels is important, since this
is how we will distinguish certain groups of pods from others.
```go
podList := &corev1.PodList{}
listOpts := []client.ListOption{
  client.InNamespace(memcached.Namespace),
  client.MatchingLabels(labelsForMemcached(memcached.Name)),
}
```
The filters we set in the previous `ListOpts` variable are passed into the List function, as a way to actually 
see which pods are currently in our namespace and also match the same labels as our custom resource. 

The [List](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Reader.List) function
will update the list which you pass into it, which in our 
case is the `podList` which we pass in.

```go
if err = r.List(ctx, podList, listOpts...); err != nil {
  log.Error(err, "Failed to list pods", "Memcached.Namespace", memcached.Namespace, "Memcached.Name", memcached.Name)
  return ctrl.Result{}, err
}
```

It will also give the `podList` variable a `.Items` field, which we will pass into getPodNames below.

GetPodNames converts the PodList returned from our List function into a string array, since that 
is what our `MemcachedStatus` struct is expecting, as we have defined it in our `memcached_types.go` file.

```go
podNames := getPodNames(podList.Items)

func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}
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

To summarize: the [Update](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Writer) function is a very important step in changing the actual state of the cluster. This is used if there is an object already existing in the cluster. 

## 6. Understanding KubeBuilder Markers

Now, one more thing to understand before we deploy our operator. Let's discuss the [Kubebuilder markers](https://book.kubebuilder.io/reference/markers.html) which you can see at the top of the file:

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

**Kubebuilder markers are tricky since they are written in comments.**

KubeBuilder markers, i.e. single-line comments which start with a plus, followed by a marker name enable config and code generation.
They are extremely important, especially when used for RBAC (role-based access control). The `controller-gen` utility, listed in 
your `bin` directory, is what actually code and YAML files from these markers. 

```go
// generate rbac to get, list, watch, create, update and patch the memcached status the nencached resource
// +kubebuilder:rbac:groups=cache.example.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
```

For example, the marker above tells us the following - for any `memcacheds` resources, within the `cache.example.com` API Group, 
the operator is able to get, list, watch, create, update, path, and delete these resources. Once we run `make manifests`, the `controller-gen` utility will see that we have a new KubeBuilder marker, and will update the rbac yaml files in our `config/rbac` directory to change our 
RBAC configuration.

For example, if our memcached resource didn't have the `List` verb listed in the kubebuilder marker, we would not be able to use r.List() on our memcached resource - we would get a permissions error such as `Failed to list *v1.Pod`. Once we change these markers and add the `list` command, we have to run `make generate` and `make manifests` and that will in turn apply the changes from our kubebuilder commands into our `config/rbac` yaml files. To 
learn more about KubeBuilder markers, see the docs [here](https://book.kubebuilder.io/reference/markers/rbac.html).


# License

This code pattern is licensed under the Apache Software License, Version 2.  Separate third party code objects invoked within this code pattern are licensed by their respective providers pursuant to their own separate licenses. Contributions are subject to the [Developer Certificate of Origin, Version 1.1 (DCO)](https://developercertificate.org/) and the [Apache Software License, Version 2](https://www.apache.org/licenses/LICENSE-2.0.txt).

[Apache Software License (ASL) FAQ](https://www.apache.org/foundation/license-faq.html#WhatDoesItMEAN)
