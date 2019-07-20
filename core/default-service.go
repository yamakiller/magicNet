package core

import (
	"magicNet/engine/logger"
	"magicNet/engine/util"
	"magicNet/library"
	"magicNet/service"
)

// DefaultService : 默认服务系统
type DefaultService struct {
	monitorSrv *service.MonitorService
}

// InitService : 初始化服务模块
func (ds *DefaultService) InitService() error {
	/**ds.monitorSrv = service.Make("monitor/HTTP", func() service.IService {
		srv := &service.MonitorService{}
		return srv
	}).(*service.MonitorService)**/
	logger.Info(0, "service start:")
	monitorEnv := util.GetEnvMap(util.GetEnvRoot(), "monitor")
	if monitorEnv != nil {
		monitorName := util.GetEnvString(monitorEnv, "name", "")
		monitorProto := util.GetEnvString(monitorEnv, "protocol", "http")
		monitorAddr := util.GetEnvString(monitorEnv, "address", "127.0.0.1")
		monitorPort := util.GetEnvString(monitorEnv, "port", "8001")
		logger.Info(0, "%s->%s://%s:%s", monitorName, monitorProto, monitorAddr, monitorPort)
		ds.monitorSrv = service.Make(monitorName, func() service.IService {
			srv := &service.MonitorService{Proto: monitorProto, Addr: monitorAddr + ":" + monitorPort}
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
				ds.monitorSrv.OAuto2 = oauto2Value
			}
			ds.monitorSrv.MakerMethod = library.NewHTTPSrvMethodJS //带JS功能分配器
			ds.monitorSrv.Regiter = func() {
				if ds.monitorSrv.OAuto2 != nil {
					//注册授权客户端
					clientEnv := util.GetEnvArray(oauto2Env, "auth-clients")
					if clientEnv != nil {
						for _, v := range clientEnv {
							cvalue := util.ToEnvMap(v)
							ID := util.GetEnvString(cvalue, "ID", "")
							Secret := util.GetEnvString(cvalue, "Secret", "")
							Domain := util.GetEnvString(cvalue, "Domain", "")
							UserID := util.GetEnvString(cvalue, "UserID", "")
							ds.monitorSrv.OAuto2.RegisterClient(ID, Secret, Domain, UserID)
						}
					}

					//注册服务
				}
			}

			return srv
		}).(*service.MonitorService)
	}
	return nil
}

// CloseService : 关闭服务系统
func (ds *DefaultService) CloseService() {

}

func (ds *DefaultService) spawnMonitorService() error {

	return nil
}
