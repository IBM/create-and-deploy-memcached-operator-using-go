# Intro to Operators

This article will be the first in a series of articles and tutorials on learning how to build and deploy 
a Kubernetes Operator. This article will assume you have no knowledge on Kubernetes Operators, and will 
give you all of the basic knowledge needed to understand the different components and concepts of developing a
a Golang based operator. If you are already familiar with operators, you can either skim this article or 
simply skip ahead to the `SIMPLE_OPERATOR.md` file which shows how to develop and deploy your first operator to 
an OpenShift cluster.

## Expectations (What you have)
* You have little to no experience developing operators
* You have little to no knowledge on Kubernetes Operators concepts

## Expectations (What you want)
* You want to learn the basic concepts and steps needed to develop a Golang based operator to manage Kubernetes resources

## Estimated time
* This article should take roughly 15-30 minutes to complete, depending on how long you spend reading through 
documentation.

Before we start diving into operators, you need to know basic knoweldge about Kubernetes itself, since operators are just 
automating parts of app delivery and deployment on Kubernetes. 

### What is Kubernetes (a.k.a. K8s)
* A Kubernetes cluster is a collection of computers, called nodes 
  * All cluster work runs on one or more of the nodes
* The basic unit of work, and replication, is called a pod
  * A pod is one or more linux containers with common resources like networking, storage, and access to shared memory 
* A k8s cluster can be divided into two planes but sometimes these planes have some overlap. 
  * <b>The control plane</b>, which is k8s itself and implements the K8s API and cluster orchestration logic
  * The <b>app plane, or data plane</b>, is everything else. This is the group of pods where the app pods run
* One or more nodes are usually dedicated to running applications, while one or more nodes are reserved only for the control plane 
* Multiple replicas of control plane components can run on multiple pods to provide redundancy

To read more about Kubernetes, go to the official [Kuberenetes documentation](https://kubernetes.io/docs/concepts/overview/what-is-kubernetes/).

### How Kubernetes works
* Kubernetes automates the lifecycle of an app, such as a static web server 
* Without state, any of the app instances are interchangeable 
  * Because the server is not tracking the state, or storing input or data of any kind, when once instance fails, Kubernetes can replace it with another. 
* These instances are called replicas which are just copies of an app running on a cluster 
* The controllers of the control loop implement logic to automatically check for the difference between the actual state and the desired state 
  * When the two diverge, the controller takes action and makes sure the two match

## What are operators?
Operators are "software extensions to Kubernetes that make use of custom resources to amange applications and their 
components". You can read more about the operator pattern [here](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).
 


## What do operators do?
The main idea is that when the desired state and the actual state of the cluster diverge, operators have custom logic that will 
enable the app to be automatically installed, upgraded, recovered, analyzed, and scaled.* Operators wrap any necessary logic for deploying and operating a Kubernetes app using Kubernetes constructs. Here are a few more details you should understand about operators:

* They provision and manage the resources that you would normally need to use manually and since it is provided with enough authorization in the cluster, it can do cluster-management for you, such as rescheduling pods as they fail, or scaling the replica sets as needed. 
* They can help you in the selection of cloud resources from your cloud environment
* They can automatically provision storage, volume, and any other infrastructure you may need
* Operators are clients of the Kubernetes API that act as controllers for a custom resource
  * Operators are the app specific combo of a CR and a custom controller that does know all the details about starting, scaling, recovering, etc
  * The operators operand is what we call the app, service, or whatever resource an operator manage 
* CRDs are one of two essential building blocks for the most basic description of the operator pattern: a custom controller managing CRs 