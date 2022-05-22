package db

import (
	"context"
	"testing"
	"time"

	"simplebank/util"

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
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
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

	arg := ListAccountParams{
		Limit:  10,
		Offset: 1,
	}
	for i := int64(0); i < int64(arg.Limit); i++ {
		createRandomAccounts(t)
	}

	accounts, err := testQueries.ListAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 10)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}
