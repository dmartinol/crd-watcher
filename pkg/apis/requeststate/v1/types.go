package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type RequestState struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	Status RequestStateStatus `json:"status,omitempty"`

	Spec RequestStateSpec `json:"spec,omitempty"`
}

type RequestStateSpec struct {
	RequestUid string `json:"request-uid,omitempty"`
	Job        string `json:"job,omitempty"`
	State      string `json:"state,omitempty"`
}

type RequestStateStatus struct {
	History []RequestStateHistory `json:"history,omitempty"`
}

type RequestStateHistory struct {
	Job       string `json:"job,omitempty"`
	State     string `json:"state,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type RequestStateList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `son:"metadata,omitempty"`

	Items []RequestState `json:"items"`
}
