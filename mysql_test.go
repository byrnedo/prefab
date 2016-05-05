package prefab

import "testing"

func TestStartMysqlContainer(t *testing.T) {
	id, url := StartMysqlContainer(func(cnf SetupOpts)SetupOpts{return cnf})
	if url == "" {
		t.Error("Didn't get url")
	}
	t.Log(id)
	t.Log(url)
}
