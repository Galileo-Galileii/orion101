package render

import (
	"context"
	"fmt"
	"strings"

	"github.com/orion101-ai/nah/pkg/router"
	"github.com/orion101-ai/orion101/apiclient/types"
	v1 "github.com/orion101-ai/orion101/pkg/storage/apis/orion101.orion101.ai/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type WorkflowOptions struct {
	Step             *types.Step
	ManifestOverride *types.WorkflowManifest
	Input            string
}

func IsExternalTool(tool string) bool {
	return strings.ContainsAny(tool, ".\\/")
}

func ResolveToolReference(ctx context.Context, c kclient.Client, toolRefType types.ToolReferenceType, ns, name string) (string, error) {
	if IsExternalTool(name) {
		return name, nil
	}

	var tool v1.ToolReference
	if err := c.Get(ctx, router.Key(ns, name), &tool); apierror.IsNotFound(err) {
		return name, nil
	} else if err != nil {
		return "", err
	}
	if toolRefType != "" && tool.Spec.Type != toolRefType {
		return name, fmt.Errorf("tool reference %s is not of type %s", name, toolRefType)
	}
	if tool.Status.Reference == "" {
		return "", fmt.Errorf("tool reference %s has no reference", name)
	}
	if toolRefType == types.ToolReferenceTypeTool {
		return fmt.Sprintf("%s as %s", tool.Status.Reference, name), nil
	}
	return tool.Status.Reference, nil
}

func Workflow(ctx context.Context, c kclient.Client, wf *v1.Workflow, opts WorkflowOptions) (*v1.Agent, error) {
	agentManifest := wf.Spec.Manifest.AgentManifest
	if opts.ManifestOverride != nil {
		agentManifest = opts.ManifestOverride.AgentManifest
	}

	agent := v1.Agent{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: wf.Namespace,
		},
		Spec: v1.AgentSpec{
			Manifest:            agentManifest,
			Credentials:         wf.Spec.Manifest.Credentials,
			CredentialContextID: wf.Name,
		},
		Status: v1.AgentStatus{
			WorkspaceName:     wf.Status.WorkspaceName,
			KnowledgeSetNames: wf.Status.KnowledgeSetNames,
		},
	}

	for _, env := range wf.Spec.Manifest.Env {
		if env.Name == "" {
			continue
		}
		if env.Value == "" {
			agent.Spec.Credentials = append(agent.Spec.Credentials,
				fmt.Sprintf(`github.com/gptscript-ai/credential as %s with "%s" as message and "%s" as env and %s as field`,
					env.Name, env.Description, env.Name, env.Name))
		} else {
			agent.Spec.Env = append(agent.Spec.Env, fmt.Sprintf("%s=%s", env.Name, env.Value))
		}
	}

	if step := opts.Step; step != nil {
		if step.Cache != nil {
			agent.Spec.Manifest.Cache = step.Cache
		}
		if step.Temperature != nil {
			agent.Spec.Manifest.Temperature = step.Temperature
		}

		agent.Spec.Manifest.Tools = append(agent.Spec.Manifest.Tools, step.Tools...)
		agent.Spec.Manifest.Agents = append(agent.Spec.Manifest.Agents, step.Agents...)
		agent.Spec.Manifest.Workflows = append(agent.Spec.Manifest.Workflows, step.Workflows...)
		if step.Template != nil && step.Template.Name != "" {
			name, err := ResolveToolReference(ctx, c, types.ToolReferenceTypeStepTemplate, wf.Namespace, step.Template.Name)
			if err != nil {
				return nil, err
			}
			agent.Spec.InputFilters = append(agent.Spec.InputFilters, name)
		}
	}

	if agent.Spec.Manifest.Prompt == "" {
		agent.Spec.Manifest.Prompt = v1.DefaultWorkflowAgentPrompt
	}

	if opts.Input != "" {
		agent.Spec.Manifest.Prompt = fmt.Sprintf("WORKFLOW INPUT: %s\nEND WORKFLOW INPUT\n\n%s", opts.Input, agent.Spec.Manifest.Prompt)
	}

	return &agent, nil
}