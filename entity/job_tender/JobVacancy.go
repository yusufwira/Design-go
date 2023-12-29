package job_tender

import (
	"time"

	"gorm.io/gorm"
)

type JobVacancy struct {
	Id              int       `json:"id" gorm:"primary_key"`
	PositionName    string    `json:"position_name"`
	KompTitle       string    `json:"komp_title"`
	JobDesc         string    `json:"job_desc"`
	MaxGrade        string    `json:"max_grade"`
	MinGrade        string    `json:"min_grade"`
	Klasifikasi     string    `json:"klasifikasi"`
	KlasifikasiUmum string    `json:"klasifikasi_umum"`
	ValidFrom       time.Time `json:"valid_from" gorm:"default:null"`
	ValidTo         time.Time `json:"valid_to" gorm:"default:null"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (JobVacancy) TableName() string {
	return "mobile.job_tender_vacancy"
}

type JobVacancyRepo struct {
	DB *gorm.DB
}

func GetJobVacancyRepo(db *gorm.DB) *JobVacancyRepo {
	return &JobVacancyRepo{DB: db}
}

func (repo JobVacancyRepo) FindJobByID(id int) (JobVacancy, error) {
	var job JobVacancy
	err := repo.DB.First(&job, id).Error
	return job, err
}

func (t JobVacancyRepo) Find(id int) (JobVacancy, error) {
	var job JobVacancy
	err := t.DB.Where("id=?", id).Take(&job).Error
	return job, err
}
