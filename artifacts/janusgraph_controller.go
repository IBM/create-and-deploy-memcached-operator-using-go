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
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	graphv1alpha1 "github.com/example/janusgraph-operator/api/v1alpha1"
)

// JanusgraphReconciler reconciles a Janusgraph object
type JanusgraphReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=graph.example.com,resources=janusgraphs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=graph.example.com,resources=janusgraphs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=graph.example.com,resources=janusgraphs/finalizers,verbs=update
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

	// fetch Service resource
	serviceFound := &corev1.Service{}
	log.Info("Checking for service")
	//check for Service resources in our namespace, and with a "JanusGraph" name prefix
	err = r.Get(ctx, types.NamespacedName{Name: janusgraph.Name + "-service", Namespace: janusgraph.Namespace}, serviceFound)
	if err != nil && errors.IsNotFound(err) {
		srv := r.serviceForJanusgraph(janusgraph)
		log.Info("Creating a new headless service", "Service.Namespace", srv.Namespace, "Service.Name", srv.Name)
		err = r.Create(ctx, srv)
		if err != nil {
			log.Error(err, "Failed to create new service", "service.Namespace", srv.Namespace, "service.Name", srv.Name)
			return ctrl.Result{}, err
		}
		// Service created successfully - return and requeue
		log.Info("Janusgraph service created, requeuing")
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get service")
		return ctrl.Result{}, err
	}

	// look for a resource of type StatefulSet
	found := &appsv1.StatefulSet{}
	// Check if the StatefulSet already exists in our namespace, if not create a new one
	err = r.Get(ctx, types.NamespacedName{Name: janusgraph.Name, Namespace: janusgraph.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new StatefulSet
		statefulSet := r.statefulSetForJanusgraph(janusgraph)
		log.Info("Creating a new Statefulset", "StatefulSet.Namespace", statefulSet.Namespace, "StatefulSet.Name", statefulSet.Name)
		err = r.Create(ctx, statefulSet)
		if err != nil {
			log.Error(err, "Failed to create new StatefulSet", "StatefulSet.Namespace", statefulSet.Namespace, "StatefulSet.Name", statefulSet.Name)
			return ctrl.Result{}, err
		}
		// StatefulSet created successfully - return and requeue
		log.Info("StatefulSet created, requeuing")
		return ctrl.Result{}, nil
	} else if err != nil {
		log.Error(err, "Failed to get StatefulSet")
		return ctrl.Result{}, err
	}

	// look for resource of type PodList
	podList := &corev1.PodList{}
	//create filter to check for Pods only in our Namespace with the correct matching labels
	listOpts := []client.ListOption{
		client.InNamespace(janusgraph.Namespace),
		client.MatchingLabels(labelsForJanusgraph(janusgraph.Name)),
	}
	//List all Pods that match our filter (same Namespace and matching labels)
	if err = r.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "Janusgraph.Namespace", janusgraph.Namespace, "Janusgraph.Name", janusgraph.Name)
		return ctrl.Result{}, err
	}
	//return an array of pod names
	podNames := getPodNames(podList.Items)

	// Update the status of our JanusGraph object to show Pods which were returned from getPodNames
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

// getPodNames returns a string array of Pod Names
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

// labelsForJanusgraph returns a map of string keys and string values
func labelsForJanusgraph(name string) map[string]string {
	return map[string]string{"app": "Janusgraph", "janusgraph_cr": name}
}

// serviceForJanusgraph returns a Load Balancer service for our JanusGraph object
func (r *JanusgraphReconciler) serviceForJanusgraph(m *graphv1alpha1.Janusgraph) *corev1.Service {

	//fetch labels
	ls := labelsForJanusgraph(m.Name)
	//create Service
	srv := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-service",
			Namespace: m.Namespace,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone, //"None",
			Ports: []corev1.ServicePort{{
				Port: 8182,
				Name: "janusgraph",
			},
			},
			Selector: ls,
		},
	}
	ctrl.SetControllerReference(m, srv, r.Scheme)
	return srv
}

// statefulSetForJanusgraph returns a StatefulSet for our JanusGraph object
func (r *JanusgraphReconciler) statefulSetForJanusgraph(m *graphv1alpha1.Janusgraph) *appsv1.StatefulSet {

	//fetch labels
	ls := labelsForJanusgraph(m.Name)
	//fetch the size of the JanusGraph object from the custom resource
	replicas := m.Spec.Size
	//fetch the version of JanusGraph to install from the custom resource
	version := m.Spec.Version

	//create StatefulSet
	statefulSet := &appsv1.StatefulSet{
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
							Image: "horeaporutiu/janusgraph:" + version,
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
	ctrl.SetControllerReference(m, statefulSet, r.Scheme)
	return statefulSet
}
