package profile

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Profile struct {
	ID            int       `json:"id"`
	Nik           string    `json:"nik"`
	Bio           *string   `json:"bio" gorm:"default:null"`
	LinkTwitter   *string   `json:"link_twitter" gorm:"default:null"`
	LinkInstagram *string   `json:"link_instagram" gorm:"default:null"`
	LinkWebsite   *string   `json:"link_website" gorm:"default:null"`
	LinkFacebook  *string   `json:"link_facebook" gorm:"default:null"`
	LinkTiktok    *string   `json:"link_tiktok" gorm:"default:null"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Suku          *string   `json:"suku"`
	GolonganDarah *string   `json:"golongan_darah" gorm:"default:null"`
	VisiPribadi   *string   `json:"visi_pribadi" gorm:"default:null"`
	NilaiPribadi  *string   `json:"nilai_pribadi" gorm:"default:null"`
	Interest      *string   `json:"interest" gorm:"default:null"`
	LinkLinkedin  *string   `json:"link_linkedin" gorm:"default:null"`
}

type AboutUs struct {
	ID           int       `json:"id"`
	Nik          string    `json:"nik"`
	AboutUsDesc  *string   `json:"about_us_desc" gorm:"default:null"`
	AboutUsHobby *string   `json:"about_us_hobby" gorm:"default:null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type ProfileSkill struct {
	ID            int       `json:"id"`
	IdParentSkill *int      `json:"id_parent_skill"`
	Nik           string    `json:"nik"`
	Type          string    `json:"type"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type PhotoProfile struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	EmpNo     string    `json:"emp_no"`
}

type PengalamanKerja struct {
	RiwayatJabatan string `json:"riwayat_jabatan"`
}

func (Profile) TableName() string {
	return "mobile.profile"
}

func (AboutUs) TableName() string {
	return "mobile.about_us"
}

func (ProfileSkill) TableName() string {
	return "mobile.profile_skill"
}

func (PhotoProfile) TableName() string {
	return "mobile.profile_photo"
}

type ProfileRepo struct {
	DB *gorm.DB
}

type PhotoProfileRepo struct {
	DB            *gorm.DB
	StorageClient *storage.Client
}

type AboutUsRepo struct {
	DB *gorm.DB
}

type ProfileSkillRepo struct {
	DB *gorm.DB
}

type PengalamanKerjaRepo struct {
	DB *gorm.DB
}

func NewProfileRepo(db *gorm.DB) *ProfileRepo {
	return &ProfileRepo{DB: db}

}
func NewPhotoProfileRepo(db *gorm.DB, sc *storage.Client) *PhotoProfileRepo {
	return &PhotoProfileRepo{DB: db, StorageClient: sc}
}

func NewAboutUsRepo(db *gorm.DB) *AboutUsRepo {
	return &AboutUsRepo{DB: db}
}

func NewProfileSkillRepo(db *gorm.DB) *ProfileSkillRepo {
	return &ProfileSkillRepo{DB: db}
}

func NewPengalamanKerjaRepo(db *gorm.DB) *PengalamanKerjaRepo {
	return &PengalamanKerjaRepo{DB: db}
}

func (t ProfileRepo) Create(p Profile) (Profile, error) {
	err := t.DB.Create(&p).Error
	if err != nil {
		return p, err
	}
	return p, nil
}

func (t ProfileRepo) Update(p Profile) (Profile, error) {
	err := t.DB.Where("nik=?", p.Nik).Save(&p).Error
	if err != nil {
		return p, err
	}
	return p, nil
}

func (t ProfileSkillRepo) Create(p ProfileSkill) (ProfileSkill, error) {
	err := t.DB.Create(&p).Error
	if err != nil {
		return p, err
	}
	return p, nil
}

func (t ProfileSkillRepo) CreateC(p []ProfileSkill) ([]ProfileSkill, error) {
	err := t.DB.Create(&p).Error
	if err != nil {
		return p, err
	}

	var indexA []int
	var indexB []int
	iterasiA := 0
	iterasiB := 0

	for i, data := range p {
		if data.Type == "category_skill" {
			if len(indexA) > 0 {
				iterasiA++
			}
			indexA = append(indexA, data.ID)
		}
		if data.Type == "main_skill" {
			if len(indexA) > 0 {
				data.IdParentSkill = &indexA[iterasiA]
			}
			if len(indexB) > 0 {
				iterasiB++
			}
			indexB = append(indexB, data.ID)
		}
		if data.Type == "sub_skill" {
			if len(indexB) > 0 {
				data.IdParentSkill = &indexB[iterasiB]
			}
		}
		p[i] = data
	}

	err1 := t.DB.Save(&p).Error
	if err1 != nil {
		return p, err1
	}

	return p, nil
}

func (t ProfileSkillRepo) Update(p ProfileSkill) (ProfileSkill, error) {
	err := t.DB.Where("id=?", p.ID).Save(&p).Error
	if err != nil {
		fmt.Println("ERROR")
		return p, err
	}

	return p, nil
}

func (t ProfileSkillRepo) UpdateC(p []ProfileSkill) ([]ProfileSkill, error) {
	var indexA []int
	var indexB []int
	iterasiA := 0
	iterasiB := 0

	for i, data := range p {
		if data.Type == "category_skill" {
			if len(indexA) > 0 {
				iterasiA++
			}
			indexA = append(indexA, data.ID)
		}
		if data.Type == "main_skill" {
			if len(indexA) > 0 {
				data.IdParentSkill = &indexA[iterasiA]
			}
			if len(indexB) > 0 {
				iterasiB++
			}
			indexB = append(indexB, data.ID)
		}
		if data.Type == "sub_skill" {
			if len(indexB) > 0 {
				data.IdParentSkill = &indexB[iterasiB]
			}
		}
		p[i] = data
	}

	err := t.DB.Save(&p).Error
	if err != nil {
		return p, err
	}

	return p, nil
}

func (t ProfileSkillRepo) DeleteC(p []ProfileSkill) ([]ProfileSkill, error) {
	err := t.DB.Delete(&p).Error
	if err != nil {
		return p, err
	}

	return p, nil
}

func (t AboutUsRepo) Create(au AboutUs) (AboutUs, error) {
	err := t.DB.Create(&au).Error
	if err != nil {
		return au, err
	}
	return au, nil
}

func (t AboutUsRepo) Update(au AboutUs) (AboutUs, error) {
	err := t.DB.Where("nik=?", au.Nik).Save(&au).Error
	if err != nil {
		fmt.Println("ERROR")
		return au, err
	}
	return au, nil
}

func (t ProfileRepo) FindProfile(nik string) (Profile, error) {
	var profile Profile
	err := t.DB.Where("nik=?", nik).First(&profile).Error
	if err != nil {
		fmt.Println("ERROR2")
		return profile, err
	}
	return profile, nil
}
func (t ProfileSkillRepo) FindProfileCategorySkill(id int) (ProfileSkill, error) {
	var profile ProfileSkill
	err := t.DB.Where("id=?", id).First(&profile).Error
	if err != nil {
		return profile, err
	}
	return profile, nil
}

func (t ProfileSkillRepo) FindProfileCategorySkillArr(nik string) ([]ProfileSkill, error) {
	var profile []ProfileSkill
	err := t.DB.Where("nik=? AND id_parent_skill IS NULL", nik).Find(&profile).Error
	if err != nil {
		return profile, err
	}
	return profile, nil
}

func (t ProfileSkillRepo) FindProfileSkillArr(id int) ([]ProfileSkill, error) {
	var profile []ProfileSkill
	err := t.DB.Where("id_parent_skill=?", id).Find(&profile).Error
	if err != nil {
		return profile, err
	}
	return profile, nil
}

func (t ProfileSkillRepo) GetProfileSkillArr(nik string, typeSkill string) ([]ProfileSkill, error) {
	var results []ProfileSkill
	var ps string
	if typeSkill == "category_skill" {
		ps = "ps"
	}
	if typeSkill == "main_skill" {
		ps = "ps2"
	}
	if typeSkill == "sub_skill" {
		ps = "ps3"
	}
	query := fmt.Sprintf(`select distinct(%s.*)
                          from mobile.profile_skill ps
                          left join mobile.profile_skill ps2 on ps2.id_parent_skill = ps.id
                          left join mobile.profile_skill ps3 on ps3.id_parent_skill = ps2.id
                          where ps.nik = ? and ps.type ='category_skill'`, ps)

	err := t.DB.Raw(query, nik).
		Scan(&results).Error

	if err != nil {
		return results, err
	}

	return results, nil
}

func (t ProfileSkillRepo) FindProfileSkill(id int, parent int) (ProfileSkill, error) {
	var profile ProfileSkill
	err := t.DB.Where("id=? AND id_parent_skill=?", id, parent).First(&profile).Error
	if err != nil {
		return profile, err
	}
	return profile, nil
}

func (t ProfileSkillRepo) FindProfileSkillIndiv(id int) (ProfileSkill, error) {
	var profile ProfileSkill
	err := t.DB.Where("id=?", id).First(&profile).Error
	if err != nil {
		return profile, err
	}
	return profile, nil
}

func (t AboutUsRepo) FindProfileAboutUs(nik string) (AboutUs, error) {
	var au AboutUs
	err := t.DB.Where("nik=?", nik).First(&au).Error
	if err != nil {
		return au, err
	}
	return au, nil
}
func (t PengalamanKerjaRepo) FindRiwayatJabatan(nik string) ([]PengalamanKerja, error) {
	var pk []PengalamanKerja

	err := t.DB.Raw("SELECT dbo.get_riwayat_jabatan_v1(?) as riwayat_jabatan", nik).Find(&pk).Error
	if err != nil {
		return pk, err
	}
	return pk, nil
}

func (t PhotoProfileRepo) FindPhotoProfile(nik string) (PhotoProfile, error) {
	var pp PhotoProfile
	// err := t.DB.Raw(`select pp.url from mobile.profile_photo pp where emp_no =?`, nik).Scan(&pp).Error
	err := t.DB.Where(`emp_no =?`, nik).First(&pp).Error
	if err != nil {
		return pp, err
	}
	return pp, nil
}

func (t PhotoProfileRepo) Create(pp PhotoProfile) (PhotoProfile, error) {
	err := t.DB.Create(&pp).Error
	if err != nil {
		return pp, err
	}
	return pp, nil
}

func (t PhotoProfileRepo) Update(pp PhotoProfile) (PhotoProfile, error) {
	err := t.DB.Save(&pp).Error
	if err != nil {
		return pp, err
	}
	return pp, nil
}

func (t PhotoProfileRepo) UploadFilePhotoProfile(nik string, objName string, files multipart.File) (string, string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}
	// gcsFile := "serviceAccount.json"

	// gcs, err := ioutil.ReadFile(gcsFile)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// cfg, err := google.JWTConfigFromJSON(gcs)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// expires := time.Now().Add(time.Second * 60)

	ctx := context.Background() // Create a new context

	bckt := t.StorageClient.Bucket(os.Getenv("GC_LUMEN_BUCKET"))
	folderName1 := "PhotoProfile"
	time := time.Now()
	tahun := time.Format("2006")
	emp_no := nik
	tanggal := time.Format("02012006")
	jam := time.Format("150405")
	fileName := tanggal + "_" + jam + "_" + objName
	fd, _ := createFolderPhotosProfile(bckt, ctx, folderName1, tahun, emp_no)

	location := fd + fileName
	object := bckt.Object(location)
	wc := object.NewWriter(ctx)

	// set cache control so the image will be served fresh by browsers
	// To do this with the object handle, you'd first have to upload, then update
	wc.ObjectAttrs.CacheControl = "Cache-Control:no-cache, max-age=0"

	// multipart.File has a reader!
	if _, err := io.Copy(wc, files); err != nil {
		log.Printf("Unable to write a file to Google Cloud Storage: %v\n", err)
		return "", " ", err
	}

	if err := wc.Close(); err != nil {
		return "", " ", fmt.Errorf("Writer.Close: %v", err)
	}

	// Set the object's ACL to public read access
	if err := makePublic(ctx, object); err != nil {
		log.Fatalf("Failed to make the object public: %v", err)
	}

	// opts := &storage.SignedURLOptions{
	// 	GoogleAccessID: cfg.Email,
	// 	PrivateKey:     cfg.PrivateKey,
	// 	Method:         "GET",
	// 	Expires:        expires,
	// }

	// url, err := bckt.SignedURL(fd+fileName, opts)
	// if err != nil {
	// 	return "", fmt.Errorf("Bucket(%v).SignedURL: %w", bckt, err)
	// }

	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", os.Getenv("GC_LUMEN_BUCKET"), location)

	return imageURL, fileName, nil
}

func createFolderPhotosProfile(bucketName *storage.BucketHandle, ctx context.Context, folderName1 string, tahun string, emp_no string) (string, error) {
	// Create an empty object (blob) with the folder name as the object name
	foldername := folderName1 + "/" + tahun + "/" + emp_no + "/"

	writer := bucketName.Object(foldername).NewWriter(ctx)
	if _, err := writer.Write([]byte("")); err != nil {
		return foldername, err
	}
	if err := writer.Close(); err != nil {
		return foldername, err
	}

	return foldername, nil
}

// makePublic makes the object publicly accessible by setting its ACL.
func makePublic(ctx context.Context, object *storage.ObjectHandle) error {
	// Create a new ACL rule to allow public read access
	rule := storage.ACLRule{
		Entity: storage.AllUsers,
		Role:   storage.RoleReader,
	}

	// Add the ACL rule to the object's ACL
	if err := object.ACL().Set(ctx, rule.Entity, rule.Role); err != nil {
		return err
	}

	return nil
}
