package defaultmanager

import (
	"context"

	"github.com/uor-framework/uor-client-go/nodes/collection"
	v2 "github.com/uor-framework/uor-client-go/nodes/descriptor/v2"
	"github.com/uor-framework/uor-client-go/registryclient"
)

func (d DefaultManager) List(ctx context.Context, reference string, remote registryclient.Remote) (*collection.Collection, string, error) {
	loadCollection, err := remote.LoadCollection(ctx, reference)
	if err != nil {
		return nil, "", err
	}

	root, err := loadCollection.Root()
	if err != nil {
		return nil, "", err
	}
	var digest string
	rootDesc, ok := root.(*v2.Node)
	if ok {
		digest = rootDesc.Descriptor().Digest.String()
	}

	return &loadCollection, digest, nil
}
