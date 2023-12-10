# RedLock
> [distributed-locks](https://redis.io/docs/manual/patterns/distributed-locks/)

透過 github.com/redis/go-redis/v9 簡單實作 RedLock

如果你需要使用到`RedLock`,請使用 [redsync](https://github.com/go-redsync/redsync)

因為更新`Database`的資料前,可能遇到`stop the world gc`,所以還需要搭配
樂觀鎖或其他方式來避免`stop the world gc`之後鎖過期被取走的問題:

```go
	rl := redlock.New(
		&redlock.Options{
			Addr: "localhost:6379",
		},
		&redlock.Options{
			Addr: "localhost:6380",
		},
	)

	m := rl.NewMutex("test1", 10*time.Second)
	err := m.Lock()
	if err != nil {
		log.Default().Println(err)
		m.Unlock()
		return
	}

	mysql := NewMysql()
	data := mysql.GetTestById(1)
	version := data.Version
	newVersion := version + 1
	mysql.UpdateTest.
		Where("id", 1).
		Where("version", version).
		Set("version", newVersion).
		Set("value", "new value").
		Exec()

	m.Unlock()
```

## Redlock分析

Martin Kleppmann的[分析](https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html)

[與此分析相反的觀點](http://antirez.com/news/101)
