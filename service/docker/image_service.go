package service

import (
	"DYCLOUD/global"
	global2 "DYCLOUD/global/docker"
	model "DYCLOUD/model/docker"
	"DYCLOUD/utils/docker/docker"
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"sort"
	"strings"
	"sync"
)

type ImageService struct {
	store sync.Map
}

type AuthConfig struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	ServerHost string `json:"serveraddress"`
	Email      string `json:"email"`
}

// ListImage 获取镜像列表
func (e *ImageService) ListImage(host string, req model.SearchImage) (res []model.Image, err error) {

	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return nil, err
	}

	option := types.ImageListOptions{
		All: true,
	}

	if strings.Trim(req.Name, "") != "" {
		option.Filters = filters.NewArgs(filters.KeyValuePair{Key: "reference", Value: req.Name})
	}

	var imageSummaryList []types.ImageSummary

	imageSummaryList, err = cli.ImageList(context.TODO(), option)
	if err != nil {
		return nil, err
	}

	val, ok := e.store.Load("pull-image-list")
	if ok {
		imageSummaryList = append(imageSummaryList, val.([]types.ImageSummary)...)
	}

	imageList := make([]model.Image, 0)
	for _, re := range imageSummaryList {
		tag := "<none>:<none>"
		if len(re.RepoTags) > 0 {
			tag = re.RepoTags[0]
		}

		imageList = append(imageList, model.Image{
			Tag:     tag,
			Id:      re.ID,
			Created: docker.TimeFormat(re.Created),
			Size:    docker.UnitFormat(re.Size),
		})
	}

	sort.Slice(imageList, func(i, j int) bool {
		if docker.StringToTimestamp(imageList[i].Created) > docker.StringToTimestamp(imageList[j].Created) {
			return true
		} else {
			return false
		}
	})
	return imageList, nil
}

// PullImage 下载镜像
func (e *ImageService) PullImage(host string, c model.PullImage) error {

	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return err
	}

	imagePullOptions := types.ImagePullOptions{}

	if c.Auth.Enable {
		marshal, _ := json.Marshal(AuthConfig{
			Username:   c.Auth.Username,
			Password:   c.Auth.Password,
			ServerHost: c.Auth.ServerAddress,
		})
		auth := base64.StdEncoding.EncodeToString(marshal)
		imagePullOptions.RegistryAuth = auth
	}

	res, err := cli.ImagePull(context.Background(), c.Name, imagePullOptions)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(res)
	lastLine := ""
	for scanner.Scan() {
		lastLine = scanner.Text()
		global.DYCLOUD_LOG.Info(lastLine)
	}

	return nil
}

// RemoveImage 删除镜像
func (e *ImageService) RemoveImage(host string, ids []string) error {

	cli, err := global2.DockerClient.Load(host)
	if err != nil {
		return err
	}

	imageRemoveOptions := types.ImageRemoveOptions{}
	for _, id := range ids {
		_, err := cli.ImageRemove(context.Background(), id, imageRemoveOptions)
		if err != nil {
			return err
		}
	}
	return nil
}
