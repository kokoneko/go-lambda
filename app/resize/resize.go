package event

import (
	"fmt"
	"log"
	"os"
	"strings"
	"bytes"
	"image/jpeg"
	"image/png"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/disintegration/imaging"
)

const (
	squarePrefix = "sq"
	squareSize = 300
)

func ExecResize(bucketName string, objectKey string) {
	log.Println(fmt.Sprintf("画像リサイズ開始。対象オブジェクト: %s", objectKey))

	sess := session.Must(session.NewSession())

	s3Sdk := s3.New(sess)
	obj, err := s3Sdk.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(objectKey),
	})
	if err != nil {
		log.Fatal(err)
	}

	body := obj.Body
	defer body.Close()

	// 元画像読み込み
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(body)
	if err != nil {
		log.Fatal(err)
	}

	// 画像の向きを自動調整してデコード
	img, err := imaging.Decode(buf, imaging.AutoOrientation(true))
	if err != nil {
		log.Fatal(err)
	}

	// 短辺の方に合わせて正方形にリサイズ
    imgRectangle := img.Bounds()
	size := imgRectangle.Max.X
	if imgRectangle.Max.Y < size {
		size = imgRectangle.Max.Y
	}

	triming := imaging.CropAnchor(img, size, size, imaging.Center)
	resizedImg := imaging.Resize(triming, squareSize, squareSize, imaging.NearestNeighbor)

	// 画像のエンコード（書き込み）
	ext := filepath.Ext(objectKey)
	switch ext {
	case "jpeg", "jpg":
		if err := jpeg.Encode(buf, resizedImg, &jpeg.Options{Quality: 100}); err != nil {
			log.Fatal(err)
		}
	case "png":
		if err := png.Encode(buf, resizedImg); err != nil {
			log.Fatal(err)
		}
	default:
		if err := png.Encode(buf, resizedImg); err != nil {
			log.Fatal(err)
		}
	}

	// S3リサイズ画像用フォルダにアップロード
	uploader := s3manager.NewUploader(sess)
	uploadKey := squarePrefix + strings.Replace(objectKey, os.Getenv("READ_PREFIX"), "", 1)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key: aws.String(os.Getenv("UPLOAD_PREFIX") + uploadKey),
		Body: buf,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("画像リサイズ完了。実施オブジェクト: %s", uploadKey))
}