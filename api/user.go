package api

import (
	"net/http"
	"time"

	db "github.com/ghkdqhrbals/simplebank/db/sqlc"
	"github.com/ghkdqhrbals/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"` // Do not expose this Hashed Password
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// User은 username, email이 unique key로 설정되어있음.
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 임시로 주석처리(user_test.go 에서 Testing하기 위해)
	// 실제로는 주석해제해야함
	// response := createUserResponse{
	// 	Username:          user.Username,
	// 	FullName:          user.FullName,
	// 	Email:             user.Email,
	// 	PasswordChangedAt: user.PasswordChangeAt,
	// 	CreatedAt:         user.CreatedAt,
	// }

	ctx.JSON(http.StatusOK, user)
}
