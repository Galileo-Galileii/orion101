package controller

import (
	"github.com/orion101-ai/nah/pkg/handlers"
	"github.com/orion101-ai/orion101/pkg/controller/handlers/agents"
	"github.com/orion101-ai/orion101/pkg/controller/handlers/alias"
	"github.com/orion101-ai/orion101/pkg/controller/handlers/cleanup"
	"github.com/orion101-ai/orion101/pkg/controller/handlers/cronjob"
	"github.com/orion101-ai/orion101/pkg/controller/handlers/knowledgefile"
	"github.com/orion101-ai/orion101/pkg/controller/handlers/knowledgeset"
	"github.com/orion101-ai/orion101/pkg/controller/handlers/knowledgesource"
	"github.com/orion101-ai/orion101/pkg/controller/handlers/oauthapp"
	"github.com/orion101-ai/orion101/pkg/controller/handlers/runs"
	"github.com/orion101-ai/orion101/pkg/controller/handlers/threads"
	"github.com/orion101-ai/orion101/pkg/controller/handlers/toolreference"
	"github.com/orion101-ai/orion101/pkg/controller/handlers/webhook"
	"github.com/orion101-ai/orion101/pkg/controller/handlers/workflow"
	"github.com/orion101-ai/orion101/pkg/controller/handlers/workflowexecution"
	"github.com/orion101-ai/orion101/pkg/controller/handlers/workflowstep"
	"github.com/orion101-ai/orion101/pkg/controller/handlers/workspace"
	v1 "github.com/orion101-ai/orion101/pkg/storage/apis/orion101.orion101.ai/v1"
)

