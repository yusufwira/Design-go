package users

import (
	"html"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id uint `json:"id" gorm:"primary_key"`
	// Name          string    `json:"name"`
	// Email         string    `json:"email"`
	// Nik           string    `json:"nik"`
	Password string `json:"password"`
	// RememberToken string    `json:"remember_token"`
	Username string `json:"username"`
	// Type          string    `json:"type"`
	// UserType      string    `json:"user_type"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "public.users"
}

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (u UserRepo) Create(user User) {
	u.DB.Create(&user)
}

func (user *User) BeforeCreate(*gorm.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	return nil
}

func (u UserRepo) GetAll(user User) []User {
	var data []User
	u.DB.Find(&data)
	return data
}

func (u UserRepo) GetUsersID(user User, id string, ctx *gin.Context) []User {
	var data []User
	if err := u.DB.Where("id = ?", id).Take(&data).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return data
	}
	u.DB.Where("id = ?", id).Take(&data)
	return data
}

func (u UserRepo) DelUsersID(user User, id string, ctx *gin.Context) []User {
	var data []User
	if err := u.DB.Where("id = ?", id).Take(&data).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return data
	}
	u.DB.Where("id = ?", id).Delete(&data)
	ctx.JSON(http.StatusOK, gin.H{"data": true})
	return data
}

func (u UserRepo) UpUsersID(user User, id string, ctx *gin.Context) []User {
	var data []User
	if err := u.DB.Where("id = ?", id).Take(&data).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return data
	}
	u.DB.Save(&data)
	ctx.JSON(http.StatusOK, &data)
	return data
}

func (u UserRepo) LoginCheck(username string, password string) (User, error) {
	var user User

	err_username := u.DB.Where("username=?", username).Take(&user).Error
	if err_username == nil {
		err_pw := user.ValidatePassword(password)
		if err_pw == nil {
			return user, nil
		}
		return User{}, err_pw
	}
	return User{}, err_username
}

func (user *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}
