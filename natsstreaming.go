package prefab

import (
	"fmt"
	"time"
	"errors"
	nurl "net/url"
)

const (
	NatsStreamingImage = "mattmastersitb/nats-streaming-server:latest"
)

func StartNatsStreamingContainer(clientOpts ...ConfOverrideFunc) (id string, url string) {

	var confFunc = func(baseOpts *SetupOpts){
		baseOpts.Image = NatsStreamingImage
		baseOpts.ExposedPort = 4222
		for _, clientOpt := range clientOpts {
			clientOpt(baseOpts)
		}
	}

	con, ip, port, err := startStandardContainer(confFunc)
	if err != nil {
		panic(err.Error())
	}
	return con.ID, fmt.Sprintf("nats://%s:%d/", ip, port)
}

func WaitForNatsStreaming(url string, timeout time.Duration) error {
	u, err := nurl.Parse(url)
	if err != nil {
		return errors.New("Failed to parse url: " + url)
	}
	return WaitForPort(u.Host, timeout)
}
