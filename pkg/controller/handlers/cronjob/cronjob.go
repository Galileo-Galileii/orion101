package cronjob

import (
	"fmt"
	"time"

	"github.com/orion101-ai/nah/pkg/router"
	"github.com/orion101-ai/orion101/apiclient/types"
	"github.com/orion101-ai/orion101/pkg/alias"
	"github.com/orion101-ai/orion101/pkg/storage/apis/orion101.orion101.ai/v1"
	"github.com/orion101-ai/orion101/pkg/system"
	"github.com/robfig/cron/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Run(req router.Request, resp router.Response) error {
	cj := req.Object.(*v1.CronJob)
	lastRun := cj.Status.LastRunStartedAt
	if lastRun.IsZero() {
		lastRun = &cj.CreationTimestamp
	}

	sched, err := cron.ParseStandard(cj.Spec.Schedule)
	if err != nil {
		return fmt.Errorf("failed to parse schedule: %w", err)
	}

	if until := time.Until(sched.Next(lastRun.Time)); until > 0 {
		resp.RetryAfter(until)
		return nil
	}

	var workflow v1.Workflow
	if err := alias.Get(req.Ctx, req.Client, &workflow, cj.Namespace, cj.Spec.Workflow); err != nil {
		return err
	}

	if err = req.Client.Create(req.Ctx,
		&v1.WorkflowExecution{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: system.WorkflowExecutionPrefix,
				Namespace:    req.Namespace,
			},
			Spec: v1.WorkflowExecutionSpec{
				WorkflowName: workflow.Name,
				Input:        cj.Spec.Input,
				CronJobName:  cj.Name,
			},
		},
	); err != nil {
		return err
	}

	cj.Status.LastRunStartedAt = &[]metav1.Time{metav1.Now()}[0]

	return nil
}

func (h *Handler) SetSuccessRunTime(req router.Request, _ router.Response) error {
	cj := req.Object.(*v1.CronJob)

	var workflowExecutions v1.WorkflowExecutionList
	if err := req.List(&workflowExecutions, &kclient.ListOptions{
		FieldSelector: fields.SelectorFromSet(map[string]string{"spec.cronJobName": cj.Name}),
		Namespace:     cj.Namespace,
	}); err != nil {
		return err
	}

	for _, execution := range workflowExecutions.Items {
		if execution.Status.State == types.WorkflowStateComplete && (cj.Status.LastSuccessfulRunCompleted == nil || cj.Status.LastSuccessfulRunCompleted.Before(execution.Status.EndTime)) {
			cj.Status.LastSuccessfulRunCompleted = execution.Status.EndTime
		}
	}

	return nil
}
