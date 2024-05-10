package domain_test

import (
	"encoder/domain"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestNewJob(t *testing.T) {
	video := domain.NewVideo()

	video.ID = uuid.NewV4().String()
	video.ResourceID = "any_resourceID"
	video.FilePath = "any_path"
	video.CreatedAt = time.Now()

	job, err := domain.NewJob(video.FilePath, "Converted", video)

	require.Nil(t, err)
	require.NotNil(t, job)
}
