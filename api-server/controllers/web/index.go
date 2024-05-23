package web

import (
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/danielpickens/yeti/api-server/config"
	"github.com/dnaielpickens/yeti/common/scookie"
)

var (
	indexContent  []byte
	indexLoadOnce sync.Once
)

func Index(ctx *gin.Context) {
	indexLoadOnce.Do(func() {
		var err error
		indexContent, err = os.ReadFile(path.Join(config.GetUIDistDir(), "index.html"))
		if err != nil {
			logrus.Panicf("failed to read index.html:%s", err.Error())
		}
	})
	ctx.Data(200, "text/html; charset=utf-8", indexContent)
}

func Logout(ctx *gin.Context) {
	_ = scookie.DeleteUsernameFromCookie(ctx)
	ctx.Redirect(http.StatusFound, "/login")
}
