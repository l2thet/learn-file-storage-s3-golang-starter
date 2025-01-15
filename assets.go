package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func getAssetPath(mediaType string) (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic("Couldn't generate random bytes")
	}

	fileName := base64.RawURLEncoding.EncodeToString(randomBytes)

	ext, err := parseImageContentTypeToExtension(mediaType)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", fileName, ext), nil
}

func (cfg apiConfig) getAssetDiskPath(assetPath string) string {
	return filepath.Join(cfg.assetsRoot, assetPath)
}

func (cfg apiConfig) getLocalAssetURL(assetPath string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, assetPath)
}

func (cfg apiConfig) getS3AssetURL(assetPath string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", cfg.s3Bucket, cfg.s3Region, assetPath)
}

func parseImageContentTypeToExtension(contentType string) (string, error) {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", errors.Wrap(err, "parsing Content-Type")
	}

	parsedContentType := strings.Split(mediaType, "/")
	if len(parsedContentType) < 2 {
		return "", errors.Errorf("invalid Content-Type: %s", contentType)
	}

	switch parsedContentType[1] {
	case "jpeg":
		return ".jpg", nil
	case "png":
		return ".png", nil
	case "mp4":
		return ".mp4", nil
	default:
		return "", errors.Errorf("unsupported image format: %s", contentType)
	}
}
