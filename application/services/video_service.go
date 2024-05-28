package services

import (
	"encoder/application/repositories"
	"encoder/domain"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type VideoService struct {
	Video            *domain.Video
	VideoRepository  repositories.VideoRepository
	localStoragePath string
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucketName string) error {
	// Create a file to write the S3 Object contents to.
	fileName := os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".mp4"

	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file %q, %v", fileName, err)
	}

	// The session the S3 Downloader will use
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))

	_, errCre := sess.Config.Credentials.Get()
	if errCre != nil {
		return fmt.Errorf("failed to credential, %v", err)
	}

	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	// Write the contents of S3 Object to the file
	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(v.Video.FilePath),
	})

	if err != nil {
		return fmt.Errorf("failed to download file, %v", err)
	}

	fmt.Printf("file downloaded, %d bytes\n", n)

	defer f.Close()

	return nil
}

func (v *VideoService) Fragment() error {

	err := os.Mkdir(os.Getenv("localStoragePath")+"/"+v.Video.ID, os.ModePerm)
	if err != nil {
		return err
	}

	source := v.getLocalStoragePath(v.Video.ID, "mp4")
	target := v.getLocalStoragePath(v.Video.ID, "frag")

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("=====> Output: %s\n", string(out))
	}
}

func (v *VideoService) getLocalStoragePath(videoId string, videoType string) string {
	if v.localStoragePath == "" {
		v.localStoragePath = os.Getenv("localStoragePath")
	}

	return v.localStoragePath + "/" + videoId + "." + videoType
}
