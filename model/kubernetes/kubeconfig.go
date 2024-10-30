package kubernetes

import (
	"DYCLOUD/utils"
	v1 "k8s.io/api/authorization/v1"
)

var (
	ProxyVerify     = utils.Rules{"Path": {utils.NotEmpty()}}
	TerminalVerify  = utils.Rules{"Name": {utils.NotEmpty()}, "PodName": {utils.NotEmpty()}, "Namespace": {utils.NotEmpty()}}
	RoleTypeVerify  = utils.Rules{"RoleType": {utils.NotEmpty()}}
	ApiGroupsVerify = utils.Rules{"ApiType": {utils.NotEmpty()}}
	RoleVerify      = utils.Rules{"Rules": {utils.NotEmpty()}, "Metadata": {utils.NotEmpty()}}
	PodVerify       = utils.Rules{"ClusterId": {utils.NotEmpty()}}
)

type Kubeconfig struct {
	APIVersion     string         `yaml:"apiVersion" json:"apiVersion"`
	Kind           string         `yaml:"kind" json:"kind"`
	Clusters       []ClusterEntry `yaml:"clusters" json:"clusters"`
	Contexts       []ContextEntry `yaml:"contexts" json:"contexts"`
	CurrentContext string         `yaml:"current-context" json:"current-context"`
	Preferences    struct{}       `yaml:"preferences" json:"preferences"`
	Users          []UserEntry    `yaml:"users" json:"users"`
}

type KubeCluster struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data" json:"certificate-authority-data"`
	Server                   string `yaml:"server" json:"server"`
}

type ClusterEntry struct {
	Name    string      `yaml:"name" json:"name"`
	Cluster KubeCluster `yaml:"cluster" json:"cluster"`
}

type Context struct {
	Cluster string `yaml:"cluster" json:"cluster"`
	User    string `yaml:"user" json:"user"`
}

type ContextEntry struct {
	Name    string  `yaml:"name" json:"name"`
	Context Context `yaml:"context" json:"context"`
}

type KubeUser struct {
	ClientCertificateData string `yaml:"client-certificate-data" json:"client-certificate-data"`
	ClientKeyData         string `yaml:"client-key-data" json:"client-key-data"`
}

type UserEntry struct {
	Name string   `yaml:"name" json:"name"`
	User KubeUser `yaml:"user" json:"user"`
}

type PermissionCheckResult struct {
	Resource v1.ResourceAttributes
	Allowed  bool
}

type ApiGroupOption struct {
	Group     string              `json:"group"`
	Resources []ApiResourceOption `json:"resources"`
}
type ApiResourceOption struct {
	Resource string   `json:"resource"`
	Verbs    []string `json:"verbs"`
}
