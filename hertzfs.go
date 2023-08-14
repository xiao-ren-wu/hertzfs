package hertzfs

import (
	"context"
	"embed"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"net/http"
	"strings"
)

func HertzFS(hertz *server.Hertz, fs embed.FS, ops ...FSOption) {
	var conf = staticFSConf{CacheControl: 36500}
	for _, op := range ops {
		op(&conf)
	}
	matchPath := fmt.Sprintf("/%s/*filepath", conf.BasePath)
	cacheControl := fmt.Sprintf("max-age=%d", conf.CacheControl)
	hertz.Any(matchPath,
		func(ctx context.Context, c *app.RequestContext) {
			c.Header("CacheControl-Control", cacheControl)
			reqPath := string(c.Request.URI().Path())
			if strings.Trim(reqPath, "/") == conf.BasePath {
				reqPath = fmt.Sprintf("/%s/index.htm", conf.BasePath)
			}
			request, err := http.NewRequest(string(c.Request.Method()), reqPath, c.RequestBodyStream())
			if err != nil {
				c.JSON(400, map[string]interface{}{
					"message": "bad request",
				})
				return
			}
			staticServer := http.FileServer(http.FS(fs))
			staticServer.ServeHTTP(NewCustomerWriter(c), request)
		},
	)
}

type CustomerWriter struct {
	c      *app.RequestContext
	header http.Header
	clear  bool
}

func NewCustomerWriter(c *app.RequestContext) *CustomerWriter {
	return &CustomerWriter{
		c:      c,
		header: map[string][]string{},
	}
}

func (c *CustomerWriter) Header() http.Header {
	c.clear = false
	return c.header
}

func (c *CustomerWriter) Write(bytes []byte) (int, error) {
	if !c.clear {
		for k, vals := range c.header {
			c.c.Header(k, vals[0])
		}
		c.clear = true
		c.header = map[string][]string{}
	}
	return c.c.Write(bytes)
}

func (c *CustomerWriter) WriteHeader(statusCode int) {
	c.c.Status(statusCode)
}
