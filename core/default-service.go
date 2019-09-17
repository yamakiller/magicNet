package core

import (
	"strings"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/library"
	"github.com/yamakiller/magicNet/service"
	"github.com/yamakiller/magicNet/util"

	"github.com/yamakiller/magicNet/engine/logger"
)

// DefaultService : 默认服务系统
type DefaultService struct {
	monitorSrv *service.MonitorService
}

// InitService : 初始化服务模块
func (ds *DefaultService) InitService() error {
	return ds.spawnMonitorService()
}

// CloseService : 关闭服务系统
func (ds *DefaultService) CloseService() {
	if ds.monitorSrv == nil {
		return
	}

	ds.monitorSrv.Shutdown()
}

func (ds *DefaultService) spawnMonitorService() error {
	monitorEnv := util.GetEnvMap(util.GetEnvRoot(), "monitor")
	if monitorEnv != nil {
		monitorName := util.GetEnvString(monitorEnv, "name", "")
		monitorProto := util.GetEnvString(monitorEnv, "protocol", "http")
		monitorAddr := util.GetEnvString(monitorEnv, "address", "127.0.0.1")
		monitorPort := util.GetEnvString(monitorEnv, "port", "8001")
		logger.Info(0, "%s->%s://%s:%s", monitorName, monitorProto, monitorAddr, monitorPort)
		ds.monitorSrv = service.Make(monitorName, func() service.IService {
			srv := &service.MonitorService{Proto: monitorProto, Addr: monitorAddr + ":" + monitorPort}
			srv.Init()
			if monitorProto == "https" {
				srv.CertFile = util.GetEnvString(monitorEnv, "cert-file", "")
				srv.KeyFile = util.GetEnvString(monitorEnv, "key-file", "")
			}

			//var oauto2 *library.OAuth2
			oauto2Env := util.GetEnvMap(util.GetEnvRoot(), "oauth2")
			if oauto2Env != nil {
				oauto2Value := &library.OAuth2{TokenExp: util.GetEnvInt(oauto2Env, "auth-token-exp", 120),
					RefreshTokenExp:   util.GetEnvInt(oauto2Env, "auth-refresh-token-exp", 120),
					IsGenerateRefresh: util.GetEnvBoolean(oauto2Env, "auth-is-generate-refresh-token", true),
					S256Key:           util.GetEnvString(oauto2Env, "auth-signing-256-key", ""),
					AccessURI:         util.GetEnvString(oauto2Env, "auth-access-uri", "/login/access_token")}
				srv.OAuto2 = oauto2Value
			}
			srv.MakerMethod = library.NewHTTPSrvMethodJS //带JS功能分配器
			srv.Regiter = func(pid *actor.PID, m library.IHTTPSrvMethod) {
				if srv.OAuto2 != nil {
					//注册授权客户端
					clientEnv := util.GetEnvArray(oauto2Env, "auth-clients")
					if clientEnv != nil {
						for _, v := range clientEnv {
							cvalue := util.ToEnvMap(v)
							ID := util.GetEnvString(cvalue, "ID", "")
							Secret := util.GetEnvString(cvalue, "Secret", "")
							Domain := util.GetEnvString(cvalue, "Domain", "")
							UserID := util.GetEnvString(cvalue, "UserID", "")
							srv.OAuto2.RegisterClient(ID, Secret, Domain, UserID)
						}
					}

					//注册服务
					oauto2MetodEnv := util.GetEnvArray(oauto2Env, "auth-server-set")
					if oauto2MetodEnv != nil {
						//Method get|put|post
						//Pattern    方法
						//JSFile JS文件名
						for _, v := range oauto2MetodEnv {
							cvalue := util.ToEnvMap(v)
							Method := util.GetEnvString(cvalue, "Method", "")
							Pattern := util.GetEnvString(cvalue, "Pattern", "")
							JSFile := util.GetEnvString(cvalue, "JSFile", "")
							if Method == "" ||
								Pattern == "" ||
								JSFile == "" {
								continue
							}
							logger.Info(pid.ID, "OAuth2 Register Method[%s] Pattern:[%s] MapFile:[%s]", strings.ToUpper(Method), Pattern, JSFile)
							srv.OAuto2.RegisterAuth2MethodJS(m, Pattern, Method, JSFile)
						}
					}
				}
			}

			return srv
		}).(*service.MonitorService)
	}

	return nil
}
