package feedsController

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Feed represents the entity structure
type Feed struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	NIK         string    `gorm:"column:nik"`
	Description string    `gorm:"column:desc"`
	Hashtag     string    `gorm:"column:hastag"`
	URL         string    `gorm:"column:url"`
	IsBroadcast int32     `gorm:"column:is_broadcast"`
	IsPin       int32     `gorm:"column:is_pin"`
	CompanyCode string    `gorm:"column:comp_code"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
	IsModerate  int16     `gorm:"column:is_moderate"`
	Category    string    `gorm:"column:kategori"`
	IsClickable int16     `gorm:"column:is_clickable"`
	ReferenceID string    `gorm:"column:referensi_id"`
}

func (Feed) TableName() string {
	return "mobile.feeds"
}

type FeedRepo struct {
	db *gorm.DB
}

func GetFeedRepo(db *gorm.DB) *FeedRepo {
	return &FeedRepo{db: db}
}

// Handler to create a new feed
func (r FeedRepo) createFeed(c *gin.Context) {
	var feed Feed
	if err := c.ShouldBindJSON(&feed); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	r.db.Create(&feed)
	c.JSON(201, feed)
}

// Handler to get a feed by ID
func (r FeedRepo) getFeed(c *gin.Context) {
	var feed Feed
	id := c.Param("id")
	if err := r.db.First(&feed, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(200, feed)
}

// Handler to update a feed by ID
func (r FeedRepo) updateFeed(c *gin.Context) {
	var feed Feed
	id := c.Param("id")
	if err := r.db.First(&feed, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Record not found!"})
		return
	}
	if err := c.ShouldBindJSON(&feed); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	r.db.Save(&feed)
	c.JSON(200, feed)
}

// Handler to delete a feed by ID
func (r FeedRepo) deleteFeed(c *gin.Context) {
	var feed Feed
	id := c.Param("id")
	if err := r.db.First(&feed, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Record not found!"})
		return
	}
	r.db.Delete(&feed)
	c.JSON(204, gin.H{})
}
