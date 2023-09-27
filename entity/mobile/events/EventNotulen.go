package events

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"gorm.io/gorm"
)

type EventNotulen struct {
	IdNotulen int       `json:"id_notulen" gorm:"primary_key"`
	IdEvent   int       `json:"id_event"`
	Deskripsi string    `json:"deskripsi" gorm:"default:null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type EventNotulenFile struct {
	IdNotulenFile int       `json:"id_notulen_file" gorm:"primary_key"`
	IdNotulen     int       `json:"id_notulen"`
	FileName      string    `json:"file_name" gorm:"default:null"`
	FileUrl       string    `json:"file_url" gorm:"default:null"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (EventNotulen) TableName() string {
	return "mobile.event_notulen"
}

type EventNotulenRepo struct {
	DB            *gorm.DB
	StorageClient *storage.Client
}

func (t EventNotulenRepo) CreateNotulen(en EventNotulen) (EventNotulen, error) {
	err := t.DB.Create(&en).Error
	if err != nil {
		return en, err
	}
	return en, nil
}

func (t EventNotulenRepo) CreateNotulenFiles(en EventNotulenFile) (EventNotulenFile, error) {
	err := t.DB.Table("mobile.event_notulen_file").Create(&en).Error
	if err != nil {
		return en, err
	}
	return en, nil
}

func NewEventNotulenRepo(db *gorm.DB, sc *storage.Client) *EventNotulenRepo {
	return &EventNotulenRepo{DB: db, StorageClient: sc}
}

func (t EventNotulenRepo) FindEventNotulenK(idEvent int) (*EventNotulen, error) {
	var ev_notulen EventNotulen
	err := t.DB.Where("id_event=?", idEvent).Take(&ev_notulen).Error
	if err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	// Return nil and nil error to indicate that no record was found
		// 	ev_notulen.Deskripsi = nil
		// 	return nil, nil
		// }
		return nil, err
	}
	return &ev_notulen, nil
}

func (t EventNotulenRepo) GetDataNotulenFile(idNotulen int) ([]EventNotulenFile, error) {
	var ev_notulen_file []EventNotulenFile
	err := t.DB.Table("mobile.event_notulen_file").Where("id_notulen=?", idNotulen).Find(&ev_notulen_file).Error
	if err != nil {
		return nil, err
	}
	return ev_notulen_file, nil
}

func (t EventNotulenRepo) DeleteEventNotulen(id int) error {
	var ev_rb EventNotulen
	err := t.DB.Table("mobile.event_notulen").Where("id_event=?", id).Take(ev_rb).Error
	if err == nil {
		t.DB.Table("mobile.event_notulen").Where("id_event= ?", id).Delete(&ev_rb)
		t.DeleteEventNotulenFile(ev_rb.IdNotulen)
		return nil
	}
	return err
}

func (t EventNotulenRepo) DeleteEventNotulenFile(id int) (EventNotulenFile, error) {
	var ev_rb EventNotulenFile
	err := t.DB.Table("mobile.event_notulen_file").Where("id_notulen=?", id).First(&ev_rb).Error
	if err == nil {
		t.DB.Table("mobile.event_notulen_file").Where("id_notulen= ?", id).Delete(&ev_rb)
		return ev_rb, nil
	}
	return ev_rb, err
}

func (t EventNotulenRepo) DeleteEventNotulenFiles(id int) (EventNotulenFile, error) {
	var ev_rb EventNotulenFile
	err := t.DB.Table("mobile.event_notulen_file").Where("id_notulen_file=?", id).First(&ev_rb).Error
	if err == nil {
		t.DB.Table("mobile.event_notulen_file").Where("id_notulen_file= ?", id).Delete(&ev_rb)
		return ev_rb, nil
	}
	return ev_rb, err
}

