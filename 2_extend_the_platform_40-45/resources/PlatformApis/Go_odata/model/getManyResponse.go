package model

import (
	"encoding/json"
)

type GetManyResponse[T any] struct {
	Next string `json:"@odata.nextLink"`
	Data []T    `json:"value"`
}

func NewGetManyResponseFromJson[T any](responseJson []byte) *GetManyResponse[T] {
	res := GetManyResponse[T]{}
	json.Unmarshal(responseJson, &res)
	return &res
}
