package prefab

import (
	"testing"
	"time"
	"github.com/jlaffaye/ftp"
	"bytes"
	"io/ioutil"
)

func TestSetupFtpContainer(t *testing.T) {
	id, url := StartFtpContainer()
	t.Log(url)
	if url == "" {
		t.Fatal("url empty")
	}
	if id == "" {
		t.Fatal("id empty")
	}

	if err := WaitForFtp(url, 30 * time.Second); err != nil {
		t.Fatal(err.Error())
	}

	if c, err := ftp.DialTimeout(url, 30*time.Second); err != nil {
		t.Fatal(err.Error())
	} else {
		defer c.Quit()
		if err := c.Login(ftpUser,ftpPassword); err != nil {
			t.Fatal(err.Error())
		}

		b := bytes.NewBufferString("test data")

		if err := c.Stor("./test.txt", b); err != nil {
			t.Fatal(err)
		}

		if r, err := c.Retr("./test.txt"); err != nil {
			t.Fatal(err)
		} else {
			defer r.Close()
			if b, err := ioutil.ReadAll(r); err != nil {
				t.Fatal(err)
			} else if string(b) != "test data" {
				t.Fatal("Not same data:", string(b))
			}
		}

	}


	Remove(id)
}
