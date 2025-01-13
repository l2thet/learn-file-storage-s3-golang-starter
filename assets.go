package main

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func getAssetPath(videoID uuid.UUID, mediaType string) (string, error) {
	ext, err := parseContentTypeToExtension(mediaType)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", videoID, ext), nil
}

func (cfg apiConfig) getAssetDiskPath(assetPath string) string {
	return filepath.Join(cfg.assetsRoot, assetPath)
}

func (cfg apiConfig) getAssetURL(assetPath string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, assetPath)
}

func parseContentTypeToExtension(contentType string) (string, error) {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", errors.Wrap(err, "parsing Content-Type")
	}

	parsedContentType := strings.Split(mediaType, "/")
	if len(parsedContentType) < 2 || parsedContentType[0] != "image" {
		return "", errors.Errorf("invalid Content-Type: %s", contentType)
	}

	switch parsedContentType[1] {
	case "jpeg":
		return ".jpg", nil
	case "png":
		return ".png", nil
	default:
		return "", errors.Errorf("unsupported image format: %s", contentType)
	}
}
