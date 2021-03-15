# Intro to Operators

This article will be the first in a series of articles and tutorials on learning how to build and deploy 
a Kubernetes Operator. 

This article assumes you are familiar with how Kubernetes works but that you have no knowledge on Kubernetes Operators, and will 
give you all of the basic knowledge needed to develop an
operator implemented using Golang. If you are already familiar with operators, you can either skim this article or 
simply skip ahead to the [develop and deploy a Memcached Operator on OpenShift](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md) which shows how to develop and deploy your first operator to 
an OpenShift cluster, or go back to one of the other articles in the [learning path](https://github.ibm.com/TT-ISV-org/operator#kubernetes-operators-learning-path).

## Expectations (What you have)
* You have a basic understanding of Kubernetes concepts and how to install a workload
* You have little to no knowledge on Kubernetes Operators concepts
* You have little to no experience developing operators

## Expectations (What you want)
* You want to learn the basic concepts and steps needed to develop a Golang-based operator to manage Kubernetes resources

## Estimated time
* This article should take roughly 15-30 minutes to complete, depending on how long you spend reading through 
documentation.

## Outline
1. [What are operators](#1-What-are-operators)
1. [What do operators do](#2-what-do-operators-do)
1. [Operator SDK](#3-Operator-SDK)
1. [Operator Capability Levels](#4-operator-capability-levels)
1. [Operator Hub](#5-Operator-Hub)

## 1. What are operators
[Operators](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) are extensions to Kubernetes that make use of custom resources 
to manage Kubernetes applications and their components. Operators are used to automate software configuration and maintenance activities 
that are typically performed by human operators. That's why they are called operators! Kubernetes is great at managing 
stateless applications, but when you need more complex configuration details for a stateful application, such as a 
database, that is when operators are very useful. An operator can also automate other more complex lifecycle management tasks such as version
upgrades, failure recovery, and scaling. 

Operators extend the Kubernetes control plane with specialized functionality to manage a workload on behalf of a Kubernetes admin. An operator includes these components:
- A [custom resource definition](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/) (CRD) that defines a schema of settings available for configuring the workload
- A [custom resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) (CR) that specifies values for the settings defined by the CRD, values that describe the configuration of a workload instance
- A [controller](https://kubernetes.io/docs/concepts/architecture/controller/) customized for the workload that configures the current state of the workload to match the desired state represented by the values in the CR

Operators have the following features:

* The user provides configuration settings within a CR, and then the operator translates the configuration into low-level actions,
using the operator's custom controller logic to implement the translation.
* An operator introduces new object types through its custom resource definition. These objects can be handled by the Kubernetes API just like
native Kubernetes objects, including interaction via Kubernetes client tools and inclusion in role-based access control policies (RBAC).

The article [What is a Kubernetes operator?](https://www.redhat.com/en/topics/containers/what-is-a-kubernetes-operator) by Red Hat explains more details about operators.


## 2. What do operators do

In Kubernetes, controllers in the
[control plane](https://kubernetes.io/docs/concepts/overview/components/) run in a control loop that repeatedly compares the desired state of the cluster to its current state. If the states don't match,
then the controller takes action to adjust the current state to more closely match the desired state. Similarly, the controller in an operator watches a specific CR type and takes application-specific actions to make the workload's current state match the desired state expressed in the CR.

This diagram illustrates how the control plane runs the controllers in a loop, where some controllers are built into Kubernetes and some are part of operators:

![Control Loops with Operators](./images/operator-reconciliation.png)

The controllers in the control plane are optimized for stateless workloads and one set of controllers works for all stateless workloads because they're all very similar. The controller in an operator is customized for one particular stateful workload. Each stateful workload has its own operator with its own controller that knows how to manage this workload.

## Why does Kubernetes need operators?

Kubernetes needs operators in order to automate tasks that are normally performed manually by IT operations personnel. Statefulness changes how a workload needs to be installed, upgrades to a new version, recovers from failures, needs to be monitored, and scales out and back in again. The operator will take care of everything needed to make sure the service is up and running. A workload with an operator manages itself better, making it easier for application teams to use it with less effort from the operations team.

## 3. Operator SDK

The [Operator SDK](https://sdk.operatorframework.io/) is an open source toolkit that provides tools to build, test, and package operators. The SDK CLI enables you to scaffold a project and also provides commands to generate code. The SDK uses [make](https://en.wikipedia.org/wiki/Make_(software)), a build automation tool that runs commands configured in a Makefile to generate executable code and libraries from source code. The SDK includes pre-built make commands that we will use to develop our operator, such as:

* `make manifests` -- generates yaml manifests based on `kubebuilder` markers
* `make install` -- compiles source code into executables
* `make generate` -- updates the generated code based on an operator's API schema
* `make docker-build` -- builds the operator's Docker container image
* `make docker-push` -- pushes the Docker image
* `make deploy` -- deploys all of the operator's resources to the cluster
* `make undeploy` -- deletes all of the operator's deployed resources from the cluster

These commands in the SDK greatly simplify implementing an operator.

### Operator Lifecycle Manager

The SDK also enables you to install the [Operator Lifecycle Manager](https://olm.operatorframework.io/) (OLM) using the `operator-sdk olm install` command. The OLM is a set of cluster resources that manage the lifecycle of an operator. Once installed, you can get the status of the OLM using `operator-sdk olm status`, which verifies whether the SDK can successfully communicate with the OLM components in the cluster.

### Terminology

In addition to operator, custom resource, and custom resource definition, the Operator SDK adds the following [terminology](https://operatorframework.io/operator-capabilities/), _operand_ and _managed resource_:
- Operand - the managed workload provided by the Operator as a service
- Managed resources - the Kubernetes objects or off-cluster services the Operator uses to constitute an Operand (also known as secondary resources)


## 4. Operator Capability Levels

Some operators are more sophisticated at managing their operand's lifecycle than others. The Operator Capability Levels model defines five levels of sophistication, as illustrated here:

![Operator Capability Levels](./images/operator-capability-level.png)

This model aims to provide 
guidance for the features that users can expect from a particular operator. As the picture shows, only
Ansible and Go can be used to achieve all five capability levels. Helm can only be used to achieve the first two levels. <b>Capability levels build on top of one another. That means that if an operator has level 3 capabilities, then it should also have all of the capabilities required from level 1 and level 2.</b>

Before an operator can even achieve level 1, first you must install the operator itself. The operator is a Kubernetes workload and so can be installed the way any workload is installed, such as by using the Kubernetes CLI or by using a Helm chart. It can also be installed by an operator repository such as the Operator Hub or through the Operator Lifecycle Manager.

Let's examine the capabilities of a level 1 operator in more detail.

### Operator Capability Level 1 - Basic Install

In level 1, your operator can provision an application as described by a custom resource. The CR specifies all of the configuration
details. Avoid making the user create and manage configuration files outside
of Kubernetes, that's what the CR is for.

### Level 1 Example - Installing the Operand

When a custom resource is created, that triggers the operator, which responds by creating and installing the operand. If the custom resource is deleted, then the operator removes the operand.

To install an operand, the controller creates the managed resources for that operand and installs them, which causes the cluster to install the operand. These managed resources are typical Kubernetes workload kinds such as a `Deployment` and a `Service`, as well as other kinds like a `ConfigMap`, a `Secret`, and a `PersistentVolumeClaim`.

The configuration of these managed resources is specified in the custom resource's specification (i.e. the `spec` section of its yaml (more on this soon)). A simple controller may simply copy values out of the CR's specification and into the appropriate fields of the managed resources, maybe transforming the values as needed. Once the cluster has installed the operand, the controller gathers the status of those managed resources and updates the custom resource's status (i.e. the `status` section of its yaml (more on this soon)). Status is how the CR remembers the managed resources that were created for its operand.

### Level 1 Example - Managing the Operand 

Now, let's say that you want to increase the capacity of the operand. How would you do this using the operator?

This should be done by updating the custom resource's specification--perhaps it has a `size` setting and you increase it. When we update the custom resource's specification, we are specifying a different configuration for the operand. The controller notices the changes in the custom resource and responds by changing the configuration of the managed resources--perhaps to increase the number of pods or create new persistent volume claims.

Once the controller applies these changes to the managed resources, the cluster responds by applying them the operand, thereby scaling the operand.

[Operator Capability Levels](https://sdk.operatorframework.io/docs/advanced-topics/operator-capabilities/operator-capabilities/) in the Operator SDK documentation also describes the behavior of the other capability levels.

## 5. Operator Hub

[OperatorHub.io](https://operatorhub.io/) is a repository for operators, a public website where you can find and share operators. Its homepage looks like this:

![Operator Hub](./images/operatorHub.png)

There
are more than 180 different operators to choose from on OperatorHub.io. OperatorHub.io is very important since this is where 
you can use other operators to automate the configuration of your Kubernetes applications, and submit your own operator to be published online. All of the details of how to package, test, preview, and submit your operator for addition to the Hub can be found in [How to contribute an Operator](https://operatorhub.io/contribute). 

## Conclusion
In this article, we learned about how operators can extend the base Kubernetes functionality 
using custom controllers and custom resources. We've also learned that the Operator SDK offers code scaffolding 
tools to enable you to write your operator more easily, and offers guidelines for the capability levels of an operator. Lastly,
we learned that we can browse existing operators and submit our own on OperatorHub.io.

In the [next article](https://github.ibm.com/TT-ISV-org/operator/blob/main/articles/demystified.md), we will dive deeper 
into the Kubernetes architecture that enables operators to work. 

If you would rather go straight to developing an operator, go to the intermediate level tutorial [Develop and deploy and operator to OpenShift](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md).
