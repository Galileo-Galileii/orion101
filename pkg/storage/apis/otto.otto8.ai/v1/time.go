package v1

import (
	"github.com/orion101-ai/orion101/apiclient/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewTime(t *metav1.Time) *types.Time {
	if t == nil {
		return nil
	}
	return types.NewTime(t.Time)
}
