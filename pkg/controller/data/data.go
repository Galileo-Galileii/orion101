package data

import (
	"context"
	_ "embed"

	"github.com/orion101-ai/nah/pkg/apply"
	v1 "github.com/orion101-ai/orion101/pkg/storage/apis/orion101.orion101.ai/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

//go:embed orion101.yaml
var orion101Data []byte

//go:embed default-models.yaml
var defaultModelsData []byte

//go:embed default-model-aliases.yaml
var defaultModelAliasesData []byte

func Data(ctx context.Context, c kclient.Client) error {
	var defaultModels v1.ModelList
	if err := yaml.Unmarshal(defaultModelsData, &defaultModels); err != nil {
		return err
	}

	for _, model := range defaultModels.Items {
		var existing v1.Model
		if err := c.Get(ctx, kclient.ObjectKey{Namespace: model.Namespace, Name: model.Name}, &existing); err == nil {
			// If the usage is different, update the existing model.
			if model.Spec.Manifest.Usage != existing.Spec.Manifest.Usage {
				existing.Spec.Manifest.Usage = model.Spec.Manifest.Usage
				if err := c.Update(ctx, &existing); err != nil {
					return err
				}
			}
		} else if !apierrors.IsNotFound(err) {
			return err
		} else {
			if err = kclient.IgnoreAlreadyExists(c.Create(ctx, &model)); err != nil {
				return err
			}
		}
	}

	var defaultModelAliases v1.DefaultModelAliasList
	if err := yaml.Unmarshal(defaultModelAliasesData, &defaultModelAliases); err != nil {
		return err
	}

	for _, alias := range defaultModelAliases.Items {
		var existing v1.DefaultModelAlias
		if err := c.Get(ctx, kclient.ObjectKey{Namespace: alias.Namespace, Name: alias.Name}, &existing); apierrors.IsNotFound(err) {
			if err := c.Create(ctx, &alias); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	var orion101 v1.Agent
	if err := yaml.Unmarshal(orion101Data, &orion101); err != nil {
		return err
	}

	return apply.Ensure(ctx, c, &orion101)
}