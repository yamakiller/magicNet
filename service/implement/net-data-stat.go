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
func (nst *NetStat) UpdateRead(tts uint64, bytes uint64) {
	nst.read.Update(tts, bytes)
}

//UpdateWrite Update write data status
func (nst *NetStat) UpdateWrite(tts uint64, bytes uint64) {
	nst.write.Update(tts, bytes)
}

//UpdateOnline Update online time information
func (nst *NetStat) UpdateOnline(tts uint64) {
	nst.online = tts
}

//GetRead returns read status object
func (nst *NetStat) GetRead() NetDataStat {
	return nst.read
}

//GetWrite returns write status object
func (nst *NetStat) GetWrite() NetDataStat {
	return nst.write
}

//GetOnline returns online time last
func (nst *NetStat) GetOnline() uint64 {
	return nst.online
}
