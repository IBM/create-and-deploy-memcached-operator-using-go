# Kubernetes Operators Learning Path

Welcome to IBM Developer's Kubernetes Operators learning path! In this series of articles and tutorials, you will learn how to create and deploy a Golang based operator. You will also 
learn about all of the foundational Kubernetes knowledge needed to understand how to develop
an operator from scratch. By the end of this learning path you should be able to develop
your own operator, deploy it to your cluster, and develop an operator with all five of 
the [operator capibility levels](https://sdk.operatorframework.io/docs/advanced-topics/operator-capabilities/operator-capabilities/).

The learning path has several levels of tutorials, starting from beginner, to intermediate, 
to advanced.

## Beginner level

1. [Environment Setup](https://github.ibm.com/TT-ISV-org/operator/blob/main/installation.md): 
This article will guide you through all of the prerequisites and environment setup needed to develop an Operator using the Operator-SDK.

2. [Intro to Operators](https://github.ibm.com/TT-ISV-org/operator/blob/main/INTRO_TO_OPERATORS.md): This article does a summary into Kubernetes concepts such as workloads, architecture, controllers, and custom resources. It explains the control loop and the declaritive API that is 
at the heart of Kubernetes, and how the operator pattern works.

3. [Develop and Deploy a Memcached Operator on OpenShift Container Platform](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md): 
In this tutorial we will be creating a simple Go-based Memcached operator using operator-sdk, and then deploying it onto the OpenShift Container Platform. 

## Intermediate level
1. [Deep dive into Memcached Operator Code](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md): In this tutorial we will build upon the [Memcached Operator tutorial](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md) and deep-dive into the code to understand what the operator is doing in the custom controller code, and why it is doing it.
