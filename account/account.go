package account

import (
	"github.com/samber/lo"
	"go.chrastecky.dev/fio-api/fio"
	"go.chrastecky.dev/fio-client/fioclient/model"
)

type Account interface {
}

type account struct {
	api          fio.Client
	accountModel model.Account
}

func New(model model.Account) Account {
	return &account{
		api:          lo.Must(fio.NewClient(model.ApiKey)),
		accountModel: model,
	}
}
