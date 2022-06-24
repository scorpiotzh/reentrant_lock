package reentrant_lock

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type LockHandleRedis struct {
	client   *redis.Client
	key      string
	lockTime time.Duration
}

var ErrDistributedLockPreemption = errors.New("distributed lock preemption")

func (j *LockHandleRedis) Lock() error {
	if ok, err := j.client.SetNX(j.key, "", j.lockTime).Result(); err != nil {
		return fmt.Errorf("redis set nx err: %s", err.Error())
	} else if !ok {
		return ErrDistributedLockPreemption
	}
	return nil
}

func (j *LockHandleRedis) Watch() error {
	return j.client.Expire(j.key, j.lockTime).Err()
}

func (j *LockHandleRedis) Unlock() error {
	if err := j.client.Del(j.key).Err(); err != nil {
		return fmt.Errorf("redis del err: %s", err.Error())
	}
	return nil
}
