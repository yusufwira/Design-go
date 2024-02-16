package feed

import (
	"github.com/gin-gonic/gin"
	"github.com/yusufwira/lern-golang-gin/entity/dbo/pihc"
	"gorm.io/gorm"
)

type FeedController struct {
	PihcMasterKaryRtDbRepo *pihc.PihcMasterKaryRtDbRepo
	PihcMasterCompanyRepo  *pihc.PihcMasterCompanyRepo
}

func getFeedController(Db *gorm.DB) *FeedController {
	return &FeedController{
		PihcMasterKaryRtDbRepo: pihc.NewPihcMasterKaryRtDbRepo(Db),
		PihcMasterCompanyRepo:  pihc.NewPihcMasterCompanyRepo(Db)}
}

func (c *FeedController) getFeed(ctx *gin.Context) {
	// return
}
