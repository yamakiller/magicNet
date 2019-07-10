package frame

type defaultService struct {
}

func (ds *defaultService) InitService() bool {
	return true
}

func (ds *defaultService) CloseService() {

}
