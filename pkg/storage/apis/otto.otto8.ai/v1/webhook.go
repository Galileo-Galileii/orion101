package v1

import (
	"github.com/orion101-ai/orion101/apiclient/types"
	"github.com/orion101-ai/orion101/pkg/system"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	_ Aliasable = (*Webhook)(nil)
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Webhook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebhookSpec   `json:"spec,omitempty"`
	Status WebhookStatus `json:"status,omitempty"`
}

func (w *Webhook) GetAliasName() string {
	return w.Spec.WebhookManifest.Alias
}

func (w *Webhook) SetAssigned(assigned bool) {
	w.Status.AliasAssigned = assigned
}

func (w *Webhook) IsAssigned() bool {
	return w.Status.AliasAssigned
}

func (w *Webhook) GetAliasObservedGeneration() int64 {
	return w.Status.AliasObservedGeneration
}

func (w *Webhook) SetAliasObservedGeneration(gen int64) {
	w.Status.AliasObservedGeneration = gen
}

func (*Webhook) GetColumns() [][]string {
	return [][]string{
		{"Name", "Name"},
		{"Alias", "Spec.Alias"},
		{"Workflow", "Spec.Workflow"},
		{"Created", "{{ago .CreationTimestamp}}"},
		{"Last Success", "{{ago .Status.LastSuccessfulRunCompleted}}"},
		{"Description", "Spec.Description"},
	}
}

func (w *Webhook) DeleteRefs() []Ref {
	if system.IsWebhookID(w.Spec.Workflow) {
		return []Ref{
			{ObjType: new(Workflow), Name: w.Spec.Workflow},
		}
	}
	return nil
}

type WebhookSpec struct {
	types.WebhookManifest `json:",inline"`
	TokenHash             []byte `json:"tokenHash,omitempty"`
}

type WebhookStatus struct {
	AliasAssigned              bool         `json:"aliasAssigned,omitempty"`
	LastSuccessfulRunCompleted *metav1.Time `json:"lastSuccessfulRunCompleted,omitempty"`
	AliasObservedGeneration    int64        `json:"aliasProcessed,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type WebhookList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Webhook `json:"items"`
}