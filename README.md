# Kubernetes Operators Learning Path

Welcome to IBM Developer's Kubernetes Operators learning path! In this series of articles and tutorials, you will learn how to create and 
deploy a Golang based operator. You will also learn all of the foundational Kubernetes knowledge needed to understand how to develop
a Golang based operator from scratch. 

## Beginner level
1. [Intro to Kubernetes Operators](https://github.ibm.com/TT-ISV-org/operator/blob/main/INTRO_TO_OPERATORS.md): This article explains what operators 
are and why they are useful.

2. [Anatomy of an operator, demystified](https://github.ibm.com/TT-ISV-org/operator/blob/main/articles/demystified.md): In this article we will build upon the [Intro to Operators](https://github.ibm.com/TT-ISV-org/operator/blob/main/INTRO_TO_OPERATORS.md) article and explore Kubernetes concepts such as workloads, controllers, custom resources, and the control loop. This article will explain how operators extend
Kubernetes functionality. 

## Intermediate level

1. [Develop and Deploy a Memcached Operator on OpenShift Container Platform](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md): 
In this tutorial we will start by ensuring we have our [environment setup](https://github.ibm.com/TT-ISV-org/operator/blob/main/installation.md) in order to be able to use the Operator-SDK. Next, we create a simple Go-based Memcached operator using operator-sdk, and then deploy it onto the OpenShift Container Platform. 

2. [Deep dive into Memcached Operator Code](https://github.ibm.com/TT-ISV-org/operator/blob/main/INTERMEDIATE_TUTORIAL.md): In this article we will build upon the [Memcached Operator tutorial](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md). We will deep-dive into the reconcile function, the KubeBuilder markers, and the low-level APIs that enable an operator to update Kubernetes resources.

## Advanced level

1. [Develop and Deploy a Level 1 Janusgraph Operator on OpenShift Container Platform](https://github.ibm.com/TT-ISV-org/operator/blob/main/articles/level-1-operator.md): 
In this tutorial, we will discuss how to develop and deploy a Level 1 operator on the OpenShift Container Platform. We will use the 
[Operator SDK Capability Levels](https://operatorframework.io/operator-capabilities/) as our guidelines for what is considered a 
level 1 operator. In part 1 of the tutorial we will deploy JanusGraph using the default (BerkeleyDB) 
backend storage. This will be a simple approach, and only recommended for testing purposes. Once we've gotten the default configuration working, we will move on to part 2 which will feature 
Cassandra as the backend storage for JanusGraph.



