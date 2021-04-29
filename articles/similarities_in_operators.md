# The Operator Cookbook: How to make an operator 

In this article, we will discuss common building blocks for level 1 operators, and what logic a service vendor would need to write themselves in order
to build a level 1 operator. We will use the 
[Operator SDK Capability Levels](https://operatorframework.io/operator-capabilities/) as our guidelines for what is considered a 
level 1 operator.

By developing and deploying the [Memcached Operator](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md) 
and the [JanusGraph Operator](https://github.ibm.com/TT-ISV-org/operator/blob/main/articles/level-1-operator.md) we can 
analyze the similarities in the controller code, and think about what each operator must do at a high-level.

# Characteristics of an operator

An operator ensures that certain Kubernetes resources (the ones that are required to run your service) are created, and configured 
properly. It also relays status information back to the user, to communicate when certain resources are running.

In the Memcached example, we created a Deployment resource for the manager, which is the operator itself. And then, once we 
applied our custom resource using `kubectl` we created a Memcached Deployment, which is the operand, or the application we 
are deploying. Similarly, in the JanusGraph operator we create a StatefulSet, instead of a Deployment, and then create a service. 

Below are the main main characteristics of a level 1 operator that we will cover:

1. [Define the API](https://github.ibm.com/TT-ISV-org/operator/blob/main/articles/similarities_in_operators.md#the-api)
2. [Create Kubernetes resources if they do not exist](https://github.ibm.com/TT-ISV-org/operator/blob/main/articles/similarities_in_operators.md#the-main-logic-for-your-operator)
3. [Update replicas in your controller code](https://github.ibm.com/TT-ISV-org/operator/blob/main/articles/similarities_in_operators.md#replicas-should-be-set-in-the-cr-and-updated-in-the-controller-code)
4. [Update the status](https://github.ibm.com/TT-ISV-org/operator/blob/main/articles/similarities_in_operators.md#update-the-status)
5. [Scale up and down via custom resource](https://github.ibm.com/TT-ISV-org/operator/blob/main/articles/similarities_in_operators.md#ensure-operator-can-scale-up-and-down-via-the-custom-resource)

# The API
When building an operator, the easiest way to get started is by using the [Operator SDK](https://sdk.operatorframework.io/). Once you've 
finished the first steps such as using the [`operator sdk init`](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md#1-create-a-new-project-using-operator-sdk) and [`operator sdk create api`](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md#2-create-api-and-custom-controller) commands, you'll want to update the API.

This is where we design the structure of our custom resource. For simple cases, you'll likely use something like the `Size` and `Version` fields 
in the `Spec` section of your custom resource.

The Operator SDK generates the following code for your API:

```go
package v1alpha1
import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
type ExampleSpec struct {
	Foo string `json:"foo,omitempty"`
}
type ExampleStatus struct {
}
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Example struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MemcachedSpec   `json:"spec,omitempty"`
	Status MemcachedStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type ExampleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Example `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Example{}, &ExampleList{})
}
```

First, we will update the `Spec` section, like so:

```go
// ExampleSpec defines the desired state of Example database
type ExampleSpec struct {
	Size    int32  `json:"size"`
	Version string `json:"version"`
}
```

Next, we will update the `Status` section like so: 

```go
// ExampleStatus defines the observed state of Example database
type ExampleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Nodes []string `json:"nodes"`
}
```

And then for the last part of specifying the fields of your `Example` custom resource:

```go
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Example is the Schema for the example API
type Example struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExampleSpec   `json:"spec,omitempty"`
	Status ExampleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ExampleList contains a list of Example
type ExampleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Example `json:"items"`
}
```

That's it for your API.

# The main logic for your operator

The main logic when checking for different types of Kubernetes resources (Service, StatefulSet, Deployment, etc.) is as follows.

First, you will get a reference to a certain type of Kubernetes resource that you want to create:

```go	
found := &appsv1.Deployment{}
```
Then you use the Get function to find resources of that type in your namespace.

```go	
err = r.Get(ctx, req.NamespacedName, found)
```

The main logic is shown below, and this is similar no matter what resource you want to ensure is running (Deployment, StatefulSet, Service, PVC).

## Check if a resource exists, create one if it does not

First, we check that the error is not nil. If there is no error, that implies that the resource we want to create is already created, so we do not 
need to create another one. 

Next, we check for a `IsNotFound` error. This means that this resource doesn't exist at all, so we should create one. 

After we create, and the deployment or statefulset has been created successfully, then you can return and requeue. Otherwise, we 
return an error.

```go

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

```

 

## Replicas should be set in the CR and updated in the controller code

From our memcached example, we can can see that we set a variable to be what the `size` is from the custom resource. From there
we will check if the deployment's spec section has the same number of replicas as what is specified in the custom resource. 
If the numbers don't match, then we will update the Replicas. 

```go
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
```
## Update the status
The last thing we need to do in any operator is to update the status. This can be done by using the
reconciler [`Status().Update()`](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/client#StatusWriter.Update) function. We 
will show this below, but first, we need to format our Pods in a way in which we can quickly compare the current state with the desired 
state. In this example we use the reconciler `List` function to retrieve the pods which match our labels and are in our namespace. 

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
We then have a function which will return an array of strings. This is done so we can easily compare Pods in the current state 
versus the desired state.
```go
podNames := getPodNames(podList.Items)
```

Next, we update the status if we need to. We check if the podNames we've retrieved from the `List` function match the 
custom resource's status. The status which is defined in the API is as follows: 

```go
Nodes []string `json:"nodes"`
```
These nodes are the pod names which are currently in the cluster. If those nodes are different than the ones we've retrieved 
from the `List` function, then we will update the status using the reconciler [`Status().Update()`](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/client#StatusWriter.Update) function. 

```
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
Once we've updated the status, we are ready to test our operator!

## Ensure operator can scale up and down via the custom resource

The last thing to check is to make sure that you can scale up and scale down your operand via the custom resource. 

You can do this by changing the `size` value in your custom resource:

```go
apiVersion: cache.example.com/v1alpha1
kind: Memcached
metadata:
  name: memcached-sample
spec:
  size: 3

```
Change the size From 3, to 1.

```go
apiVersion: cache.example.com/v1alpha1
kind: Memcached
metadata:
  name: memcached-sample
spec:
  size: 1
```
Once you issue a `kubectl apply -f` command on the custom resource, you should see two pods terminating. As long as your 
application continues to work and is able to scale up and down via the custom resource, then you have a properly working Level 1 
operator. 

## Conclusion

Let's recap. To build an operator for your Kubernetes service, you will need to implement three main tasks.

1. Implement functions to check if the desired resource exists, and then create if it does not exist.
2. Set your replicas in your custom resource, and update them within your controller code. 
3. Update your status. This will communicate with the user what state the Pods are in. 

After that, you'll want to test your operator by scaling it up and down via the custom resource. If it can scale up and down 
successfully via custom resource, and your application still runs smoothly, then you are done in terms of a level 1 operator. 
**Congratulations!!** You understand the main concepts behind building a level 1 operator. Stay tuned for subsequent tutorials 
which will cover level 2 operators.  


