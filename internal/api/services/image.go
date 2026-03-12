package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rangodisco/yhar/internal/api/models"
	"github.com/rangodisco/yhar/internal/api/repositories"
)

type ImageService struct {
	repo *repositories.ImageRepository
}

func NewImageService(repo *repositories.ImageRepository) *ImageService {
	return &ImageService{repo: repo}
}

// GetOrCreate looks for the url in database, if it doesn't exist, creates, persists and returns the image
func (s *ImageService) GetOrCreate(ctx context.Context, url string) (*models.Image, error) {
	if url == "" {
		return nil, fmt.Errorf("url is empty")
	}

	existingImage, err := s.repo.FindActiveByUrl(ctx, url)
	if err == nil && existingImage.Path != "" {
		return existingImage, nil
	}
	filename, err := s.SaveLocally(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("unable to save image locally: %w", err)
	}

	model := &models.Image{
		Path:       filename,
		Type:       "local",
		ContentURL: url,
	}

	err = s.repo.Persist(ctx, model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

// SaveLocally saves an distant image to our local public/img dir
func (s *ImageService) SaveLocally(ctx context.Context, url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("unable to fetch image : %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("unable to close body: %w", err)
		}
	}(res.Body)

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read image : %w", err)
	}

	// Probably overkill, just needed a random string
	rawHash := sha256.Sum256(data)
	hash := hex.EncodeToString(rawHash[:])[:16]
	extension := s.getExtension(res.Header.Get("Content-Type"))
	filename := fmt.Sprintf("%s.%s", hash, extension)

	err = os.MkdirAll("public/img", 0755)
	if err != nil {
		return "", fmt.Errorf("unable to create dir: %w", err)
	}

	path := fmt.Sprintf("public/img/%s", filename)

	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("unable to create file: %w", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("unable to close file: %w", err)
		}
	}(file)

	_, err = file.Write(data)
	if err != nil {
		return "", fmt.Errorf("unable to write file: %w", err)
	}

	return filename, nil
}

// getExtension returns the file extension based on the content type
func (s *ImageService) getExtension(contentType string) string {
	switch contentType {
	case "image/jpg", "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	case "image/webp":
		return "webp"
	default:
		return "jpg"
	}
}
