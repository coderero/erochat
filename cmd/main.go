package main

import (
	"log"
	"time"

	"github.com/coderero/erochat-server/api/handler"
	apiMiddleware "github.com/coderero/erochat-server/api/middleware"
	"github.com/coderero/erochat-server/api/service"
	"github.com/coderero/erochat-server/api/utils"
	"github.com/coderero/erochat-server/db/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Configuration variables.
	const (
		// MySQL DSN.
		dsn            = "root:secretpassword@tcp(localhost:3306)/erochat?parseTime=true"
		maxConnections = 10
	)

	// Create a new connection pool.
	db, err := mysql.NewConnectionPool(dsn, maxConnections)
	if err != nil {
		panic(err)
	}

	// Close the connection pool when the main function returns.
	obj, err := db.Get()
	if err != nil {
		panic(err)
	}
	defer obj.Close()

	// Get RSA keys for JWT from the certificate files.
	privKey, err := utils.GetFile("certs/app.rsa.key")
	if err != nil {
		panic(err)
	}
	pubKey, err := utils.GetFile("certs/app.rsa.pub")
	if err != nil {
		panic(err)
	}

	// Create a new token service.
	tokenService, err := service.NewJWTService(privKey, pubKey, time.Hour*24, time.Hour*24*7)
	if err != nil {
		panic(err)
	}

	// Echo and HTTP server Configuration variables.

	var (
		// Echo instance.
		app = echo.New()

		// Echo middleware.
		recover = middleware.Recover()
		logger  = middleware.Logger()

		// Echo Routes.
		_         = app.Group("/api/v1")
		apiAuthV1 = app.Group("/api/auth/v1")

		// Service initialization.
		passService     = service.NewScryptService(1<<14, 8, 1, 32, 22)
		jwtTokenService = tokenService

		// Middleware initialization.
		auth = apiMiddleware.JWTMiddleware(jwtTokenService)

		// Handler initialization.
		authHandler = handler.NewAuthHandler(mysql.NewUserStore(db), passService, jwtTokenService)
	)

	// Use middleware.
	app.Use(recover)
	app.Use(logger)

	// Custom HTTP error handler.
	app.HTTPErrorHandler = utils.CustomHTTPErrorHandler(app)

	// Custom Middleware.
	apiAuthV1.Use(auth)

	// Routes.
	apiAuthV1.POST("/login", authHandler.Login)
	apiAuthV1.POST("/register", authHandler.Register)
	apiAuthV1.POST("/logout", authHandler.Logout)

	if err := app.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
