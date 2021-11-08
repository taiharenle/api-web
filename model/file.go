package model

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

var (
	endpoint        = os.Getenv("MINIO_ENDPOINT")      // URL to object storage service.
	bucketName      = os.Getenv("MINIO_BUCKET_NAME")   // 桶名称
	accessKeyID     = os.Getenv("MINIO_ROOT_USER")     // 用户名
	secretAccessKey = os.Getenv("MINIO_ROOT_PASSWORD") // 密码
	location        = os.Getenv("MINIO_LOCATION")      // minio时区
)

type Storage struct {
	gorm.Model
	Hash string
	Name string
	Type string
	Size int64
	Data string
}

func (s *Storage) Create(hash, name, types, data string, size int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		s.Hash = hash
		s.Name = name
		s.Type = types
		s.Data = data
		s.Size = size
		return tx.Create(s).Error
	})
}

type Hash struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
}

func GetUUID() string {
	u2 := uuid.NewV4()
	return u2.String()
}

func UploadFile(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file[]"]
	hashes := make([]Hash, 0)
	for _, file := range files {
		objectName := GetUUID()
		hash, err := UploadToMinio(objectName, file, "")
		if err != nil {
			log.Println("Failed upload file")
		}
		hashes = append(hashes, Hash{Name: file.Filename, Hash: hash})
	}
	c.JSON(http.StatusOK, gin.H{
		"hashes": hashes,
	})
}

func UploadToMinio(objectName string, file *multipart.FileHeader, contentType string) (hash string, err error) {
	ctx := context.Background()
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		// Use http or https
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Make a new bucket.
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	src, err := file.Open()
	if err != nil {
		log.Fatalln(err)
	}
	defer src.Close()

	// 使用PutObject上传文件
	info, err := minioClient.PutObject(ctx, bucketName, objectName, src, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	hash = objectName
	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)

	fileNameWithSuffix := path.Base(file.Filename)
	fileType := path.Ext(fileNameWithSuffix)

	if err := DB.Create(
		&Storage{
			Hash: objectName,
			Name: file.Filename,
			Type: "application/" + fileType[1:],
			Size: info.Size,
		}).Error; err != nil {
		log.Printf("Insert file info to Storage err")
	}

	return hash, err
}

func DownloadFile(c *gin.Context) {
	// url中的hash
	fileName := c.Param("hash")
	s := Storage{}
	if err := DB.Where("hash = ?", fileName).First(&s).Error; err != nil {
		log.Printf("This hash has no corresponding attachment")
	}

	ctx := context.Background()
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		// Use http or https
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// 取minio地址
	expires := time.Second * 24 * 60 * 60
	reqParams := make(url.Values)
	fileURL, err := minioClient.PresignedGetObject(ctx, bucketName, fileName, expires, reqParams)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := http.Get(fmt.Sprintf("%s", fileURL))
	if err != nil {
		log.Fatalln(err)
	}

	c.Writer.Header().Add("Content-Disposition", "filename=\""+s.Name+"\"")
	io.Copy(c.Writer, res.Body)

	return
}
