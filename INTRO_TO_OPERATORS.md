# Intro to Kubernetes Operators

This article is the first in a series of articles and tutorials that show you how to build and deploy a Kubernetes Operator. In this article, learn what Kubernetes Operators are, what they do, and key features including the Operator SDK, Operator capability levels and Operator Hub. You will learn all of the basic knowledge needed to develop an operator implemented using Golang. 

If you are already familiar with operators, you can either skim this article or skip ahead to the next Operator tutorial: [Develop and deploy a Memcached Operator on OpenShift](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md) which shows how to develop and deploy your first operator to 
an OpenShift cluster. Alternately, read one of the other articles in the [learning path](https://github.ibm.com/TT-ISV-org/operator#kubernetes-operators-learning-path).

## Prerequisites

This article assumes you want to learn the basic concepts about Kubernetes Operators and steps needed to develop a Golang-based operator to manage Kubernetes resources.

* You have a basic understanding of Kubernetes concepts and how to install a workload
* You have little to no knowledge about Kubernetes Operators concepts
* You have little to no experience developing operators

## Estimated time

This article should take roughly 15-30 minutes to complete.

## What are operators?

[Operators](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) are extensions to Kubernetes that use custom resources to manage Kubernetes applications and their components. Operators are used to automate software configuration and maintenance activities that are typically performed by human operators. That's why they are called operators! 

While Kubernetes is great at managing stateless applications, when you need more complex configuration details for a stateful application, such as a database, operators are very useful. An operator can also automate other more complex lifecycle management tasks such as version upgrades, failure recovery, and scaling. 

Operators extend the Kubernetes control plane with specialized functionality to manage a workload on behalf of a Kubernetes admin. An operator includes these components:

- A [custom resource definition](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/) (CRD) that defines a schema of settings available for configuring the workload
- A [custom resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) (CR) is the Kubernetes API extension created by a CRD
    - A custom resource instance sets specific values for the settings defined by the CRD to describe the configuration of a workload
- A [controller](https://kubernetes.io/docs/concepts/architecture/controller/) customized for the workload that configures the current state of the workload to match the desired state represented by the values in the CR

Operators have the following features:

* The user provides configuration settings within a CR, and then the operator translates the configuration into low-level actions, using the operator's custom controller logic to implement the translation
* An operator introduces new object types through its custom resource definition. These objects can be handled by the Kubernetes API just like native Kubernetes objects, including interaction via Kubernetes client tools and inclusion in role-based access control policies (RBAC)

The article [What is a Kubernetes operator?](https://www.redhat.com/en/topics/containers/what-is-a-kubernetes-operator) by Red Hat explains more details about operators.

## What do operators do?

In Kubernetes, controllers in the [control plane](https://kubernetes.io/docs/concepts/overview/components/) run in a control loop that repeatedly compares the desired state of the cluster to its current state. If the states don't match,
then the controller takes action to adjust the current state to more closely match the desired state. 

Similarly, the controller in an operator watches a specific CR type and takes application-specific actions to make the workload's current state match the desired state expressed in the CR.

This diagram illustrates how the control plane runs the controllers in a loop, where some controllers are built into Kubernetes and some are part of operators:

![Control Loops with Operators](./images/operator-reconciliation.png)

The controllers in the control plane are optimized for stateless workloads and one set of controllers works for all stateless workloads because they're all very similar. The controller in an operator is customized for one particular stateful workload. Each stateful workload has its own operator with its own controller that knows how to manage this workload.

### Why does Kubernetes need operators?

Kubernetes needs operators to automate tasks that are normally performed manually by IT operations personnel. A stateful workload is more difficult to manage than a stateless workload. The state in a workload changes how a workload: needs to be installed, upgrades to a new version, recovers from failures, needs to be monitored, and scales out and back in again. The operator will take care of everything needed to make sure the service keeps running successfully. An operator makes its service more self-managing, enabling an application team to spend less effort managing the service so that it can spend more effort using the service. 

### How do I use an operator?

To use an operator, you must first deploy the operator as a workload in your cluster. You can deploy it like 
any other workload, via the Kubernetes CLI or using a Helm chart. Additionally, you can make use of Operator Hub 
or deploy your operator through the Operator Lifecycle Manager.
 
## Operator SDK

The [Operator SDK](https://sdk.operatorframework.io/) is a set of open source tools to build, test, and package operators. The SDK CLI enables you to scaffold a project and also provides commands to generate code. The SDK uses [make](https://en.wikipedia.org/wiki/Make_(software)), a build automation tool that runs commands configured in a Makefile to generate executable code and libraries from source code. 

The SDK includes pre-built make commands that we will use to develop our operator, such as:

* `make manifests` -- generates YAML manifests based on `kubebuilder` markers
* `make install` -- compiles source code into executables
* `make generate` -- updates the generated code based on an operator's API schema
* `make docker-build` -- builds the operator's Docker container image
* `make docker-push` -- pushes the Docker image
* `make deploy` -- deploys all of the operator's resources to the cluster
* `make undeploy` -- deletes all of the operator's deployed resources from the cluster

These commands in the SDK greatly simplify implementing an operator.

### Operator Lifecycle Manager

The SDK also enables you to install the [Operator Lifecycle Manager](https://olm.operatorframework.io/) (OLM) using the `operator-sdk olm install` command. The OLM is a set of cluster resources that manage the lifecycle of an operator. Once installed, you can get the status of the OLM using `operator-sdk olm status`, which verifies whether the SDK can successfully communicate with the OLM components in the cluster.

### Additional Operator SDK concepts

In addition to operator, custom resource, and custom resource definition, you need to understand what we mean by the Operator SDK _operand_ and _managed resource_:

- Operand - the managed workload provided by the Operator as a service
- Managed resources - the Kubernetes objects or off-cluster services the Operator uses to constitute an Operand (also known as secondary resources)

## Operator capability Levels

Some operators are more sophisticated at managing their operand's lifecycle than others. The [Operator Capability Levels model](https://operatorframework.io/operator-capabilities/) defines five levels of sophistication, as illustrated here:

![Operator Capability Levels](./images/operator-capability-level.png)

This model aims to provide guidance for the features that users can expect from a particular operator. As the picture shows, you can only use Ansible and Go to achieve all five capability levels. Helm can only be used to achieve the first two levels. 

<b>Capability levels build on top of one another. That means that if an operator has level 3 capabilities, then it should also have all of the capabilities required from level 1 and level 2.</b>

Let's examine the capabilities of a level 1 operator in more detail.

### Operator Capability Level 1 - Basic install

In level 1, your Operator can provision an application as described by a custom resource. The CR specifies all of the configuration details. Avoid making the user create and manage configuration files outside of Kubernetes &mdash; that's what the CR is for.

### Level 1 Example - Install the Operand

When a custom resource is created, that triggers the Operator, which responds by creating and installing the Operand. If the custom resource is deleted, then the Operator removes the Operand.

To install an Operand, the controller creates the managed resources for that Operand and installs them, which causes the cluster to install the Operand. These managed resources are typical Kubernetes workload kinds such as a `Deployment` and a `Service`, as well as other kinds like a `ConfigMap`, a `Secret`, and a `PersistentVolumeClaim`.

You can specify the configuration of these managed resources in the custom resource's specification (that is, the `spec` section of its YAML); I cover this in greater detail in the [next article](https://github.ibm.com/TT-ISV-org/operator/blob/main/articles/demystified.md#2-kubernetes-architecture). 

A simple controller may simply copy values out of the CR's specification and into the appropriate fields of the managed resources, maybe transforming the values as needed. Once the cluster installs the Operand, the controller gathers the status of those managed resources and updates the custom resource's status (that is, the `status` section of its YAML); I cover this in greater detail in my [next article](https://github.ibm.com/TT-ISV-org/operator/blob/main/articles/demystified.md#2-kubernetes-architecture). Status is how the CR remembers the managed resources that were created for its Operand.

### Level 1 Example - Manage the Operand 

Now, let's say that you want to increase the capacity of the Operand. How would you do this using the Operator?

You can do this by updating the custom resource's specification &mdash; perhaps it has a `size` setting, and you increase it. When you update the custom resource's specification, you are specifying a different configuration for the Operand. The controller notices the changes in the custom resource and responds by changing the configuration of the managed resources &mdash; perhaps to increase the number of pods or create new persistent volume claims.

Once the controller applies these changes to the managed resources, the cluster responds by applying them the Operand, thereby scaling the Operand.

We will only cover the first level of the capability levels in this article, but you can read more about the other four levels at 
[Operator Capability Levels](https://sdk.operatorframework.io/docs/advanced-topics/operator-capabilities/operator-capabilities/) in the Operator SDK documentation.

## Operator Hub

[OperatorHub.io](https://operatorhub.io/) is a web-based application that enables finding, installing, and publishing operators. It is a one stop shop for Kubernetes operators, as shown below:

![Operator Hub](./images/operatorHub.png)

Check out [How to contribute an Operator](https://operatorhub.io/contribute) for details of how to package, test, preview, and submit your operator for addition to the Hub. 

## How to implement an operator

Operators are usually implemented in three main technologies: Go, Helm, and Ansible. You can read more about the pros and cons of each approach from
this [Red Hat article](https://www.openshift.com/blog/build-your-kubernetes-operator-with-the-right-tool). For the rest of the learning path, we will use Go to 
implement a Kubernetes operator, as 71% of operators from OperatorHub.io are written in Go. Also, Kubernetes itself is written in Go, so this makes writing the 
operator which interacts with the Kubernetes API much easier. 

## Conclusion

In this article, you learned how operators can extend the base Kubernetes functionality using custom controllers and custom resources. You also learned that the Operator SDK offers code scaffolding tools to enable you to write your Operator more easily, and offers guidelines for the capability levels of an operator. Lastly, I introduced you to OperatorHub.io where you can browse existing operators and submit your own.

The next article, [Anatomy of an Operator, demystified](https://github.ibm.com/TT-ISV-org/operator/blob/main/articles/demystified.md), takes a closer look at the Kubernetes architecture that enables Operators to work. 

If you would rather go straight to developing an Operator, go to the [Develop and deploy a Memcached Operator](https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md) tutorial.
