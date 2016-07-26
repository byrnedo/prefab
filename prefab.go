package prefab

import (
	"errors"
	gDoc "github.com/fsouza/go-dockerclient"
	"github.com/satori/go.uuid"
	"os"
	"strconv"
	"net"
	"time"
	"strings"
)

var (
	dockCli *gDoc.Client
	runId   string
)

func init() {
	var err error

	dockCli, err = gDoc.NewClientFromEnv()
	if err != nil {
		panic("Failed to connect to docker:" + err.Error())
	}

	runId = uuid.NewV1().String()
}

func start(image string, portB map[gDoc.Port][]gDoc.PortBinding, envs []string, forcePull bool, privileged bool) (*gDoc.Container, error) {

	if imgs, err := dockCli.ListImages(gDoc.ListImagesOptions{Filter: image}); (err != nil || len(imgs) == 0) || forcePull {
		if err := dockCli.PullImage(gDoc.PullImageOptions{Repository: image, OutputStream: os.Stdout}, gDoc.AuthConfiguration{}); err != nil {

			return nil, err
		}
	}

	con, err := dockCli.CreateContainer(gDoc.CreateContainerOptions{
		Config: &gDoc.Config{
			Labels: map[string]string{
				"com.byrnedo.prefab.id": image + ":" + runId,
			},
			Env:   envs,
			Image: image,
		},
		HostConfig: &gDoc.HostConfig{
			PortBindings: portB,
			Privileged: privileged,
		},
	})
	if err != nil {
		return nil, err
	}

	if err := dockCli.StartContainer(con.ID, nil); err != nil {
		return nil, err
	}
	return dockCli.InspectContainer(con.ID)
}

type SetupOpts struct {
	Image         string
	ExposedPort   int
	PublishedPort int
	Protocol      string
	HostIp        string
	ForcePull     bool
	Envs          []string
	Privileged    bool
	ExtraPorts    map[gDoc.Port][]gDoc.PortBinding
}

type ConfOverrideFunc func(*SetupOpts)

func startStandardContainer(cnfOverride ConfOverrideFunc) (*gDoc.Container, string, int, error) {
	var (
		con           *gDoc.Container
		hostPortStr   string
		containerPort string
		defaultCnf    *SetupOpts
	)

	defaultCnf = &SetupOpts{
		Protocol:  "tcp",
		HostIp:    "127.0.0.1",
		ForcePull: false,
	}

	cnfOverride(defaultCnf)

	if defaultCnf.ExposedPort != 0 {
		containerPort = strconv.Itoa(defaultCnf.ExposedPort)
	}

	if defaultCnf.PublishedPort != 0 {
		hostPortStr = strconv.Itoa(defaultCnf.PublishedPort)
	}

	dockerExposedPort := gDoc.Port(containerPort + "/" + defaultCnf.Protocol)

	if id, err := Running(defaultCnf.Image); err != nil || len(id) < 1 {
		ports := map[gDoc.Port][]gDoc.PortBinding{
			dockerExposedPort: []gDoc.PortBinding{
				gDoc.PortBinding{
					HostIP:   defaultCnf.HostIp,
					HostPort: hostPortStr,
				},
			},
		}

		for exposed, binding := range defaultCnf.ExtraPorts {
			ports[exposed] = binding
		}

		if con, err = start(defaultCnf.Image, ports, defaultCnf.Envs, defaultCnf.ForcePull, defaultCnf.Privileged); err != nil {
			return nil, "", 0, errors.New("Error starting container:" + err.Error())
		}
	}
	port, found := con.NetworkSettings.Ports[dockerExposedPort]
	if !found {
		panic("No port mapping found")
	}
	if len(port) != 1 {
		panic("No port mapping found")
	}
	hostPort, _ := strconv.Atoi(port[0].HostPort)

	return con, port[0].HostIP, hostPort, nil
}

func StartCustom(image string, portB map[gDoc.Port][]gDoc.PortBinding, envs []string, forcePull bool, privileged bool) (string, error) {
	var (
		con *gDoc.Container
		err error
	)

	if con, err = start(image, portB, envs, forcePull, privileged); err != nil {
		return "", err
	}

	return con.ID, nil
}

func Running(image string) (string, error) {
	cons, err := dockCli.ListContainers(gDoc.ListContainersOptions{
		Filters: map[string][]string{
			"label": []string{"com.byrnedo.prefab.id=" + image + ":" + runId},
		},
	})
	if err != nil {
		return "", err
	}

	if len(cons) == 0 {
		return "", nil
	}
	return cons[0].ID, nil
}

func Remove(id string) (error) {

	if err := dockCli.RemoveContainer(gDoc.RemoveContainerOptions{
		RemoveVolumes: true,
		Force: true,
		ID:    id,
	}); err != nil {
		return err
	}
	return nil

}

func RemoveByImage(image string) (error) {
	var (
		id  string
		err error
	)
	if id, err = Running(image); err != nil && len(id) > 0 {
		return err
	}

	return Remove(id)

}

func WaitForPort(addr string, timeout time.Duration) error {

	var (
		proto = "tcp"
		timedOut = time.Now().Add(timeout)
		buff = make([]byte, 10)
	)

	if i := strings.Index(addr, "("); i > 0 {
		proto = addr[0:i]
		addr = addr[i+1: len(addr)-1]
	}

	for {
		if time.Now().After(timedOut) {
			return errors.New("Timed out waiting to connect")
		}
		time.Sleep(1 * time.Second)

		c, err := net.DialTimeout(proto, addr, 300 * time.Millisecond)
		if err !=  nil {
			continue
		}

		c.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		c.SetDeadline(time.Now().Add(100 * time.Millisecond))
		if readFromSocket(c, buff) {
			return nil
		}
	}

}

func readFromSocket(c net.Conn, buffer []byte) bool {
	_, err := c.Read(buffer)
	c.Close()
	if err != nil {
		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			return true
		}
		return false
	}
	return true
}
