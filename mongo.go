package prefab

import (
	"fmt"
	"time"
	nurl "net/url"
	"errors"
)

const (
	MongoImage = "mongo:3"
	MongoTmpfsImage = "byrnedo/mongo-tmpfs"
)

func StartMongoTmpfsContainer(optsFuncs ...ConfOverrideFunc) (id string, url string) {


	var confFunc = func(baseOpts *SetupOpts) {
		baseOpts.Image = MongoTmpfsImage
		baseOpts.ExposedPort = 27017
		baseOpts.Privileged = true
		for _, optFunc := range optsFuncs {
			optFunc(baseOpts)
		}
	}

	con, ip, port, err := startStandardContainer(confFunc)
	if err != nil {
		panic(err.Error())
	}
	return con.ID, fmt.Sprintf("mongodb://%s:%d/", ip, port)
}


func StartMongoContainer(optsFuncs ...ConfOverrideFunc) (id string, url string) {


	var confFunc = func(baseOpts *SetupOpts) {
		baseOpts.Image = MongoImage
		baseOpts.ExposedPort = 27017
		for _, optFunc := range optsFuncs {
			optFunc(baseOpts)
		}
	}

	con, ip, port, err := startStandardContainer(confFunc)
	if err != nil {
		panic(err.Error())
	}
	return con.ID, fmt.Sprintf("mongodb://%s:%d/", ip, port)
}

func WaitForMongo(url string, timeout time.Duration) error {
	u, err := nurl.Parse(url)
	if err != nil {
		return errors.New("Failed to parse url: " + url)
	}
	return WaitForPort(u.Host, timeout)
}

