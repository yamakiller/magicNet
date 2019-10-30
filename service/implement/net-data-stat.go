package implement

//NetDataStat Network data status information
type NetDataStat struct {
	lastTime  uint64
	lastBytes uint64
}

//Update Update Network data status information
func (ndst *NetDataStat) Update(tts uint64, bytes uint64) {
	ndst.lastTime = tts
	ndst.lastBytes += bytes
}

//GetTime returns last time
func (ndst *NetDataStat) GetTime() uint64 {
	return ndst.lastTime
}

//GetBytes returns count bytes
func (ndst *NetDataStat) GetBytes() uint64 {
	return ndst.lastBytes
}

//NetStat network status
type NetStat struct {
	online uint64
	read   NetDataStat
	write  NetDataStat
}

//UpdateRead Update read data status
func (slf *NetStat) UpdateRead(tts uint64, bytes uint64) {
	slf.read.Update(tts, bytes)
}

//UpdateWrite Update write data status
func (slf *NetStat) UpdateWrite(tts uint64, bytes uint64) {
	slf.write.Update(tts, bytes)
}

//UpdateOnline Update online time information
func (slf *NetStat) UpdateOnline(tts uint64) {
	slf.online = tts
}

//GetRead returns read status object
func (slf *NetStat) GetRead() NetDataStat {
	return slf.read
}

//GetWrite returns write status object
func (slf *NetStat) GetWrite() NetDataStat {
	return slf.write
}

//GetOnline returns online time last
func (slf *NetStat) GetOnline() uint64 {
	return slf.online
}
