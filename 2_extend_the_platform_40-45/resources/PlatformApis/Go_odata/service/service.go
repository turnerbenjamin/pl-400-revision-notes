package service

// import (
// 	"github.com/turnerbenjamin/go_odata/model"
// 	"github.com/turnerbenjamin/go_odata/msal"
// )

// type EntityList[T any] interface {
// 	GetData() []T
// 	HasNext() bool
// 	GetNext() (EntityList[T], error)
// 	HasPrevious() bool
// 	GetPrevious() EntityList[T]
// }

// type entityList[T any] struct {
// 	svc      msal.DataverseService
// 	data     []T
// 	next     string
// 	previous EntityList[T]
// }

// func (el *entityList[T]) GetData() []T {
// 	return el.data
// }

// func (el *entityList[T]) HasNext() bool {
// 	return el.next != ""
// }

// func (el *entityList[T]) GetNext() (EntityList[T], error) {
// 	r, err := el.svc.GetNext(el.next)
// 	if err != nil {
// 		return nil, err
// 	}
// 	gmr := model.NewGetManyResponseFromJson[T](*r.Body)
// 	return &entityList[T]{
// 		svc:      el.svc,
// 		data:     gmr.Data,
// 		next:     gmr.Next,
// 		previous: el,
// 	}, nil
// }

// func (el *entityList[T]) HasPrevious() bool {
// 	return el.previous != nil
// }

// func (el *entityList[T]) GetPrevious() EntityList[T] {
// 	return el.previous
// }
