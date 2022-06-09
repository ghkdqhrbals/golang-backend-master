package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/ghkdqhrbals/simplebank/db/sqlc"
	"github.com/ghkdqhrbals/simplebank/token"
	"github.com/gin-gonic/gin"
)

// oneof = item1, itme2 -> value should be one of next items
type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`       // gt = N -> the number is greater than N
	Currency      string `json:"currency" binding:"required,currency"` // doesn't require any oneof methods, because we already have currency validator(reflection)
}

func (server *Server) createTrasnfer(ctx *gin.Context) {
	var req transferRequest
	// http.requests into req by unmarshelling json struct information
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	// check each account currency whether they have the same currency and also same with Req.Currency
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account dosen't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	// Create transferParams
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	results, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, results)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Accounts, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
