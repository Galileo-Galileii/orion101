package create

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func OrGet(ctx context.Context, c kclient.Client, obj kclient.Object) error {
	err := c.Create(ctx, obj)
	if apierrors.IsAlreadyExists(err) {
		return c.Get(ctx, kclient.ObjectKeyFromObject(obj), obj)
	}
	return err
}

func IfNotExists(ctx context.Context, c kclient.Client, obj kclient.Object) error {
	err := c.Get(ctx, kclient.ObjectKeyFromObject(obj), obj)
	if apierrors.IsNotFound(err) {
		return OrGet(ctx, c, obj)
	}
	return err
}
