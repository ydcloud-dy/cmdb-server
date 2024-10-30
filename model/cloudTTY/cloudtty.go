package cloudTTY

type CloudTTY struct {
	ClusterId int `json:"cluster_id"`
}
type PodMessage struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Container string `json:"container"`
}
