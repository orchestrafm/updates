package objstore

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spidernest-go/logger"
)

func Upload(s *session.Session, f io.Reader, fname string) (string, error) {
	uploader := s3manager.NewUploader(s)

	out, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String(os.Getenv(fname)),
		Body:   f,
	})

	if err != nil {
		logger.Error().
			Err(err).
			Msg("File was unable to be put into Object Storage.")

		return "", err
	}

	return out.Location, nil
}
