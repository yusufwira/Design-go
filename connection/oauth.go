package connection

import (
	"time"

	ginserver "github.com/go-oauth2/gin-server"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

func Middleware() {
	manager := manage.NewDefaultManager()
	// token store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// SetClientTokenCfg set the client grant token config
	manager.SetClientTokenCfg(&manage.Config{
		AccessTokenExp:    time.Duration(2000),
		RefreshTokenExp:   time.Duration(2000),
		IsGenerateRefresh: true,
	})

	// SetAuthorizeCodeExp set the authorization code expiration time
	manager.SetAuthorizeCodeExp(time.Duration(2000))

	// SetRefreshTokenCfg set the refreshing token config
	manager.SetRefreshTokenCfg(&manage.RefreshingConfig{
		AccessTokenExp:     time.Duration(2000),
		RefreshTokenExp:    time.Duration(2000),
		IsGenerateRefresh:  true,
		IsResetRefreshTime: true,
	})

	// client store
	clientStore := store.NewClientStore()
	clientStore.Set("000000", &models.Client{
		ID:     "000000",
		Secret: "999999",
		Domain: "http://localhost",
	})
	manager.MapClientStorage(clientStore)

	// Initialize the oauth2 service
	ginserver.InitServer(manager)
	ginserver.SetAllowGetAccessRequest(true)
	ginserver.SetClientInfoHandler(server.ClientFormHandler)
}
