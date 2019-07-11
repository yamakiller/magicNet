package frame

// DefaultService : 默认服务系统
type DefaultService struct {
}

// InitService : 初始化服务模块
func (ds *DefaultService) InitService() error {
	return nil
}

// CloseService : 关闭服务系统
func (ds *DefaultService) CloseService() {

}
