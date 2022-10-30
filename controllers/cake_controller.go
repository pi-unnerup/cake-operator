package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	tutorialsv1 "github.com/pi-unnerup/cake-operator/api/v1"
)

// CakeReconciler reconciles a Cake object
type CakeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tutorials.piunnerup.com,resources={cakes,secrets,serviceaccounts,services},verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tutorials.piunnerup.com,resources=cakes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tutorials.piunnerup.com,resources=cakes/finalizers,verbs=update
// Custom RBAC to allow the operator to interact with mandatory resources
//+kubebuilder:rbac:groups="",resources={secrets,serviceaccounts,services},verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

func (r *CakeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)
	cake := &tutorialsv1.Cake{}
	existingNginxDeployment := &appsv1.Deployment{}
	existingService := &corev1.Service{}

	log.Info("⚡️ Event received! ⚡️")
	log.Info("Request: ", "req", req)

	// CR deleted : check if  the Deployment and the Service must be deleted
	err := r.Get(ctx, req.NamespacedName, cake)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Cake resource not found, check if a deployment must be deleted.")

			// Delete Deployment
			err = r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, existingNginxDeployment)
			if err != nil {
				if errors.IsNotFound(err) {
					log.Info("Nothing to do, no deployment found.")
					return ctrl.Result{}, nil
				} else {
					log.Error(err, "❌ Failed to get Deployment")
					return ctrl.Result{}, err
				}
			} else {
				log.Info("☠️ Deployment exists: delete it. ☠️")
				r.Delete(ctx, existingNginxDeployment)
			}

			// Delete Service
			err = r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, existingService)
			if err != nil {
				if errors.IsNotFound(err) {
					log.Info("Nothing to do, no service found.")
					return ctrl.Result{}, nil
				} else {
					log.Error(err, "❌ Failed to get Service")
					return ctrl.Result{}, err
				}
			} else {
				log.Info("☠️ Service exists: delete it. ☠️")
				r.Delete(ctx, existingService)
				return ctrl.Result{}, nil
			}
		}
	} else {
		log.Info("ℹ️  CR state ℹ️", "cake.Name", cake.Name, " cake.Namespace", cake.Namespace, "cake.Spec.ReplicaCount", cake.Spec.ReplicaCount, "cake.Spec.Port", cake.Spec.Port)

		// Check if the deployment already exists, if not: create a new one.
		err = r.Get(ctx, types.NamespacedName{Name: cake.Name, Namespace: cake.Namespace}, existingNginxDeployment)
		if err != nil && errors.IsNotFound(err) {
			// Define a new deployment
			newNginxDeployment := r.createDeployment(cake)
			log.Info("✨ Creating a new Deployment", "Deployment.Namespace", newNginxDeployment.Namespace, "Deployment.Name", newNginxDeployment.Name)

			err = r.Create(ctx, newNginxDeployment)
			if err != nil {
				log.Error(err, "❌ Failed to create new Deployment", "Deployment.Namespace", newNginxDeployment.Namespace, "Deployment.Name", newNginxDeployment.Name)
				return ctrl.Result{}, err
			}
		} else if err == nil {
			// Deployment exists, check if the Deployment must be updated
			var replicaCount int32 = cake.Spec.ReplicaCount
			if *existingNginxDeployment.Spec.Replicas != replicaCount {
				log.Info("🔁 Number of replicas changes, update the deployment! 🔁")
				existingNginxDeployment.Spec.Replicas = &replicaCount
				err = r.Update(ctx, existingNginxDeployment)
				if err != nil {
					log.Error(err, "❌ Failed to update Deployment", "Deployment.Namespace", existingNginxDeployment.Namespace, "Deployment.Name", existingNginxDeployment.Name)
					return ctrl.Result{}, err
				}
			}
		} else if err != nil {
			log.Error(err, "Failed to get Deployment")
			return ctrl.Result{}, err
		}

		// Check if the service already exists, if not: create a new one
		err = r.Get(ctx, types.NamespacedName{Name: cake.Name, Namespace: cake.Namespace}, existingService)
		if err != nil && errors.IsNotFound(err) {
			// Create the Service
			newService := r.createService(cake)
			log.Info("✨ Creating a new Service", "Service.Namespace", newService.Namespace, "Service.Name", newService.Name)
			err = r.Create(ctx, newService)
			if err != nil {
				log.Error(err, "❌ Failed to create new Service", "Service.Namespace", newService.Namespace, "Service.Name", newService.Name)
				return ctrl.Result{}, err
			}
		} else if err == nil {
			// TODO: update for env??
			// Service exists, check if the port has to be updated.
			var port int32 = cake.Spec.Port
			if existingService.Spec.Ports[0].Port != port {
				log.Info("🔁 Port number changes, update the service! 🔁")
				existingService.Spec.Ports[0].Port = port
				err = r.Update(ctx, existingService)
				if err != nil {
					log.Error(err, "❌ Failed to update Service", "Service.Namespace", existingService.Namespace, "Service.Name", existingService.Name)
					return ctrl.Result{}, err
				}
			}
		} else if err != nil {
			log.Error(err, "Failed to get Service")
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// Create a Deployment for the Nginx server.
func (r *CakeReconciler) createDeployment(cakeCR *tutorialsv1.Cake) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cakeCR.Name,
			Namespace: cakeCR.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &cakeCR.Spec.ReplicaCount,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "cake-server"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "cake-server"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "unnerup/muffin-time:v0.0.3",
						Name:  "cake",
						Env: []corev1.EnvVar{
							{
								Name:  "COLOUR",
								Value: cakeCR.Spec.COLOUR,
							},
							{
								Name:  "MESSAGE",
								Value: cakeCR.Spec.MESSAGE,
							},
							{
								Name:  "DECORATION",
								Value: cakeCR.Spec.DECORATION,
							},
						},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Name:          "http",
							Protocol:      "TCP",
						}},
					}},
				},
			},
		},
	}
	return deployment
}

// Create a Service for the Nginx server.
func (r *CakeReconciler) createService(cakeCR *tutorialsv1.Cake) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cakeCR.Name,
			Namespace: cakeCR.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "cake-server",
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       cakeCR.Spec.Port,
					TargetPort: intstr.FromInt(80),
				},
			},
			Type: corev1.ServiceTypeLoadBalancer,
		},
	}

	return service
}

// SetupWithManager sets up the controller with the Manager.
func (r *CakeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tutorialsv1.Cake{}).
		Watches(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForObject{}, builder.WithPredicates(predicate.Funcs{
			// Check only delete events for a service
			UpdateFunc: func(e event.UpdateEvent) bool {
				return false
			},
			CreateFunc: func(e event.CreateEvent) bool {
				return false
			},
			DeleteFunc: func(e event.DeleteEvent) bool {
				return true
			},
		})).
		Complete(r)
}
