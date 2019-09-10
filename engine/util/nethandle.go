package util

//----------------------------------
//  NetHandle |  uint64   |
//----------------------------------------------------------------------
//  7 bit(128)          |  9 bit(512)   |    16 bit(65535)  |  32 bit   |
//---------------------------------------------------------------------
//  gateway server id　 |  world id     |     sreach id     |  socket id |
//-----------------------------------------------------------------------

//NetHandle : net handle
type NetHandle struct {
	value uint64
}

const (
	//constNetHandleMax net handle max bit
	constNetHandleMax = 64
	//constNetHandleIDMask nethandle id mask
	constNetHandleIDMask = 0xFF
	//constNetHandleGatewayIDMask gateway id mask
	constNetHandleGatewayIDMask = 0x7F
	//constNetHandleWorldIDMask world id mask
	constNetHandleWorldIDMask = 0x1FF
	//constNetHandleSocketIDMask socket id mask
	constNetHandleSocketIDMask = 0xFFFF
	//constNetHandleGatewayIDBit  gateway id bit number
	constNetHandleGatewayIDBit = 7
	//constNetHandleWorldIDBit  world id bit number
	constNetHandleWorldIDBit = 9
	//constNetHandleIDBit  handle id bit number
	constNetHandleIDBit = 16
	//constNetHandleSocketIDBit socket id bit number
	constNetHandleSocketIDBit = 32
	//constNetHandleGatewayIDShift  gateway id shift
	constNetHandleGatewayIDShift = constNetHandleMax - constNetHandleGatewayIDBit
	//constNetHandleWorldIDShift world id is bit number
	constNetHandleWorldIDShift = constNetHandleGatewayIDShift - constNetHandleWorldIDBit
	//constNetHandleIDShift nethandle id is bit number
	constNetHandleIDShift = constNetHandleWorldIDShift - constNetHandleIDBit
)

// Generate : Generate Handle
func (nh *NetHandle) Generate(gatewayID int32, worldID int32, handleID int32, sock int32) {
	nh.value = ((uint64(gatewayID) & constNetHandleGatewayIDMask) << constNetHandleGatewayIDShift) &
		((uint64(worldID) & constNetHandleWorldIDMask) << constNetHandleWorldIDShift) &
		((uint64(handleID) & constNetHandleIDMask) << constNetHandleIDShift) &
		(uint64(sock) & constNetHandleSocketIDMask)
}

// GatewayID : Get Gateway ID
func (nh *NetHandle) GatewayID() int32 {
	return int32((nh.value >> constNetHandleGatewayIDShift) & constNetHandleGatewayIDMask)
}

// WorldID ： Get World ID
func (nh *NetHandle) WorldID() int32 {
	return int32((nh.value >> constNetHandleWorldIDShift) & constNetHandleWorldIDMask)
}

// HandleID : Get Handle ID
func (nh *NetHandle) HandleID() int32 {
	return int32((nh.value >> constNetHandleSocketIDBit) & constNetHandleIDMask)
}

// SocketID : Get Socket ID
func (nh *NetHandle) SocketID() int32 {
	return int32(nh.value & constNetHandleSocketIDMask)
}

// GetValue : Get Handle Value
func (nh *NetHandle) GetValue() uint64 {
	return nh.value
}

// IsEmpty : is empty
func (nh *NetHandle) IsEmpty() bool {
	if nh.value == 0 {
		return true
	}
	return false
}

// Rest :
func (nh *NetHandle) Rest() {
	nh.value = 0
}
