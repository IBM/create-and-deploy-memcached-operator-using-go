# Develop and Deploy a Level 1 JanusGraph Operator on OpenShift Container Platform
In this article, we will discuss how to develop and deploy a Level 1 operator on the OpenShift Container Platform. We will use the 
[Operator SDK Capability Levels](https://operatorframework.io/operator-capabilities/) as our guidelines for what is considered a 
level 1 operator.

## Expectations (What you have)
* You have some experience developing operators.
* You've finished the beginner and intermediate tutorials in this learning path, including  [Develop and Deploy a Memcached Operator on OpenShift Container Platform](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md).
* You've read articles and blogs on the basic idea of a Kubernetes Operators, and you know the basic Kubernetes resource types.

## Expectations (What you want)
* You want deep technical knowledge of how to implement a Level 1 operator.

## Outline
1. [What is a Level 1 Operator](#1-What-is-a-Level-1-Operator?)
1. [How should my operator deploy the operand?](#2-How-should-my-operator-deploy-the-operand?)

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

**Important: JanusGraph is just an example. The main ideas learned from JanusGraph are meant to be applied to any application or service you want to create an operator for.** 

## JanusGraph operator  

With that aside, let's understand what the JanusGraph operator must to do to successfully run JanusGraph on OpenShift. More specifically, we will show how 
to implement the below changes in the controller code which will run each time a change to the custom resource is observed. 

1. Create a service if one does not exist.
2. Create a deployment (or statefulset) if ones does not exist.
3. Create persistent volume and or persistent volume claims if it does not exist 

## Create the JanusGraph project and API  

Let's dive into each a bit more deeply. At this point, we are familiar with using the Operator SDK to scaffold an operator for us. 

First, let's create our project directory: 

```bash
mkdir $HOME/projects/memcached-operator
cd $HOME/projects/memcached-operator
```


First, let's create our project:

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

## Update the JanusGraph API

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

All we've done is add the Size and Version fields to the `Spec`. We've also added the spec and status fields to the `Janusgraph` struct. This 
should be familiar to you if you've completed the [Develop and Deploy a Memcached Operator on OpenShift Container Platform](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md) tutorial. If you have not, that tutorial will offer more details about using the Operator SDK.

## Controller logic - creating a service

Now, let's take a look at the heart of the Level 1 operator - the controller code.

