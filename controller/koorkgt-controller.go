package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/yusufwira/lern-golang-gin/entity/tjsl/koorKgt"
	"gorm.io/gorm"
)

type KoorKgtController struct {
	KegiatanKoordinatorRepo *koorKgt.KegiatanKoordinatorRepo
}

func NewKoorKgtController(db *gorm.DB) *KoorKgtController {
	return &KoorKgtController{KegiatanKoordinatorRepo: koorKgt.NewKegiatanKoordinatorRepo(db)}
}

func (c *KoorKgtController) ListKoordinator(ctx *gin.Context) {

}

func (c *KoorKgtController) StoreKoordinator(ctx *gin.Context) {

}

func (c *KoorKgtController) ShowDetailKoordinator(ctx *gin.Context) {

}

func (c *KoorKgtController) DeleteKoordinator(ctx *gin.Context) {

}
