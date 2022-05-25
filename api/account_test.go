package api

import (
	"testing"

	mockdb "github.com/ghkdqhrbals/simplebank/db/mock"
	db "github.com/ghkdqhrbals/simplebank/db/sqlc"
	"github.com/ghkdqhrbals/simplebank/util"
	"github.com/golang/mock/gomock"
)

func TestGetACcountAPI(t *testing.T) {
	account := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)
}

func randomAccount() db.Accounts {
	return db.Accounts{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
