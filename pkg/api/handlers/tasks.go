package handlers

import (
	"net/http"
	"slices"

	"github.com/orion101-ai/orion101/apiclient/types"
	"github.com/orion101-ai/orion101/pkg/api"
	"github.com/orion101-ai/orion101/pkg/events"
	"github.com/orion101-ai/orion101/pkg/invoke"
	v1 "github.com/orion101-ai/orion101/pkg/storage/apis/orion101.orion101.ai/v1"
	"github.com/orion101-ai/orion101/pkg/system"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type TaskHandler struct {
	invoker *invoke.Invoker
	events  *events.Emitter
}

func NewTaskHandler(invoker *invoke.Invoker, events *events.Emitter) *TaskHandler {
	return &TaskHandler{
		invoker: invoker,
		events:  events,
	}
}

func (t *TaskHandler) Events(req api.Context) error {
	var (
		follow = req.URL.Query().Get("follow") == "true"
	)

	workflow, err := t.getTask(req)
	if err != nil {
		return err
	}

	var thread v1.Thread
	if err := req.Get(&thread, req.PathValue("thread_id")); kclient.IgnoreNotFound(err) != nil {
		return err
	}

	if thread.Spec.WorkflowName != workflow.Name {
		return types.NewErrHttp(http.StatusForbidden, "thread does not belong to the task")
	}

	_, events, err := t.events.Watch(req.Context(), req.Namespace(), events.WatchOptions{
		History:                  true,
		MaxRuns:                  100,
		ThreadName:               thread.Name,
		Follow:                   true,
		FollowWorkflowExecutions: follow,
	})
	if err != nil {
		return err
	}

	return req.WriteEvents(events)
}

func (t *TaskHandler) Run(req api.Context) error {
	var (
		threadID = req.Request.URL.Query().Get("thread")
		stepID   = req.Request.URL.Query().Get("step")
	)

	workflow, err := t.getTask(req)
	if err != nil {
		return err
	}

	resp, err := t.invoker.Workflow(req.Context(), req.Storage, workflow, "", invoke.WorkflowOptions{
		ThreadName: threadID,
		StepID:     stepID,
	})
	if err != nil {
		return err
	}

	return req.WriteCreated(map[string]any{
		"threadID": resp.Thread.Name,
	})
}

func (t *TaskHandler) Delete(req api.Context) error {
	workflow, err := t.getTask(req)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	return req.Delete(workflow)
}

func (t *TaskHandler) Update(req api.Context) error {
	workflow, err := t.getTask(req)
	if err != nil {
		return err
	}

	_, manifest, err := t.getAssistantAndManifestFromRequest(req)
	if err != nil {
		return err
	}

	workflow.Spec.Manifest = manifest
	if err := req.Update(workflow); err != nil {
		return err
	}

	return req.Write(convertTask(*workflow))
}

func (t *TaskHandler) getAssistantAndManifestFromRequest(req api.Context) (*v1.Agent, types.WorkflowManifest, error) {
	assistantID := req.PathValue("assistant_id")

	assistant, err := getAssistant(req, assistantID)
	if err != nil {
		return nil, types.WorkflowManifest{}, err
	}

	thread, err := getUserThread(req, assistantID)
	if err != nil {
		return nil, types.WorkflowManifest{}, err
	}

	var manifest types.TaskManifest
	if err := req.Read(&manifest); err != nil {
		return nil, types.WorkflowManifest{}, err
	}

	if manifest.Name == "" {
		manifest.Name = "New Task"
	}

	return assistant, toWorkflowManifest(assistant, thread, manifest), nil
}

func (t *TaskHandler) Create(req api.Context) error {
	assistant, workflowManifest, err := t.getAssistantAndManifestFromRequest(req)
	if err != nil {
		return err
	}

	workflow := v1.Workflow{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: system.WorkflowPrefix,
			Namespace:    req.Namespace(),
		},
		Spec: v1.WorkflowSpec{
			AgentName: assistant.Name,
			UserID:    req.User.GetUID(),
			Manifest:  workflowManifest,
		},
	}

	if err := req.Create(&workflow); err != nil {
		return err
	}

	return req.WriteCreated(convertTask(workflow))
}

