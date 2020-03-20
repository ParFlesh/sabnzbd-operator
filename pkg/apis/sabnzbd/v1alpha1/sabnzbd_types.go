package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SABnzbdSpec defines the desired state of SABnzbd
type SABnzbdSpec struct {
	// SABnzbd Config Volume mounted to /config
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Config Volume"
	// +optional
	ConfigVolume corev1.VolumeSource `json:"configVolume,omitempty"`

	// Container image capable of running SABnzbd (Default: quay.io/parflesh/sabnzbd:latest)
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Container Image"
	// +optional
	Image string `json:"image,omitempty"`

	// Additional Volumes to be mounted in SABnzbd container
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Additional Volumes"
	// +optional
	AdditionalVolumes []corev1.Volume `json:"additionalVolumes,omitempty"`

	// Stop automatic updates when hash for image tag changes
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Disable Image Updates"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:booleanSwitch"
	// +optional
	DisableUpdates bool `json:"disableUpdates,omitempty"`

	// Image pull secret for private container images
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Image Pull Secret"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:io.kubernetes:Secret"
	// +optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecret,omitempty"`

	// Time to wait between checking resource status (Default: 1m)
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Watch Frequency"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text"
	// +optional
	WatchFrequency string `json:"watchFrequency,omitempty"`

	// Priority Class Name
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="Priority Class Name"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text"
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`

	// Run as User Id
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="User ID"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:number"
	// +optional
	RunAsUser int64 `json:"runAsUser,omitempty"`

	// Run as Group Id
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.displayName="GroupID"
	// +operator-sdk:gen-csv:customresourcedefinitions.specDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:number"
	// +optional
	RunAsGroup int64 `json:"runAsGroup,omitempty"`
}

// SABnzbdStatus defines the observed state of SABnzbd
type SABnzbdStatus struct {
	// Desired Image hash for container
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.displayName="Image"
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors.x-descriptors="urn:alm:descriptor:com.tectonic.ui:text"
	Image string `json:"image,omitempty"`

	// Phase
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	Phase string `json:"phase,omitempty"`

	// Reason
	// +operator-sdk:gen-csv:customresourcedefinitions.statusDescriptors=true
	Reason string `json:"reason,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SABnzbd is the Schema for the sabnzbds API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=sabnzbds,scope=Namespaced
// +operator-sdk:gen-csv:customresourcedefinitions.displayName="SABnzbd"
// +operator-sdk:gen-csv:customresourcedefinitions.resources=`Deployment,v1,"sabnzbd-operator"`
// +operator-sdk:gen-csv:customresourcedefinitions.resources=`Service,v1,"sabnzbd-operator"`
type SABnzbd struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SABnzbdSpec   `json:"spec,omitempty"`
	Status SABnzbdStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SABnzbdList contains a list of SABnzbd
type SABnzbdList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SABnzbd `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SABnzbd{}, &SABnzbdList{})
}
