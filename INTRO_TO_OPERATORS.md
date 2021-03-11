# Intro to Operators

This article will be the first in a series of articles and tutorials on learning how to build and deploy 
a Kubernetes Operator. 

This article will assume you have no knowledge on Kubernetes Operators, and will 
give you all of the basic knowledge needed to develop a
Golang based operator. If you are already familiar with operators, you can either skim this article or 
simply skip ahead to the [develop and deploy a Memcached Operator on OpenShift](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md) which shows how to develop and deploy your first operator to 
an OpenShift cluster, or go back to one of the other articles in the [learning path](https://github.ibm.com/TT-ISV-org/operator#kubernetes-operators-learning-path).

## Expectations (What you have)
* You have little to no experience developing operators
* You have little to no knowledge on Kubernetes Operators concepts
* You have a basic understanding of Kubernetes concepts

## Expectations (What you want)
* You want to learn the basic concepts and steps needed to develop a Golang based operator to manage Kubernetes resources

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
that are typically done by human operators. That's why they are called operators. Kubernetes is great at managing 
stateless applications, but when you need more complex configuration details for a stateful application, such as a 
database, that is when operators are very useful. Other more complex lifecycle management tasks such as patches and minor
upgrades can be automated using an operator. 

Operators are application-specific controllers which extend the functionality of the Kubernetes API to manage instances of complex applications, on behalf of a Kubernetes admin. The [custom resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)(CR) is the mechanism through which the Kubernetes API is extended. 
A Custom Resource Definition (CRD) lists out all of the configuration available to users of the operator. 

In Kubernetes, controllers of the
[control plane](https://kubernetes.io/docs/reference/glossary/?all=true#term-control-plane) implement control loops that repeatedly compare the desired state of the cluster to its actual state. If the states don't match,
then the controller takes action to fix the problem. Similarly, a Kubernetes operator watches a specific CR type and takes application-specific actions to make the current state match the desired state in that resource.

Operators have the following features:

* The user provides configuration and settings within a CR, and then the operator translates the configuration into low-level actions,
based on the logic defined in the operator's custom controller logic.
* Operator introduce new object types through its custom resource definition. These objects can be handled by the Kubernetes API just like
native Kubernetes objects, including interaction via `kubectl` and inclusion in role-based access control policies.


Read more about operators from this [Red Hat blog](https://www.redhat.com/en/topics/containers/what-is-a-kubernetes-operator).

## 2. What do operators do

![Alt text](./images/operator-reconciliation.png)

In this image, you can see that the when we create a custom controller for our operator, it is treated as a workload (more on this later), and it is added to the 
worker node. Much like many of the Kubernetes native controllers, each of which watch for a specific resource (such as Deployments, or Jobs), the 
custom controllers we will develop in this series will watch for custom resources (i.e. the resources which our operators manages). Instead of an admin using 
`kubectl` commands to change the desired state, we will instead specify our desired state in the custom resource's `Spec` section and the operator will take care 
of ensuring the current state of the cluster reaches the desired state.

### Operators vs. Operands
An operator is the combination of CRs and a custom controller that extends Kubernetes functionality
to enable the starting, scaling, and recovering of a specific application or service. <b>The `operand`, on the 
other hand, is what we call the resources an operator manages, i.e. the workload. </b>
 
## Why does Kubernetes need operators?

Kubernetes needs operators in order to automate tasks which are normally done manually by a 
SRE. Instead of having to set up multiple deployments, configmaps, secrets, and services, as 
an end user, you can just deploy your operator instead. Your operator will take care of everything
needed to make sure your service is up and running. The approach of using an operator is 
inherently easier, and scales better, than creating all of the deployments, configmaps, secrets, and services manually. 

## 3. Operator SDK

Operator SDK is an open source toolkit that provides tools to build, test and package operators. The SDK cli allows you to scaffold a project and also provides commands to generate code. Also operator SDK makes use of `make`, a build automation tool, to build, test, package and deploy your operator through series of `make` commands that is provided in generated `Makefile`. The `Makefile` comes with pre-built commands like below which we will be using in our project.

* `make manifests` generates manifests `yaml` definitions based on `kubebuilder` markers.
* `make install` compiles your code and create executables.
* `make generate` updates the generated code based on your operator API schema.
* `make docker-build` builds the operator docker image.
* `make docker-push` pushes the operator docker image.
* `make deploy` deploys all the resources to the cluster.
* `make undeploy` deletes all the deployed resources from the cluster.

Operator SDK also allows you to install OLM (operator lifecycle manager) using `operator-sdk olm install` command. OLM is a set of cluster resources that manage the lifecycle of an Operator. Once installed, you can get the status of the OLM using `operator-sdk olm status`, to make sure all the resources in the cluster are in `installed` status.

## 4. Operator Capability Levels

![Alt text](./images/operator-capability-level.png)

Operators come in different maturity levels in regards to their lifecycle management capabilities. This model aims to provide 
guidance in terms of what features users can expect from a particular operator. As you can see from the picture above, only
Ansible and Go can be used to achieve all five capability levels. Helm can only be used to achieve seamless upgrades and basic install. <b>Capability levels build on top of one another. That means if you have level 3 capabilities, then you should have all capabilities required from Level 1 and Level 2.</b>

Let's take a look at level one in more detail:

### Operator Capability Level 1 - Basic Install 

In level 1, your operator can provision an application through a custom resource. All of the configuration
details are specified in the CR. You should also be able to install your operator in multiple ways (`kubectl`, Operator Hub, 
or through the Operator Lifecycle Manager). Avoid the practice of making the user create / manage configuration files outside
of Kubernetes.

### Level 1 Example - Installing the Workload

The operator enables the deployment of a database by ensuring that `Deployment`, `ServiceAccount`, `RoleBinding`, `ConfigMap`, `PersistentVolumeClaim`,
and `Secret` resources are created. These resources are specified in the custom resource's `spec` section (more on this soon). Once the custom resource is 
created, the custom resource will install the workload (or operand). If the custom resource is deleted, then the workload (or operand) is 
removed. 

Once the custom resource installs the workload, it then initializes an empty database schema, and alerts the user when the database is ready to accept requests by updating the `status` section of the custom resource.


### Level 1 Example - Managing the Workload 

Now, let's say that you want to increase the capacity of your underlying database. How would you do this using operators?
This should be done by resizing the `PersistentVolumeClaim` resources within the `Spec` section of the Custom Resource. When we update the `Spec` section of the custom resource, we are configuring the workload (or operand). Once 
these changes are applied, the custom resource will take care of scaling the underlying `PersistentVolumeClaim` resource to match 
what was declared in the `Spec` section of the Custom Resource.

To read more about the other capability levels, read this article from the [Operator SDK documentation](https://sdk.operatorframework.io/docs/advanced-topics/operator-capabilities/operator-capabilities/).

## 5. Operator Hub

![Alt text](./images/operatorHub.png)

[OperatorHub.io](https://operatorhub.io/) is where you can find and share Operators. As you can see in the picture above, there
are more than 180 different operators to choose from on OperatorHub.io. OperatorHub.io is very important since this is where 
you can use other operators to automate the configuration of your Kubernetes applications, and submit your own operator to be published online. All of the details of how to package, test, preview, and submit your operator can be found in [this article](https://operatorhub.io/contribute). 

## Conclusion
In this article, we learned about how operators can extend the base Kubernetes functionality 
by the use of custom controllers and custom resources. We've also learned that the Operator SDK offers code scaffolding 
tools to enable you to write your operator faster, and offers guidelines for the capability levels of an operator. Lastly,
we learned that we can view operators and submit our own on OperatorHub.io.

In the [next article](https://github.ibm.com/TT-ISV-org/operator/blob/main/articles/demystified.md), we will dive deeper 
into the Kubernetes architecture that enables operators to work. 

If you would rather go straight to developing an operator, you can go to the intermediate level tutorial [develop and deploy and operator to OpenShift](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md) tutorial instead.