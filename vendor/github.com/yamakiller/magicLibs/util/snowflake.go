package util

import (
	"fmt"
	"sync"
	"time"
)

const (
	constWorkerIDBits     int64 = 5
	constDatacenterIDBits int64 = 5
	constSequenceBits     int64 = 12

	constMaxWorkerID     int64 = -1 ^ (-1 << uint64(constWorkerIDBits))
	constMaxDatacenterID int64 = -1 ^ (-1 << uint64(constDatacenterIDBits))
	constMaxSequence     int64 = -1 ^ (-1 << uint64(constSequenceBits))

	timeLeft uint8 = 22
	dataLeft uint8 = 17
	workLeft uint8 = 12

	twepoch int64 = 1525705533000
)

//NewSnowFlake doc
//@Summary new snowflake object
//@Method NewSnowFlake
//@Param int64 worker id
//@Param int64 id
//@Return *SnowFlake
func NewSnowFlake(workerID int64, id int64) *SnowFlake {
	return &SnowFlake{_laststamp: -1, _workerid: workerID, _datacenterid: id, _sequence: 1}
}

//SnowFlake doc
//@Summary snowflake object
//@Struct SnowFlake
//@Member int64 last stamp
//@Member int64 worker id
//@Member int64 data center id
//@Member int64 sequence
//@Member Mutex
type SnowFlake struct {
	_laststamp    int64
	_workerid     int64
	_datacenterid int64
	_sequence     int64
	_sync         sync.Mutex
}

//NextID doc
//@Summary spawn id
//@Method NextID
//@Return int64 id
//@Return error
func (slf *SnowFlake) NextID() (int64, error) {
	slf._sync.Lock()
	defer slf._sync.Unlock()

	timestamp := slf.getCurrentTime()
	if timestamp < slf._laststamp {
		return 0, fmt.Errorf("time stamp abnormal")
	}

	if slf._laststamp == timestamp {
		slf._sequence = (slf._sequence + 1) & constMaxSequence
		if slf._sequence == 0 {
			for timestamp <= slf._laststamp {
				timestamp = slf.getCurrentTime()
			}
		}
	} else {
		slf._sequence = 0
	}
	slf._laststamp = timestamp
	return ((timestamp - twepoch) << timeLeft) | (slf._datacenterid << dataLeft) | (slf._workerid << workLeft) | slf._sequence, nil
}

func (slf *SnowFlake) getCurrentTime() int64 {
	return time.Now().UnixNano() / 1e6
}
