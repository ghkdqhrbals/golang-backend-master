package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/ghkdqhrbals/golang-backend-master/db/sqlc"
	"github.com/ghkdqhrbals/golang-backend-master/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"` // Do not expose this Hashed Password
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

// remove hashPassword
func newUserResponse(user db.Users) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangeAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) testCreateUser(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, user)
}

// Synchronous api handler
func (server *Server) createUser_synchronous(ctx *gin.Context) {
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
	response := newUserResponse(user)
	ctx.JSON(http.StatusOK, response)
}

// Asynchronous api handler
func (server *Server) createUser_asynchronous(ctx *gin.Context) {
	counter := server.avaliable
	server.avaliable++
	result := make(chan int)
	errb := make(chan error)
	responseb := make(chan userResponse)
	go func(ctx *gin.Context, counter int) {
		counter++
		logrus.WithFields(logrus.Fields{
			"thread_number": counter,
		}).Info("Start create user api")

		var req createUserRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			result <- http.StatusBadRequest
			errb <- err
			logrus.WithFields(logrus.Fields{
				"Status": "bad request",
				"thread": counter,
			}).Warn("Finished with error")
			return
		}

		hashedPassword, err := util.HashPassword(req.Password)
		if err != nil {
			result <- http.StatusInternalServerError
			errb <- err
			logrus.WithFields(logrus.Fields{
				"Status":        "StatusInternalServerError",
				"thread_number": counter,
			}).Warn("Finished with error")
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
					result <- http.StatusForbidden
					errb <- err
					logrus.WithFields(logrus.Fields{
						"Status":        "unique_violation",
						"thread_number": counter,
					}).Warn("Finished with error")
					return
				}
			}
			result <- http.StatusInternalServerError
			errb <- err
			logrus.WithFields(logrus.Fields{
				"Status":        "StatusInternalServerError",
				"thread_number": counter,
			}).Warn("Finish")
			return
		}
		response := newUserResponse(user)
		result <- http.StatusOK
		responseb <- response
		logrus.WithFields(logrus.Fields{
			"Status":        "StatusOK",
			"username":      req.Username,
			"email":         req.Email,
			"thread_number": counter,
		}).Info("Successfully created")

	}(ctx.Copy(), counter)
	msg := <-result
	if msg == http.StatusOK {
		ctx.JSON(msg, <-responseb)
	} else {
		ctx.JSON(msg, errorResponse(<-errb))
	}
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"` // Do not expose this Hashed Password
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// get user's hasedpassword
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// check password
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	// create token
	accessToken, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}
