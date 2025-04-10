package service

import (
	"github.com/turnerbenjamin/go_odata/model"
	"github.com/turnerbenjamin/go_odata/msal"
)

type EntityList[T any] interface {
	GetData() []T
	HasNext() bool
	GetNext() (EntityList[T], error)
}

type entityList[T any] struct {
	res *model.GetManyResponse[T]
	svc msal.DataverseService
}

func (el *entityList[T]) GetData() []T {
	return el.res.Data
}

func (el *entityList[T]) HasNext() bool {
	return el.res.HasNext()
}

func (el *entityList[T]) GetNext() (EntityList[T], error) {
	r, err := el.svc.GetNext(el.res.Next)
	if err != nil {
		return nil, err
	}
	gmr := model.NewGetManyResponseFromJson[T](*r.Body)
	return &entityList[T]{
		res: gmr,
		svc: el.svc,
	}, nil

}

// type Service[T any] interface {
// 	Create(*T) (*T, error)
// 	List(string, int)
// 	RetrieveByName(name string)
// }
