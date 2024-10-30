package response

import "DYCLOUD/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
