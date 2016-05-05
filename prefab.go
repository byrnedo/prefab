package prefab

import (
	"encoding/json"
	"errors"
	"fmt"
	gDoc "github.com/fsouza/go-dockerclient"
	"github.com/satori/go.uuid"
	"os"
	"strconv"
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

func start(image string, portB map[gDoc.Port][]gDoc.PortBinding, envs []string, forcePull bool) (*gDoc.Container, error) {

	if imgs, err := dockCli.ListImages(gDoc.ListImagesOptions{Filter: image}); err != nil || len(imgs) == 0 {
		if forcePull {
			if err := dockCli.PullImage(gDoc.PullImageOptions{Repository: image, OutputStream: os.Stdout}, gDoc.AuthConfiguration{}); err != nil {

				return nil, err
			}

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
}

func startStandardContainer(cnfOverride func(SetupOpts) SetupOpts) (*gDoc.Container, string, int, error) {
	var (
		con           *gDoc.Container
		hostPortStr   string
		containerPort string
		opts          SetupOpts
	)

	defaultCnf := SetupOpts{
		Protocol:  "tcp",
		HostIp:    "127.0.0.1",
		ForcePull: false,
	}

	opts = cnfOverride(defaultCnf)

	if opts.ExposedPort != 0 {
		containerPort = strconv.Itoa(opts.ExposedPort)
	}

	if opts.PublishedPort != 0 {
		hostPortStr = strconv.Itoa(opts.PublishedPort)
	}

	dockerExposedPort := gDoc.Port(containerPort + "/" + opts.Protocol)

	if id, err := Running(opts.Image); err != nil || len(id) < 1 {
		if con, err = start(opts.Image, map[gDoc.Port][]gDoc.PortBinding{
			dockerExposedPort: []gDoc.PortBinding{gDoc.PortBinding{
				HostIP:   opts.HostIp,
				HostPort: hostPortStr,
			}},
		}, opts.Envs, opts.ForcePull); err != nil {
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
	v, _ := json.Marshal(con)
	fmt.Println(string(v))
	hostPort, _ := strconv.Atoi(port[0].HostPort)

	return con, port[0].HostIP, hostPort, nil
}

func StartCustom(image string, portB map[gDoc.Port][]gDoc.PortBinding, envs []string, forcePull bool) (string, error) {
	var (
		con *gDoc.Container
		err error
	)

	if con, err = start(image, portB, envs, forcePull); err != nil {
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

func Remove(id string) (bool, error) {

	if err := dockCli.RemoveContainer(gDoc.RemoveContainerOptions{
		Force: true,
		ID:    id,
	}); err != nil {
		return false, err
	}
	return true, nil

}

func RemoveByImage(image string) (bool, error) {
	var (
		id  string
		err error
	)
	if id, err = Running(image); err != nil && len(id) > 0 {
		return false, err
	}

	return Remove(id)

}
