package core

//DefaultService doc
//@type default Initialization
type DefaultService struct {
}

//InitService doc
//@Method InitService doc
//@Return (error)
func (slf *DefaultService) InitService() error {
	return nil
}

//CloseService doc
//@Method CloseService doc
func (slf *DefaultService) CloseService() {
}
