package sabnzbd

import (
	"context"
	"fmt"
	"github.com/parflesh/sabnzbd-operator/defaults"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	sabnzbdv1alpha1 "github.com/parflesh/sabnzbd-operator/pkg/apis/sabnzbd/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_sabnzbd")

// Add creates a new SABnzbd Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSABnzbd{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("sabnzbd-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource SABnzbd
	err = c.Watch(&source.Kind{Type: &sabnzbdv1alpha1.SABnzbd{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sabnzbdv1alpha1.SABnzbd{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sabnzbdv1alpha1.SABnzbd{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileSABnzbd implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileSABnzbd{}

// ReconcileSABnzbd reconciles a SABnzbd object
type ReconcileSABnzbd struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

func (r *ReconcileSABnzbd) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling SABnzbd")

	// Fetch the SABnzbd instance
	instance := &sabnzbdv1alpha1.SABnzbd{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	newStatus := instance.Status

	err = r.reconcileSpec(instance)
	if err != nil {
		err := r.client.Update(context.TODO(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}

	/*imageManifest, err := r.imageInspector.GetImageLabels(ctx, instance.Spec.Image)
	if err != nil {
		return reconcile.Result{}, err
	}*/

	newStatus.Image = instance.Spec.Image
	if newStatus.Image != instance.Status.Image {
		instance.Status.Image = newStatus.Image
		if err := r.client.Status().Update(context.TODO(), instance); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}

	newDep, err := r.newDeployment(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	foundDep := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), request.NamespacedName, foundDep)
	if err != nil && errors.IsNotFound(err) {
		err := r.client.Create(context.TODO(), newDep)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	newSvc, err := r.newService(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	foundSvc := &corev1.Service{}
	err = r.client.Get(context.TODO(), request.NamespacedName, foundSvc)
	if err != nil && errors.IsNotFound(err) {
		err := r.client.Create(context.TODO(), newSvc)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileSABnzbd) reconcileSpec(cr *sabnzbdv1alpha1.SABnzbd) error {
	if cr.Spec.Image == "" {
		cr.Spec.Image = defaults.SABnzbdImage
		return fmt.Errorf("image not set")
	}
	if cr.Spec.WatchFrequency == "" {
		cr.Spec.WatchFrequency = defaults.OperatorRequeuTime
		return fmt.Errorf("watch frequency not set")
	}
	return nil
}

func (r *ReconcileSABnzbd) newDeployment(cr *sabnzbdv1alpha1.SABnzbd) (*appsv1.Deployment, error) {
	labels := r.labelsForCR(cr)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &[]int32{1}[0],
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name:         "config",
							VolumeSource: cr.Spec.ConfigVolume,
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "sabnzbd",
							Image: cr.Status.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8080,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Resources: corev1.ResourceRequirements{},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/config",
								},
							},
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "",
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 8080,
											StrVal: "",
										},
										Scheme: corev1.URISchemeHTTP,
									},
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "",
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 8080,
											StrVal: "",
										},
										Scheme: corev1.URISchemeHTTP,
									},
								},
							},
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
					RestartPolicy:     corev1.RestartPolicyAlways,
					SecurityContext:   &corev1.PodSecurityContext{},
					ImagePullSecrets:  cr.Spec.ImagePullSecrets,
					PriorityClassName: cr.Spec.PriorityClassName,
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			RevisionHistoryLimit: &[]int32{5}[0],
		},
	}

	if cr.Spec.RunAsUser != int64(0) {
		dep.Spec.Template.Spec.Containers[0].SecurityContext.RunAsUser = &cr.Spec.RunAsUser
	}

	if cr.Spec.RunAsGroup != int64(0) {
		dep.Spec.Template.Spec.Containers[0].SecurityContext.RunAsUser = &cr.Spec.RunAsUser
	}

	err := controllerutil.SetControllerReference(cr, dep, r.scheme)
	if err != nil {
		return dep, err
	}
	return dep, nil
}

func (r *ReconcileSABnzbd) newService(cr *sabnzbdv1alpha1.SABnzbd) (*corev1.Service, error) {
	labels := r.labelsForCR(cr)

	dep := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Protocol: corev1.ProtocolTCP,
					Port:     8080,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 8080,
						StrVal: "",
					},
				},
			},
			Selector: labels,
		},
	}

	err := controllerutil.SetControllerReference(cr, dep, r.scheme)
	if err != nil {
		return dep, err
	}

	return dep, nil
}

func (r *ReconcileSABnzbd) labelsForCR(cr *sabnzbdv1alpha1.SABnzbd) map[string]string {
	return map[string]string{
		"sabnzbd": cr.Name,
	}
}
