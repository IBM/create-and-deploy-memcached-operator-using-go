/*
Copyright 2021.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"reflect"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.ibm.com/TT-ISV-org/janusgraph-operator/api/v1alpha1"
	graphv1alpha1 "github.ibm.com/TT-ISV-org/janusgraph-operator/api/v1alpha1"
)

// JanusgraphReconciler reconciles a Janusgraph object
type JanusgraphReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=graph.ibm.com,resources=janusgraphs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=graph.ibm.com,resources=janusgraphs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=graph.ibm.com,resources=janusgraphs/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=pods;deployments;statefulsets;services;persistentvolumeclaims;persistentvolumes;,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods;services;persistentvolumeclaims;persistentvolumes;,verbs=get;list;create;update;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Janusgraph object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *JanusgraphReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("janusgraph", req.NamespacedName)

	// Fetch the Janusgraph instance
	janusgraph := &graphv1alpha1.Janusgraph{}
	err := r.Get(ctx, req.NamespacedName, janusgraph)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Janusgraph resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Janusgraph")
		return ctrl.Result{}, err
	}

	serviceFound := &corev1.Service{}
	log.Info("Checking for service")
	err = r.Get(ctx, types.NamespacedName{Name: janusgraph.Name + "-service", Namespace: janusgraph.Namespace}, serviceFound)
	if err != nil && errors.IsNotFound(err) {
		srv := r.serviceForJanusgraph(janusgraph)
		log.Info("Creating a new headless service", "Service.Namespace", srv.Namespace, "Service.Name", srv.Name)
		err = r.Create(ctx, srv)
		if err != nil {
			log.Error(err, "Failed to create new service", "service.Namespace", srv.Namespace, "service.Name", srv.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		log.Info("Janusgraph service created, requeuing")
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get service")
		return ctrl.Result{}, err
	}

	// deployment
	found := &appsv1.StatefulSet{}
	// Check if the deployment already exists, if not create a new one
	err = r.Get(ctx, types.NamespacedName{Name: janusgraph.Name, Namespace: janusgraph.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForJanusgraph(janusgraph)
		log.Info("Creating a new Statefulset", "StatefulSet.Namespace", dep.Namespace, "StatefulSet.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new StatefulSet", "StatefulSet.Namespace", dep.Namespace, "StatefulSet.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		log.Info("Deployment created, requeuing")
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// Status check
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(janusgraph.Namespace),
		client.MatchingLabels(labelsForJanusgraph(janusgraph.Name)),
	}
	if err = r.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "Janusgraph.Namespace", janusgraph.Namespace, "Janusgraph.Name", janusgraph.Name)
		return ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, janusgraph.Status.Nodes) {
		janusgraph.Status.Nodes = podNames
		err := r.Status().Update(ctx, janusgraph)
		if err != nil {
			log.Error(err, "Failed to update Janusgraph status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

// SetupWithManager sets up the controller with the Manager.
func (r *JanusgraphReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&graphv1alpha1.Janusgraph{}).
		Complete(r)
}

func labelsForJanusgraph(name string) map[string]string {
	return map[string]string{"app": "Janusgraph", "janusgraph_cr": name}
}

func (r *JanusgraphReconciler) serviceForJanusgraph(m *v1alpha1.Janusgraph) *corev1.Service {
	ls := labelsForJanusgraph(m.Name)
	srv := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-service",
			Namespace: m.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeLoadBalancer,
			Ports: []corev1.ServicePort{
				{
					Port: 8182,
					TargetPort: intstr.IntOrString{
						IntVal: 8182,
					},
					NodePort: 30184,
				},
			},
			Selector: ls,
		},
	}
	ctrl.SetControllerReference(m, srv, r.Scheme)
	return srv
}

func (r *JanusgraphReconciler) deploymentForJanusgraph(m *v1alpha1.Janusgraph) *appsv1.StatefulSet {
	ls := labelsForJanusgraph(m.Name)
	replicas := m.Spec.Size
	version := m.Spec.Version

	dep := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			ServiceName: m.Name + "-service",
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
					Name:   "janusgraph",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "sanjeevghimire/janusgraph:" + version,
							Name:  "janusgraph",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8182,
									Name:          "janusgraph",
								},
							},
							Env: []corev1.EnvVar{},
						}},
					RestartPolicy: corev1.RestartPolicyAlways,
				},
			},
		},
	}
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}
