package library

import (
	"net/http"
	"time"

	"github.com/yamakiller/magicNet/logger"

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
func (slf *OAuth2) Init(method IHTTPSrvMethod) {
	oauth2CodeTokenCfg := &manage.Config{
		AccessTokenExp:    time.Minute * time.Duration(slf.TokenExp),
		RefreshTokenExp:   time.Minute * time.Duration(slf.RefreshTokenExp),
		IsGenerateRefresh: slf.IsGenerateRefresh,
	}

	oauth2ClientTokenCfg := &manage.Config{AccessTokenExp: time.Minute * time.Duration(slf.TokenExp)}

	slf.m = manage.NewDefaultManager()
	slf.m.SetAuthorizeCodeTokenCfg(oauth2CodeTokenCfg)
	slf.m.SetClientTokenCfg(oauth2ClientTokenCfg)

	slf.m.MustTokenStorage(store.NewMemoryTokenStore())

	if slf.S256Key != "" {
		slf.m.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte(slf.S256Key), jwt.SigningMethodHS256))
	}

	slf.c = store.NewClientStore()

	slf.m.MapClientStorage(slf.c)

	slf.s = server.NewDefaultServer(slf.m)

	slf.s.SetAllowGetAccessRequest(true)
	slf.s.SetAllowedGrantType(oauth2.ClientCredentials, oauth2.Refreshing)
	slf.s.SetClientInfoHandler(server.ClientFormHandler)
	slf.m.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

	method.RegisterMethod(slf.AccessURI, "get|put|post", HTTPSrvFunc(func(w http.ResponseWriter, r *http.Request) {
		err := slf.s.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}))
}

// SetErrorHandle : xxx
func (slf *OAuth2) SetErrorHandle(owner uint32) {
	slf.s.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		logger.Error(owner, "internal error:%s", err.Error())
		return
	})

	slf.s.SetResponseErrorHandler(func(re *errors.Response) {
		logger.Error(owner, "response error:%s", re.Error.Error())
	})
}

//RegisterClient : 注册授权的客户端
func (slf *OAuth2) RegisterClient(id string, secret string, domain string, userid string) {
	slf.c.Set(id, &models.Client{ID: id, Secret: secret, Domain: domain, UserID: userid})
}

// RegisterAuth2Method : 注册受保护的方法
func (slf *OAuth2) RegisterAuth2Method(method IHTTPSrvMethod, pattern string, httpMetod string, f HTTPSrvFunc) {
	method.RegisterMethod(pattern, httpMetod, validateToken(method, f, slf.s))
}

// RegisterAuth2MethodJS : 注册受保护的JS方法
func (slf *OAuth2) RegisterAuth2MethodJS(method IHTTPSrvMethod, pattern string, httpMetod string, f string) {
	method.RegisterMethod(pattern, httpMetod, validateToken(method, f, slf.s))
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
			if jsm, jsmok := method.(*HTTPSrvMethodJS); jsmok {
				jsm.runJs(v, w, r)
			}
		}

		if v, ok := f.(func(http.ResponseWriter, *http.Request)); ok {
			v(w, r)
		}
	})
}
