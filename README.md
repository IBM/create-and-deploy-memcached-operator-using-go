# Kubernetes Operators Learning Path

Welcome to IBM Developer's Kubernetes Operators learning path! In this series of articles and tutorials, you will learn how to create and deploy a Golang based operator. You will also 
learn about all of the foundational Kubernetes knowledge needed to understand how to develop
an operator from scratch. 

By the end of this learning path you should be able to: 

* Understand what operators are, why they are useful, and the five [operator capability levels](https://sdk.operatorframework.io/docs/advanced-topics/operator-capabilities/operator-capabilities/). 
* Setup your local machine to be able to develop an operator using the Operator-SDK.
* Develop and deploy a Memcached Operator on OpenShift.
* Understand Kubernetes concepts such as workloads, controllers, custom resources, and the control plane. 
* Understand how Operators extend Kubernetes. 
* Understand the reconcile function within an operator's custom controller.
* Understand the low-level APIs that enable an operator to update Kubernetes resources.  

## Beginner level
1. [Intro to Operators](https://github.ibm.com/TT-ISV-org/operator/blob/main/INTRO_TO_OPERATORS.md): This article does a summary into Kubernetes concepts such as workloads, architecture, controllers, and custom resources. It explains the control loop and the declaritive API that is 
at the heart of Kubernetes, and how the operator pattern works.

2. [Anatomy of an operator, demystified](https://github.ibm.com/TT-ISV-org/operator/blob/main/articles/demystified.md): In this article we will build upon the [Intro to Operators](https://github.ibm.com/TT-ISV-org/operator/blob/main/INTRO_TO_OPERATORS.md) article and explore
the Kubernetes architecture which enable the custom functionality of operators.

## Intermediate level

1. [Develop and Deploy a Memcached Operator on OpenShift Container Platform](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md): 
In this tutorial we will start by ensuring we have our [environment setup](https://github.ibm.com/TT-ISV-org/operator/blob/main/installation.md) in order to be able to use the Operator-SDK. Next, we create a simple Go-based Memcached operator using operator-sdk, and then deploy it onto the OpenShift Container Platform. 

2. [Deep dive into Memcached Operator Code](https://github.ibm.com/TT-ISV-org/operator/blob/main/INTERMEDIATE_TUTORIAL.md): In this article we will build upon the [Memcached Operator tutorial](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md) and deep-dive into the code to understand what the operator is doing in the custom controller code, and why it is doing it.