func (t EventNotulenRepo) UploadFile(objName string, files multipart.File) (string, string, error) {
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
	folderName1 := "Event"
	folderName2 := "Notulen"
	time := time.Now()
	tahun := time.Format("2006")
	tanggal := time.Format("02012006")
	jam := time.Format("150405")
	fileName := tanggal + "_" + jam + "_" + objName
	fmt.Println(fileName)
	fd, _ := createFolderNotulen(bckt, ctx, folderName1, folderName2, tahun)

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

func createFolderNotulen(bucketName *storage.BucketHandle, ctx context.Context, folderName1 string, folderName2 string, tahun string) (string, error) {
	// Create an empty object (blob) with the folder name as the object name
	foldername := folderName1 + "/" + folderName2 + "/" + tahun + "/"

	writer := bucketName.Object(foldername).NewWriter(ctx)
	if _, err := writer.Write([]byte("")); err != nil {
		return foldername, err
	}
	if err := writer.Close(); err != nil {
		return foldername, err
	}

	return foldername, nil
}

func (t EventNotulenRepo) RenameFileGCS(objName string, newObjName string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading: %v", err)
	}
	ctx := context.Background() // Create a new context

	bckt := t.StorageClient.Bucket(os.Getenv("GC_IMAGE_BUCKET"))
	objectIterator := bckt.Objects(ctx, nil)
	var newObjectName string
	notFound := false

	for {
		attrs, err := objectIterator.Next()
		if err == iterator.Done {
			notFound = true
			break
		}
		if err != nil {
			log.Fatalf("error iterating over objects: %v", err)
		}

		// Jika nama objek mengandung kata kunci pencarian
		if strings.Contains(attrs.Name, objName) {
			objName += filepath.Ext(attrs.Name)

			if attrs.Name == objName {
				fmt.Printf("Found matching object: %s\n", attrs.Name)

				// // Ubah nama objek
				newObjectName = newObjName + filepath.Ext(attrs.Name)

				// Create a Copier to copy the object
				copier := bckt.Object(newObjectName).CopierFrom(bckt.Object(attrs.Name))
				if _, err := copier.Run(ctx); err != nil {
					log.Fatalf("error copying object: %v", err)
				}

				// Hapus objek lama jika perlu
				if err := bckt.Object(attrs.Name).Delete(ctx); err != nil {
					log.Fatalf("error deleting object: %v", err)
				}
				fmt.Printf("Object renamed to: %s\n", newObjectName)
				break
			}
		}
	}

	if notFound {
		return "", fmt.Errorf("object not found")
	}

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(time.Second * 60),
	}

	url, err := bckt.SignedURL(newObjectName, opts)
	if err != nil {
		return "", fmt.Errorf("bucket(%v).SignedURL: %w", bckt, err)
	}

	return url, nil
}

func (t EventNotulenRepo) DeleteFileGCS(objName string) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading: %v", err)
	}
	ctx := context.Background() // Create a new context

	bckt := t.StorageClient.Bucket(os.Getenv("GC_IMAGE_BUCKET"))
	objectIterator := bckt.Objects(ctx, nil)
	notFound := false

	for {
		attrs, err := objectIterator.Next()
		if err == iterator.Done {
			notFound = true
			break
		}
		if err != nil {
			log.Fatalf("error iterating over objects: %v", err)
		}

		// Jika nama objek mengandung kata kunci pencarian
		if strings.Contains(attrs.Name, objName) {
			objName += filepath.Ext(attrs.Name)

			if attrs.Name == objName {
				fmt.Printf("Found matching object: %s\n", attrs.Name)

				// Hapus objek lama jika perlu
				if err := bckt.Object(attrs.Name).Delete(ctx); err != nil {
					log.Fatalf("error deleting object: %v", err)
				}
				break
			}
		}
	}

	if notFound {
		return fmt.Errorf("object not found")
	}

	return nil
}

func (t EventNotulenRepo) GetFileExtensionFromUrl(rawUrl string) (string, error) {
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
