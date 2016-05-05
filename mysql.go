package prefab

import "fmt"

const (
	MysqlImage = "mysql:latest"
	MysqlUser = "user"
	MysqlPassword = "pass"
	MysqlRootPassword = "toor"
	MysqlDatabase = "test"

)

func StartMysqlContainer(clientOpts func(SetupOpts)SetupOpts) (id string, url string) {

	var confFunc = func(baseOpts SetupOpts)SetupOpts{
		baseOpts.Image = MysqlImage
		baseOpts.ExposedPort = 3306
		baseOpts.Envs = append(baseOpts.Envs, "MYSQL_ROOT_PASSWORD="+MysqlRootPassword)
		baseOpts.Envs = append(baseOpts.Envs, "MYSQL_PASSWORD="+MysqlPassword)
		baseOpts.Envs = append(baseOpts.Envs, "MYSQL_USER="+MysqlUser)
		baseOpts.Envs = append(baseOpts.Envs, "MYSQL_DATABASE="+MysqlDatabase)
		return baseOpts
	}

	con, ip, port, err := startStandardContainer(confFunc)
	if err != nil {
		panic(err.Error())
	}

	return con.ID, fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s", MysqlUser, MysqlPassword, ip, port, MysqlDatabase)
}
