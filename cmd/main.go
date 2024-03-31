package main

import (
	"log"
	"time"

	"github.com/coderero/erochat-server/api/handler"
	apiMiddleware "github.com/coderero/erochat-server/api/middleware"
	"github.com/coderero/erochat-server/api/service"
	"github.com/coderero/erochat-server/api/utils"
	"github.com/coderero/erochat-server/db/mysql"
	"github.com/go-playground/validator/v10"
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
		cors    = middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		})

		// Echo Routes.
		apiV1     = app.Group("/api/v1")
		apiAuthV1 = app.Group("/api/auth/v1")

		// Service initialization.
		passService     = service.NewScryptService(1<<14, 8, 1, 32, 22)
		jwtTokenService = tokenService

		// Middleware initialization.
		auth = apiMiddleware.JWTMiddleware(jwtTokenService)

		// Store initialization.
		user    = mysql.NewUserStore(db)
		profile = mysql.NewProfileStore(db)
		status  = mysql.NewStatusStore(db)
		friend  = mysql.NewFriendStore(db)

		// Validator initialization.
		validator = validator.New()

		// Handler initialization.
		authHandler       = handler.NewAuthHandler(validator, user, passService, jwtTokenService)
		profileHandler    = handler.NewProfileHandler(validator, profile, user)
		statusHandler     = handler.NewUserStatusHandler(validator, user, status)
		friendshipHandler = handler.NewUserFriendShipHandler(validator, user, friend)
	)

	// Use middleware.

	/* Main App */
	app.Use(recover)
	app.Use(logger)
	app.Use(cors)

	/* API V1 */
	apiV1.Use(auth)

	// Echo configration
	app.HTTPErrorHandler = utils.CustomHTTPErrorHandler(app)

	// Validator configuration.
	validator.RegisterTagNameFunc(utils.ValidatorTagFunc)

	// Routes.

	/* Auth routes. */
	apiAuthV1.POST("/login", authHandler.Login)
	apiAuthV1.POST("/register", authHandler.Register)
	apiAuthV1.POST("/logout", authHandler.Logout)

	/* User routes. */
	apiV1.GET("/user/profile", profileHandler.GetProfile)
	apiV1.POST("/user/profile", profileHandler.CreateProfile)
	apiV1.PUT("/user/profile", profileHandler.UpdateProfile)
	apiV1.GET("/user/profile/:uid", profileHandler.GetProfileByID)
	apiV1.POST("/user/profile/:uid", profileHandler.AddFriend)
	apiV1.DELETE("/user/profile", profileHandler.DeleteProfile)
	apiV1.PATCH("/user/profile/reactivate", profileHandler.ReactivateProfile)

	/* Status routes. */
	apiV1.GET("/user/status", statusHandler.GetStatus)
	apiV1.POST("/user/status", statusHandler.CreateStatus)
	apiV1.DELETE("/user/status/:uid", statusHandler.DeleteStatus)

	/* Friend routes. */
	apiV1.GET("/user/friends/details", friendshipHandler.GetFriends)
	apiV1.GET("/user/friends/details/:uid", friendshipHandler.GetFriend)
	apiV1.DELETE("/user/friends/details/:uid", friendshipHandler.DeleteFriend)
	apiV1.GET("/user/friends/requests", friendshipHandler.GetFriendRequests)
	apiV1.GET("/user/friends/requests/:uid", friendshipHandler.GetFriendRequest)
	apiV1.PATCH("/user/friends/requests/:uid", friendshipHandler.AcceptFriendRequest)
	apiV1.DELETE("/user/friends/requests/:uid", friendshipHandler.DeleteFriendRequest)
	apiV1.GET("/user/friends/status", friendshipHandler.GetFriendsStatus)
	apiV1.GET("/user/friends/status/:uid", friendshipHandler.GetFriendStatus)

	/* Start the HTTP server. */

	if err := app.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
