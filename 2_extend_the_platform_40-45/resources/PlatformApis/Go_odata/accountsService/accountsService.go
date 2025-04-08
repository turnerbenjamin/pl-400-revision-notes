package accountsService

import (
	"github.com/turnerbenjamin/go_odata/model"
)

type AccountsService interface {
	Create(*model.Account)
	RetrieveByName(name string)
}
