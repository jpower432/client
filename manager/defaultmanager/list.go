package defaultmanager

import (
	"context"

	"github.com/uor-framework/uor-client-go/nodes/collection"
	"github.com/uor-framework/uor-client-go/registryclient"
)

func (d DefaultManager) List(ctx context.Context, reference string, remote registryclient.Remote) (*collection.Collection, error) {
	loadCollection, err := remote.LoadCollection(ctx, reference)
	if err != nil {
		return nil, err
	}
	return &loadCollection, nil
}
