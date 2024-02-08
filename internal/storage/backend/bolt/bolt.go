package bolt

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"go.etcd.io/bbolt"

	"github.com/xmapst/osreapi/internal/logx"
	"github.com/xmapst/osreapi/internal/storage/backend"
	"github.com/xmapst/osreapi/internal/storage/types"
)

type Bolt struct {
	db *bbolt.DB
}

func (b *Bolt) Get(bucket, key string) (value []byte, err error) {
	err = b.db.View(func(tx *bbolt.Tx) error {
		if buk := tx.Bucket([]byte(bucket)); buk != nil {
			value = buk.Get([]byte(key))
		}
		return nil
	})
	if value == nil {
		return nil, bbolt.ErrBucketNotFound
	}
	return
}

func (b *Bolt) Set(bucket, key string, value []byte) (err error) {
	err = b.db.Update(func(tx *bbolt.Tx) error {
		buk, e := tx.CreateBucketIfNotExists([]byte(bucket))
		if e != nil {
			return e
		}
		err = buk.Put([]byte(key), value)
		return err
	})
	return
}

func (b *Bolt) Del(bucket, key string) (err error) {
	err = b.db.Update(func(tx *bbolt.Tx) error {
		if buk := tx.Bucket([]byte(bucket)); buk != nil {
			return buk.Delete([]byte(key))
		}
		return nil
	})

	return
}

func (b *Bolt) Prefix(bucket, prefix string) (values [][]byte, err error) {
	err = b.db.View(func(tx *bbolt.Tx) error {
		if buk := tx.Bucket([]byte(bucket)); buk != nil {
			c := buk.Cursor()
			for key, val := c.Seek([]byte(prefix)); key != nil && bytes.HasPrefix(key, []byte(prefix)); key, val = c.Next() {
				values = append(values, backend.SafeCopy(nil, val))
			}
		}

		return nil
	})
	if values == nil {
		return nil, bbolt.ErrBucketNotFound
	}
	return
}

func (b *Bolt) Suffix(bucket, suffix string) (values [][]byte, err error) {
	err = b.db.View(func(tx *bbolt.Tx) error {
		if buk := tx.Bucket([]byte(bucket)); buk != nil {
			err = buk.ForEach(func(k, v []byte) error {
				if bytes.HasSuffix(k, []byte(suffix)) {
					values = append(values, backend.SafeCopy(nil, v))
				}
				return nil
			})
		}

		return nil
	})
	if values == nil {
		return nil, bbolt.ErrBucketNotFound
	}
	return
}

func (b *Bolt) Range(bucket, start, limit string) (values [][]byte, err error) {
	err = b.db.View(func(tx *bbolt.Tx) error {
		if buk := tx.Bucket([]byte(bucket)); buk != nil {
			c := buk.Cursor()
			_start := filepath.ToSlash(filepath.Join(bucket, start))
			_limit := filepath.ToSlash(filepath.Join(bucket, limit))
			for k, v := c.Seek([]byte(_start)); k != nil && bytes.Compare([]byte(_start), k) <= 0; k, v = c.Next() {
				if bytes.Compare([]byte(_limit), k) > 0 {
					values = append(values, backend.SafeCopy(nil, v))
				} else {
					break
				}
			}
		}

		return nil
	})
	if values == nil {
		return nil, bbolt.ErrBucketNotFound
	}
	return
}

func (b *Bolt) BatchSet(bucket string, kvs map[string][]byte) error {
	if err := b.db.Update(func(tx *bbolt.Tx) (err error) {
		_, err = tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	}); err != nil {
		return err
	}
	return b.db.Batch(func(tx *bbolt.Tx) (err error) {
		buk := tx.Bucket([]byte(bucket))
		for k, v := range kvs {
			if err = buk.Put([]byte(k), v); err != nil {
				return err
			}
		}
		return
	})
}

func (b *Bolt) Name() string {
	return "bolt"
}

func (b *Bolt) Close() error {
	logx.Infoln("wait for data removal to complete")
	return b.db.Close()
}

func New(path string) (*Bolt, error) {
	db, err := bbolt.Open(filepath.Join(path, "database.db"), os.ModePerm, &bbolt.Options{})
	if err != nil {
		return nil, err
	}

	b := &Bolt{
		db: db,
	}
	return b, nil
}

