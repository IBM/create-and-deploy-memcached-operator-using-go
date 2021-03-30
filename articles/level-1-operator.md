# Develop and Deploy a Level 1 JanusGraph Operator on OpenShift Container Platform
In this article, we will discuss how to develop and deploy a Level 1 operator on the OpenShift Container Platform. We will use the 
[Operator SDK Capability Levels](https://operatorframework.io/operator-capabilities/) as our guidelines for what is considered a 
level 1 operator.

### Develop and Deploy a Level 1 JanusGraph Operator using BerkeleyDB - Part 1
In part 1 of the tutorial, we will deploy JanusGraph using the default backend storage (BerkeleyDB). 
This is much more simple to deploy, since it will only use one Pod, and doesn't use any persistent volumes to our cluster. This approach is really only recommended for 
testing purposes, as noted in the [JanusGraph docs](https://docs.janusgraph.org/storage-backend/bdb/).
This approach will be much easier and quicker to get up and running, so we will start with this approach. Once we've gotten this approach working, we will move on to a more advanced approach (Part 2), using Cassandra as the backend storage.

### Develop and Deploy a Level 1 JanusGraph Operator using Cassandra - Part 2
In part 2 of the tutorial, we will update our operator to use Cassandra as the backend storage. 
This will enable use to replicate our data across multiple Pods, and give us high availability.  

## Expectations (What you have)
* You have some experience developing operators.
* You've finished the beginner and intermediate tutorials in this learning path, including  [Develop and Deploy a Memcached Operator on OpenShift Container Platform](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md).
* You've read articles and blogs on the basic idea of a Kubernetes Operators, and you know the basic Kubernetes resource types.

## Expectations (What you want)
* You want deep technical knowledge of how to implement a Level 1 operator.

## Outline 
<!-- 1. [What is a Level 1 Operator](#1-What-is-a-Level-1-Operator?)
1. [How should my operator deploy the operand?](#2-How-should-my-operator-deploy-the-operand?) -->

## What is a Level 1 Operator? 

According to the Operator Capability Levels, a Level 1 Operator is one which has ["automated application provisioning and configuration management"](https://sdk.operatorframework.io/docs/advanced-topics/operator-capabilities/operator-capabilities/#level-1---basic-install). 

Your operator should have the following features to be qualified as a Level 1 Operator
* Provision an application through a custom resource
* Allow **all** installation configuration details to be specified in the `spec` section of the CR 
* Should be possible to install the operator in multiple ways (kubectl, OLM, OperatorHub)
* All configuration files should be able to be created within Kubernetes 

## How should my operator deploy the operand? 

Your operator should have the following features when deploying an operand:

* The operator must wait for the operand to reach a healthy state
* The operator must use the `status` subresource of the custom resource to communicate with the user when the operand or application has reconciled.

## JanusGraph example 

Now that we understand at a high-level what an operator must do to be considered level 1, 
let's dive into our example. 

In our article, we will use JanusGraph as the example of the service we want to to create an
operator for. Currently, there is no JanusGraph operator on OperatorHub (as of March 23, 2021).
[JanusGraph](https://janusgraph.org/) is a distributed, open source, scalable graph database. 

**Important: JanusGraph is an example. The main ideas learned from JanusGraph are meant to be applied to any application or service you want to create an operator for.** 

## JanusGraph operator - Part 1

With that aside, let's understand what the JanusGraph operator must to do to successfully run JanusGraph on OpenShift. More specifically, we will show how 
to implement the below changes in the controller code which will run each time a change to the custom resource is observed. 

1. Create a Service if one does not exist.
2. Create a StatefulSet if ones does not exist.

These are the only two resources that our operator must create in order to get the default 
JanusGraph configuration (using BerkeleyDB) up and running. The reason that we create a 
headless service first is that our [StatefulSet needs to have a headless service](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#limitations) to be responsible for the network identity of the Pods. 

### What is a StatefulSet
A [StateFulSet](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) is the object that is used to manage stateful applications. Similar to a Deployment, a 
StatefulSet manages pods that are based on an identical container spec. The difference is that in a 
Deployment, pods are interchangeable. But in a StatefulSet, they are not - each has a unique identifier that
is maintained across any rescheduling. We will get into why this is important in part 2 of the tutorial. 

### Create the JanusGraph project and API  

At this point, we are familiar with using the Operator SDK to scaffold an operator for us. 

First, let's create our project directory: 

```bash
mkdir $HOME/projects/memcached-operator
cd $HOME/projects/memcached-operator
```

Next, let's create our project:

```bash
operator-sdk init --domain=example.com --repo=github.com/example/janusgraph-operator
```

For Go Modules to work properly, make sure you activate GO module support by running the following command:

```bash
export GO111MODULE=on
```

Now, create the api, with the `kind` being `Janusgraph`:

```bash
operator-sdk create api --group=graph --version=v1alpha1 --kind=Janusgraph --controller --resource

..
..
Writing scaffold for you to edit...
api/v1alpha1/memcached_types.go
controllers/memcached_controller.go
```

### Update the JanusGraph API

Next, let's update the API. Your `janusgraph_types.go` file should look like the following:

```go
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// JanusgraphSpec defines the desired state of Janusgraph
type JanusgraphSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Janusgraph. Edit Janusgraph_types.go to remove/update
	Size    int32  `json:"size"`
	Version string `json:"version"`
}

// JanusgraphStatus defines the observed state of Janusgraph
type JanusgraphStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Janusgraph is the Schema for the janusgraphs API
type Janusgraph struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JanusgraphSpec   `json:"spec,omitempty"`
	Status JanusgraphStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// JanusgraphList contains a list of Janusgraph
type JanusgraphList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Janusgraph `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Janusgraph{}, &JanusgraphList{})
}
```

As shown above, we've added the `Size` and `Version` fields to the `Spec`. We've also added the `Spec` and `Status` fields to the `Janusgraph` struct. This 
should be familiar to you if you've completed the [Develop and Deploy a Memcached Operator on OpenShift Container Platform](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md) tutorial. If you have not, that tutorial will offer more details about using the Operator SDK.

### Controller Logic - Creating a Service

Now, let's take a look at the heart of the Level 1 operator - the controller code. The first thing we must do at
a high-level to create an operator for JanusGraph, is to create a [headless service](https://kubernetes.io/docs/concepts/services-networking/service/#headless-services). A headless service is a service in which 
you do not specify the cluster IP. The service is used to control the network domain.

The first thing we will do 
in the controller code is to fetch the `Janusgraph` instance from our cluster.

```go
janusgraph := &graphv1alpha1.Janusgraph{}
err := r.Get(ctx, req.NamespacedName, janusgraph)
```

If we get any errors back from the `Get` request, such as errors reading the object, or a resource not found error, we will return (and requeue if we get errors reading the object). Otherwise, we will keep going and check for a service:

```go
serviceFound := &corev1.Service{}
err = r.Get(ctx, types.NamespacedName{Name: janusgraph.Name + "-service", Namespace: janusgraph.Namespace}, serviceFound)
```

We will use the [`errors.IsNotFound(err)`](https://pkg.go.dev/k8s.io/apimachinery@v0.19.2/pkg/api/errors#IsNotFound) function
to see if the service resource exists. If it does not, we will create one using the `serviceForJanusgraph(janusgraph)` function.

```go
if err != nil && errors.IsNotFound(err) {
    srv := r.serviceForJanusgraph(janusgraph)
    ...
}
```

### Service for JanusGraph

Let's look at the `serviceForJanusgraph(janusgraph)` function in more detail. The function signature is the following:

`func (r *JanusgraphReconciler) serviceForJanusgraph(m *v1alpha1.Janusgraph) *corev1.Service` which means that 
we will pass in a JanusGraph object, and return a `corev1.Service`. 

Below, you can see the full `serviceForJanusgraph` function:

```go
func (r *JanusgraphReconciler) serviceForJanusgraph(m *v1alpha1.Janusgraph) *corev1.Service {
	ls := labelsForJanusgraph(m.Name)
	srv := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-service",
			Namespace: m.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeLoadBalancer,
			Ports: []corev1.ServicePort{
				{
					Port: 8182,
					TargetPort: intstr.IntOrString{
						IntVal: 8182,
					},
					NodePort: 30184,
				},
			},
			Selector: ls,
		},
	}
	ctrl.SetControllerReference(m, srv, r.Scheme)
	return srv
}
```
### Labels for JanusGraph

Let's take it step by step. First, we create labels by calling the `labelsForJanusgraph` function:

```go
func labelsForJanusgraph(name string) map[string]string {
	return map[string]string{"app": "Janusgraph", "janusgraph_cr": name}
}
```
This function returns a map which looks like this:

```json
{
	"app": "Janusgraph",
	"janusgraph_cr": "<name>"
}
```
The way that a service works is that it will target any Pod with the `"app": "Janusgraph"` and `"janusgraph_cr": "<name>"` label, that is on the port 8182 (as shown in the code above).

### Configuring the service

Once we've created our labels, we will create the service using the [corev1.Service](https://pkg.go.dev/k8s.io/api/core/v1#Service) package. 

The service looks like the following: 

```go
srv := &corev1.Service{
	ObjectMeta: metav1.ObjectMeta{
		Name:      m.Name + "-service",
		Namespace: m.Namespace,
	},
	Spec: corev1.ServiceSpec{
		Type: corev1.ServiceTypeLoadBalancer,
		Ports: []corev1.ServicePort{
			{
				Port: 8182,
				TargetPort: intstr.IntOrString{
					IntVal: 8182,
				},
				NodePort: 30184,
			},
		},
		Selector: ls,
	},
}
```

Notice that at the top, we've filled out the [`ObjectMeta`](https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#ObjectMeta) 
field with the name and namespace of our custom resource. The `ObjectMeta` field is the metadata that we 
want to create with our service. Next, we fill out the heart of the service, which is the `Spec` field. In the `Spec` field, the 
package is expecting a [`corev1.ServiceSpec`](https://pkg.go.dev/k8s.io/api/core/v1#ServiceSpec), which contains the required 
fields of `Ports` and the optional `Selector` and `Type` fields. 

For the `Selector` field, we want to make sure to target only Pods that are part of our Janusgraph StatefulSet, so we do so by using the map returned from our `labelsForJanusgraph` function.

For our `Type` we 
create a [`ServiceTypeLoadBalancer`](https://pkg.go.dev/k8s.io/api/core/v1#ServiceType). Load balancers have an extra `NodePort`
field, which is set to `30184` in our case. 

Once we've finished configuring the service, we will return it the service, i.e. we will return a `corev1.Service` object.  

```go 
ctrl.SetControllerReference(m, srv, r.Scheme)
return srv
```

### Updating the cluster state

Once we've successfully created our service, we will use the `Create` function to save the `service` resources to our cluster. 

```go
srv := r.serviceForJanusgraph(janusgraph)
log.Info("Creating a new headless service", "Service.Namespace", srv.Namespace, "Service.Name", srv.Name)
err = r.Create(ctx, srv)
```

If we failed to create a service, we return an error. 

```go
if err != nil {
	log.Error(err, "Failed to create new service", "service.Namespace", srv.Namespace, "service.Name", srv.Name)
	return ctrl.Result{}, err
}
```

Otherwise, we return and requeue. 

```go
// Service created successfully - return and requeue
log.Info("Janusgraph service created, requeuing")
return ctrl.Result{Requeue: true}, nil
```

### StatefulSet for JanusGraph

Next, we will create a [StatefulSet](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) for JanusGraph. You will see that the code is very similar to that 
of creating a service for JanusGraph, other than some minor details with creating the StatefulSet object itself.
Note that instead of a deployment, we will use a StatefulSet, but this same logic can be applied to the deployment 
object.

First, we check to see if there are any StatefulSets in our cluster by using the `Get` function:

```go
found := &appsv1.StatefulSet{}
err = r.Get(ctx, types.NamespacedName{Name: janusgraph.Name, Namespace: janusgraph.Namespace}, found)
```

Next, we check for errors, as before. We want to make sure that no other StatefulSet resources exist in the 
cluster. If they do, then we do not need to create any, so we can return: 

```go
return ctrl.Result{}, nil
```

If there are no StatefulSet resources in the cluster, then we can go ahead and create one. We will call the 
`deploymentForJanusgraph(janusgraph)` function to create our deployment. 

### Understanding the deploymentForJanusgraph function

Let's dive into the `deploymentForJanusgraph(janusgraph)` function. It looks like the following:

```go
func (r *JanusgraphReconciler) deploymentForJanusgraph(m *v1alpha1.Janusgraph) *appsv1.StatefulSet {
	ls := labelsForJanusgraph(m.Name)
	replicas := m.Spec.Size
	version := m.Spec.Version

	dep := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			ServiceName: m.Name + "-service",
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
					Name:   "janusgraph",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "horeaporutiu/janusgraph:" + version,
							Name:  "janusgraph",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8182,
									Name:          "janusgraph",
								},
							},
							Env: []corev1.EnvVar{},
						}},
					RestartPolicy: corev1.RestartPolicyAlways,
				},
			},
		},
	}
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}
```

First, we get the labels as we did before with our service:

```go
ls := labelsForJanusgraph(m.Name)
```

Next, we grab the values from the `Spec`. These will determine how many pods to create, and which version of JanusGraph to deploy.

```go
replicas := m.Spec.Size
version := m.Spec.Version
```

Next, we have the heart of the function. This is when we will use the `appsv1` package to create our StatefulSet:

`dep := &appsv1.StatefulSet{`

We will create the metadata for the object as we did for the service:

```go
ObjectMeta: metav1.ObjectMeta{
	Name:      m.Name,
	Namespace: m.Namespace,
},
```

Next, in the `Spec` section of our StatefulSet, we use the `replicas` which
we set earlier. The replicas are coming from the user-entered values from the custom resource. We also use the labels which we've generated from our 
`labelsForJanusgraph` function and pass those into our `Selector` field. 
The `Selector` field defines how the StatefulSet finds which Pods to manage:

```go
Replicas: &replicas,
Selector: &metav1.LabelSelector{
	MatchLabels: ls,
},
```

Next, we pass in the service which we created earlier and create use the `corev1.PodTemplateSpec` to create our 
Pods:

```go
ServiceName: m.Name + "-service",
Template: corev1.PodTemplateSpec{
...
```

Then, we use the `corev1.PodSpec` and `corev1.Container`
package to specify which image we want to use to 
create our container. We will use our forked version of 
the JanusGraph Docker image, called `horeaporutiu/jansugraph`:

```go
Spec: corev1.PodSpec{
	Containers: []corev1.Container{
		{
			Image: "horeaporutiu/janusgraph:" + version,
			Name:  "janusgraph",
			...
		}
	}	
```

Then, we specify the container port, which is `8182` in this case. We found this from the JanusGraph documentation.

```go
Ports: []corev1.ContainerPort{
	{
		ContainerPort: 8182,
		Name:          "janusgraph",
	},
},
```

After we've specified all of the details of the 
StatefulSet, we return the object:

```go
ctrl.SetControllerReference(m, dep, r.Scheme)
return dep
```


### Updating the cluster state with our StatefulSet

Once we've successfully created our StatefulSet, we will use the `Create` function to save the `StatefulSet` resources to our cluster. 

```go
dep := r.deploymentForJanusgraph(janusgraph)
err = r.Create(ctx, dep)
```

If we failed to create a service, we return an error. 

```go
if err != nil {
	return ctrl.Result{}, err
}
```

Otherwise, we return and requeue. 

```go
// Deployment created successfully - return and requeue
log.Info("Deployment created, requeuing")
return ctrl.Result{}, nil
```

At the end of the reconcile function, we 
return and requeue:

```go
return ctrl.Result{}, nil
```

## Compile, Build and Push

Now, we will go ahead and login to our OpenShift cluster. 
You can follow the steps described in the previous 
tutorial. After you've logged in, go ahead and 
create a new project:

```bash
$ oc new-project JanusGraph-demo-project
Now using project "JanusGraph-demo-project" on server "https://c116-e.us-south.containers.cloud.ibm.com:31047".
```

### Edit the manager.yaml file

The `manager.yaml` file defines a Deployment manifest used to deploy the operator. That manifest includes a security context that tells Kubernetes to run the pods as a specific user (uid=65532). OpenShift already manages the users employed to run pods which is behavior the manifest should not override, so we will remove that from the manifest.

To do this, we can modify the `config/manager/manager.yaml` file to remove the following line:

```
runAsUser: 65532
```

### Create CRD and RBAC

Now that we have our controller code and API implemented, run the following command to implement the required Go type interfaces:

```bash
$ make generate
```

Once we've generated the code for our custom resource, we can use the `make manifests` command to generate CRD manifests and RBAC from KubeBuilder Markers:

```bash
$ make manifests
```

### Compile your Operator

To compile the code run the following command in the terminal from your project root:
```bash
$ make install
```

### Set the Operator Namespace

Next, we need to make sure to update our config to tell our operator to run in our own project namespace. Do this by issuing the following Kustomize 
commands:

```bash
$ export IMG=docker.io/<username>/janusgraph-operator:<version>
$ export NAMESPACE=<oc-project-name>

$ cd config/manager
$ kustomize edit set image controller=${IMG}
$ kustomize edit set namespace "${NAMESPACE}"
$ cd ../../

$ cd config/default
$ kustomize edit set namespace "${NAMESPACE}"
$ cd ../../
```

### Build and Push your Image

**Note:** You will need to have an account to a image repository like Docker Hub to be able to push your 
operator image. Use `Docker login` to login.

To build the Docker image run the following command:

```bash
$ make docker-build IMG=$IMG
```
and push the docker image to your registry using following from your terminal:

```bash
$ make docker-push IMG=$IMG
```

### Deploy the operator to OpenShift cluster

To Deploy the operator run the following command from your terminal:

```bash
$ make deploy IMG=$IMG
```

To make sure everything is working correctly, use the `oc get pods` command.

```bash
$ oc get pods

NAME                                                     READY   STATUS    RESTARTS   AGE
janusgraph-operator-controller-manager-54c5864f7b-znwws   2/2     Running   0          14s
```

This means your operator is up and running. Great job!


### Create the Custom Resource

Next, let's create the custom resource.

Update your custom resource, by modifying the `config/samples/graph_v1alpha1_janusgraph.yaml` file
to look like the following:

```yaml
apiVersion: graph.ibm.com/v1alpha1
kind: Janusgraph
metadata:
  name: janusgraph-sample
spec:
  # Add fields here
  size: 1
  version: latest
``` 
In the above code, we set the replicas to 1, and the 
version to `latest`. We aren't using the version
parameter currently in the controller code, but we will
in a later part of the tutorial.

And finally create the custom resources using the following command:

```bash
$ kubectl apply -f config/samples/graph_v1alpha1_janusgraph.yaml
```

#### Verify that resources are Running

From the terminal run `kubectl get all` or `oc get all` to make sure that controllers, managers and pods have been successfully created and is in `Running` state with the right number of pods as defined in the spec.

```bash
$ kubectl get all 
```

You should see one `janusgraph-sample` pod running.

### Part 1 Conclusion

**Congratulations!!**. You've just created a level 1 
operator for JanusGraph, using the default 
BerkeleyDB configuration. Great job! In the next 
section of the tutorial, we will show how to create a 
more complex level 1 operator for JanusGraph, using 
Cassandra as the backend storage. We will also show
how to scale the JanusGraph application up and down, 
using the custom resource. 

