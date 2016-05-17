package prefab

import (
	"fmt"
	"time"
	"strings"
)

const (
	mysqlImage = "mysql:latest"
	mysqlTmpfsImage = "theasci/docker-mysql-tmpfs:latest"
	mysqlUser = "user"
	mysqlPassword = "pass"
	mysqlRootPassword = "toor"
	mysqlDatabase = "test"

)

func StartMysqlContainer(clientOpts ...ConfOverrideFunc) (id string, url string) {

	var confFunc = func(baseOpts *SetupOpts){
		baseOpts.Image = mysqlImage
		baseOpts.ExposedPort = 3306
		baseOpts.Envs = append(baseOpts.Envs, "MYSQL_ROOT_PASSWORD="+ mysqlRootPassword)
		baseOpts.Envs = append(baseOpts.Envs, "MYSQL_PASSWORD="+ mysqlPassword)
		baseOpts.Envs = append(baseOpts.Envs, "MYSQL_USER="+ mysqlUser)
		baseOpts.Envs = append(baseOpts.Envs, "MYSQL_DATABASE="+ mysqlDatabase)

		for _, clientOpt := range clientOpts {
			clientOpt(baseOpts)
		}
	}

	con, ip, port, err := startStandardContainer(confFunc)
	if err != nil {
		panic(err.Error())
	}

	return con.ID, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mysqlUser, mysqlPassword, ip, port, mysqlDatabase)
}

func StartMysqlTmpfsContainer(clientOpts ...ConfOverrideFunc) (id string, url string) {

	var confFunc = func(baseOpts *SetupOpts){
		baseOpts.Image = mysqlTmpfsImage
		baseOpts.ExposedPort = 3306
		baseOpts.Privileged = true
		for _, clientOpt := range clientOpts {
			clientOpt(baseOpts)
		}
	}
	con, ip, port, err := startStandardContainer(confFunc)
	if err != nil {
		panic(err.Error())
	}

	return con.ID, fmt.Sprintf("testrunner:testrunner@tcp(%s:%d)/", ip, port)
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
