package conf

import (
	"sync"

	"github.com/AdRoll/goamz/aws"
	"github.com/AdRoll/goamz/s3"
	"github.com/timehop/goth/env"
)

var (
	once sync.Once

	TimehopUploadsS3Bucket *s3.Bucket
)

func SetupAll() {
	once.Do(func() {
		awsCredentials := aws.Auth{
			AccessKey: env.MandatoryVar("AWS_ACCESS_KEY_ID"),
			SecretKey: env.MandatoryVar("AWS_SECRET_ACCESS_KEY"),
		}
		TimehopUploadsS3Bucket = s3.New(awsCredentials, aws.USEast).Bucket("timehop.uploads")
	})
}
