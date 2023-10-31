package connection

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	ginserver "github.com/go-oauth2/gin-server"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

// var (
// 	dumpvar     bool
// 	idvar       string
// 	userNamevar string
// 	userPassvar string
// 	secretvar   string
// 	domainvar   string
// 	// portvar     int
// )

// func init() {
// 	flag.BoolVar(&dumpvar, "d", true, "Dump requests and responses")
// 	flag.StringVar(&idvar, "i", "000000", "The client id being passed in")
// 	flag.StringVar(&secretvar, "s", "999999", "The client secret being passed in")
// 	flag.StringVar(&userNamevar, "y", "test", "The username being passed in")
// 	flag.StringVar(&userPassvar, "z", "test", "The password being passed in")
// 	flag.StringVar(&domainvar, "r", "http://localhost:9094", "The domain of the redirect url")
// 	flag.IntVar(&portvar, "p", 9096, "the base port for the server")
// }

// Define a function to reset the OAuth2 manager and related configurations
func resetOAuth2Manager(idvar, secretvar, domainvar string) *manage.Manager {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	// manager.MustTokenStorage(store.NewFileTokenStore("data.db"))
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// generate jwt access token
	manager.MapAccessGenerate(generates.NewAccessGenerate())

	// client store
	clientStore := store.NewClientStore()
	// cs := oauth2.Config{
	// 	ClientID:     idvar,
	// 	ClientSecret: secretvar,
	// 	RedirectURL:  "http://localhost:9096/oauth2/api",
	// 	Scopes:       []string{"profile", "email"},
	// }
	// id_client, cs := ClientInfo(idvar, secretvar)
	clientInfo := &models.Client{
		ID:     idvar,
		Secret: secretvar,
		Domain: domainvar,
	}
	clientStore.Set(idvar, clientInfo)
	manager.MapClientStorage(clientStore)

	// Define a custom grant type handler
	// manager.SetAuthorizeCodeTokenCfg(&manage.Config{AccessTokenExp: 0})

	// SetClientTokenCfg set the client grant token config
	// manager.SetClientTokenCfg(&manage.Config{
	// 	AccessTokenExp:    time.Hour * 24 * 365 * 10, // 10 years
	// 	RefreshTokenExp:   time.Hour * 24 * 365 * 10, // 10 years
	// 	IsGenerateRefresh: true,
	// })

	// SetPasswordTokenCfg the access token and refresh token configuration
	manager.SetPasswordTokenCfg(&manage.Config{
		AccessTokenExp:    time.Hour * 24 * 365 * 10, // 10 years
		RefreshTokenExp:   time.Hour * 24 * 365 * 10, // 10 years
		IsGenerateRefresh: true,
	})

	return manager
}

