package db

import (
	"context"
	"testing"
	"time"

	"github.com/ghkdqhrbals/simplebank/util"

	"github.com/stretchr/testify/require"
)

// func TestCreateAccounts(t *testing.T) {
// 	arg := CreateAccountParams{
// 		Owner:    util.RandomOwner(),
// 		Balance:  util.RandomMoney(),
// 		Currency: util.RandomCurrency(),
// 	}

// 	account, err := testQueries.CreateAccount(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, account)

// 	require.Equal(t, arg.Owner, account.Owner)
// 	require.Equal(t, arg.Balance, account.Balance)
// 	require.Equal(t, arg.Currency, account.Currency)

// 	require.NotEmpty(t, account.ID)
// 	require.NotEmpty(t, account.CreatedAt)
// }
func createRandomAccounts(t *testing.T) Accounts {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotEmpty(t, account.ID)
	require.NotEmpty(t, account.CreatedAt)
	return account
}
func TestCreateCcount(t *testing.T) {
	createRandomAccounts(t)
}
func TestGetAccounts(t *testing.T) {
	account1 := createRandomAccounts(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	t.Log(account2.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, 1*time.Second)
}

func TestUpdateAccounts(t *testing.T) {
	account1 := createRandomAccounts(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: 100,
	}
	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, 1*time.Second)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Accounts

	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccounts(t)
	}

	arg := ListAccountParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}
	accounts, err := testQueries.ListAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner) // List account의 Owner, lastaccount owner 확인
	}

}
