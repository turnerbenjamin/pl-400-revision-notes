package model

import "github.com/turnerbenjamin/go_odata/view"

type EntityListService interface {
	GetNext()
}

type entityList[T view.Entity] struct {
	fetchNext func(string) (*GetManyResponse[T], error)
	data      []T
	next      string
	previous  view.EntityList[T]
}

func CreateEntityList[T view.Entity](getManyResponse GetManyResponse[T], fetchNext func(string) (*GetManyResponse[T], error)) view.EntityList[T] {
	return &entityList[T]{
		fetchNext: fetchNext,
		data:      getManyResponse.Data,
		next:      getManyResponse.Next,
		previous:  nil,
	}
}

func (el *entityList[T]) Data() []T {
	return el.data
}

func (el *entityList[T]) HasNext() bool {
	return el.next != ""
}

func (el *entityList[T]) Next() (view.EntityList[T], error) {
	r, err := el.fetchNext(el.next)
	if err != nil {
		return nil, err
	}
	// gmr := NewGetManyResponseFromJson[T](*r.Body)
	return &entityList[T]{
		fetchNext: el.fetchNext,
		data:      r.Data,
		next:      r.Next,
		previous:  el,
	}, nil
}

func (el *entityList[T]) HasPrevious() bool {
	return el.previous != nil
}

func (el *entityList[T]) Previous() view.EntityList[T] {
	return el.previous
}
