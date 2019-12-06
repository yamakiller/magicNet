package implement

//----------------------------------
//  NetHandle |  uint64   |
//----------------------------------------------------------------------
//  8 bit(256)      |   24 bit(16 777 216)  |  32 bit   |
//---------------------------------------------------------------------
//   server id　    |    sreach id          |  socket id |
//-----------------------------------------------------------------------

//NetHandle : net handle
type NetHandle struct {
	_value uint64
}

const (
	//constNetHandleMax net handle max bit
	constNetHandleMax = 64
	//constNetHandleIDMask nethandle id mask
	constNetHandleIDMask = 0xFF
	//constNetHandleServiceIDMask gateway id mask
	constNetHandleServiceIDMask = 0xFF
	//constNetHandleSocketIDMask socket id mask
	constNetHandleSocketIDMask = 0xFFFF
	//constNetHandleGatewayIDBit  gateway id bit number
	constNetHandleServiceIDBit = 8
	//constNetHandleIDBit  handle id bit number
	constNetHandleIDBit = 24
	//constNetHandleSocketIDBit socket id bit number
	constNetHandleSocketIDBit = 32
	//constNetHandleServiceIDShift  service id shift
	constNetHandleServiceIDShift = constNetHandleMax - constNetHandleServiceIDBit
	//constNetHandleIDShift nethandle id is bit number
	constNetHandleIDShift = constNetHandleServiceIDShift - constNetHandleIDBit
)

// Generate : Generate Handle
func (slf *NetHandle) Generate(serviceID int32, handleID int32, sock int32) {
	slf._value = ((uint64(serviceID) & constNetHandleServiceIDMask) << constNetHandleServiceIDShift) |
		((uint64(handleID) & constNetHandleIDMask) << constNetHandleIDShift) |
		(uint64(sock) & constNetHandleSocketIDMask)
}

// GetServiceID :  Return to the server ID
func (slf *NetHandle) GetServiceID() int32 {
	return int32((slf._value >> constNetHandleServiceIDShift) & constNetHandleServiceIDMask)
}

// GetHandle : Returns the handle ID of the allocated resource
func (slf *NetHandle) GetHandle() int32 {
	return int32((slf._value >> constNetHandleSocketIDBit) & constNetHandleIDMask)
}

// GetSocket : Return socket ID
func (slf *NetHandle) GetSocket() int32 {
	return int32(slf._value & constNetHandleSocketIDMask)
}

// GetValue : Get Handle Value
func (slf *NetHandle) GetValue() uint64 {
	return slf._value
}

// SetValue : Set Handle Value
func (slf *NetHandle) SetValue(v uint64) {
	slf._value = v
}

// IsEmpty : is empty
func (slf *NetHandle) IsEmpty() bool {
	if slf._value == 0 {
		return true
	}
	return false
}

// Rest :
func (slf *NetHandle) Rest() {
	slf._value = 0
}
