package ws

type TerminalRequest struct {
	Name      string `json:"name"  form:"name"`
	PodName   string `json:"pod_name" form:"pod_name"`
	Namespace string `json:"namespace" form:"namespace"`
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	XToken    string `json:"x-token" form:"x-token"`
	Cols      int    `json:"cols" form:"cols"`
	Rows      int    `json:"rows" form:"rows"`
}