func Middleware(ctx *gin.Context) {
	idvar := ctx.Query("client_id")
	secretvar := ctx.Query("client_secret")
	userNamevar := ctx.Query("username")
	userPassvar := ctx.Query("password")
	dumpvar := true
	domainvar := "http://localhost:9094"

	// manager := manage.NewManager()
	manager := resetOAuth2Manager(idvar, secretvar, domainvar)

	// Initialize the oauth2 service
	server.NewServer(server.NewConfig(), manager)
	ginserver.InitServer(manager)
	ginserver.SetAllowGetAccessRequest(dumpvar)
	ginserver.SetClientInfoHandler(server.ClientFormHandler)
	grantType := oauth2.GrantType(oauth2.PasswordCredentials.String())
	ginserver.SetAllowedGrantType(grantType)
	// SetPasswordAuthorizationHandlers(userNamevar, userPassvar, idvar)
	ginserver.SetPasswordAuthorizationHandler(func(ctx context.Context, ClientID, username, password string) (userID string, err error) {
		// Implement your authentication logic here and return the userID if valid
		if username == userNamevar && password == userPassvar {
			fmt.Println("XXXX")
			userID = ClientID
		}
		return
	})

	// Handle the token request
	ginserver.HandleTokenRequest(ctx)

	// manager := manage.NewDefaultManager()
	// manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// // token store
	// manager.MustTokenStorage(store.NewMemoryTokenStore())

	// // generate jwt access token
	// // manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	// manager.MapAccessGenerate(generates.NewAccessGenerate())

	// // client store
	// clientStore := store.NewClientStore()
	// clientInfo := &models.Client{
	// 	ID:     idvar,
	// 	Secret: secretvar,
	// 	Domain: domainvar,
	// }
	// clientStore.Set(idvar, clientInfo)

	// manager.MapClientStorage(clientStore)

	// // manager.SetPasswordTokenCfg({
	// // 	AccessTokenExp:    time.Hour,
	// // 	RefreshTokenExp:   30 * 24 * time.Hour,
	// // 	IsGenerateRefresh: true,
	// // })

	// manager.SetPasswordTokenCfg(&manage.Config{
	// 	AccessTokenExp:    time.Hour,
	// 	RefreshTokenExp:   30 * 24 * time.Hour,
	// 	IsGenerateRefresh: true,
	// })

	// srv.SetUserAuthorizationHandler(userAuthorizeHandler)

	// srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
	// 	log.Println("Internal Error:", err.Error())
	// 	return
	// })

	// srv.SetResponseErrorHandler(func(re *errors.Response) {
	// 	log.Println("Response Error:", re.Error.Error())
	// })

	// // Initialize the oauth2 service
	// ginserver.InitServer(manager)
	// ginserver.SetAllowGetAccessRequest(true)
	// ginserver.SetClientInfoHandler(server.ClientFormHandler)

	// // // SetClientTokenCfg set the client grant token config
	// // manager.SetClientTokenCfg(&manage.Config{
	// // 	AccessTokenExp:    time.Duration(2000),
	// // 	RefreshTokenExp:   time.Duration(2000),
	// // 	IsGenerateRefresh: true,
	// // })

	// // // SetAuthorizeCodeExp set the authorization code expiration time
	// // manager.SetAuthorizeCodeExp(time.Duration(2000))

	// // // SetRefreshTokenCfg set the refreshing token config
	// // manager.SetRefreshTokenCfg(&manage.RefreshingConfig{
	// // 	AccessTokenExp:     time.Duration(2000),
	// // 	RefreshTokenExp:    time.Duration(2000),
	// // 	IsGenerateRefresh:  true,
	// // 	IsResetRefreshTime: true,
	// // })

	// // // Initialize the oauth2 service
	// // ginserver.InitServer(manager)
	// // ginserver.SetAllowGetAccessRequest(true)
	// // ginserver.SetClientInfoHandler(server.ClientFormHandler)
}

// func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
// 	if dumpvar {
// 		_ = dumpRequest(os.Stdout, "userAuthorizeHandler", r) // Ignore the error
// 	}
// 	store, err := session.Start(r.Context(), w, r)
// 	if err != nil {
// 		return
// 	}

// 	uid, ok := store.Get("LoggedInUserID")
// 	if !ok {
// 		if r.Form == nil {
// 			r.ParseForm()
// 		}

// 		store.Set("ReturnUri", r.Form)
// 		store.Save()

// 		w.Header().Set("Location", "/login")
// 		w.WriteHeader(http.StatusFound)
// 		return
// 	}

// 	userID = uid.(string)
// 	store.Delete("LoggedInUserID")
// 	store.Save()
// 	return
// }

// func dumpRequest(writer io.Writer, header string, r *http.Request) error {
// 	data, err := httputil.DumpRequest(r, true)
// 	if err != nil {
// 		return err
// 	}
// 	writer.Write([]byte("\n" + header + ": \n"))
// 	writer.Write(data)
// 	return nil
// }

// func ClientInfo(id_client string, secret_client string) (string, *models.Client) {
// 	clientInfo := &models.Client{
// 		ID:     id_client,
// 		Secret: secret_client,
// 		Domain: domainvar,
// 	}

// 	return id_client, clientInfo
// }

// func SetPasswordAuthorizationHandlers(u_name string, u_pass string, id_client string) {
// 	ginserver.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string, err error) {
// 		// Implement your authentication logic here and return the userID if valid
// 		if username == u_name && password == u_pass {
// 			userID = id_client
// 		}
// 		return
// 	})
// }
