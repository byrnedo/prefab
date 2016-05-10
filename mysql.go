package prefab

import (
	"fmt"
	"time"
	"strings"
)

const (
	MysqlImage = "mysql:latest"
	MysqlUser = "user"
	MysqlPassword = "pass"
	MysqlRootPassword = "toor"
	MysqlDatabase = "test"

)

func StartMysqlContainer(clientOpts ...ConfOverrideFunc) (id string, url string) {

	var confFunc = func(baseOpts *SetupOpts){
		baseOpts.Image = MysqlImage
		baseOpts.ExposedPort = 3306
		baseOpts.Envs = append(baseOpts.Envs, "MYSQL_ROOT_PASSWORD="+MysqlRootPassword)
		baseOpts.Envs = append(baseOpts.Envs, "MYSQL_PASSWORD="+MysqlPassword)
		baseOpts.Envs = append(baseOpts.Envs, "MYSQL_USER="+MysqlUser)
		baseOpts.Envs = append(baseOpts.Envs, "MYSQL_DATABASE="+MysqlDatabase)

		for _, clientOpt := range clientOpts {
			clientOpt(baseOpts)
		}
	}

	con, ip, port, err := startStandardContainer(confFunc)
	if err != nil {
		panic(err.Error())
	}

	return con.ID, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", MysqlUser, MysqlPassword, ip, port, MysqlDatabase)
}

func WaitForMysql(url string, timeout time.Duration) error {

	var (
		addr string
		addrStart int
		addrEnd int
	)

	if atInd:= strings.Index(url, "@"); atInd > 0 {
		addrStart = atInd + 1
	}

	if slashInd := strings.Index(url, "/"); slashInd > 0 {
		addrEnd = slashInd
	}

	addr = url[addrStart:addrEnd]

	return WaitForPort(addr, timeout)
}
