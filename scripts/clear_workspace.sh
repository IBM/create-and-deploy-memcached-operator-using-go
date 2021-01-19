#!/bin/bash
oc delete -f config/samples/cache_v1alpha1_memcached.yaml
oc delete deployments,service -l control-plane=controller-manager
oc delete role,rolebinding --all
# make undeploy
# make uninstall
