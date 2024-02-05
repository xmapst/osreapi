package types

import (
	"time"
)

const (
	TableTask = "task"
	TableStep = "step"
	TableLog  = "log"
)

type MetaData struct {
	HardWareID   string `json:"hard_ware_id" yaml:"HardWareID" toml:"hard_ware_id"`
	VMInstanceID string `json:"vm_instance_id" yaml:"VMInstanceID" toml:"vm_instance_id"`
}

type TaskState struct {
	ID       string        `json:"id" yaml:"ID" toml:"id"`
	State    int64         `json:"state" yaml:"State" toml:"state"`
	Count    int64         `json:"count" yaml:"Count" toml:"count"`
	EnvVars  []string      `json:"env_vars" yaml:"EnvVars" toml:"env_vars"`
	Timeout  time.Duration `json:"timeout" yaml:"Timeout" toml:"timeout"`
	Message  string        `json:"message" yaml:"Message" toml:"message"`
	OldState int64         `json:"old_state" yaml:"OldState" toml:"old_state"`
	MetaData MetaData      `json:"meta_data" yaml:"MetaData" toml:"meta_data"`
	Times    *Times        `json:"times" yaml:"Times" toml:"times"`
}

type TaskStepState struct {
	ID             string        `json:"id" yaml:"ID" toml:"id"`
	State          int64         `json:"state" yaml:"State" toml:"state"`
	Code           int64         `json:"code" yaml:"Code" toml:"code"`
	Message        string        `json:"message" yaml:"Message" toml:"message"`
	EnvVars        []string      `json:"env_vars" yaml:"EnvVars" toml:"env_vars"`
	Timeout        time.Duration `json:"timeout" yaml:"Timeout" toml:"timeout"`
	DependsOn      []string      `json:"depends_on" yaml:"DependsOn" toml:"depends_on"`
	CommandType    string        `json:"command_type" yaml:"CommandType" toml:"command_type"`
	CommandContent string        `json:"command_content" yaml:"CommandContent" toml:"command_content"`
	Times          Times         `json:"times" yaml:"Times" toml:"times"`
}

type Times struct {
	ST int64 `json:"st" yaml:"ST" toml:"st"` // 开始时间
	ET int64 `json:"et" yaml:"ET" toml:"et"` // 结束时间
}

type TaskStepLog struct {
	Timestamp int64  `json:"timestamp" yaml:"Timestamp" toml:"timestamp"`
	Line      int64  `json:"line" yaml:"Line" toml:"line"`
	Content   string `json:"content" yaml:"Content" toml:"content"`
}

type TaskStates []*TaskState

func (l TaskStates) Len() int           { return len(l) }
func (l TaskStates) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l TaskStates) Less(i, j int) bool { return l[i].Times.ST < l[j].Times.ST }

type TaskStepStates []*TaskStepState

func (e TaskStepStates) Len() int           { return len(e) }
func (e TaskStepStates) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e TaskStepStates) Less(i, j int) bool { return e[i].ID < e[j].ID }

type TaskStepLogs []*TaskStepLog

func (e TaskStepLogs) Len() int           { return len(e) }
func (e TaskStepLogs) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e TaskStepLogs) Less(i, j int) bool { return e[i].Line < e[j].Line }
