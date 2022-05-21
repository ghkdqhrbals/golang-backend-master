package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTransferTest(t *testing.T, fromID int64, toID int64, amount int64) Transfers {
	//
	arg := CreateTransferParams{
		FromAccountID: fromID,
		ToAccountID:   toID,
		Amount:        amount,
	}
	tx3, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tx3)

	fromAccount, _ := testQueries.GetAccount(context.Background(), arg.FromAccountID)
	toAccount, _ := testQueries.GetAccount(context.Background(), arg.ToAccountID)

	// 송신자 빼기
	arg1 := UpdateAccountParams{
		ID:      arg.FromAccountID,
		Balance: fromAccount.Balance - arg.Amount,
	}
	tx, err := testQueries.UpdateAccount(context.Background(), arg1)
	require.NoError(t, err)
	require.NotEmpty(t, tx)

	// 수신자 더하기
	arg2 := UpdateAccountParams{
		ID:      arg.ToAccountID,
		Balance: toAccount.Balance + arg.Amount,
	}

	tx2, err := testQueries.UpdateAccount(context.Background(), arg2)
	require.NoError(t, err)
	require.NotEmpty(t, tx2)

	fromAccount, _ = testQueries.GetAccount(context.Background(), arg.FromAccountID)
	toAccount, _ = testQueries.GetAccount(context.Background(), arg.ToAccountID)

	return tx3
}

func TestCreateTransfer(t *testing.T) {
	createTransferTest(t, 7, 8, 10)
}

func TestDeleteTransfer(t *testing.T) {
	err := testQueries.DeleteTransfer(context.Background(), 1)

	require.NoError(t, err)
}

func TestListTransfer(t *testing.T) {
	var listOfTransfer []Transfers

	arg := ListTransferParams{
		Limit:  10,
		Offset: 1,
	}

	listOfTransfer, err := testQueries.ListTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, listOfTransfer, 10)

	for _, transfer := range listOfTransfer {
		require.NotEmpty(t, transfer)
	}

}
