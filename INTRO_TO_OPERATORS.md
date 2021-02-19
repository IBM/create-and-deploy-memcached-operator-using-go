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

## What is Kubernetes (a.k.a. K8s)
* A Kubernetes cluster is a collection of computers, called nodes 
  * All cluster work runs on one or more of the nodes
* The basic unit of work, and replication, is called a pod
  * A pod is one or more linux containers with common resources like networking, storage, and access to shared memory 
* A k8s cluster can be divided into two planes but sometimes these planes have some overlap. 
  * <b>The control plane</b>, which is k8s itself and implements the K8s API and cluster orchestration logic
  * The <b>app plane, or data plane</b>, is everything else. This is the group of pods where the app pods run
* One or more nodes are usually dedicated to running applications, while one or more nodes are reserved only for the control plane 
* Multiple replicas of control plane components can run on multiple pods to provide redundancy

To read more about Kubernetes including the official definition, go to the official [Kuberenetes documentation](https://kubernetes.io/docs/concepts/overview/what-is-kubernetes/).

## How Kubernetes works
* Kubernetes automates the lifecycle of an app, such as a static web server 
* Without state, any of the app instances are interchangeable 
  * Because the server is not tracking the state, or storing input or data of any kind, when once instance fails, Kubernetes can replace it with another. 
* These instances are called replicas which are just copies of an app running on a cluster 
* The controllers of the control loop implement logic to automatically check for the difference between the actual state and the desired state 
  * When the two diverge, the controller takes action and makes sure the two match

To learn more about how Kuberentes works, read [this blog](https://sensu.io/blog/how-kubernetes-works#:~:text=Kubernetes%20keeps%20track%20of%20your,storage%2C%20and%20CPU%20when%20necessary.).

## What are operators?
Operators are "software extensions to Kubernetes that make use of custom resources to amange applications and their 
components". You can read more about the operator pattern [here](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).
 
## What do operators do?
The main idea is that when the desired state and the actual state of the cluster diverge, operators have custom logic that will 
enable the app to be automatically installed, upgraded, recovered, analyzed, and scaled. 

<b> The advantage of operators lies in their 
automation. Usually, a SRE (Site reliability engineer) would have to take care of recoving an application if it crashes, or upgrading to 
a later version of an application. But with an operator, all of this can be automated. </b>

Operators wrap any necessary logic for deploying and operating a Kubernetes app using Kubernetes constructs. Here are a few more details you should understand about operators:

* They provision and manage the resources that you would normally need to use manually and since it is provided with enough authorization in the cluster, it can do cluster-management for you, such as rescheduling pods as they fail, or scaling the replica sets as needed. 
* They can help you in the selection of cloud resources from your cloud environment
* They can automatically provision storage, volume, and any other infrastructure you may need
* Operators are clients of the Kubernetes API that act as controllers for a custom resource
  * Operators are the app specific combo of a CR and a custom controller that does know all the details about starting, scaling, recovering, etc
  * The operators operand is what we call the app, service, or whatever resource an operator manage 
* CRDs are one of two essential building blocks for the most basic description of the operator pattern: a custom controller managing CRs

To learn more, read this [article from Red Hat](https://www.redhat.com/en/topics/containers/what-is-a-kubernetes-operator#:~:text=A%20Kubernetes%20operator%20is%20a,and%20managing%20a%20Kubernetes%20application.&text=A%20Kubernetes%20operator%20is%20an,behalf%20of%20a%20Kubernetes%20user.) outlining what a kubernetes operator is, and what they do.

### Stateful vs. Stateless Apps
* In stateless deployments, the order of deploying pods, their labels, network address or port, storage class, or volume are not important. You keep them if they are healthy and serving, you dispose of them when they become unhealthy, outdated, or just no longer needed and replace them as necessary. <b>You do not need an operator for stateless applications.</b>
* In stateful apps, some order is necessary. You also need to add storage and persistent volume so that the state is saved, and the cluster admin has to manage that. 
* <b>The majority of applications are stateful. This is where Kubernetes Operators are helpful.</b>

To learn more about Stateful vs. Stateless apps, read [this article](https://www.redhat.com/en/topics/cloud-native-apps/stateful-vs-stateless) from Red Hat.

## Why does Kubernetes need operators?

Kuberenetes needs operators for stateful deployments. This is because we can automate manual tasks such as setting configuration flags, 
and changing runtime configuration that is needed for many stateful applications. Read more about why Kubernetes needs operators in this [blog](https://kublr.com/blog/understanding-kubernetes-operators/).

## Operator Code - the Controller and the API
Now, let's start exploring the heart of the operator - the controller code. But before we do that we must understand custom 
resources, and custom resource definitions, since that is what we will use to create our operator.

### Custom Resources(CRs) - Custom API endpoints
A resource is an endpoint in the Kubernetes API that stores a collection of API objects of a certain kind; for example, the built-in pods resource contains a collection of Pod objects. A 
custom resource definition is a schema that defines the fileds and types of fields within an
instance of a Custom Resource. We will use a Custom Resource Definition as one of the initial 
steps when developing our operator. Here are some more useful details about Custom Resources:

* CRs hold structured data, and the Kubernetes API server provides a mechanism for reading and setting their files as you would those in a native resource, by using kubectl or another API client. 
* CRs are most useful when they are watched by <b>custom controller code that creates, updates, deletes other cluster objects or even resources outside of the cluster</b>
* CRs obey k8s conventions, like the resources .spec and .status 
  * A CR is a named endpoint in the Kubernetes API. Users define a CR by providing a custom-resource definition or CRD.
  * The CR’s group, version, and kind together form the fully qualified name of a Kuberentes resource type. That name must be unique across a cluster. 

To learn more about Custom Resources, refer to the official [Kubernetes documentation](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/).

### Custom Controllers (the code that watches your Custom Resource)
When you combine the custom resource with a custom controller that watches that custom resource, 
you get what is called a "declaritive API" by Kubernetes. 

* A declaritive API is one in which the user can "declare" or "specify" (keep in mind <b>Spec</b> for later, since that is important) the
desired state of the cluster, and Kubernetes will try to change whatever resources are needed to 
be changed in order to reach the "desired state" of the cluster.
  * For example if we "specify" or "declare" that we want 4 replicas of a certain resource, and currently we only have 3, Kuberentes will automatically create that 4th replica for us.

To learn more about custom controller, and if you should use one, read the official Kubernetes 
documentation on custom controller, [here](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/#custom-controllers).


## Operator SDK
The Operator SDK is a framework that uses the controller-runtime library to make writing operators easier by providing:

* High level APIs and abstractions to write the operational logic more intuitively
* Tools for scaffolding and code generation to bootstrap a new project fast
* Extensions to cover common Operator use cases

For the purposes of our learning path, the operator SDK will be used to scaffold our code.




<!-- What are operators? -->
<!-- What do operators do? Explain the advantages promised by using operators -->

<!-- Why does Kubernetes need operators? Explain why we need operators -->
TODO: 
<!-- 1. Describe the code in an operator – controller and API (what it does, not how to implement it) -->
2. Introduction to operator capability levels
<!-- 3. Kubernetes Operator SDK -->

The information in this article can be found in a few different sources:

* Kubernetes Operators by Jason Doies and Joshua Wood (O'Reilly)

* kublr.com/blog/understanding-kubernetes-operators

* kubernetes.io/docs/concepts/extend-kubernetes/operator