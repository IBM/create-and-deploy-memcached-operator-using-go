Certifying a RedHat Openshift Operator image

 At a high level the Operator certification steps include:
 1. Confirming that the application’s container images and Operator image is Red Hat certified. Certified images must be:
    1. Built using Red Hat Enterprise Linux or Red Hat Universal Base Image as the base layer.
    2. Comply with Red Hat Container Certification requirements.
    3. Submit to Red Hat to verify and publish.
    4. Provide product information for publishing in the Red Hat Container Catalog.
2. Then the Operator and Application container images are packaged and submitted through the Red Hat Partner Connect website along with additional metadata. The Operator metadata allow it to be deployed using the Operator Lifecycle Manager and to render the Operator’s detail page in OperatorHub.
3. The Operator and Application images then undergo a security health check.
4. Once they complete the health checks the RH team will post the Operator in the embedded OperatorHub.

Check out the detailed guide to OpenShift Operator Certification
- [Partner Guide for OpenShift Operator and Container Certification](https://redhat-connect.gitbook.io/partner-guide-for-red-hat-openshift-and-container/)

In general, to have a Certified Operator, you must first complete the Container Application Certification for all the applications you will be deploying with your operator. 

    ![Certification workflow](https://github.ibm.com/TT-ISV-org/operator/blob/operator-certification/images/Certification-workflow.png)

The certification of an operator as shown in the above image is done 2 stages as follows:

1. Operator image certification process:

    . Create a Container Image project in https://connect.redhat.com/projects

    . Update operator files for certification
    
    . Build and test the operator image
    
    . Upload the operator image to scan.connect.redhat.com

2. Operator bundle image certification process:

    . Create an Operator Bundle Image project in https://connect.redhat.com/projects

    . Bundle the operator

    . Update operator bundle files for certification
    
    . Build and test the bundle image

    . Upload the operator bundle image to scan.connect.redhat.com

How to create the Operator images:
There are 3 kind of Operators available according to the Operator maturity model as follows:

1. Building a Helm Operator: 
https://redhat-connect.gitbook.io/certified-operator-guide/helm-operators/building-a-helm-operator

2. Building an Ansible Operator: 
https://redhat-connect.gitbook.io/certified-operator-guide/ansible-operators/building-an-ansible-operator

3. Building a Golang Operator:
https://github.ibm.com/TT-ISV-org/operator/blob/main/BEGINNER_TUTORIAL.md#develop-and-deploy-a-memcached-operator-on-openshift-container-platform

Red Hat awards OpenShift Certification Badges to Kubernetes Operators built and tested for specific-use, cloud-native cases—like networking and storage—and comply with industry-standard specifications or domain best practices. Current OpenShift Certification Badges are:

• Container Storage Interface (CSI)—for providing and supporting a block/file persistent storage backend for OpenShift.

• Container Networking Interface (CNI)—for delivery of networking services through a pluggable framework.

• Cloud-Native Network Functions (CNF)—for implementation of Telco functions deployed as containers.

Other references:
- [Red Hat OpenShift Operator Certification](https://www.openshift.com/blog/red-hat-openshift-operator-certification) (blog)
- [Red Hat OpenShift Operator Certification for Red Hat Partners](https://connect.redhat.com/en/partner-with-us/red-hat-openshift-operator-certification)
- [Red Hat OpenShift Operator Certification](https://redhat-connect.gitbook.io/red-hat-partner-connect-general-guide/certification-offerings/red-hat-openshift-operator-certification-1) (Git Book)
- [Red Hat OpenShift Operator Certification Data Sheet](https://connect.redhat.com/sites/default/files/2020-07/RH-OpenShift-Operator-Cert-Datasheet-US_062020F4.pdf)
- [Red Hat Ecosystem Catalog: Certified OpenShift Operators](https://catalog.redhat.com/software/operators/explore)
- [OpenShift Certification Badges](https://www.openshift.com/blog/badge-announcement-blog) (blog)
- [Certify your Operator Image](https://redhat-connect.gitbook.io/partner-guide-for-red-hat-openshift-and-container/certify-your-operator/creating-an-operator-project)
- [Certified oprator build guide](https://redhat-connect.gitbook.io/certified-operator-guide/)