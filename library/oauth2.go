package library

import (
	"fmt"
	"net/http"
	"time"

	"github.com/yamakiller/magicNet/engine/logger"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
)

// OAuth2 : xxx
type OAuth2 struct {
	TokenExp          int //Minute one
	RefreshTokenExp   int //Minute one
	IsGenerateRefresh bool
	S256Key           string
	AccessURI         string

	m *manage.Manager
	c *store.ClientStore
	s *server.Server
}

// Init : xxx
func (oa *OAuth2) Init(method IHTTPSrvMethod) {
	oauth2CodeTokenCfg := &manage.Config{
		AccessTokenExp:    time.Minute * time.Duration(oa.TokenExp),
		RefreshTokenExp:   time.Minute * time.Duration(oa.RefreshTokenExp),
		IsGenerateRefresh: oa.IsGenerateRefresh,
	}

	oauth2ClientTokenCfg := &manage.Config{AccessTokenExp: time.Minute * time.Duration(oa.TokenExp)}

	oa.m = manage.NewDefaultManager()
	oa.m.SetAuthorizeCodeTokenCfg(oauth2CodeTokenCfg)
	oa.m.SetClientTokenCfg(oauth2ClientTokenCfg)

	oa.m.MustTokenStorage(store.NewMemoryTokenStore())

	if oa.S256Key != "" {
		oa.m.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte(oa.S256Key), jwt.SigningMethodHS256))
	}

	oa.c = store.NewClientStore()

	oa.m.MapClientStorage(oa.c)

	oa.s = server.NewDefaultServer(oa.m)

	oa.s.SetAllowGetAccessRequest(true)
	oa.s.SetAllowedGrantType(oauth2.ClientCredentials, oauth2.Refreshing)
	oa.s.SetClientInfoHandler(server.ClientFormHandler)
	oa.m.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

	method.RegisterMethod(oa.AccessURI, "get|put|post", func(w http.ResponseWriter, r *http.Request) {
		err := oa.s.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})
}

// SetErrorHandle : xxx
func (oa *OAuth2) SetErrorHandle(owner uint32) {
	oa.s.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		logger.Error(owner, "internal error:%s", err.Error())
		return
	})

	oa.s.SetResponseErrorHandler(func(re *errors.Response) {
		logger.Error(owner, "response error:%s", re.Error.Error())
	})
}

//RegisterClient : 注册授权的客户端
func (oa *OAuth2) RegisterClient(id string, secret string, domain string, userid string) {
	oa.c.Set(id, &models.Client{ID: id, Secret: secret, Domain: domain, UserID: userid})
}

// RegisterAuth2Method : 注册受保护的方法
func (oa *OAuth2) RegisterAuth2Method(method IHTTPSrvMethod, pattern string, httpMetod string, f HTTPSrvFunc) {
	method.RegisterMethod(pattern, httpMetod, validateToken(method, f, oa.s))
}

// RegisterAuth2MethodJS : 注册受保护的JS方法
func (oa *OAuth2) RegisterAuth2MethodJS(method IHTTPSrvMethod, pattern string, httpMetod string, f string) {
	method.RegisterMethod(pattern, httpMetod, validateToken(method, f, oa.s))
}

//HTTPSrvFunc
func validateToken(method IHTTPSrvMethod, f interface{}, srv *server.Server) HTTPSrvFunc {
	return HTTPSrvFunc(func(w http.ResponseWriter, r *http.Request) {
		/*_, err := srv.ValidationBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}*/

		if v, ok := f.(string); ok {
			fmt.Println("run js 1")
			if jsm, jsmok := method.(*HTTPSrvMethodJS); jsmok {
				fmt.Println("run js 2")
				jsm.runJs(v, w, r)
			}
		}

		if v, ok := f.(func(http.ResponseWriter, *http.Request)); ok {
			v(w, r)
		}
	})
}
