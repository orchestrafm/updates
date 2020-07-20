package objstore

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spidernest-go/logger"
)

var session_ *session.Session

func Login() {
	sess, err := session.NewSession(
		&aws.Config{
			Region:   aws.String(os.Getenv("AWS_S3_REGION")),
			Endpoint: aws.String(os.Getenv("AWS_ENDPOINT")),
		})
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("Object Storage session failed to initialize.")
	}

	session_ = sess
}
