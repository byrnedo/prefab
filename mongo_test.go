package prefab

import (
	"testing"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TestSetupMongoContainer(t *testing.T) {
	id, url := StartMongoContainer(func(c SetupOpts)SetupOpts{return c})
	t.Log(url)
	if url == "" {
		t.Error("url empty")
	}
	if id == "" {
		t.Error("id empty")
	}

	ses, err := mgo.Dial(url)
	if err != nil {
		t.Error(err)
	}
	if ses == nil {
		t.Error("nil session")
	}

	if err := ses.Copy().DB("test").C("some-collection").Insert(bson.M{"hey":"you"}); err != nil {
		t.Error(err)
	}

	Remove(id)
}