func (c *Controller) setupRoutes() error {
	root := c.router

	workflowExecution := workflowexecution.New(c.services.Invoker)
	workflowStep := workflowstep.New(c.services.Invoker)
	toolRef := toolreference.New(c.services.GPTClient, c.services.ToolRegistryURL)
	workspace := workspace.New(c.services.GPTClient, c.services.WorkspaceProviderType)
	knowledgeset := knowledgeset.New(c.services.AIHelper, c.services.Invoker)
	knowledgesource := knowledgesource.NewHandler(c.services.Invoker, c.services.GPTClient)
	knowledgefile := knowledgefile.New(c.services.Invoker, c.services.GPTClient, c.services.KnowledgeSetIngestionLimit)
	runs := runs.New(c.services.Invoker)
	webHooks := webhook.New()
	cronJobs := cronjob.New()
	oauthLogins := oauthapp.NewLogin(c.services.Invoker, c.services.ServerURL)

	// Runs
	root.Type(&v1.Run{}).HandlerFunc(runs.MigrateRemoveRunFinalizer) // to be removed
	root.Type(&v1.Run{}).HandlerFunc(cleanup.Cleanup)
	root.Type(&v1.Run{}).HandlerFunc(runs.Resume)

	// RunStates
	root.Type(&v1.RunState{}).HandlerFunc(cleanup.Cleanup)

	// Threads
	root.Type(&v1.Thread{}).HandlerFunc(cleanup.Cleanup)
	root.Type(&v1.Thread{}).HandlerFunc(threads.CreateWorkspaces)
	root.Type(&v1.Thread{}).HandlerFunc(threads.CreateKnowledgeSet)
	root.Type(&v1.Thread{}).HandlerFunc(threads.WorkflowState)

	// Workflows
	root.Type(&v1.Workflow{}).HandlerFunc(workflow.WorkspaceObjects)
	root.Type(&v1.Workflow{}).HandlerFunc(workflow.EnsureIDs)
	root.Type(&v1.Workflow{}).HandlerFunc(workflow.CreateWorkspaceAndKnowledgeSet)
	root.Type(&v1.Workflow{}).HandlerFunc(workflow.BackPopulateAuthStatus)

	// WorkflowExecutions
	root.Type(&v1.WorkflowExecution{}).HandlerFunc(cleanup.Cleanup)
	root.Type(&v1.WorkflowExecution{}).HandlerFunc(workflowExecution.Run)

	// Agents
	root.Type(&v1.Agent{}).HandlerFunc(agents.CreateWorkspaceAndKnowledgeSet)
	root.Type(&v1.Agent{}).HandlerFunc(agents.BackPopulateAuthStatus)

	// Uploads
	root.Type(&v1.KnowledgeSource{}).HandlerFunc(cleanup.Cleanup)
	root.Type(&v1.KnowledgeSource{}).FinalizeFunc(v1.KnowledgeSourceFinalizer, knowledgesource.Cleanup)
	root.Type(&v1.KnowledgeSource{}).HandlerFunc(knowledgesource.Reschedule)
	root.Type(&v1.KnowledgeSource{}).HandlerFunc(knowledgesource.Sync)

	// ToolReference
	root.Type(&v1.ToolReference{}).HandlerFunc(toolRef.Populate)

	// Reference
	root.Type(&v1.Agent{}).HandlerFunc(alias.AssignAlias)
	root.Type(&v1.Workflow{}).HandlerFunc(alias.AssignAlias)
	root.Type(&v1.Model{}).HandlerFunc(alias.AssignAlias)
	root.Type(&v1.DefaultModelAlias{}).HandlerFunc(alias.AssignAlias)

	// Knowledge files
	root.Type(&v1.KnowledgeFile{}).HandlerFunc(cleanup.Cleanup)
	root.Type(&v1.KnowledgeFile{}).FinalizeFunc(v1.KnowledgeFileFinalizer, knowledgefile.Cleanup)
	root.Type(&v1.KnowledgeFile{}).HandlerFunc(knowledgefile.IngestFile)
	root.Type(&v1.KnowledgeFile{}).HandlerFunc(knowledgefile.Unapproved)

	// Workspaces
	root.Type(&v1.Workspace{}).HandlerFunc(cleanup.Cleanup)
	root.Type(&v1.Workspace{}).FinalizeFunc(v1.WorkspaceFinalizer, workspace.RemoveWorkspace)
	root.Type(&v1.Workspace{}).HandlerFunc(workspace.CreateWorkspace)

	// KnowledgeSets
	root.Type(&v1.KnowledgeSet{}).HandlerFunc(cleanup.Cleanup)
	root.Type(&v1.KnowledgeSet{}).FinalizeFunc(v1.KnowledgeSetFinalizer, knowledgeset.Cleanup)
	// Also cleanup the dataset when there is no content.
	// This will allow the user to switch the embedding model implicitly.
	root.Type(&v1.KnowledgeSet{}).HandlerFunc(knowledgeset.Cleanup)
	root.Type(&v1.KnowledgeSet{}).HandlerFunc(knowledgeset.GenerateDataDescription)
	root.Type(&v1.KnowledgeSet{}).HandlerFunc(knowledgeset.CreateWorkspace)
	root.Type(&v1.KnowledgeSet{}).HandlerFunc(knowledgeset.CheckHasContent)
	root.Type(&v1.KnowledgeSet{}).HandlerFunc(knowledgeset.SetEmbeddingModel)

	// Webhooks
	root.Type(&v1.Webhook{}).HandlerFunc(cleanup.Cleanup)
	root.Type(&v1.Webhook{}).HandlerFunc(alias.AssignAlias)
	root.Type(&v1.Webhook{}).HandlerFunc(webHooks.SetSuccessRunTime)

	// Cronjobs
	root.Type(&v1.CronJob{}).HandlerFunc(cleanup.Cleanup)
	root.Type(&v1.CronJob{}).HandlerFunc(cronJobs.SetSuccessRunTime)
	root.Type(&v1.CronJob{}).HandlerFunc(cronJobs.Run)

	// OAuthApps
	root.Type(&v1.OAuthApp{}).HandlerFunc(cleanup.Cleanup)
	root.Type(&v1.OAuthApp{}).HandlerFunc(alias.AssignAlias)

	// OAuthAppLogins
	root.Type(&v1.OAuthAppLogin{}).HandlerFunc(cleanup.Cleanup)
	root.Type(&v1.OAuthAppLogin{}).HandlerFunc(oauthLogins.RunTool)

	// Alias
	root.Type(&v1.Alias{}).HandlerFunc(alias.UnassignAlias)

	// WorkflowSteps
	steps := root.Type(&v1.WorkflowStep{})
	steps.HandlerFunc(cleanup.Cleanup)
	steps.HandlerFunc(handlers.GCOrphans)

	running := steps.Middleware(workflowStep.Preconditions)
	running.HandlerFunc(workflowStep.RunInvoke)
	running.HandlerFunc(workflowStep.RunIf)
	running.HandlerFunc(workflowStep.RunWhile)
	steps.HandlerFunc(workflowStep.RunSubflow)

	c.toolRefHandler = toolRef
	return nil
}