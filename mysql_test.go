package prefab

import (
	"testing"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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
		t.Fatal("Got error waiting on url " + url + ": " + err.Error())
	}
	t.Log(time.Now())

	db, err := sql.Open("mysql", url)
	if err != nil {
		t.Fatal(err)
	}

	r, err := db.Query("SELECT 1")
	if err != nil {
		t.Fatal(err)
	}
	r.Close()
	db.Close()

	Remove(id)
}

func TestStartMysqlTmpfsContainer(t *testing.T) {
	id, url := StartMysqlTmpfsContainer()
	if url == "" {
		t.Fatal("Didn't get url")
	}
	t.Log(id)


	t.Log(time.Now())
	if err := WaitForMysql(url, 20 * time.Second); err != nil {
		t.Fatal("Got error waiting on url " + url + ": " + err.Error())
	}
	t.Log(time.Now())

	db, err := sql.Open("mysql", url)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("create database mysql_test")
	if err != nil {
		t.Fatal(err)
	}
	db.Close()

	Remove(id)
}
