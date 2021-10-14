package florm

import (
	"time"

	"github.com/godruoyi/go-snowflake"
)

const timestampMoveLength = 22
const machineIDMoveLength = 12

type SID snowflake.SID

func (id *SID) CalculateID() uint64 {
	id.ID = uint64((id.Timestamp << timestampMoveLength) | (id.MachineID << machineIDMoveLength) | uint64(id.Sequence))
	return id.ID
}

func (id *SID) SetTime(t time.Time) {
	st := sfkCfg.StartTime.UTC().UnixNano() / 1e6
	tst := t.UnixNano()/1e6 - st

	id.Timestamp = uint64(tst)
}

type SnowflakeCfg struct {
	StartTime   time.Time `json:"start_time"`
	UseMacineID bool      `json:"use_machine_id"`
}

var sfkCfg SnowflakeCfg = SnowflakeCfg{
	StartTime:   time.Unix(0, 0),
	UseMacineID: true,
}

func ConfigSnowflake(cfg SnowflakeCfg) {
	sfkCfg = cfg

	snowflake.SetStartTime(cfg.StartTime)

	if cfg.UseMacineID {
		snowflake.SetMachineID(snowflake.PrivateIPToMachineID())
	}
}
