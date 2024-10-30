package api

type ApiGroup struct {
	ContainerApi
	ImageApi
	NetworkApi
	VolumeApi
	SystemApi
	HostApi
}

var ApiGroupApp = new(ApiGroup)
