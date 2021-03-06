package objstore

import (
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spidernest-go/logger"
)

func Upload(f io.Reader, fname string, acl string, cdn bool) (string, error) {
	uploader := s3manager.NewUploader(session_)

	out, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:          aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:             aws.String(fname),
		Body:            f,
		ACL:             aws.String(acl),
		ContentEncoding: aws.String("brotli"),
	})

	if err != nil {
		logger.Error().
			Err(err).
			Msg("File was unable to be put into Object Storage.")

		return "", err
	}

	//HACK: there should be a better way to get the cdn url
	if cdn == true {
		out.Location = strings.Replace(out.Location,
			"nyc3.digitaloceanspaces.com",
			"nyc3.cdn.digitaloceanspaces.com",
			1)
	}

	return out.Location, nil
}
