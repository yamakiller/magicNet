package preset_function

import (
  "time"
  "net/http"
  "magicNet/engine/util"
  "magicNet/engine/monitor"
  "magicNet/engine/logger"
  "gopkg.in/oauth2.v3"
  "gopkg.in/oauth2.v3/errors"
  "gopkg.in/oauth2.v3/models"
  "gopkg.in/oauth2.v3/manage"
  "gopkg.in/oauth2.v3/server"
  "gopkg.in/oauth2.v3/store"
)


var managerAuth2 *manage.Manager
var serverAuth2 *server.Server

//暂支持Clien模式及刷i新令牌
func InitializeAuth2() {
  logger.Info(0, "monitor service initialize auth2")

  autoEnv := util.GetEnvInstance().GetMap("auto2")
  authTokenExp := util.GetEnvInt(autoEnv, "auth-token-exp", 10)
  authRefreshTokenExp := util.GetEnvInt(autoEnv, "auth-refresh-token-exp", 20)
  authIsGenerateRefresh := util.GetEnvBool(autoEnv, "auth-is-generate-refresh-token", true)

  authCodeTokenCfg := &manage.Config{
  	AccessTokenExp: time.Minute * time.Duration(authTokenExp),
  	RefreshTokenExp: time.Minute * time.Duration(authRefreshTokenExp),
  	IsGenerateRefresh: authIsGenerateRefresh,
  }

  authClientTokenCfg := &manage.Config{AccessTokenExp: time.Minute * time.Duration(authTokenExp)}

  managerAuth2 = manage.NewDefaultManager()
  managerAuth2.SetAuthorizeCodeTokenCfg(authCodeTokenCfg)
  managerAuth2.SetClientTokenCfg(authClientTokenCfg)

  //后续提供优化方案-----------------------------------------------
  managerAuth2.MustTokenStorage(store.NewMemoryTokenStore())
  if autoEnv["auth-clients"].Exists() {
    csStore := store.NewClientStore()
    for _, cs := range autoEnv["auth-clients"].Array() {
        csId := cs.Map()["ID"].String()
        csSecret := cs.Map()["Secret"].String()
        csDomain := cs.Map()["Domain"].String()

        csStore.Set(csId, &models.Client{ID : csId,
                    Secret : csSecret,
                    Domain : csDomain})
    }
    managerAuth2.MapClientStorage(csStore)
  }
  //--------------------------------------------------

  serverAuth2 = server.NewDefaultServer(managerAuth2)

  serverAuth2.SetAllowGetAccessRequest(true)
  serverAuth2.SetAllowedGrantType(oauth2.ClientCredentials, oauth2.Refreshing)
  serverAuth2.SetClientInfoHandler(server.ClientFormHandler)
  managerAuth2.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

  serverAuth2.SetInternalErrorHandler(func(err error)(re *errors.Response) {
    logger.Error(0, "internal error:%s", err.Error())
    return
  })

  serverAuth2.SetResponseErrorHandler(func(re *errors.Response) {
    logger.Error(0, "response error:%s", re.Error.Error())
  })

  monitor.RegisterHttpMethod("/login/oauth/acces_token", func(w http.ResponseWriter, r *http.Request){
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
