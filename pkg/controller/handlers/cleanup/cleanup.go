package cleanup

import (
	"github.com/orion101-ai/nah/pkg/router"
	"github.com/orion101-ai/nah/pkg/uncached"
	"github.com/orion101-ai/orion101/logger"
	v1 "github.com/orion101-ai/orion101/pkg/storage/apis/orion101.orion101.ai/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var log = logger.Package()

func Cleanup(req router.Request, _ router.Response) error {
	toDelete := req.Object.(v1.DeleteRefs)

	for _, ref := range toDelete.DeleteRefs() {
		if ref.Name == "" {
			continue
		}

		namespace := req.Namespace
		if namespace == "" && ref.Namespace != "" {
			namespace = ref.Namespace
		}

		objType := ref.ObjType
		if ref.Kind != "" {
			o, err := req.Client.Scheme().New(schema.GroupVersionKind{
				Group:   objType.GetObjectKind().GroupVersionKind().Group,
				Version: objType.GetObjectKind().GroupVersionKind().Version,
				Kind:    ref.Kind,
			})
			if err != nil {
				return err
			}
			objType = o.(kclient.Object)
		}

		if err := req.Get(objType, namespace, ref.Name); apierrors.IsNotFound(err) {
			if err := req.Get(uncached.Get(objType), namespace, ref.Name); apierrors.IsNotFound(err) {
				log.Infof("Deleting %s/%s due to missing %s", namespace, req.Name, ref.Name)
				return req.Delete(req.Object)
			}
		}
	}

	return nil
}