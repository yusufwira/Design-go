package tjsl

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"gorm.io/gorm"
)

type KegiatanPhotos struct {
	Id            int       `json:"id" gorm:"primary_key"`
	KegiatanId    int       `json:"kegiatan_id"` //id_kegiatan_karyawan
	IsKoordinator int       `json:"is_koordinator"`
	OriginalName  string    `json:"original_name"`
	Url           string    `json:"url"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Extendtion    string    `json:"extendtion"`
}

func (KegiatanPhotos) TableName() string {
	return "tjsl.kegiatan_photos"
}

type KegiatanPhotosRepo struct {
	DB *gorm.DB
}

func NewKegiatanPhotosRepo(db *gorm.DB) *KegiatanPhotosRepo {
	return &KegiatanPhotosRepo{DB: db}
}

func (t KegiatanPhotosRepo) LastString(ss []string) string {
	return ss[len(ss)-1]
}

func (t KegiatanPhotosRepo) Create(kp KegiatanPhotos) KegiatanPhotos {
	t.DB.Create(&kp)
	return kp
}

func (t KegiatanPhotosRepo) Update(kp KegiatanPhotos) (KegiatanPhotos, error) {
	err := t.DB.Save(&kp).Error
	return kp, err
}

func (t KegiatanPhotosRepo) FindData(id int) []KegiatanPhotos {
	var kgtn_phto []KegiatanPhotos
	t.DB.Where("kegiatan_id=? AND is_koordinator=?", id, 0).Find(&kgtn_phto)
	return kgtn_phto
}

func (t KegiatanPhotosRepo) GetFileExtensionFromUrl(rawUrl string) (string, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}
	pos := strings.LastIndex(u.Path, ".")
	if pos == -1 {
		return "", errors.New("couldn't find a period to indicate a file extension")
	}
	return u.Path[pos+1 : len(u.Path)], nil
}

func (t KegiatanPhotosRepo) DelPhotosID(kegiatan_id int) ([]KegiatanPhotos, error) {
	var data []KegiatanPhotos
	err := t.DB.Where("kegiatan_id = ?", kegiatan_id).First(&data).Error
	if err == nil {
		t.DB.Where("kegiatan_id = ?", kegiatan_id).Delete(&data)
		return data, nil
	}
	return data, err
}

func (t KegiatanPhotosRepo) DelPhotosIDLama(kegiatan_id int, list_id []int) {
	t.DB.Where("kegiatan_id = ? AND id not in(?)", kegiatan_id, list_id).Delete(&KegiatanPhotos{})
}

func (t KegiatanPhotosRepo) FindDataPhotosID(id int, is_koordinator int) []KegiatanPhotos {
	var kgtn_phtos []KegiatanPhotos
	t.DB.Where("kegiatan_id=? AND is_koordinator=?", id, is_koordinator).Find(&kgtn_phtos)
	return kgtn_phtos
}
