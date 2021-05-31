package gosnowflake

import (
	"sync"
	"time"
)

type Snowflake struct {
	MachineId int16
	lastTime  int64
	sequence  int64

	locker sync.Locker
}

const machineBit = 9
const sequenceBit = 13

const maxMachine = (1 << (machineBit)) - 1
const maxSequence = (1 << (sequenceBit)) - 1

const sequenceLeft = 0
const machineLeft = sequenceLeft + sequenceBit
const timestampLeft = machineLeft + machineBit

// 2006/01/02 15:04:05 起始时间,雪花算法可以支持69年 大概到2075年超出精度
const startTimestamp = 1136185445000

func NewSnowFlake(machineId int) *Snowflake {
	// 取后十位
	time.Now().Format("20060102150405")
	sn := &Snowflake{MachineId: int16(machineId & maxMachine), locker: &sync.Mutex{}}
	return sn
}

func (s *Snowflake) NextId() int64 {

	defer s.locker.Unlock()
	s.locker.Lock()

	now := now()
	if now < s.lastTime {
		panic("时钟回拨,拒绝生成id")
	}

	if now == s.lastTime {
		s.sequence = s.sequence + 1&maxSequence
		if s.sequence == 0 {
			now = s.getNextMill()
			s.sequence = 0
		}
	} else {
		s.sequence = 0
	}
	s.lastTime = now
	return s.generateId(now, s.sequence)
}

func (s *Snowflake) generateId(now int64, seq int64) int64 {
	return ((now - startTimestamp) << timestampLeft) | (int64(s.MachineId) << machineLeft) | seq
}

func (s *Snowflake) getNextMill() int64 {
	n := now()
	for n = now(); !(n > s.lastTime); {
		return n
	}
	return n
}

func now() int64 {
	n := time.Now().UnixNano() / 1e6
	return n
}
