package router

type RouterGroup struct {
	CloudPlatformRouter
	CloudRegionRouter
	CloudVirtualMachineRouter
	CloudLoadBalancerRouter
	CloudRDSRouter
}

var RouterGroupApp = new(RouterGroup)
