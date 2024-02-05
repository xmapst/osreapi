package types

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/xmapst/osreapi/internal/dag"
	"github.com/xmapst/osreapi/internal/logx"
	"github.com/xmapst/osreapi/internal/server/config"
	"github.com/xmapst/osreapi/internal/storage/types"
	"github.com/xmapst/osreapi/internal/utils"
)

const (
	XTaskID    = "X-Task-ID"
	XTaskState = "X-Task-STATE"
)

var reg = regexp.MustCompile("[^a-zA-Z\\p{Han}0-9]")

type TaskDetailRes struct {
	Name           string    `json:"name" yaml:"Name" toml:"name"`
	Code           int64     `json:"code" yaml:"Code" toml:"code"`
	State          int64     `json:"state" yaml:"State" toml:"state"`
	Manager        string    `json:"manager" yaml:"Manager" toml:"manager"`
	Workspace      string    `json:"workspace" yaml:"Workspace" toml:"workspace"`
	Message        string    `json:"msg" yaml:"Message" toml:"message"`
	Timeout        string    `json:"timeout" yaml:"Timeout" toml:"timeout"`
	DependsOn      []string  `json:"depends_on" yaml:"DependsOn" toml:"depends_on"`
	CommandType    string    `json:"command_type" yaml:"CommandType" toml:"command_type"`
	EnvVars        []string  `json:"env_vars" yaml:"EnvVars" toml:"env_vars"`
	CommandContent string    `json:"command_content" yaml:"CommandContent" toml:"command_content"`
	Times          *StrTimes `json:"times" yaml:"Times" toml:"times"`
}

type TaskListRes struct {
	Total int         `json:"total" yaml:"Total" toml:"total"`
	Tasks []*TaskList `json:"tasks" yaml:"Tasks" toml:"tasks"`
}

type TaskList struct {
	ID        string   `json:"id" yaml:"ID" toml:"id"`
	Code      int64    `json:"code" yaml:"Code" toml:"code"`
	State     int64    `json:"state" yaml:"State" toml:"state"`
	Manager   string   `json:"manager" yaml:"Manager" toml:"manager"`
	Workspace string   `json:"workspace" yaml:"Workspace" toml:"workspace"`
	Message   string   `json:"msg" yaml:"Message" toml:"message"`
	EnvVars   []string `json:"env_vars" yaml:"EnvVars" toml:"env_vars"`
	Timeout   string   `json:"timeout" yaml:"Timeout" toml:"timeout"`
	Count     int64    `json:"count" yaml:"Count" toml:"count"`
	Times     StrTimes `json:"times" yaml:"Times" toml:"times"`
}

type StrTimes struct {
	ST string `json:"st" yaml:"ST" toml:"st"` // 开始时间
	ET string `json:"et" yaml:"ET" toml:"et"` // 结束时间
}

type Step struct {
	Name           string   `json:"name" form:"name" yaml:"Name" toml:"name" example:"script.ps1"`
	CommandType    string   `json:"command_type" form:"command_type" yaml:"CommandType" toml:"command_type" binding:"required" example:"powershell"`
	CommandContent string   `json:"command_content" form:"command_content" yaml:"CommandContent" toml:"command_content" binding:"required" example:"sleep 10"`
	EnvVars        []string `json:"env_vars" form:"env_vars" yaml:"EnvVars" toml:"env_vars" example:"env1=value1,env2=value2"`
	DependsOn      []string `json:"depends_on" form:"depends_on" yaml:"DependsOn" toml:"depends_on" example:""`
	Timeout        string   `json:"timeout" form:"timeout" yaml:"Timeout" toml:"timeout" example:"3m"`
	Notify         *Notifys `json:"notify" form:"notify" yaml:"Notify" toml:"notify"`

	TimeoutDuration time.Duration `json:"-" form:"-" yaml:"-" toml:"-"`
}

type Notifys []*Notify

type Notify struct {
	Type   string `json:"type" form:"type" yaml:"Type" toml:"type" binding:"required" example:"webhook"`
	Action string `json:"action" form:"action" yaml:"Action" toml:"action" binding:"required" example:"before"` // or after
}

