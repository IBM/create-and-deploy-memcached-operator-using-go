# Notes from Remko's Operator materials 

Remko's materials can be found [here](https://ibm.github.io/kubernetes-operators/lab1/).

A few interesting ideas in his materials:

## Environment with pre-requisites installed for developer workshop
He has a [web-based command line](https://github.com/IBMAppModernization/web-terminal) for developer workshops which can quickly setup terminal sessions for participants with all necessary software such as Docker, OC and kubectl cli, go, operator-sdk already pre-installed. That's not much use right now, but if we want to deliver a workshop in the future, that may be very useful.

## Explanation of CRs
He also has a nice explanation of what a custom resource is, and has a nice example of creating one [here](https://ibm.github.io/kubernetes-operators/lab1/). This may be useful to include in the intro of our "develop a simple operator" tutorial, issue #48

## Explanation of Operator development workflow
His "About Operators and About the Operator Framework" intro is also useful, and we may be able to recycle the workflow for a new Go operator, i.e.:

The following workflow is for a new Go operator:

1. Create a new operator project using the SDK Command Line Interface(CLI)
2. Define new resource APIs by adding Custom Resource Definitions(CRD)
3. Define Controllers to watch and reconcile resources
4. Write the reconciling logic for your Controller using the SDK and controller-runtime APIs
5. Use the SDK CLI to build and generate the operator deployment manifests

## Create an Operator of Type Go using the Operator SDK
The rest of the [Go-based operator tutorial](https://ibm.github.io/kubernetes-operators/lab2/) features the old operator-sdk commands, so that is not very useful, but we may want to use his 
structure, since it is pretty simple and quick.

 Also, there is an operator from [Helm tutorial](https://ibm.github.io/kubernetes-operators/lab3/) too.