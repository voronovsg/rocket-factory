package part

import (
	def "github.com/voronovsg/rocket-factory/inventory/internal/repository"
	repoModel "github.com/voronovsg/rocket-factory/inventory/internal/repository/model"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	parts map[string]repoModel.Part
}

func NewRepository() *repository {
	return &repository{
		parts: make(map[string]repoModel.Part),
	}
}