func toWorkflowManifest(agent *v1.Agent, thread *v1.Thread, manifest types.TaskManifest) types.WorkflowManifest {
	workflowManifest := types.WorkflowManifest{
		AgentManifest: agent.Spec.Manifest,
	}

	for _, tool := range thread.Spec.Manifest.Tools {
		if !slices.Contains(workflowManifest.Tools, tool) {
			workflowManifest.Tools = append(workflowManifest.Tools, tool)
		}
	}

	workflowManifest.Steps = toWorkflowSteps(manifest.Steps)
	workflowManifest.Name = manifest.Name
	workflowManifest.Description = manifest.Description
	return workflowManifest
}

func toWorkflowSteps(steps []types.TaskStep) []types.Step {
	workflowSteps := make([]types.Step, 0, len(steps))
	for _, step := range steps {
		workflowSteps = append(workflowSteps, types.Step{
			ID:   step.ID,
			Step: step.Step,
			If:   toWorkflowIf(step.If),
		})
	}
	return workflowSteps
}

func toWorkflowIf(ifStep *types.TaskIf) *types.If {
	if ifStep == nil {
		return nil
	}
	return &types.If{
		Condition: ifStep.Condition,
		Steps:     toWorkflowSteps(ifStep.Steps),
		Else:      toWorkflowSteps(ifStep.Else),
	}
}

func (t *TaskHandler) Get(req api.Context) error {
	task, err := t.getTask(req)
	if err != nil {
		return err
	}

	return req.Write(convertTask(*task))
}

func (t *TaskHandler) getTask(req api.Context) (*v1.Workflow, error) {
	assistantID := req.PathValue("assistant_id")

	var workflow v1.Workflow
	if err := req.Get(&workflow, req.PathValue("id")); err != nil {
		return nil, err
	}

	assistant, err := getAssistant(req, assistantID)
	if err != nil {
		return nil, err
	}

	if workflow.Spec.AgentName != assistant.Name || workflow.Spec.UserID != req.User.GetUID() {
		return nil, types.NewErrHttp(http.StatusForbidden, "task does not belong to the user")
	}

	return &workflow, nil
}

func (t *TaskHandler) List(req api.Context) error {
	assistant, err := getAssistant(req, req.PathValue("assistant_id"))
	if err != nil {
		return err
	}

	var workflows v1.WorkflowList
	if err := req.List(&workflows, kclient.MatchingFields{
		"spec.agentName": assistant.Name,
		"spec.userID":    req.User.GetUID(),
	}); err != nil {
		return err
	}

	taskList := types.TaskList{Items: make([]types.Task, 0, len(workflows.Items))}

	for _, workflow := range workflows.Items {
		taskList.Items = append(taskList.Items, convertTask(workflow))
	}

	return req.Write(taskList)
}

func convertTask(workflow v1.Workflow) types.Task {
	task := types.Task{
		Metadata: MetadataFrom(&workflow),
		TaskManifest: types.TaskManifest{
			Name:        workflow.Spec.Manifest.Name,
			Description: workflow.Spec.Manifest.Description,
		},
	}
	task.Steps = toTaskSteps(workflow.Spec.Manifest.Steps)
	return task
}

func toTaskSteps(steps []types.Step) []types.TaskStep {
	taskSteps := make([]types.TaskStep, 0, len(steps))
	for _, step := range steps {
		taskSteps = append(taskSteps, types.TaskStep{
			ID:   step.ID,
			Step: step.Step,
			If:   toIf(step.If),
		})
	}
	return taskSteps
}

func toIf(ifStep *types.If) *types.TaskIf {
	if ifStep == nil {
		return nil
	}
	return &types.TaskIf{
		Condition: ifStep.Condition,
		Steps:     toTaskSteps(ifStep.Steps),
		Else:      toTaskSteps(ifStep.Else),
	}
}
