# Explanation of Memcached operator code

This article offers a detailed look at how the Memcached custom controller code works, describing the logic of the custom controller code from the [Develop and Deploy a Memcached Operator on OpenShift Container Platform](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md) tutorial. If you understand the low-level functions needed to write your own operator, you will be able to develop a complex operator yourself. <!--EM:  Link to advanced article where we show them how to write their own logic for an operator.-->

Read this article to gain deep technical knowledge about:

* The code that enables operators to run
* How the Reconcile loop works and how you can use it to manage Kubernetes resources
* Basic Get, Update, and Create functions used to save resources to your Kubernetes cluster
* KubeBuilder markers and how to use them to set role-based access control (RBAC).

As a reminder, a *controller* is the core part of Kubernetes that ensures that an object's actual state matches the object's desired state.

## Prerequisites

* Read the accompanying article, [Develop and deploy a Memcached operator on Red Hat OpenShift Container Platform](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md).

## Outline
1. [Reconcile function overview](#1-reconcile-function-overview)
1. [Understanding the Get function](#2-understanding-the-get-function)
1. [Understanding the Reconcile function return types](#3-understanding-the-reconcile-function-return-types)
1. [Create deployment](#4-create-deployment)
1. [Understanding the Update function](#5-understanding-the-update-function)
1. [Understanding KubeBuilder Markers](#6-understanding-kubebuilder-markers)

## Examine the code

This article details the custom controller code for the Memcached Operator, found in our [GitHub repo](https://github.ibm.com/TT-ISV-org/operator/blob/main/artifacts/memcached_controller.go). The complete code is shown below, too, for convenience:

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

<!--EM: One thing I'm wondering here & throughout is: Do we need to chunk up that code listing above and break it down here for the reader? So here should we start with , grab the specific block of related code from the above long one, and post it here. What do you think?-->

The controller's [Reconcile](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/reconcile#Reconciler) method contains the logic responsible for monitoring and applying the requested state for specific deployments. The Reconciler sends client requests to Kubernetes APIs and runs every time a custom resource is modified by a user or changes state (for example, if a pod fails). If the Reconcile method fails, it can be re-queued to run again.

After scaffolding your controller via the `operator-sdk`, you have an empty Reconciler function.

In this example, the Reconciler should:<!--EM: In what example? from the code listing above?-->

1. Check for an existing memcached deployment, and create one if it does not exist.
2. Retrieve the current state of the memcached deployment and compare it to the desired state. More specifically, the method compares the memcached deployment `ReplicaSet` value to the `Size` parameter that is defined in the `memcached_types.go` file.
3. Ensure the `ReplicaSet` value matches the `Size` parameter. If the number of pods in the deployment `ReplicaSet` does not match the provided size, then the Reconciler updates the `ReplicaSet` value and re-queues the Reconciler until the desired state is achieved.

In the code, logic is added to the empty Reconciler function. First, reference the instance you want to observe. In this code, it's the `Memcached` object defined in our `api/v1alpha1/memcached_types.go` file. Do this by retrieving the Memcached CRD from the `cachev1alpha1` object, which is listed in the import statements. Note that the trailing endpoint of the URL maps to the files in the `/api/v1alpha1/` directory.

```go
import (
  ...
  cachev1alpha1 "github.com/example/memcached-operator/api/v1alpha1"  
)
```

Here, `cachev1alpha1.<Object>{}` is used to reference any of the defined objects within that `memcached_types.go` file.

```go
memcached := &cachev1alpha1.Memcached{}
```

## 2. Get function overview
<!--EM: again, I wonder if it would help/hurt/make it too long to have the code listing part that you're talking about here referenced-->

Use the [`Get` function](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Reader.Get) to confirm that the `Memcached` resource is defined within your namespace. This function retrieves an object from a Kubernetes cluster based on the arguments that are passed in.

The function definition is the following: <b>Get(ctx context.Context, key types.NamespacedName, obj client.Object)</b>.

### Understanding the Get function's context in Go

<!--EM: is the "context" referenced here referring to the "context" line of code (line 8) in the big code listing?-->
The `Get` function expects the objects and the  [context](https://pkg.go.dev/context#Context) as arguments. *Context* refers to the object key that is the namespace and the name of the object. These context arguments are in many function calls in the controller code, so let's take a closer look.

The context carries a deadline, a cancellation signal, and other values across API boundaries. The context takes into account the identity of the end user, auth tokens, and the request's deadline.

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
<!--I think it might be helpful to wrap this part up by identifying what the output code above shows as it relates to the definition of the context that you had described above this code: "The context carries a deadline, a cancellation signal, and other values across API boundaries. The context takes into account the identity of the end user, auth tokens, and the request's deadline"-->
Read more about context in Golang [here](https://blog.golang.org/context).

### Understanding objects in Go

An object passed into the `Get` function must implement the [Object interface](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Object), which means that it needs to embed both [runtime.Object](https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime#Object), and [metav1.Object](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Object). This object is written via YAML and then created via `Kubectl create`. As such, the object is treated like a Kubernetes-native object.

In the later parts of the code listing above<!--EM: Should we show the snippet here in context?-->, we passed in a different type of resource (a Deployment) to the `Get` function. Because the `Get` function accepts any Kubernetes object that implements the object interface, it doesn't matter if our object is a custom resource (Memcached) or a native Kubernetes resource like a [`Deployment`](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/).

The Reconcile function produces two things:

* the context (`ctx`)
* the request (`req`).

The `request` parameter includes the information needed to reconcile a Kubernetes object. In this code example, that is the `memcached` object. More specifically, the `req` struct contains the `NamespacedName` field which is the name and the namespace of the object to reconcile. This `NamespacedName` is what gets passed into the `Get` function.

If the resource doesn't exist, you receive an error like the following.

```go
err := r.Get(ctx, req.NamespacedName, memcached)
```

If the Memcached object does not exist in the namespace yet, the Reconciler will return an error and try again.

```go
return ctrl.Result{}, err
```

## 3. Reconcile function return types

The [Reconcile function](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/reconcile#Reconciler) can produce various return types.

The function definition is <b>Reconcile(ctx context.Context, req ctrl.Request) (Result, error)</b>.

The reconcile function returns a `(Result, err)`.

The first field <!--EM first field in what? The reconcile function? or the return types??-->is the [`Result` struct](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/reconcile#Result) which has two fields, `Requeue` and `RequeueAfter`.

* `Requeue` is a boolean data type that tells the reconcile function to requeue again. This data type defaults to "false".
* `RequeueAfter` expects a `time.Duration` that tells the reconciler to requeue after a specific amount of time.

For example the following code requeues after 30 seconds.

```go
return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
```

Furthermore, the controller requeues the request again if the error is not `nil` or `Result.Requeue` is true.

### Most common return types

Three of the most common return types include:

1. `return ctrl.Result{Requeue: true}, nil` often occurs when the state of the cluster or spec is updated. This type returns and requeues the request.
2. `return ctrl.Result{}, err` occurs when there is an error and requeues the request.
3. `return ctrl.Result{}, nil` occurs when the function is successful and the function doesn't need to requeue. This type occurs at the bottom of the reconcile loop, when the observed state of the cluster matches the desired state. In our code, this happens when the `MemcachedSpec` is the same as the `MemcachedStatus`.

To summarize, if the `Reconcile` function returns an error or if the state of the cluster is updated, the process requeues. If the current state is the same as the desired state, there is no need to requeue.

At this point, the Reconciler function above looks like:

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

## 4. Create Deployment <!--EM: Should this be "Memcached Deployment" or "Observe your memcached deployment"-->

If the resource is defined, you can Observe the state of your Memcached Deployment. *Memcached Deployment*
refers to the standard `Deployment` Kubernetes resource. In OpenShift, the custom resource creates these deployments, instead of a SRE or Kubernetes administrator.

First, use the [k8s.io/api/apps/v1](https://godoc.org/k8s.io/api/apps/v1#Deployment) package, defined in your import statement, to confirm that a Memcached deployment exists within the namespace: <!--EM This is where it starts to feel like a tutorial more than a code deep dive. I tried ot rewrite it, but now I'm worried that the parts I changed were in fact meant to be something the reader typed/inputted-->

```go
import (
	appsv1 "k8s.io/api/apps/v1"
  ...
)
```

The `apps` package references a [Deployment object](https://pkg.go.dev/k8s.io/api/apps/v1#Deployment). Note that a deployment object is a Kubernetes object which implements the [Object interface](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Object).

The reconciler `Get` function checks whether the Memcached deployment exists with the provided name within your namespace.

```go
found := &appsv1.Deployment{}
err = r.Get(ctx, req.NamespacedName, found)
```

If a deployment is not found, use the `Deployment` definition within the the `apps` package to create a new one using the reconciler [`Create`](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Writer) method:

```
go
if err != nil && errors.IsNotFound(err) {
  dep := r.deploymentForMemcached(memcached)
  log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
  err = r.Create(ctx, dep)
  ...
  // if successful, return and re-queue Reconciler method
  return ctrl.Result{Requeue: true}, nil
```

For improved readability, the deployment definition is in a different function called [`deploymentForMemcached`](https://github.ibm.com/TT-ISV-org/operator/blob/main/artifacts/memcached_controller.go#L134). This function includes the pod runtime specs (ports, startup command, image name), and the `Memcached.Spec.Size` value to determine how many replicas should be deployed. This function returns the deployment resource -- a Kubernetes object.

```
go
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

Creating a deployment and, more specifically, creating a [`PodSpec`](https://pkg.go.dev/k8s.io/api/core/v1#PodSpec) is extremely important. The [`Image`](https://kubernetes.io/docs/concepts/containers/images/) and Ports field are important.

The code above uses the Docker Hub's Official [`Memcached Image`](https://hub.docker.com/_/memcached) and version 1.4.36-alpine and exposes container port 11211 in the PodSpec.

### Use the Create function to save a new object to the cluster

After creating the deployment, using the [`r.Create(ctx context.Context, obj client.Object)`](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Writer) function saves the object in the Kubernetes cluster. This function is only used if this object does not exist yet. If the object *does* exist, using the `Update()` function saves any changes.

The `r.Create(ctx context.Context, obj client.Object)` function takes in the context (which is passed into the `Reconcile` function) and the Kubernetes object that needs to be saved (the deployment we just created) in the `deploymentForMemcached` function:

```go
dep := r.deploymentForMemcached(memcached)
log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
err = r.Create(ctx, dep)
```

Since an update was made to our cluster, the function requeues:

```go
return ctrl.Result{Requeue: true}, nil
```

To summarize, using the `Create` function changes the current state of the cluster by creating an object. `Update` is used to update an already-created object.

## 5. Overview of the Update function

The next part of the code adds logic to our method to adjust the number of replicas in our deployment whenever the `Size` parameter is adjusted. This assumes a deployment already exists in our namespace. Specifically, this changes the desired state of our cluster to match the desired state of the Custom Resource.

### Use Update() to save the state after modifying an existing object

First, request the `Size` field from our Memcached Custom Resource and then compare the desired size to the number of replicas running in the deployment. If the numbers of replicas isn't the same as the desired `Size` from our Memcached Spec, we'll use the `Update` method to adjust the amount of replicas in the deployment to be the same as the desired `Size` from our Memcached Spec.

The Update(ctx context.Context, obj Object) function has a similar function definition to Create(), except that we must pass in a struct pointer to the object we want to update. In our case, this is the Memcached Deployment resource we created in the `deploymentForMemcached` function.  

```go
found := &appsv1.Deployment{}
...
size := memcached.Spec.Size
if *found.Spec.Replicas != size {  
  found.Spec.Replicas = &size
  err = r.Update(ctx, found)
  ...
}

```

In this snippet of code, the CR effectively changes the desired state by setting the deployment's replicas value to match the value that the user sets in the CR. This changes the desired state of the cluster to match the desired state of the CR.

If all goes well, the spec is updated, and the [] requeues<!--EM: What requees? We don't. But what does--the method? the code? the function? -->. Otherwise, an error is returned.
You always want to requeue after you update the state of the cluster. If the actual state is equal to the desired state,
then we do not have to requeue.

```go
if err != nil {
  log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
  return ctrl.Result{}, err
}
// Spec updated - return and requeue
return ctrl.Result{Requeue: true}, nil
```

### Update the Status to save the current state of the cluster

To save the current state of the cluster, modify the `Status` subresource of
our Memcached object using the [`StatusClient`](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#StatusClient.Status) interface.

First let's review what type our status subresource is, according to the `api` which we created.<!--EM: Where was the API? -->

The `Status` struct in our code looks like the following:

```go
type MemcachedStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Nodes []string `json:"nodes"`
}
```

In this listing, the `Status` subresource expects an array of strings which represent the current list of pods in our namespace.

Use the [`List`](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Reader.List) function to retrieve the list of pods in a specific namespace.

<!-- The r.List function will create a `.Items` field in the
ObjectList we pass in which will be populated with the objects for a given namespace. -->

This code is important because it uses the [ListOption](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#ListOption) package which offers options for filtering results. In our case, filter all the the pods which
are in our given namespace and have the same labels as our Memcached CR. Matching labels is important because it  distinguishes certain groups of pods from others.

```go
podList := &corev1.PodList{}
listOpts := []client.ListOption{
  client.InNamespace(memcached.Namespace),
  client.MatchingLabels(labelsForMemcached(memcached.Name)),
}
```

The filters we set in the previous `ListOpts` variable are passed into the `List` function to show which pods are currently in our namespace and also match the same labels as our CR.

The [List](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Reader.List) function
takes in a context, a list, and the list options. In our example, we pass in the `podList` and `listOpts` objects
which return the list of pods in our namespace that have the same labels as our Memcached resource.

```go
if err = r.List(ctx, podList, listOpts...); err != nil {
  log.Error(err, "Failed to list pods", "Memcached.Namespace", memcached.Namespace, "Memcached.Name", memcached.Name)
  return ctrl.Result{}, err
}
```

After the `List` function returns, it creates an `.Items` field in our `podList` object. We pass that field
into our `getPodNames` function, as shown below.

`getPodNames` converts the `podList` returned from our `List` function into a string array, since that
is how we defined the `MemcachedStatus` struct in the `memcached_types.go` file.

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

Lastly, we check if the `podNames` that we listed from `r.List` are the same
as the `memcached.Status.Nodes`. If they are not the same, we use [`Update(ctx context.Context, obj Object)` function](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Writer) to update the `MemcachedStatus` struct:

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

Updating the status updates the current state of the cluster. To reiterate:

* When the `Spec` is updated, the desired state is updated.
* When the `Status` is updated, the current state of the cluster is updated.

If all goes well, the function runs without an error. This means that the current state of the cluster is the same as the desired state, so no reconciliation is needed until the desired state changes again.

```go
return ctrl.Result{}, nil
```

In summary, the `Update` function is an important step in changing the state of the cluster. This function allows you to save the desired state when you update the Spec and it allows you to save the current state when you update the Status.

## 6. Understanding KubeBuilder Markers

Finally, let's discuss the [KubeBuilder markers](https://book.kubebuilder.io/reference/markers.html) which you can see at the top of the file:

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

*KubeBuilder markers* are single-line comments which start with a plus, followed by a marker name to enable config and code generation.

These markers are extremely important, especially when used for RBAC (role-based access control). The `controller-gen` utility, listed in your `bin` directory, is what actually generates code and YAML files from these markers.


```go
// generate rbac to get, list, watch, create, update and patch the memcached status the nencached resource
// +kubebuilder:rbac:groups=cache.example.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
```

The marker above tells the reader: For any `memcacheds` resources within the `cache.example.com` API Group,
the operator is able to get, list, watch, create, update, path, and delete these resources.

If you run `make manifests`, the `controller-gen` utility sees the new KubeBuilder marker and updates the RBAC YAML files in the `config/rbac` directory to change the RBAC configuration.

For example, if our memcached resource didn't have the `List` verb listed in the KubeBuilder marker, we would not be able to use `r.List()` on our memcached resource. Instead, we would get a permissions error such as `Failed to list *v1.Pod`. If you change these markers and add the `list` command, you must run `make generate` and `make manifests` in order to apply the changes from your KubeBuilder commands into your `config/rbac` YAML files. To
learn more about KubeBuilder markers, see the [Kubebuilder docs](https://book.kubebuilder.io/reference/markers.html).

## Conclusion

This article gave you a better understanding of the underlying logic of the custom controller code from the [Develop and Deploy a Memcached Operator on OpenShift Container Platform](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md) tutorial.

Hopefully you have a better understanding of how to:
* Use the Go Client [Reader](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Reader.Get) and [Writer](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#Writer) interface to Get, Create, Update, and List our resources.
* Use the [StatusWriter](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/client#StatusWriter) interface to update the status of a subresource, for example, the current state.
* Automate the deployment of a Memcached service, ensure your deployment is up and that the number of replicas in that deployment is the same as the number that is listed in your CR.
* Use KubeBuilder markers to change role-based access control policies and apply those
policies to your CR.
