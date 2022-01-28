package cache

import (
	"github.com/dgraph-io/badger/v3"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/xmapst/osreapi/utils"
	"sort"
	"time"
)

// 1 Stopped, 2 Running, 3 Pending
const (
	Stopped = iota
	Running
	Pending
)

var StateMap = map[int]string{
	Stopped: "Stopped",
	Running: "Running",
	Pending: "Pending",
}

type Val struct {
	Jobs           []*Jobs       `json:"jobs,omitempty"`
	Data           []*ExecStatus `json:"data,omitempty,omitempty"`
	State          int           `json:"state,omitempty"`
	TTL            time.Duration `json:"ttl,omitempty"`
	StartTimes     int64         `json:"start_times,omitempty"`
	CompletedTimes int64         `json:"completed_times,omitempty"`
}

type Jobs struct {
	Name           string
	CommandType    string
	CommandContent string
	DependsOn      []string
	EnvVars        []string
	ExecTimeout    time.Duration
}

type ExecStatus struct {
	Name     string `json:"name"`
	Step     int    `json:"step"`
	ExitCode int    `json:"exit_code"`
	Output   string `json:"output"`
}

var (
	db   *badger.DB
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func Init() {
	var err error
	opt := badger.DefaultOptions(utils.DataPath)
	opt.WithCompactL0OnClose(true)
	opt.WithNumCompactors(10)
	opt.WithMaxLevels(8)
	opt.WithNumGoroutines(100)
	opt.Logger = logrus.StandardLogger()
	db, err = badger.Open(opt)
	if err != nil {
		logrus.Fatalln(err)
	}
}

func Close() {
	err := db.Close()
	if err != nil {
		logrus.Error(err)
	}
}

func ListAllData() (res map[string]Val, err error) {
	res = make(map[string]Val)
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			_k := string(item.Key())
			_v, err := getItemValue(item)
			if err != nil {
				return err
			}
			var val = new(Val)
			err = json.Unmarshal(_v, val)
			if err != nil {
				return err
			}
			res[_k] = *val
		}
		return nil
	})
	return
}

type ListData struct {
	ID             string        `json:"id,omitempty"`
	State          int           `json:"state,omitempty"`
	TTL            time.Duration `json:"ttl,omitempty"`
	StartTimes     int64         `json:"start_times,omitempty"`
	CompletedTimes int64         `json:"completed_times,omitempty"`
}

type ListByStartTimes []*ListData

func (l ListByStartTimes) Len() int           { return len(l) }
func (l ListByStartTimes) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l ListByStartTimes) Less(i, j int) bool { return l[i].StartTimes < l[j].StartTimes }

func GetAllByStartTimes() (res ListByStartTimes) {
	_res, err := ListAllData()
	if err != nil {
		logrus.Error(err)
		return nil
	}
	for k, v := range _res {
		res = append(res, &ListData{
			ID:             k,
			State:          v.State,
			TTL:            v.TTL,
			StartTimes:     v.StartTimes,
			CompletedTimes: v.CompletedTimes,
		})
	}
	// sort by StartTimes
	sort.Sort(res)
	return
}

type ListByCompletedTimes []*ListData

func (l ListByCompletedTimes) Len() int           { return len(l) }
func (l ListByCompletedTimes) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l ListByCompletedTimes) Less(i, j int) bool { return l[i].CompletedTimes < l[j].CompletedTimes }

func GetAllByEndTimes() (res ListByCompletedTimes) {
	_res, err := ListAllData()
	if err != nil {
		logrus.Error(err)
		return nil
	}
	for k, v := range _res {
		res = append(res, &ListData{
			ID:             k,
			State:          v.State,
			TTL:            v.TTL,
			StartTimes:     v.StartTimes,
			CompletedTimes: v.CompletedTimes,
		})
	}
	// sort by CompletedTimes
	sort.Sort(res)
	return
}

type ListByExpiredTimes []*ListData

func (l ListByExpiredTimes) Len() int           { return len(l) }
func (l ListByExpiredTimes) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l ListByExpiredTimes) Less(i, j int) bool { return l[i].TTL < l[j].TTL }

func GetAllByTTL() (res ListByExpiredTimes) {
	_res, err := ListAllData()
	if err != nil {
		logrus.Error(err)
		return nil
	}
	for k, v := range _res {
		res = append(res, &ListData{
			ID:             k,
			State:          v.State,
			TTL:            v.TTL,
			StartTimes:     v.StartTimes,
			CompletedTimes: v.CompletedTimes,
		})
	}
	// sort by ExpiredTimes
	sort.Sort(res)
	return
}

type ExecResults []*ExecStatus

func (e ExecResults) Len() int           { return len(e) }
func (e ExecResults) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e ExecResults) Less(i, j int) bool { return e[i].Step < e[j].Step }

func Get(key string) (found bool, value *Val) {
	var err error
	var item *badger.Item
	err = db.View(func(txn *badger.Txn) error {
		item, err = txn.Get([]byte(key))
		return err
	})

	if err != nil {
		if err != badger.ErrKeyNotFound {
			logrus.Error(err)
		}
		return false, nil
	}
	var _value = new(Val)
	var val []byte
	val, err = getItemValue(item)
	err = json.Unmarshal(val, _value)
	if err != nil {
		logrus.Error(err)
		return false, nil
	}
	return true, _value
}

func getItemValue(item *badger.Item) (val []byte, err error) {
	var v []byte
	err = item.Value(func(val []byte) error {
		v = append(v, val...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return v, err
}

func Del(key string) {
	bsk := []byte(key)
	err := db.Update(func(txn *badger.Txn) error {
		return txn.Delete(bsk)
	})
	if err != nil {
		logrus.Error(err)
	}
}

func Set(key string, val *Val) {
	bsv, err := json.Marshal(val)
	if err != nil {
		logrus.Error(err)
		return
	}
	bsk := []byte(key)
	err = db.Update(func(txn *badger.Txn) error {
		if val.TTL != 0 {
			return txn.SetEntry(badger.NewEntry(bsk, bsv).WithTTL(val.TTL))
		}
		return txn.Set(bsk, bsv)
	})
	if err != nil {
		logrus.Error(err)
	}
}
