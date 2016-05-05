package prefab

import (
	"testing"
	"time"
)

func TestWaitForPostgres(t *testing.T) {
	if err := WaitForPostgres("postgres://userx:passx@127.0.0.1:12345/test", 3 * time.Second); err == nil {
		t.Error("Expected timeout error")
	}

}

func TestStartPostgresContainer(t *testing.T) {
	id, url := StartPostgresContainer()
	if url == "" {
		t.Error("Didn't get url")
	}
	t.Log(id)


	t.Log(time.Now())
	if err := WaitForPostgres(url, 20 * time.Second); err != nil {
		t.Error("Got error waiting on url " + url + ": " + err.Error())
	}
	t.Log(time.Now())

	Remove(id)
}
