package domain_test

import (
	"encoder/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {
	video := domain.NewVideo()
	err := video.Validate()

	require.Error(t, err)
}

func TestVideoIdIsnotAUuid(t *testing.T) {
	video := domain.NewVideo()

	video.ID = "any_ID"
	video.ResourceID = "any_resourceID"
	video.FilePath = "any_path"
	video.CreatedAt = time.Now()

	err := video.Validate()

	require.Error(t, err)
}
