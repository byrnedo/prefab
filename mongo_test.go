package prefab

import (
	"testing"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func TestSetupMongoContainer(t *testing.T) {
	id, url := StartMongoContainer()
	t.Log(url)
	if url == "" {
		t.Fatal("url empty")
	}
	if id == "" {
		t.Fatal("id empty")
	}

	if err := WaitForMongo(url, 10 * time.Second); err != nil {
		t.Fatal(err.Error())
	}

	ses, err := mgo.Dial(url)
	if err != nil {
		t.Fatal(err)
	}
	if ses == nil {
		t.Fatal("nil session")
	}

	if err := ses.Copy().DB("test").C("some-collection").Insert(bson.M{"hey":"you"}); err != nil {
		t.Fatal(err)
	}

	Remove(id)
}

func TestSetupMongoTmpfsContainer(t *testing.T) {
	id, url := StartMongoTmpfsContainer()
	t.Log(url)
	if url == "" {
		t.Fatal("url empty")
	}
	if id == "" {
		t.Fatal("id empty")
	}

	if err := WaitForMongo(url, 10 * time.Second); err != nil {
		t.Fatal(err.Error())
	}

	ses, err := mgo.Dial(url)
	if err != nil {
		t.Fatal(err)
	}
	if ses == nil {
		t.Fatal("nil session")
	}

	if err := ses.Copy().DB("test").C("some-collection").Insert(bson.M{"hey":"you"}); err != nil {
		t.Fatal(err)
	}

	Remove(id)
}
