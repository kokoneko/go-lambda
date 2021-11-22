package event

import (
	"fmt"
	"log"

	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const uploadBacketName = "lambda-resized299"

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

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(body)
	if err != nil {
		log.Fatal(err)
	}

	img, data, err := image.Decode(buf)
	if err != nil {
		log.Fatal(err)
	}

	switch data {
	case "jpeg", "jpg":
		if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 100}); err != nil {
			log.Fatal(err)
		}
	case "png":
		if err := png.Encode(buf, img); err != nil {
			log.Fatal(err)
		}
	}

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(uploadBacketName),
		Key: aws.String(objectKey),
		Body: buf,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("画像リサイズ完了。実施オブジェクト: %s", objectKey))
}