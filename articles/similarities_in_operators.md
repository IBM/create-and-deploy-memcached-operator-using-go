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

Those are usually the main characteristics of an operator:

1. Create a Service if one does not exist
2. Create a StatefulSet (Or Deployment) if one does not exist
3. (optional) Create a PVC 
4. Update the status

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


