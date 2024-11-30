package v1

import (
	"github.com/orion101-ai/orion101/apiclient/types"
	"github.com/orion101-ai/orion101/pkg/system"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CronJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CronJobSpec   `json:"spec,omitempty"`
	Status            CronJobStatus `json:"status,omitempty"`
}

func (*CronJob) GetColumns() [][]string {
	return [][]string{
		{"Name", "Name"},
		{"Workflow", "Spec.WorkflowID"},
		{"Schedule", "Spec.Schedule"},
		{"Last Success", "{{ago .Status.LastSuccessfulRunCompleted}}"},
		{"Last Run", "{{ago .Status.LastRunStartedAt}}"},
		{"Created", "{{ago .CreationTimestamp}}"},
		{"Description", "Spec.Description"},
	}
}

func (c *CronJob) DeleteRefs() []Ref {
	if system.IsWorkflowID(c.Spec.Workflow) {
		return []Ref{
			{ObjType: new(Workflow), Name: c.Spec.Workflow},
		}
	}
	return nil
}

type CronJobSpec struct {
	types.CronJobManifest `json:",inline"`
}

type CronJobStatus struct {
	LastRunStartedAt           *metav1.Time `json:"lastRunStartedAt,omitempty"`
	LastSuccessfulRunCompleted *metav1.Time `json:"lastSuccessfulRunCompleted,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CronJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CronJob `json:"items"`
}
