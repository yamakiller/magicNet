package preset_function

import (
	"time"
	//"reflect"
	"magicNet/engine/logger"
	"magicNet/engine/monitor"
	"magicNet/engine/util"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

var managerAuth2 *manage.Manager
var serverAuth2 *server.Server

type clientConfig struct {
	ID     string
	Secret string
	Domain string
}

//暂支持Clien模式及刷i新令牌
func InitializeAuth2() {
	logger.Info(0, "monitor service initialize oauth2")
	authEnv := util.GetEnvMap(util.GetEnvRoot(), "oauth2")
	util.AssertEmpty(authEnv, "monitor serivce oauth2 config error")

	authTokenExp := util.GetEnvInt(authEnv, "auth-token-exp", 10)
	authRefreshTokenExp := util.GetEnvInt(authEnv, "auth-refresh-token-exp", 20)
	authIsGenerateRefresh := util.GetEnvBoolean(authEnv, "auth-is-generate-refresh-token", true)
	authHS256Key := util.GetEnvString(authEnv, "auth-signing-=256-key", "")

	authCodeTokenCfg := &manage.Config{
		AccessTokenExp:    time.Minute * time.Duration(authTokenExp),
		RefreshTokenExp:   time.Minute * time.Duration(authRefreshTokenExp),
		IsGenerateRefresh: authIsGenerateRefresh,
	}

	authClientTokenCfg := &manage.Config{AccessTokenExp: time.Minute * time.Duration(authTokenExp)}

	managerAuth2 = manage.NewDefaultManager()
	managerAuth2.SetAuthorizeCodeTokenCfg(authCodeTokenCfg)
	managerAuth2.SetClientTokenCfg(authClientTokenCfg)

	//后续提供优化方案-----------------------------------------------
	managerAuth2.MustTokenStorage(store.NewMemoryTokenStore())
	if strings.Compare(authHS256Key, "") != 0 {
		managerAuth2.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte(authHS256Key), jwt.SigningMethodHS256))
	}

	clients := util.GetEnvArray(authEnv, "auth-clients")
	if clients != nil {
		csStore := store.NewClientStore()
		for _, cs := range clients {
			clinetinfo := util.ToEnvMap(cs)
			csId := util.GetEnvString(clinetinfo, "ID", "")
			csSecret := util.GetEnvString(clinetinfo, "Secret", "")
			csDomain := util.GetEnvString(clinetinfo, "Domain", "")

			util.Assert(strings.Compare(csId, "") != 0 &&
				strings.Compare(csSecret, "") != 0 &&
				strings.Compare(csDomain, "") != 0,
				"Oauth2 client "+csId+" authorization configuration has security risks")

			csStore.Set(csId, &models.Client{ID: csId,
				Secret: csSecret,
				Domain: csDomain})
		}
		managerAuth2.MapClientStorage(csStore)
	}
	//--------------------------------------------------

	serverAuth2 = server.NewDefaultServer(managerAuth2)

	serverAuth2.SetAllowGetAccessRequest(true)
	serverAuth2.SetAllowedGrantType(oauth2.ClientCredentials, oauth2.Refreshing)
	serverAuth2.SetClientInfoHandler(server.ClientFormHandler)
	managerAuth2.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

	serverAuth2.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		logger.Error(0, "internal error:%s", err.Error())
		return
	})

	serverAuth2.SetResponseErrorHandler(func(re *errors.Response) {
		logger.Error(0, "response error:%s", re.Error.Error())
	})

	monitor.RegisterHttpMethod("/login/oauth/acces_token", func(w http.ResponseWriter, r *http.Request) {
		err := serverAuth2.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})
}

func RegisterAuth2Method(pattern string, f monitor.MonitorHttpFunction) {
	monitor.RegisterHttpMethod(pattern, validateToken(f, serverAuth2))
}

func validateToken(f monitor.MonitorHttpFunction, srv *server.Server) monitor.MonitorHttpFunction {
	return monitor.MonitorHttpFunction(func(w http.ResponseWriter, r *http.Request) {
		_, err := srv.ValidationBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		f(w, r)
	})
}
