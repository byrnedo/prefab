package prefab

import (
	"fmt"
	gDoc "github.com/fsouza/go-dockerclient"
	"time"
)

const (
	FtpImage = "mcreations/ftp:latest"
	ftpUser = "user"
	ftpPassword = "pass"
)

func StartFtpContainer(optsFuncs ...ConfOverrideFunc) (id string, url string) {


	var confFunc = func(baseOpts *SetupOpts) {
		baseOpts.Image = FtpImage
		baseOpts.ExposedPort = 21
		baseOpts.Envs = append(baseOpts.Envs, "FTP_USER=" + ftpUser)
		baseOpts.Envs = append(baseOpts.Envs, "FTP_PASS=" + ftpPassword)
		baseOpts.Envs = append(baseOpts.Envs, "HOST=localhost")
		baseOpts.Envs = append(baseOpts.Envs, "PASV_MIN_PORT=65000")
		baseOpts.Envs = append(baseOpts.Envs, "PASV_MAX_PORT=65005")
		baseOpts.ExtraPorts = map[gDoc.Port][]gDoc.PortBinding{
			"65000/tcp": []gDoc.PortBinding{{HostIP: "127.0.0.1",HostPort: "65000"}},
			"65001/tcp": []gDoc.PortBinding{{HostIP: "127.0.0.1",HostPort: "65001"}},
			"65002/tcp": []gDoc.PortBinding{{HostIP: "127.0.0.1",HostPort: "65002"}},
			"65003/tcp": []gDoc.PortBinding{{HostIP: "127.0.0.1",HostPort: "65003"}},
			"65004/tcp": []gDoc.PortBinding{{HostIP: "127.0.0.1",HostPort: "65004"}},
			"65005/tcp": []gDoc.PortBinding{{HostIP: "127.0.0.1",HostPort: "65005"}},
		}
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

func WaitForFtp(url string, timeout time.Duration) error {
	return WaitForPort(url, timeout)
}
