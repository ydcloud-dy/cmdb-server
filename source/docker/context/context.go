package context

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type Context struct {
}

func (ctx *Context) GetHost(c *gin.Context) (string, error) {
	fmt.Println(c.Request.Host)
	host := c.Param("host")
	if strings.Trim(host, " ") == "" {
		return "", errors.New("host is empty")
	}
	return host, nil
}
