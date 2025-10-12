package db

import (
	"context"
	"database/sql"
	"fmt"
)

// This store will provide all functions to run database functions incdividually as well as their combination within a transaction

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Here we are calling the New function with the sql.Tx (transaction) instead of sql.DB so the returned Queries object can be used for calling input callback function with the queries returned by transaction as argument for it
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err %v, rb err %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// Input paramaters for the transfer transaction
type TransferTxParams struct {
	FromAccountId int64
	ToAccountId   int64
	Amount        int64
}

// Result of the Transfer transaction
type TransferTxResult struct {
	Transfer    Transfer
	FromAccount Account
	ToAccount   Account
	FromEntry   Entry
	ToEntry     Entry
}

// create transfer, create entries, update accounts
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// NEW CODE
		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = AddMoney(ctx, q, arg.FromAccountId, -arg.Amount, arg.ToAccountId, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = AddMoney(ctx, q, arg.ToAccountId, arg.Amount, arg.FromAccountId, -arg.Amount)
			if err != nil {
				return err
			}
		}

		// OLD CODE
		// acc1, err := q.GetAccountForUpdadte(ctx, arg.FromAccountId)
		// if err != nil {
		// 	return err
		// }
		// result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.FromAccountId,
		// 	Balance: acc1.Balance - arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }
		// acc2, err := q.GetAccountForUpdadte(ctx, arg.ToAccountId)
		// if err != nil {
		// 	return err
		// }
		// result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.ToAccountId,
		// 	Balance: acc2.Balance + arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		return nil
	})

	return result, err
}

func AddMoney(ctx context.Context, q *Queries, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID:     accountID1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount2,
		ID:     accountID2,
	})
	if err != nil {
		return
	}
	return
}
