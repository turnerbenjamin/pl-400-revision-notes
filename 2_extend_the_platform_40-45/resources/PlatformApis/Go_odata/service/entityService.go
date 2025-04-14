package service

import "github.com/turnerbenjamin/go_odata/view"

type EntityService[T view.Entity] interface {
	List(searchTerm string) (view.EntityList[T], error)
	Get(guid string) (T, error)
	Create(entityToCreate T) (newEntity T, err error)
	Update(guid string, entityToUpdate T) error
	Delete(guid string) error
}
