package part

import (
	"context"

	"github.com/voronovsg/rocket-factory/inventory/internal/model"
	repoConv "github.com/voronovsg/rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, partUuid string) (model.Part, error) {
	part, ok := r.parts[partUuid]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}

	return repoConv.PartToModel(part), nil
}
