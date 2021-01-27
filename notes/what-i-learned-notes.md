# Learnings and notes from various sources

Credit: Kubernetes Operators by Jason Doies and Joshua Wood (O'Reilly)
Credit: kublr.com/blog/understanding-kubernetes-operators
Credit: kubernetes.io/docs/concepts/extend-kubernetes/operator

## Vocabulary and concepts necessary in order to understand how Operators work with K8s 

----------------------------------------------
### What is Kubernetes (i.e. K8s)
----------------------------------------------
* A k8s cluster is a collection of computers, called nodes 
  * All cluster work runs on one or more of the nodes
* The basic unit of work, and replication, is called a pod
  * A pod is one or more linux containers with common resources like networking, storage, and access to shared memory 
* A k8s cluster can be divided into two planes but sometimes these planes have some overlap. The kubelet agent which is running on every node is part of the control plane for example. 
  * The control plane, which is k8s itself and implements the K8s API and cluster orchestration logic
  * The app plane, or data plane, is everything else. This is the group of pods where the app pods run
* One or more nodes are usually dedicated to running applications, while one or more nodes are reserved only for the control plane 
* Multiple replicas of control plane components can run on multiple pods to provide redundancy

----------------------------------------------
### How K8s works
----------------------------------------------
* K8s automates the lifecycle of a stateless app, such as a static web server 
* Without state, any of the app instances are interchangeable 
  * Because the server is not tracking the state, or storing input or data of any kind, when once instance fails, k8s can replace it with another. 
* These instances are called replicas which are just copies of an app running on a cluster 
* The controllers of the control loop implement logic to automatically check for the difference between the actual state and the desired state 
  * When the two diverge, the controller takes action and makes sure the two match
* In small clusters, the control plane sometimes will share the same pod as the application pod

----------------------------------------------
### Stateful vs. Stateless Apps
----------------------------------------------
* In stateless deployments, the order of deploying pods, their labels, network address or port, storage class, or volume are not important. You keep them if they are healthy and serving, you dispose of them when they become unhealthy, outdated, or just no longer needed and replace them as necessary.  
* In stateful apps, some order is necessary. You also need to add storage and persistent volume so that the state is saved, and the cluster admin has to manage that.

----------------------------------------------
### Stateful is Hard
----------------------------------------------
* Most apps have state, and they have configs, dependence on other packages, etc
* They need to store critical and sometimes a lot of data
* It would be great to manage these apps with a uniform way, and also automate their storage, networking, and cluster connections
* K8s cannot know all about every cluster app while also remaining general, simple, and adaptable. k8s does not and should not know the config to your cloudant database 
* But it does provide an extension mechanism - which is the Custom resource and custom controller mechanism we are writing with the operator SDK

----------------------------------------------
### Standard Scaling: the ReplicaSet Resource
----------------------------------------------
* The ReplicaSet gives a sense of how resources comprise the app management database at the heart of k8s 
* Like any other resource in the k8s API, the ReplicaSet is a collection of API objects. The ReplicaSet primarily collects pod objects forming a list of the running replicas of an app
  * The ReplicaSet defines the number of replicas that should be maintained on the cluster
  * Another object spec points to a template for creating new pods when there are fewer than desired
* The ReplicaSet controller manages ReplicaSets and the pods belonging to them. The replicaSet controller creates ReplicaSets and continually monitors them. If the count of running pods
doesn't match the desired number in the `Replicas` field, the ReplicaSet controller 
starts or stops pods to get to the desired number.
* K8s replicaSets are meant to be app agnostic, so they cannot know all the intricacies of disaster recovery of a specific service.

----------------------------------------------
### Custom Resources(CRs) - Custom API endpoints
----------------------------------------------
* CRs hold structured data, and the k8s API server provides a mechanism for reading and setting their files as you would those in a native resource, by using kubectl or another API client. 
* CRs are most useful when they are watched by custom controller code that creates, updates, deletes other cluster objects or even arbitrary resources outside of the cluster
* CRs obey k8s conventions, like the resources .spec and .status 
* A CRD is a schema that defines the fields and the types of the fields within an instance of a CR. They are not part of a default k8s deployment. 
  * A YAML manifest describes a CRD. A CR is a named endpoint in the k8s API. Users define a CR by providing a custom-resource definition or CRD.
    * The CRD defines how the  should reference this new resource. 
  * The CR’s group, version, and kind together form the fully qualified name of a k8s resource type. That canonical name must be unique across a cluster. 


----------------------------------------------
### Custom Controllers
----------------------------------------------
* To provide a declarative API for a specific app running on a cluster, you need active code that captures the processes of managing that app such as maintaining desired state, as 
prescribed by the custom resource.
* Every operator has customer controller implementing app specific management logic.
* Let’s you extend the clusters behavior without having to modify the k8s code itself


----------------------------------------------
### Operator Scopes
----------------------------------------------
* A k8s cluster is divided into namespaces, which is a boundary for cluster object and resource names.
  * Names must be unique within a single namespace, but not between namespaces, which makes it easier for multiple users or teams to share a single cluster
  * Resource limits and access controls can be applied per namespace 
* An operator can be limited to a namespace or it can maintain its operand across an entire cluster 
  * Usually, restricting your operator to a single namespace makes sense and is more flexible for a cluster that is used by multiple teams.
  * An operator scoped to a name space can be upgraded independently of other instances, and this allows for some handy facilities 
  * You can test upgrades in a testing namespace, or serve older app instances from a different namespace 

----------------------------------------------
### Cluster-Scoped Operators
----------------------------------------------
* There may be a time when you want your operator to watch and manage an app or services throughout a cluster
  * For example an operator that manages a service mesh, such as istio, or one that issues TLS certs for app endpoints, like cert-manager, might be most effective when watching cluster-wide
  * Operators usually run on namespace scope, but you can change them to run on cluster-wide scope by using cluster role and cluster role binding instead of role and role binding 


----------------------------------------------
### Service Accounts
----------------------------------------------
* In k8s, regular human user accounts aren’t managed by the cluster, and there are no API resources depicting them
* Service accounts, on the other hand, are managed by k8s and can be created and manipulated thru the k8s api
* A service account is a special type of cluster user for authorizing programs instead of people. 
  * Most operators should service their access rights from a service account
  * Creating a service account is a standard step in deploying an operator
    * The service account identifies the operator, and the accounts role denotes the powers granted to the operator
* As a best practice, <b>you should give the service account a role</b>.

----------------------------------------------
### Authorization and Roles
----------------------------------------------
* The power to do things on the cluster via the API is defined in k8s by one of a few available access control systems.
  * The best one of theses is RBAC (role based access control). RBAC regulates access to system resources according to the role a system user is assigned via role-binding.
* A role is a k8s resource which defines a set of capabilities to take actions on particular API resource such as create, update, read, update or delete 
  * The capabilities described by a role are granted or bound to a user by a role binding

----------------------------------------------
### Role Binding
----------------------------------------------
* Another necessary part of RBAC is the role binding, which assigns the role to the service account for the operator.
  * Kind "RoleBinding" is another k8s native resource, and in there it will bind a specific role to a specific subject, such as a service account.
* By default, if you’re using OpenShift, your oc commands will run in the namespace `my project`. Wherever you are, the namespace value in this role binding must match the namespace on the cluster where you are working on.

----------------------------------------------
### SRE
----------------------------------------------
* A set of practices for running large systems
* A key tenet is automating systems and adding by writing software to run your software 
* Teams freed from maintenance work have more time to create new features, fix bugs, and improve their products 
* An operator is like an automated SRE for its app - It encodes its skills of an expert admin
* An operator can manage a cluster of database servers for example, and it knows the details of configuring and managing its app, and it can install a database cluster of a declared software version and number of members 
* An operator continues to monitor its app as it runs over time, automatically. 
* Operators extend kubernetes, so they can be managed via kubectl 

## What are Operators?

----------------------------------------------
### Operators
----------------------------------------------
* A software that wraps manual tasks into k8s functions for you. It aims to mirror an engineer like an SRE who manages a set of services
  * Wraps any necessary logic for deploying and operating a k8s app using k8s constructs
* Job of the operator is to provision and manage the resources that you would normally need to use manually and since it is provided with enough authorization in the cluster, it can do cluster-management for you, such as rescheduling pods as they fail, or scaling the replica sets as needed. 
  * It can help you in the selection of cloud resources from you cloud env
  * Can automatically provision storage, volume, and any other infrastructure you may need
* Operators are clients of the Kubernetes API that act as controllers for a custom resource
  * Operators are the app specific combo of a CR and a custom controller that does know all the details about starting, scaling, recovering, etc
  * The operators operand is what we call the app, service, or whatever resource an operator manage 
* CRDs are one of two essential building blocks for the most basic description of the operator pattern: a custom controller managing CRs 

----------------------------------------------
### Operator SDK
----------------------------------------------
* This is used for developing the operator
* Once you deploy it, then the operator itself will take care of the resources and things like that
* Making an operator means creating a CRD and providing a program that runs in a loop watching CRs of that kind. The Operator-SDK provides the scaffolding code to speed up development.
* What the operator does in response to changes in the CR is specific to the app that the operator manages, and is part of the custom logic in the controller code. This is all 
scaffolded for you via the Operator-SDK.

----------------------------------------------
## Operator details
----------------------------------------------
* After you've deployed your Memcached operator, you'll have a custom resource named Memcahed that you can configure in your cluster. You'll also have: 
  * You'll have a deployment that makes sure a pod is running that contains the controller part of the operator
  * A container image of the operator code
  * Controller code that queries the control pane to figure out what resources are are configured 
  * Operator code to tell the API server how to make reality match the configured resources
* For example if you have a demoDB operator deployed, and add a new demoDB, the operator helps set up PVCs to provide durable database storage. At also sets up: 
  * A statefulset to run demoDB and a job to handle initial configuration
  * A loop to see if the app is running an old database version, and if so, upgrade it for you.  
* If you delete it one of your demoDB pods, the operator takes a snapshot, and makes sure that statefulsets and volumes are also removed 
* The operator also manages regular database backups. For each demoDB resource, the operator determines when to create a pod that can connect to the database and take backups. These pods would rely on configmap and secrets to store database connection details

----------------------------------------------
### What are operators for
----------------------------------------------
* The operator pattern arose in response to SRE's waiting to extend k8s to provide features specific to their sites and software 
* Operators make it easier for cluster admins to enable, and developers to use foundation software pieces like databases and storage system with less management overhead. 
* If cloudant database server that’s perfect for your app’s backend has an operator to manage it, you can deploy cloudant without needing to become an expert in clouding DBA 
* App developers build operators to manage the apps they are delivering simplifying the dev and management experience on their customer’ k8s clusters. Infra engineers create operators to control deployed services and systems.
* The action an operator performs can include almost anything: scaling a complex app, app version upgrades, or even managing kernel modules for nodes in a computational cluster with specialized hardware 

## Operator example in practice

----------------------------------------------
### Example of Operator tasks
----------------------------------------------
* Deploying an application on demand
* Taking and restoring backups of that app’s state
* Handling upgrades of app code, or database schema changes
* Choosing a leader in a distributed app without an internal member election process


----------------------------------------------
### Example
----------------------------------------------
* We just have a simple web server with no state
  * Kubectl get pods (we have 1 pod)
* Now we declare the we should have 3 replicas, and the clusters state differs from the desired state
  * Kubectl get pods (we should have 3 now)
* Deleting one of the pods should make the control plane go ahead and create another pod
  * The web server is interchangeable with any other replica, and with any new pod that replaces one of the replicas. It doesn’t store data or maintain state in any way
  * K8s doesn’t need to make any special arrangements to replace a failed pod or to scale the application by adding or removing replicas of the server 

----------------------------------------------
### Deploying operators
----------------------------------------------
* Add the custom resource definition and its associated controller to your cluster. The controller would normally run outside of your control plane, much as you would run any containerized app
* For example, you would just run your controller in your cluster as a deployment 

----------------------------------------------
### Using an operator
----------------------------------------------
* They work by extending the Kubernetes control plane and API. 
In its simplest form, an operator adds an endpoint to the k8s api Called a custom resource (CR) along with a control plane component that monitors and maintains resource of the new type. 
* Once you’ve deployed an operator, you’d use it by creating or modifying the resource that the operator manages
  * Kubectl get sample memcached 
  * Kubectl patch force 5 size 
* That’s it, the operator will take care of keeping the service up and running

----------------------------------------------
### How do operators secure clusters
----------------------------------------------
* Not upgrading software is a common source for security vulnerabilities. Operators can take care of this for you
* They can also help with disaster recovery and periodically back up application state
* K8s app is not only designed to be deployed on k8s, but it is built to be used and operated with k8s tools and libraries

----------------------------------------------
### A common starting point
----------------------------------------------
* Etc is a distributed key-value store with roots at CoreOS.
  * Its the underlying data store at the core of k8s, and a key piece of server distributed apps
  * It provides reliable storage by implementing a protocol called Raft that guarantees consensus among a quorum of members 
* The etcd operator often serves as a kind of hello world example of the value and mechanics of the operator pattern, and we follow that tradition here
  * Using etc is super easy (just write and read key value pairs)
  * But administering an etcd cluster of 3 or more nodes requires config of endpoints auth and other concerns usually left to an etc expert (or their collection of custom shell scripts)
  * Keeping etcd running and upgraded over time requires continued admin. The etcd operator knows how to do all of this.

----------------------------------------------
### Example the etcd operator 
----------------------------------------------
* Etcd is a distributed key-vale store i.e. lightweight database
* It usually requires an admin to manage it
  * The admin must know how to join a new node to an etcd cluster, including config with its endpoints, making connections to persistent storage, and making existing members aware of it
  * Back up the cluster data and config 
  * And upgrade the etc cluster to new etc versions 
* The etc operator knows how to perform those tasks 

----------------------------------------------
### Etcd operator
----------------------------------------------
* Manages and provisions the etcd clusters on k8s
* The operator manages pod creation and deletion, failover, restoration, and much more
* Operator can optimize storage and dynamically provision volumes on its own -> this is the PVC operator by bonzai cloud


----------------------------------------------
### The case of the missing member
----------------------------------------------
* Since the etc operator understands etcd’s state, it can recover from an etc cluster members failure in the same way k8s replaced the deleted stateless web server from earlier 
  * Assume there is a 3-member cluster managed by the etc operator
  * The operator itself and the etc cluster members run as pods
  * Deleting an etc pod triggers a <b>reconciliation</b>, and the etcd operator knows how to recover to the desired state of three replicas - something kubernetes can’t do alone.
  * But unlike with the blank-state restart of a stateless web server, the operator has to arrange the new etc pod’s cluster membership, configuring it for the existing end-points and establishing it with the remaining etc members
* Now the operator will repair the etcd cluster. And the etcd API remains available due to the operator 

----------------------------------------------
## Deploying the etcd Operator
----------------------------------------------
* The operator is running in a pod, and it watches the EtcdCluster CR you defined earlier. The manifest file etcd-operator-deployment.yaml lays out the Operator pod’s specs including the container image for the operator you’re deploying and the service account which will use the operator
  * Notice that it does not define the spec for the etcd cluster
  * You’ll describe the desired etcd cluster to the deployed etcd operator in a CR once the operator is running 

----------------------------------------------
## Declaring an etcd cluster
----------------------------------------------
* Earlier, we created a CRD defining a new kind of resource, an EtcdCluster
* Now that we have an operator watching etcdCluster resources, we can declare an etcdCluster with our desired state
  * To do so, provide the two spec elements the operator recognizes
    * Size, the number of etcd cluster members
    * And version of etcd each of those members should run
    * The yaml will tell k8s that we want 
* Since EtcdCluster is now an API resource, you can get the etcd cluster spec and status directly from k8s 
* Try kubectl describe to report on the size, etcd version, and status of your etcd cluster as Kubectl describe etcdcluster/example-etcd-cluster

----------------------------------------------
## Exercising etcd
----------------------------------------------
* We are now running an etcd cluster. 
  * The etcd operator creates a k8s service in the etcd cluster’s namespace
    * A service is an endpoint where clients can obtain access to a group of pods, even though the members of the group may change
    * A service by default has a DNS name visizbe in the cluster 
    * The operator constructs the name of the service used by clients of the etcd API by appending -client to the etcd cluster name defined in the CR.
* We can also use kubectl run —rm -I —image — /bin/sh to connect to the client service and interact with the etcd API
  * From there we can use etcd specific commands to update the key-value pairs in our node

----------------------------------------------
## Scaling the etcd Cluster
----------------------------------------------
* You can grow the cluster by changing the declared size specification
* Edit cluster-cr.yaml and change the size from 3 to 4 and then do kubectl apply -f 
Checking the pods should now show 4 

----------------------------------------------
## Failure and Automated Recovery
----------------------------------------------
* Unlike a stateless program, no etcd pods run in a vacuum
  * Usually a human etcd “operator” has to notice a members failure, execute a new copy and provide it with config so it can join the etcd cluster with the remaining members. 
  * The operator understand etcd’s internal state and makes the recovery automatic 
  * The etcd operator recovers from failures in its complex, stateful app the same way k8s automates recoveries for stateless apps

----------------------------------------------
## How to upgrade etcd the hard way (without operator) 
----------------------------------------------
* Check version and health of each etcd node
* Create a snapshot of cluster state for disaster recovery
* Stop one etcd server and replace the existing version with the v3.2.13 binary, start the new version
* Repeat for each etcd cluster member - at least two more times in a three-member cluster

----------------------------------------------
## The easy way - Let the operator do it
----------------------------------------------
* With a sense of the repetitive and error-prone process of a manual upgrade, it’s easier to see the power of encoding that etcd-specific knowledge into an operator
* The operator can manage the etcd version and upgrade becomes a matter of declaring anew desired version in an EtcdCluster resource 