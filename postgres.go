package prefab

import (
	"fmt"
	"time"
	nurl "net/url"
	"errors"
)

const (
	PostgresImage = "postgres:latest"
	PostgresUser = "postgres"
	PostgresPassword = "postgres"
	PostgresDatabase = ""

)

func StartPostgresContainer(clientOpts ...ConfOverrideFunc) (id string, url string) {

	var confFunc = func(baseOpts *SetupOpts){
		baseOpts.Image = PostgresImage
		baseOpts.ExposedPort = 5432
		baseOpts.Envs = append(baseOpts.Envs, "POSTGRES_USER="+PostgresUser)
		baseOpts.Envs = append(baseOpts.Envs, "POSTGRES_PASSWORD="+PostgresPassword)
		baseOpts.Envs = append(baseOpts.Envs, "POSTGRES_DB="+PostgresDatabase)

		for _, clientOpt := range clientOpts {
			clientOpt(baseOpts)
		}
	}

	con, ip, port, err := startStandardContainer(confFunc)
	if err != nil {
		panic(err.Error())
	}

	return con.ID, fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", PostgresUser, PostgresPassword, ip, port, PostgresDatabase)
}

func WaitForPostgres(url string, timeout time.Duration) error {
	u, err := nurl.Parse(url)
	if err != nil {
		return errors.New("Failed to parse url: " + url)
	}
	return WaitForPort(u.Host, timeout)
}
