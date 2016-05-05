package prefab

import (
	"testing"
	"time"
)

func TestWaitForMysql(t *testing.T) {
	if err := WaitForMysql("userx:passx@tcp(127.0.0.1:12345)/test", 3 * time.Second); err == nil {
		t.Error("Expected timeout error")
	}

}

func TestStartMysqlContainer(t *testing.T) {
	id, url := StartMysqlContainer()
	if url == "" {
		t.Error("Didn't get url")
	}
	t.Log(id)


	t.Log(time.Now())
	if err := WaitForMysql(url, 20 * time.Second); err != nil {
		t.Error("Got error waiting on url " + url + ": " + err.Error())
	}
	t.Log(time.Now())

	Remove(id)
}
