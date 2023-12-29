package connection

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	gen "github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt"
)

var (
	clientStore   = store.NewClientStore()
	tokenStore, _ = store.NewMemoryTokenStore()
	token_exp     time.Duration
	srv           *server.Server
	// token         = models.NewToken()
)

// Define a function to reset the OAuth2 manager and related configurations
func resetOAuth2Manager(idvar, secretvar, domainvar, keypassvar string, ctx *gin.Context) *manage.Manager {
	manager := manage.NewDefaultManager()
	manager.MapClientStorage(clientStore)
	manager.MapTokenStorage(tokenStore)

	// generate jwt access token
	manager.MapAccessGenerate(gen.NewJWTAccessGenerate("HS512", []byte("secret"), jwt.SigningMethodHS512))

	if keypassvar == "MasterPassword" {
		token_exp = 0
	} else {
		token_exp = time.Hour * 24 // 1 hari
	}
	manager.SetPasswordTokenCfg(&manage.Config{
		// AccessTokenExp:    time.Hour * 24 * 365 * 10, // 10 years
		AccessTokenExp:  token_exp,
		RefreshTokenExp: token_exp,
		// RefreshTokenExp:   time.Hour * 24 * 365 * 10, // 10 years
		IsGenerateRefresh: true,
	})

	clientInfo := &models.Client{
		ID:     idvar,
		Secret: secretvar,
		Domain: domainvar,
	}
	clientStore.Set(idvar, clientInfo)

	return manager
}

func Middleware(ctx *gin.Context) {
	idvar := ctx.Query("client_id")
	secretvar := ctx.Query("client_secret")
	userNamevar := ctx.Query("username")
	userPassvar := ctx.Query("password")
	keyPassvar := ctx.Query("key")
	domainvar := "http://localhost:9096"

	// Reset the OAuth2 manager
	manager := resetOAuth2Manager(idvar, secretvar, domainvar, keyPassvar, ctx)

	// Initialize the OAuth2 service
	srv = server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	srv.SetAllowedGrantType("password")

	// Set the password authorization handler
	srv.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		if username == userNamevar && password == userPassvar {
			userID = clientID
		}
		return
	})

	// set token_type
	if keyPassvar == "MasterPassword" {
		srv.SetTokenType("TOKEN GLOBAL")
	}

	// Handle the token request
	err := srv.HandleTokenRequest(ctx.Writer, ctx.Request)
	if err != nil {
		// Handle error, e.g., return an error response
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
}

func Validation(ctx *gin.Context) {
	if srv == nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"ERROR": "Silahkan Login..",
		})
		return
	}
	ti, err := srv.ValidationBearerToken(ctx.Request)
	if err == nil {
		fmt.Println(ti)
	} else {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"ERROR": "INVALID_TOKEN",
		})
	}
}
