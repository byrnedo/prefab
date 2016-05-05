package prefab

import (
	"testing"
	"github.com/apcera/nats"
	"sync"
	"time"
)

func TestStartNatsContainer(t *testing.T) {
	id, url := StartNatsContainer()
	t.Log(url)
	if url == "" {
		t.Error("url empty")
	}
	if id == "" {
		t.Error("id empty")
	}

	if err := WaitForNats(url, 10 *time.Second) ; err != nil {
		t.Error(err)
	}

	con, err := nats.Connect(url)
	if err != nil {
		t.Error(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	con.Subscribe("test", func(m *nats.Msg){
		if string(m.Data) != "well hello" {
			t.Error("Unexpected response", string(m.Data))
		}
		wg.Done()

	})
	con.Publish("test", []byte("well hello"))

	wg.Wait()

	Remove(id)
}
