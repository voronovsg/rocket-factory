package part

import (
	"context"
	"strings"

	"github.com/voronovsg/rocket-factory/inventory/internal/model"
	repoConv "github.com/voronovsg/rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/voronovsg/rocket-factory/inventory/internal/repository/model"
)

func (r *repository) List(_ context.Context, filter model.PartsFilter) ([]model.Part, error) {
	result := make([]repoModel.Part, 0, len(r.parts))
	// Если фильтр пустой, возвращаем все элементы
	if isEmptyFilter(filter) {
		for _, part := range r.parts {
			result = append(result, part)
		}
		return repoConv.ListPartsToModel(result), nil
	}
	// Maps для быстрого поиска
	sets := buildPartsFilterSets(filter)
	// Фильтрация
	for _, part := range r.parts {
		if !matchesPartsFilter(part, sets) {
			continue
		}
		result = append(result, part)
	}

	return repoConv.ListPartsToModel(result), nil
}

type partsFilterSets struct {
	uuidSet     map[string]struct{}
	nameSet     map[string]struct{}
	categorySet map[int32]struct{}
	countrySet  map[string]struct{}
	tagSet      map[string]struct{}
}

func buildPartsFilterSets(filter model.PartsFilter) partsFilterSets {
	sets := partsFilterSets{
		uuidSet:     make(map[string]struct{}, len(filter.Uuids)),
		nameSet:     make(map[string]struct{}, len(filter.Names)),
		categorySet: make(map[int32]struct{}, len(filter.Categories)),
		countrySet:  make(map[string]struct{}, len(filter.ManufacturerCountries)),
		tagSet:      make(map[string]struct{}, len(filter.Tags)),
	}

	for _, partUuid := range filter.Uuids {
		sets.uuidSet[partUuid] = struct{}{}
	}
	for _, name := range filter.Names {
		sets.nameSet[strings.ToLower(name)] = struct{}{}
	}
	for _, category := range filter.Categories {
		sets.categorySet[category] = struct{}{}
	}
	for _, country := range filter.ManufacturerCountries {
		sets.countrySet[strings.ToLower(country)] = struct{}{}
	}
	for _, tag := range filter.Tags {
		sets.tagSet[strings.ToLower(tag)] = struct{}{}
	}

	return sets
}

func matchesPartsFilter(part repoModel.Part, sets partsFilterSets) bool {
	if len(sets.uuidSet) > 0 {
		if _, ok := sets.uuidSet[part.Uuid]; !ok {
			return false
		}
	}
	if len(sets.nameSet) > 0 {
		if _, ok := sets.nameSet[strings.ToLower(part.Name)]; !ok {
			return false
		}
	}
	if len(sets.categorySet) > 0 {
		if _, ok := sets.categorySet[part.Category]; !ok {
			return false
		}
	}
	if len(sets.countrySet) > 0 {
		if _, ok := sets.countrySet[strings.ToLower(part.Manufacturer.Country)]; !ok {
			return false
		}
	}
	if len(sets.tagSet) > 0 {
		for _, tag := range part.Tags {
			if _, ok := sets.tagSet[strings.ToLower(tag)]; ok {
				return true
			}
		}
		return false
	}

	return true
}

func isEmptyFilter(f model.PartsFilter) bool {
	return len(f.Uuids) == 0 &&
		len(f.Names) == 0 &&
		len(f.Categories) == 0 &&
		len(f.ManufacturerCountries) == 0 &&
		len(f.Tags) == 0
}
