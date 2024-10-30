package docker

import (
	"DYCLOUD/source/docker/context"
	"DYCLOUD/utils/docker/docker"
)

var Context = new(context.Context)
var DockerClient = docker.NewDockerClient()
