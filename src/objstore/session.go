package objstore

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spidernest-go/logger"
)

func Login() *session.Session {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(os.Getenv("AWS_S3_REGION")),
		})
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("Object Storage session failed to initialize.")
	}

	return sess
}
