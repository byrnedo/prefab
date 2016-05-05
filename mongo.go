package prefab

import "fmt"

const (
	MongoImage = "mongo:3"
)


func StartMongoContainer(clientOpts func(SetupOpts)SetupOpts) (id string, url string) {

	var confFunc = func(baseOpts SetupOpts)SetupOpts{
		baseOpts.Image = MongoImage
		baseOpts.ExposedPort = 27017
		return baseOpts
	}

	con, ip, port, err := startStandardContainer(confFunc)
	if err != nil {
		panic(err.Error())
	}
	return con.ID, fmt.Sprintf("mongodb://%s:%d/", ip, port)
}

