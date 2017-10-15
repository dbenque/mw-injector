package v1

import (
	v1beta2 "k8s.io/api/apps/v1beta2"
	api "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MWInjector represents a definition to inject MW in multiple namespaces
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MWInjector struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec represents the desired behaviour of the MWInjector.
	Spec MWInjectorSpec `json:"spec,omitempty"`

	// Status contains the current status off the MWInjector
	Status MWInjectorStatus `json:"status,omitempty"`
}

// MWInjectorList implements list of MWInjector.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MWInjectorList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of MWInjector
	Items []MWInjector `json:"items"`
}

// MWInjectorSpec contains MWInjector specification
type MWInjectorSpec struct {
	//MWDeployment is the deployment of the mw to injected in namespaces
	MWDeployment v1beta2.Deployment

	// Selector for created jobs (if any)
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty"`
}

// MWInjectorConditionType is the type of status's condition for MWInjector
type MWInjectorConditionType string

// These are valid conditions of a workflow.
const (
	// MWInjectorRunning means the mwinjector is running.
	MWInjectorRunning MWInjectorConditionType = "Running"
	// MWInjectorStopped means the mwinjector is Stopped.
	MWInjectorStopped MWInjectorConditionType = "Stopped"
	// MWInjectorFailed means the workflow has failed its execution.
	MWInjectorFailed MWInjectorConditionType = "Failed"
)

// MWInjectorCondition represent the condition of the MWInjector
type MWInjectorCondition struct {
	// Type of MWInjector condition
	Type MWInjectorConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status api.ConditionStatus `json:"status"`
	// Last time the condition was checked.
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// Last time the condition transited from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// (brief) reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// Human readable message indicating details about last transition.
	Message string `json:"message,omitempty"`
}

//MWInjectorStatus represent the status of the MWInjector
type MWInjectorStatus struct {
	// Conditions represent the latest available observations of an object's current state.
	Conditions []MWInjectorCondition
}