type Steps []*Step

func (s Steps) Check(async bool) error {
	for _, v := range s {
		if v.CommandType == "" {
			return errors.New("key: 'Step.CommandType' Error:Field validation for 'CommandType' failed on the 'required' tag")
		}
		if v.CommandContent == "" {
			return errors.New("key: 'Step.CommandContent' Error:Field validation for 'CommandContent' failed on the 'required' tag")
		}
	}
	s.parseDuration()
	s.fixName()
	if err := s.uniqNames(); err != nil {
		logx.Errorln(err)
		return err
	}
	if !async {
		// 非编排模式,按顺序执行
		s.fixSync()
	}
	return nil
}

func (s Steps) fixName() {
	name := ksuid.New().String()
	for k, v := range s {
		v.Name = reg.ReplaceAllString(v.Name, "")
		if v.Name == "" {
			v.Name = fmt.Sprintf("%s-%d", name, k+1)
		}
		s[k].Name = v.Name
	}
}

func (s Steps) fixSync() {
	for k := range s {
		if k == 0 {
			s[k].DependsOn = nil
			continue
		}
		s[k].DependsOn = []string{s[k-1].Name}
	}
}

func (s Steps) uniqNames() error {
	counts := make(map[string]int)
	for _, v := range s {
		counts[v.Name]++
	}
	var errs []error
	for name, count := range counts {
		if count > 1 {
			errs = append(errs, fmt.Errorf("%s repeat count %d", name, count))
		}
	}
	if errs == nil {
		return nil
	}
	return fmt.Errorf("%v", errs)
}

func (s Steps) parseDuration() {
	for k, v := range s {
		timeout, err := time.ParseDuration(v.Timeout)
		if v.Timeout == "" || err != nil {
			timeout = config.App.ExecTimeOut
		}
		s[k].TimeoutDuration = timeout
	}
}

func (s Steps) GetMetaData() (res types.MetaData) {
	// check envs
	for k, v := range s {
		var _env []string
		m := utils.SliceToStrMap(v.EnvVars)
		for k, v := range m {
			_env = append(_env, fmt.Sprintf("%s=%s", k, v))
		}
		str, ok := m["HARDWARE_ID"]
		if ok && str != "" {
			res.HardWareID = str
		}
		str, ok = m["VM_INSTANCE_ID"]
		if ok && str != "" {
			res.VMInstanceID = str
		}
		s[k].EnvVars = _env
	}
	return
}

type Task struct {
	ID      string   `query:"id" json:"id" form:"id" yaml:"ID" toml:"id" example:""`
	Timeout string   `query:"timeout" json:"timeout" form:"timeout" yaml:"Timeout" toml:"timeout" example:""`
	EnvVars []string `query:"env_vars" json:"env_vars" form:"env_vars" yaml:"EnvVars" toml:"env_vars" example:""`
	Async   bool     `query:"async" json:"async" form:"async" yaml:"Async" toml:"async" example:"false"`
	Step    Steps    `json:"step" form:"step" yaml:"Step" toml:"step"`
	Notify  Notifys  `json:"notify" form:"notify" yaml:"Notify" toml:"notify"`

	TimeoutDuration time.Duration `json:"-" form:"-" yaml:"-"`
}

func (t *Task) Check() error {
	if t.Step == nil || len(t.Step) == 0 {
		return errors.New("key: 'Task.Step' Error:Field validation for 'Step' failed on the 'required' tag")
	}
	t.ID = reg.ReplaceAllString(t.ID, "")
	if t.ID == "" {
		t.ID = ksuid.New().String()
	}
	if _, err := dag.GraphManager(t.ID); err == nil {
		return errors.New("task is running")
	}
	timeout, err := time.ParseDuration(t.Timeout)
	if err == nil {
		t.TimeoutDuration = timeout
	}
	if err := t.Step.Check(t.Async); err != nil {
		return err
	}
	return nil
}

type TaskRes struct {
	URL   string `json:"url" yaml:"URL" toml:"url"`
	ID    string `json:"id" form:"id" yaml:"ID" toml:"id"`
	Count int    `json:"count" yaml:"Count" toml:"count"`
}
