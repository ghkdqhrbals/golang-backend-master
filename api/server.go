package api

import (
	"encoding/json"
	"fmt"

	db "github.com/ghkdqhrbals/golang-backend-master/db/sqlc"
	"github.com/ghkdqhrbals/golang-backend-master/token"
	"github.com/ghkdqhrbals/golang-backend-master/util"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	avaliable  int
}

// Configuration setting 받아와서 서버 open
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		avaliable:  0,
	}

	// Check supportable currency every http request
	// reflection 이라 동적으로 관리
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency) // Register Custom Validator {tag, validator.Func}
	}

	server.setupRouter()

	// ------------------for Testing Purpose
	// recorder := httptest.NewRecorder()
	// url := "/users"
	// for i := 0; i < 10; i++ {
	// 	go func() {
	// 		request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(dataMarchal()))
	// 		server.router.ServeHTTP(recorder, request)
	// 	}()
	// }
	// ------------------for Testing Purpose

	return server, nil
}

func dataMarchal() []byte {
	data, _ := json.Marshal(
		gin.H{
			"username":  util.RandomString(6),
			"password":  "secret",
			"full_name": util.RandomString(6),
			"email":     util.RandomEmail()})

	return data
}

func NewServerForTesting(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router := gin.Default()
	router.POST("/users", server.testCreateUser)
	router.POST("/users/login", server.loginUser)
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)
	authRoutes.POST("/transfer", server.createTrasnfer)

	server.router = router

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser_asynchronous)
	router.POST("/users/login", server.loginUser)

	// Set API route and which handler will be act by its route

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)
	authRoutes.POST("/transfer", server.createTrasnfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
