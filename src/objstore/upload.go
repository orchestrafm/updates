package objstore

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spidernest-go/logger"
)

func Upload(f io.Reader, fname string) (string, error) {
	uploader := s3manager.NewUploader(session_)

	out, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String(fname),
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
