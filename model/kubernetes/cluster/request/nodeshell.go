package request

type NodeTTY struct {
	ClusterId int    `json:"cluster_id"`
	NodeName  string `json:"node_name"`
}
