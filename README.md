# Prefab

Quickly run docker containers for infrastructure when testing.

```
import "github.com/byrnedo/prefab"

id, url := StartMysqlContainer(func(c SetupOpts)SetupOpts{return c})

// Connect using url :D
```

Supports:

- Mysql
- Mongo
- Nats

Although it's easy to roll your own.


Pull requests for new containers very welcome :)
