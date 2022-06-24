# reentrant_lock

```go

addr, password := "127.0.0.1:6379", ""
client := redis.NewClient(&redis.Options{
    Addr:     addr,
    Password: password,
    DB:       0,
})

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

if err := jobTool.Unlock(); err != nil {
    return err
}
	
```