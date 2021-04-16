# Understand the building blocks of level 1 operators
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
finished the first steps such as using the [`operator sdk init`](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md#1-create-a-new-project-using-operator-sdk) and [`operator sdk create api`]() commands, 

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


