package generic

import (
	"github.com/orion101-ai/kinm/pkg/db"
	"github.com/orion101-ai/kinm/pkg/stores"
	"github.com/orion101-ai/kinm/pkg/strategy"
	"github.com/orion101-ai/orion101/pkg/storage/scheme"
	"github.com/orion101-ai/orion101/pkg/storage/tables"
	"k8s.io/apiserver/pkg/registry/rest"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type withTable struct {
	strategy.CompleteStrategy
	*strategy.TableAdapter
}

type hasColumns interface {
	GetColumns() [][]string
}

func NewStore(db *db.Factory, obj kclient.Object) (rest.Storage, rest.Storage, error) {
	storage, err := db.NewDBStrategy(obj)
	if err != nil {
		return nil, nil, err
	}

	var tableStrategy any
	if obj, ok := obj.(hasColumns); ok {
		tableStrategy, err = tables.NewConverter(obj.GetColumns())
		if err != nil {
			return nil, nil, err
		}
	}

	newStorage := withTable{
		CompleteStrategy: storage,
		TableAdapter:     strategy.NewTable(tableStrategy),
	}

	return stores.NewComplete(scheme.Scheme, newStorage), stores.NewStatus(scheme.Scheme, storage), err
}
