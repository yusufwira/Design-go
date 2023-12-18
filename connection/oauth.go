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
	// accessGenerator := gen.NewAccessGenerate()
	// globalToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiIxNDc1IiwianRpIjoiMzJlYmUwOTM1NTkxMjRiMGJhYmFiYmU1ZmVlM2RlNzI0Y2FmOTcyY2IzNTA5ZjFmYjlmM2YyMGU1ZTZhYjI2MzlmMGM2MmJiMzg1NjU0NjAiLCJpYXQiOjE2OTI3NTcyNDguNjgwOTE2LCJuYmYiOjE2OTI3NTcyNDguNjgwOTIsImV4cCI6MTcyNDM3OTY0OC42NzAxNzgsInN1YiI6IjU0MyIsInNjb3BlcyI6W119.OUYTYh6BEXWxg_powW6DlrUmhXXwYhKW30bJOjhQCLZuGL8oTxE7ymafudNI-w1jrZeToATdg0vC14cB_yBYKkYzYpAK2QDw3oFQmBoF9128PcXOeJeSlgpxJCu7cCM0bbC-1lEs6Jn4aRLt0un_D6lO2h1XI1W0-gLMD5J6C6LQ3sBVYfbeKi7nKCg0msEx9IIZLa3UyvP5foNqefe2_SvMfYISJL2-zDyZwS65b7xJQDiqjS-TWEIk2w6MfVRzAoNykEHq3GVPzOSXxeN7uO9go6i28HHSaCoJ6Z6UL5p63BFnREOTmxAnwvf9xpF4UbtiRHVr7PCTq02WwYRUuzRe18F3XPBWqMp3Co3vDRQMqMh54KkjOAOO4U9oJC4cBu1sCWPDF6oW0qSJzpPxLiaVrzFtPoxBkOWnJlws8MliMbgL0WexTFRwY5AKZT3gHACGDIhFwPZqffASxiiA8iVjaoI4t45DhSf9gYwEzXtjNi_WXZVD7GW6cXjFL_Q6LaMiLG3ZLiXFpA6JcxbllzyjUEOoaSC1jz4A08siOEtoafdJXKLJg3lV45Cxs1pAUwprexn_zzCEfL7rgTXPHAZYBFfOaV_IOWjtF5LsbLhH9N68AWcRz9jc12ysmWO9G7VGlgoWOkk6--yMbJeqZXsbkd_HZbd8sEbHCbF6C5c"
	// token = models.NewToken()
	// token.SetAccess(globalToken)
	// token.SetRefresh(globalToken)

	// // generates := gen.NewAccessGenerate().Token(context, &oauth2.GenerateBasic{}, true)
	// generates := gen.JWTAccessGenerate{
	// 	SignedKeyID:  "HS512",
	// 	SignedKey:    []byte("base64:yMmq7EsOmOJPFp+5UyK8LwIVZsmHyKFRbZ/jlO9h2Sc="),
	// 	SignedMethod: jwt.SigningMethodHS512,
	// }
	// manager.MapAccessGenerate()

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
	// generateBasic := &oauth2.GenerateBasic{
	// 	Client:    clientInfo,
	// 	TokenInfo: token,
	// }
	// access, refresh, _ := generates.Token(context.TODO(), generateBasic, true)
	// fmt.Println(access)
	// fmt.Println(refresh)
	// generates.SignedMethod.Alg()

	// accessGenerator.Token(context.TODO(), generateBasic, true)
	// accessGenerator := gen.NewJWTAccessGenerate()
	// manager.MapAccessGenerate(generates)
	// manager.MapAccessGenerate(&generates)
	clientStore.Set(idvar, clientInfo)

	return manager
}

func Middleware(ctx *gin.Context) {
	idvar := ctx.Query("client_id")
	secretvar := ctx.Query("client_secret")
	userNamevar := ctx.Query("username")
	userPassvar := ctx.Query("password")
	keyPassvar := ctx.Query("key")
	fmt.Println(userNamevar, userPassvar)
	domainvar := "http://localhost:9096"

	// Reset the OAuth2 manager
	manager := resetOAuth2Manager(idvar, secretvar, domainvar, keyPassvar, ctx)
	// clinfo, _ := manager.GetClient(ctx, idvar)
	// manager.

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
