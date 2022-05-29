package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					// mock으로 테스트할 function 이름.
					GetAccount(gomock.Any(), account.ID). // 특정 ID를 argument로 받는 store.GetAccount함수가 call 되길 예상함.
					Times(1).                             // 위의 함수가 call 되는 횟수 예상
					Return(account, nil)                  // test function의 return 값 예상.
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code) // 200 http Response 반환
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: int64(10002),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), int64(10002)). // 특정 ID를 argument로 받는 store.GetAccount함수가 call 되길 예상함.
					Times(1).                               // 위의 함수가 call 되는 횟수 예상
					Return(db.Accounts{}, sql.ErrNoRows)    // 빈 Accounts구조체 및 ErrNoRows 에러 반환 예상.
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code) // 404 ERROR http Response반환
			},
		},
		{
			name:      "BadReqeust",
			accountID: int64(-1),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), int64(-1)).
					// 애초에 id min=1로 설정해두었기에, unmarshelling할 때 오류가 뜸으로 GetAccount가 실행이 안됨.
					//따라서 0으로 설정
					Times(0)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code) // 404 ERROR http Response반환
			},
		},
		{
			name:      "InternalServerError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), account.ID).  // 특정 ID를 argument로 받는 store.GetAccount함수가 call 되길 예상함.
					Times(1).                              // 위의 함수가 call 되는 횟수 예상
					Return(db.Accounts{}, sql.ErrConnDone) // 빈 Accounts구조체 및 ErrNoRows 에러 반환 예상.
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code) // 404 ERROR http Response반환
			},
		},
		// TODO: testCase 추가
	}

	// Stubs 케이스 실행
	for i := range testCases {

		tc := testCases[i]
		fmt.Println("Total Case: ", len(testCases), "\nCase : ", tc.name, "")
		cc := tc.accountID
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// start test server and send request
			server := NewServer(store)

			// request 생성
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%d", cc)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// send API request & response in recorder
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomAccount() db.Accounts {
	return db.Accounts{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Accounts) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Accounts
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
