package implement

//NetDataStat Network data status information
type NetDataStat struct {
	_lastTime  uint64
	_lastBytes uint64
}

//Update Update Network data status information
func (ndst *NetDataStat) Update(tts uint64, bytes uint64) {
	ndst._lastTime = tts
	ndst._lastBytes += bytes
}

//GetTime returns last time
func (ndst *NetDataStat) GetTime() uint64 {
	return ndst._lastTime
}

//GetBytes returns count bytes
func (ndst *NetDataStat) GetBytes() uint64 {
	return ndst._lastBytes
}

//NetStat network status
type NetStat struct {
	_online uint64
	_read   NetDataStat
	_write  NetDataStat
}

//UpdateRead Update read data status
func (slf *NetStat) UpdateRead(tts uint64, bytes uint64) {
	slf._read.Update(tts, bytes)
}

//UpdateWrite Update write data status
func (slf *NetStat) UpdateWrite(tts uint64, bytes uint64) {
	slf._write.Update(tts, bytes)
}

//UpdateOnline Update online time information
func (slf *NetStat) UpdateOnline(tts uint64) {
	slf._online = tts
}

//GetRead returns read status object
func (slf *NetStat) GetRead() NetDataStat {
	return slf._read
}

//GetWrite returns write status object
func (slf *NetStat) GetWrite() NetDataStat {
	return slf._write
}

//GetOnline returns online time last
func (slf *NetStat) GetOnline() uint64 {
	return slf._online
}
