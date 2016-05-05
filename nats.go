package prefab

import "fmt"

const (
	NatsImage = "nats:latest"
)

func StartNatsContainer(clientOpts func(SetupOpts)SetupOpts) (id string, url string) {

	var confFunc = func(baseOpts SetupOpts)SetupOpts{
		baseOpts.Image = NatsImage
		baseOpts.ExposedPort = 4222
		return baseOpts
	}

	con, ip, port, err := startStandardContainer(confFunc)
	if err != nil {
		panic(err.Error())
	}
	return con.ID, fmt.Sprintf("nats://%s:%d/", ip, port)
}
