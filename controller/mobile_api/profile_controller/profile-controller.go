package profile_controller

import (
	"cloud.google.com/go/storage"
	"github.com/yusufwira/lern-golang-gin/entity/mobile/profile"
	"gorm.io/gorm"
)

type ProfileController struct {
	ProfileRepo *profile.ProfileRepo
}

func NewProfileController(Db *gorm.DB, StorageClient *storage.Client) *ProfileController {
	return &ProfileController{ProfileRepo: profile.NewProfileRepo(Db)}
}
