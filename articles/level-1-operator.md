# Develop and Deploy a Level 1 JanusGraph Operator on OpenShift Container Platform
In this article, we will discuss how to develop and deploy a Level 1 operator on the OpenShift Container Platform. We will use the 
[Operator SDK Capability Levels](https://operatorframework.io/operator-capabilities/) as our guidelines for what is considered a 
level 1 operator.

, based on the guidelines from the [Operator SDK Capability Levels](https://operatorframework.io/operator-capabilities/)

the low-level functions needed to write your own operator. This article builds off of the 
previous [Develop and Deploy a Memcached Operator on OpenShift Container Platform](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md) tutorial, so if you want the complete steps to develop and deploy the Memcached operator, view that tutorial. This 
article will discuss the Memcached custom controller code in depth.

## Expectations (What you have)
* You have some experience developing operators.
* You've finished the beginner and intermediate tutorials in this learning path, including  [Develop and Deploy a Memcached Operator on OpenShift Container Platform](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md).
* You've read articles and blogs on the basic idea of a Kubernetes Operators, and you know the basic Kubernetes resource types.

## Expectations (What you want)
* You want deep technical knowledge of the code which enables operators to run.
* You want to understand how the reconcile loop works, and how you can use it to manage Kubernetes resources.
* You want to learn more about the basic Get, Update, and Create functions used to save resources to your Kubernetes cluster.
* You want to learn more about KubeBuilder markers and how to use them to set role based access control.

## Outline
1. [What is a Level 1 Operator](#1-What-is-a-Level-1-Operator?)
1. [Understanding the Get function](#2-Understanding-the-get-function)

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

Let's dive into each a bit more deeply. At this point, we are familiar with using the Operator SDK to create  