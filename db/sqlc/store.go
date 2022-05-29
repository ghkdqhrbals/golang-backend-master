package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// 쿼리들을 저장
type SQLStore struct {
	*Queries //앞의 Queries 생략하면 *Queries에서 가져와 만들어짐.
	db       *sql.DB
}

// 만들어진 Store 주소반환
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// 쿼리문을 실행한다. Beigin, Rollback조건, commit로 최종실행 확인
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) // BEGIN
	if err != nil {
		return err
	}
	q := New(tx)    //
	err = fn(q)     // EXECUTE
	if err != nil { // ROLLBACK
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err:%v", err, rbErr)
		}
		return err
	}

	return tx.Commit() // COMMIT
}

type TransferTxParams struct {
	// backtick escape 문자 반영 X
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}
type TransferTxResult struct {
	Transfer    Transfers `json:"transfer"`
	FromAccount Accounts  `json:"from_account"`
	ToAccount   Accounts  `json:"to_account"`
	FromEntry   Entries   `json:from_entry`
	ToEntry     Entries   `json:to_entry`
}

var txKey = struct{}{}

// TransferTxParams[from,to,amount]를 파라미터로 받아서, Trasnfer을 생성하고, ID의 출금,입금 내역인 entry를 각각 생성
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// create transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// create entry of from account
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// create entry of to account
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// 낮은 ID부터 ADD하도록 순서정함.
		if arg.FromAccountID < arg.ToAccountID {
			// update from account balance
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})
			if err != nil {
				return err
			}

			// update to account balance
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			if err != nil {
				return err
			}
		} else {
			// update to account balance
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			if err != nil {
				return err
			}

			// update from account balance
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})
			if err != nil {
				return err
			}
		}

		// fmt.Println(txName, "get account 1")
		// account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)

		// fmt.Println(txName, "update account 1")
		// result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.FromAccountID,
		// 	Balance: account1.Balance - arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		// fmt.Println(txName, "get account 2")
		// account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		// if err != nil {
		// 	return err
		// }

		// fmt.Println(txName, "update account 2")
		// result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.ToAccountID,
		// 	Balance: account2.Balance + arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		return nil
	})
	return result, err
}
