package prefab

import "fmt"

const (
	MongoImage = "mongo:3"
)


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

