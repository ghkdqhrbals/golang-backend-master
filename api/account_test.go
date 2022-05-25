package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/ghkdqhrbals/simplebank/db/mock"
	db "github.com/ghkdqhrbals/simplebank/db/sqlc"
	"github.com/ghkdqhrbals/simplebank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mockdb.NewMockStore(ctrl)
	// build stubs
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)). // 특정 ID를 argument로 받는 store.GetAccount함수가 call 되길 예상함.
		Times(1).                                        // 위의 함수가 call 되는 횟수 예상
		Return(account, nil)

	// start test server and send request
	fmt.Println("Before")
	server := NewServer(store)
	fmt.Println("After")

	// request 생성
	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// send API request & response in recorder
	server.router.ServeHTTP(recorder, request)

	// check its response
	require.Equal(t, http.StatusOK, recorder.Code)
	fmt.Println("Finished")
}

func randomAccount() db.Accounts {
	return db.Accounts{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
