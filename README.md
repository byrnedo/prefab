# Prefab

Quickly run docker containers for infrastructure when testing.

```
import "github.com/byrnedo/prefab"

// create and run a mysql container for this session
id, url := prefab.StartMysqlContainer()

// helper to wait for port to open
if err := prefab.WaitForMysql(url, 20 * time.Second); err != nil {
    panic(err)
}

// Connect using url :D

prefab.Remove(id)
```

Supports:

- Mysql
- Postgresql
- Mongo
- Nats

Although it's easy to roll your own.


Pull requests for new containers very welcome :)
