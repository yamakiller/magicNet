package normal

import (
  "time"
  "errors"
  "strings"
  "net/http"
  "magicNet/util"
  "magicNet/monitor"
  //"github.com/tidwall/gjson"
  //"gopkg.in/oauth2.v3/models"
  "gopkg.in/oauth2.v3/manage"
  "gopkg.in/oauth2.v3/server"
  "gopkg.in/oauth2.v3/store"
)

type authAccount struct {
  userid   string
  username string
  userpass string
}

var managerAuth2 *manage.Manager
var serverAuth2 *server.Server
var usersAuth2 map[string]*authAccount

func initializeAuth2() {
  autoEnv := util.GetEnvInstance().GetMap("auto2")
  authTokenExp := util.GetEnvInt(autoEnv, "auth-token-exp", 10)
  authRefreshTokenExp := util.GetEnvInt(autoEnv, "auth-refresh-token-exp", 20)
  authIsGenerateRefresh := util.GetEnvBool(autoEnv, "auth-is-generate-refresh-token", true)

  authCodeTokenCfg := &manage.Config{
  	AccessTokenExp: time.Minute * time.Duration(authTokenExp),
  	RefreshTokenExp: time.Minute * time.Duration(authRefreshTokenExp),
  	IsGenerateRefresh: authIsGenerateRefresh,
  }

  managerAuth2 = manage.NewManager()
  managerAuth2.MustTokenStorage(store.NewMemoryTokenStore())
  managerAuth2.SetPasswordTokenCfg(authCodeTokenCfg)

  if autoEnv["auth-accounts"].Exists() {
    for _, cs := range autoEnv["auth-accounts"].Array() {
      usersAuth2[cs.Map()["username"].String()] = &authAccount{
        cs.Map()["userid"].String(),
        cs.Map()["username"].String(),
        cs.Map()["userpass"].String()}
    }
  }

  serverAuth2 = server.NewDefaultServer(managerAuth2)
  //serverAuth2.SetClientAuthorizedHandler(handler)
  serverAuth2.SetPasswordAuthorizationHandler(func(username, password string)(userID string, err error){
        user := usersAuth2[username]
        if user == nil {
          return "", errors.New("No authorized account exists")
        }

        //TODO: MD5
        if strings.Compare(user.userpass, password) != 0 {
          return "", errors.New("Password Error in Authorized Account")
        }

        return user.userid, nil
  })

  serverAuth2.SetAllowGetAccessRequest(true)
  serverAuth2.SetClientInfoHandler(server.ClientFormHandler)

  monitor.RegisterHttpMethod("/login/oauth/authorize", func(w http.ResponseWriter, r *http.Request){
      err := serverAuth2.HandleAuthorizeRequest(w, r)
      if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
      }
  })

  monitor.RegisterHttpMethod("/login/oauth/acces_token", func(w http.ResponseWriter, r *http.Request){
    err := serverAuth2.HandleTokenRequest(w, r)
    if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
    }
  })
  registerURIService()
}

func registerURIService(){

}
