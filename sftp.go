package prefab

import (
	"fmt"
	"time"
)

const (
	SftpImage = "byrnedo/sftp:1"
	SftpUser = "user"
	SftpPassword = "pass"
)

func StartSftpContainer(optsFuncs ...ConfOverrideFunc) (string, string) {


	var confFunc = func(baseOpts *SetupOpts) {
		baseOpts.Image = SftpImage
		baseOpts.ExposedPort = 22
		baseOpts.Command = []string{SftpUser + ":" + SftpPassword + ":1001"}
		for _, optFunc := range optsFuncs {
			optFunc(baseOpts)
		}
	}

	con, ip, port, err := startStandardContainer(confFunc)
	if err != nil {
		panic(err.Error())
	}
	return con.ID, fmt.Sprintf("%s:%d", ip, port)
}

func WaitForSftp(url string, timeout time.Duration) error {
	return WaitForPort(url, timeout)
}