func (b *Bolt) TaskList(table, prefix string) (res types.TaskStates, err error) {
	if prefix != "" {
		prefix = filepath.ToSlash(filepath.Join(prefix, types.TableTask))
	}
	val, err := b.Prefix(table, prefix)
	if err != nil {
		if errors.Is(err, bbolt.ErrBucketNotFound) {
			return nil, backend.ErrNotExist
		}
		logx.Warnln(err)
		return nil, err
	}
	for _, v := range val {
		var state = new(types.TaskState)
		var data = bytes.NewReader(v)
		err = gob.NewDecoder(data).Decode(state)
		if err != nil {
			logx.Warnln(err)
			continue
		}
		res = append(res, state)
	}
	sort.Sort(res)
	return
}

func (b *Bolt) TaskDetail(table, task string) (res *types.TaskState, err error) {
	val, err := b.Get(table, filepath.ToSlash(filepath.Join(task, types.TableTask)))
	if err != nil {
		if errors.Is(err, bbolt.ErrBucketNotFound) {
			return nil, backend.ErrNotExist
		}
		logx.Warnln(err)
		return
	}
	res = new(types.TaskState)
	var data = bytes.NewReader(val)
	err = gob.NewDecoder(data).Decode(res)
	if err != nil {
		logx.Warnln(err)
		return
	}
	return
}

func (b *Bolt) SetTask(table, task string, val *types.TaskState) error {
	var data bytes.Buffer
	err := gob.NewEncoder(&data).Encode(val)
	if err != nil {
		logx.Warnln(err)
		return err
	}
	return b.Set(table, filepath.ToSlash(filepath.Join(task, types.TableTask)), data.Bytes())
}

func (b *Bolt) TaskStepList(table, task string) (res types.TaskStepStates, err error) {
	val, err := b.Prefix(table, filepath.ToSlash(filepath.Join(task, types.TableTask)))
	if err != nil {
		if errors.Is(err, bbolt.ErrBucketNotFound) {
			return nil, backend.ErrNotExist
		}
		logx.Warnln(err)
		return nil, err
	}
	for _, v := range val {
		var state = new(types.TaskStepState)
		var data = bytes.NewReader(v)
		err = gob.NewDecoder(data).Decode(state)
		if err != nil {
			logx.Warnln(err)
			continue
		}
		res = append(res, state)
	}
	sort.Sort(res)
	return
}

func (b *Bolt) TaskStepDetail(table, task, step string) (res *types.TaskStepState, err error) {
	val, err := b.Get(table, filepath.ToSlash(filepath.Join(task, types.TableTask, step, types.TableStep)))
	if err != nil {
		if errors.Is(err, bbolt.ErrBucketNotFound) {
			return nil, backend.ErrNotExist
		}
		logx.Warnln(err)
		return
	}
	res = new(types.TaskStepState)
	var data = bytes.NewReader(val)
	err = gob.NewDecoder(data).Decode(res)
	if err != nil {
		logx.Warnln(err)
		return
	}
	return
}

func (b *Bolt) SetTaskStep(table, task, step string, val *types.TaskStepState) error {
	var data bytes.Buffer
	err := gob.NewEncoder(&data).Encode(val)
	if err != nil {
		logx.Warnln(err)
		return err
	}
	return b.Set(table, filepath.ToSlash(filepath.Join(task, types.TableTask, step, types.TableStep)), data.Bytes())
}

func (b *Bolt) TaskStepLogList(table, task, step string) (res types.TaskStepLogs, err error) {
	val, err := b.Prefix(table, filepath.ToSlash(filepath.Join(task, types.TableTask, step, types.TableLog)))
	if err != nil {
		if errors.Is(err, bbolt.ErrBucketNotFound) {
			return nil, backend.ErrNotExist
		}
		logx.Warnln(err)
		return nil, err
	}
	for _, v := range val {
		var state = new(types.TaskStepLog)
		var data = bytes.NewReader(v)
		err = gob.NewDecoder(data).Decode(state)
		if err != nil {
			logx.Warnln(err)
			continue
		}
		res = append(res, state)
	}
	sort.Sort(res)
	return
}

func (b *Bolt) SetTaskStepLog(table, task, step string, line int64, val *types.TaskStepLog) error {
	var data bytes.Buffer
	err := gob.NewEncoder(&data).Encode(val)
	if err != nil {
		logx.Warnln(err)
		return err
	}
	return b.Set(table, filepath.ToSlash(filepath.Join(task, types.TableTask, step, types.TableLog, strconv.FormatInt(line, 10))), data.Bytes())
}