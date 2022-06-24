package reentrant_lock

import (
	"fmt"
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func TestJobTool(t *testing.T) {

	addr, password := "127.0.0.1:6379", ""
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	jobFunc := func(index int) error {
		lockHandle := LockHandleRedis{
			client:   client,
			key:      "test-key",
			lockTime: time.Second * 20,
		}

		jobTool := NewJobTool(&lockHandle, time.Second*5, nil)
		if err := jobTool.TryLock(); err != nil {
			return err
		}

		// todo do your job
		for i := 0; i < 10; i++ {
			fmt.Println("do job:", index, i)
			time.Sleep(time.Second)
		}

		if err := jobTool.Unlock(); err != nil {
			return err
		}
		return nil
	}

	i := 0
	for {
		i++
		go func() {
			if err := jobFunc(i); err != nil {
				fmt.Println("jobFunc err:", err.Error())
			}
		}()
		time.Sleep(time.Second * 2)
	}

}
