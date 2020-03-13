package sabnzbd

import (
    "testing"

    sabnzbdv1alpha1 "github.com/parflesh/sabnzbd-operator/pkg/apis/sabnzbd/v1alpha1"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/types"
    "k8s.io/client-go/kubernetes/scheme"
    "sigs.k8s.io/controller-runtime/pkg/client/fake"
    "sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestSABnzbdController(t *testing.T) {
    var (
        name            = "sabnzbd-operator"
        namespace       = "sabnzbd"
    )
    // A SABnzbd object with metadata and spec.
    cr := &sabnzbdv1alpha1.SABnzbd{
        ObjectMeta: metav1.ObjectMeta{
            Name:      name,
            Namespace: namespace,
        },
        Spec: sabnzbdv1alpha1.SABnzbdSpec{},
    }

    // Objects to track in the fake client.
    objs := []runtime.Object{ cr }

    // Register operator types with the runtime scheme.
    s := scheme.Scheme
    s.AddKnownTypes(sabnzbdv1alpha1.SchemeGroupVersion, cr)

    // Create a fake client to mock API calls.
    cl := fake.NewFakeClient(objs...)

    // Create a ReconcileSABnzbd object with the scheme and fake client.
    r := &ReconcileSABnzbd{client: cl, scheme: s}

    // Mock request to simulate Reconcile() being called on an event for a
    // watched resource .
    req := reconcile.Request{
        NamespacedName: types.NamespacedName{
            Name:      name,
            Namespace: namespace,
        },
    }
    res, err := r.Reconcile(req)
    if err != nil {
       t.Fatalf("reconcile: (%v)", err)
    }
    // Check the result of reconciliation to make sure it has the desired state.
    if res.Requeue {
		t.Error("reconcile requeued even though all should be good")
	}
}