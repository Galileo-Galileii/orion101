package v1

import (
	"github.com/orion101-ai/orion101/apiclient/types"
	"github.com/orion101-ai/orion101/pkg/system"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	_ Aliasable = (*EmailReceiver)(nil)
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type EmailReceiver struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmailReceiverSpec   `json:"spec,omitempty"`
	Status EmailReceiverStatus `json:"status,omitempty"`
}

func (in *EmailReceiver) GetAliasName() string {
	return in.Spec.EmailReceiverManifest.User
}

func (in *EmailReceiver) SetAssigned(assigned bool) {
	in.Status.AliasAssigned = assigned
}

func (in *EmailReceiver) IsAssigned() bool {
	return in.Status.AliasAssigned
}

func (in *EmailReceiver) GetAliasObservedGeneration() int64 {
	return in.Status.AliasObservedGeneration
}

func (in *EmailReceiver) SetAliasObservedGeneration(gen int64) {
	in.Status.AliasObservedGeneration = gen
}

func (*EmailReceiver) GetColumns() [][]string {
	return [][]string{
		{"Name", "Name"},
		{"User", "Spec.User"},
		{"Workflow", "Spec.Workflow"},
		{"Created", "{{ago .CreationTimestamp}}"},
		{"Description", "Spec.Description"},
	}
}

func (in *EmailReceiver) DeleteRefs() []Ref {
	if system.IsWorkflowID(in.Spec.Workflow) {
		return []Ref{
			{ObjType: new(Workflow), Name: in.Spec.Workflow},
		}
	}
	return nil
}

type EmailReceiverSpec struct {
	types.EmailReceiverManifest `json:",inline"`
}

type EmailReceiverStatus struct {
	AliasAssigned           bool  `json:"aliasAssigned,omitempty"`
	AliasObservedGeneration int64 `json:"aliasProcessed,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type EmailReceiverList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EmailReceiver `json:"items"`
}