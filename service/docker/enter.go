package service

import (
	"sync"
)

type ServiceGroup struct {
	ContainerService
	ImageService
	NetworkService
	VolumeService
	SystemService
	HostService
}

var ServiceGroupApp = ServiceGroup{
	ContainerService{},
	ImageService{},
	NetworkService{},
	VolumeService{},
	SystemService{},
	HostService{},
}

func init() {

	//client, err := docker.NewDockerClient()
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}

	ServiceGroupApp.ImageService.store = sync.Map{}

}
