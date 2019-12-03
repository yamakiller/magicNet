package core

//DefaultService desc
//@type default Initialization
type DefaultService struct {
}

//InitService desc
//@Method InitService desc
//@Return (error)
func (slf *DefaultService) InitService() error {
	return nil
}

//CloseService desc
//@Method CloseService desc
func (slf *DefaultService) CloseService() {
}
