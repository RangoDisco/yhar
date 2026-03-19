package services

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"

	_ "image/jpeg"

	"golang.org/x/image/draw"

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

	filename, err := s.saveLocally(ctx, url)
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

// saveLocally saves a distant image to our local images dir
func (s *ImageService) saveLocally(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("unable to build req: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("unable to fetch img")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("unable to close body: %v", err)
		}
	}(res.Body)

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read image: %w", err)
	}

	// Probably overkill, just needed a random string
	rawHash := sha256.Sum256(data)
	hash := hex.EncodeToString(rawHash[:])[:16]
	filename := fmt.Sprintf("%s.%s", hash, "jpg")

	err = os.MkdirAll("images", 0755)
	if err != nil {
		return "", fmt.Errorf("unable to create dir: %w", err)
	}

	path := fmt.Sprintf("images/%s", filename)

	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("unable to create file: %w", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("unable to close file: %v", err)
		}
	}(file)

	img, err := s.resizeImg(data)
	if err != nil {
		return "", fmt.Errorf("unable to resize img: %w", err)
	}

	err = jpeg.Encode(file, img, nil)
	if err != nil {
		return "", fmt.Errorf("unable to write file: %w", err)
	}

	return filename, nil
}

// resizeImg resizes img to 500/500 (which is later encoded in jpeg)
func (s *ImageService) resizeImg(data []byte) (*image.RGBA, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("unable to decode data into image: %w", err)
	}

	dst := image.NewRGBA(image.Rect(0, 0, 500, 500))
	draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)

	return dst, nil
}
