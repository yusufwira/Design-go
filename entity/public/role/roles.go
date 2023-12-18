package role

import (
	"time"

	"gorm.io/gorm"
)

type ModelHasRole struct {
	RoleId    int `json:"role_id" gorm:"primary_key;references:ID"`
	Roles     Roles
	ModelType string `json:"model_type" gorm:"primary_key"`
	ModelId   int    `json:"model_id" gorm:"primary_key"`
	EmpNo     string `json:"emp_no"`
	CompCode  string `json:"comp_code"`
}

type Roles struct {
	ID           int       `json:"id" gorm:"primary_key"`
	Name         string    `json:"name"`
	ReadableName string    `json:"readable_name"`
	CompanyCode  string    `json:"company_code"`
	GuardName    string    `json:"guard_name"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type MyRole struct {
	Name     string  `json:"name"`
	CompCode *string `json:"comp_code"`
}

func (ModelHasRole) TableName() string {
	return "public.model_has_roles"
}

type ModelHasRoleRepo struct {
	DB *gorm.DB
}

func NewModelHasRoleRepo(db *gorm.DB) *ModelHasRoleRepo {
	return &ModelHasRoleRepo{DB: db}
}

func (Roles) TableName() string {
	return "public.roles"
}

type RolesRepo struct {
	DB *gorm.DB
}

func NewRolesRepo(db *gorm.DB) *RolesRepo {
	return &RolesRepo{DB: db}
}

func (t RolesRepo) FindRoleByUser(nik string) ([]MyRole, error) {
	var role []MyRole
	err := t.DB.Table("public.roles").
		Select("public.roles.name, public.model_has_roles.comp_code").
		Joins("JOIN public.model_has_roles ON public.model_has_roles.role_id = public.roles.id").
		Where("public.model_has_roles.emp_no = ?", nik).
		Scan(&role).
		Error
	if err != nil {
		return role, err
	}
	return role, nil
}
func (t ModelHasRoleRepo) FindRoleByUser(nik string) ([]ModelHasRole, error) {
	var test []ModelHasRole

	err := t.DB.
		Preload("ModelHasRole").
		Find(&test).Error

	if err != nil {
		return test, err
	}
	return test, nil
}
