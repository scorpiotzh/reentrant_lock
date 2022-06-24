package reentrant_lock

import (
	"context"
	"fmt"
	"github.com/scorpiotzh/mylog"
	"time"
)

type LockHandle interface {
	Lock() error
	Watch() error
	Unlock() error
}

type JobTool struct {
	lockHandle          LockHandle
	watchTickerDuration time.Duration
	log                 *mylog.Logger

	cancelWatchFunc context.CancelFunc
}

func NewJobTool(lockHandle LockHandle, watchTickerDuration time.Duration, log *mylog.Logger) *JobTool {
	if log == nil {
		log = mylog.NewLogger("reentrant_lock", mylog.LevelDebug)
	}
	return &JobTool{
		lockHandle:          lockHandle,
		watchTickerDuration: watchTickerDuration,
		log:                 log,
		cancelWatchFunc:     nil,
	}
}

func (t *JobTool) TryLock() (err error) {
	if t.lockHandle == nil {
		return fmt.Errorf("lockHandle is nil")
	}

	if err = t.lockHandle.Lock(); err != nil {
		return fmt.Errorf("lockHandle.Lock() err: %s", err.Error())
	}
	t.log.Info("lockHandle.Lock() ok")

	ctx, cancelFunc := context.WithCancel(context.Background())
	t.cancelWatchFunc = cancelFunc
	t.watch(ctx)

	return nil
}

func (t *JobTool) watch(ctx context.Context) {
	ticker := time.NewTicker(t.watchTickerDuration)
	count := 0
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := t.lockHandle.Watch(); err != nil {
					t.log.Error("lockHandle.Watch() err: ", err.Error())
				}
				count++
				t.log.Info("watch:", count)
			case <-ctx.Done():
				t.log.Info("watch done:", count)
				return
			}
		}
	}()
}

func (t *JobTool) Unlock() (err error) {
	if t.cancelWatchFunc != nil {
		t.cancelWatchFunc()
	}
	if err := t.lockHandle.Unlock(); err != nil {
		return fmt.Errorf("lockHandle.Unlock() err: %s", err.Error())
	}
	t.log.Info("lockHandle.Unlock() ok")

	return nil
}
