package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	// "github.com/cloudinary/cloudinary-go/v2/api"
	// "github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cld *cloudinary.Cloudinary
var cloudName string
var apiSecret string
var apiKey string
var baseFolder = "daterrr/"

func init() {
	configPath := "../"

	cfg, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("CRITICAL ERROR: Failed to load config from '%s'. Error: %v", configPath, err)
	}
	cloudName = cfg.CloudinaryCloudName
	apiKey = cfg.CloudinaryAPIKey
	apiSecret = cfg.CloudinaryAPISecret
	// Add config
	cld, err = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Printf("Failed to initialize Cloudinary client %v", err)
		panic(err)
	}

}

// Returns a pointer to a boolean variable
func BoolAddress(b bool) *bool {
	boolVar := b
	return &boolVar
}

// In utils/utils.go
func UploadImage(file io.Reader, key string, folder string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // Ensure the cancel function is called to release resources

	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:  key,
		Overwrite: BoolAddress(true),
		Faces:     BoolAddress(true),
	})

	log.Printf("%v", resp)
    if err != nil {
        log.Printf("Cloudinary upload failed for key %s: %v", key, err)
        return "", fmt.Errorf("cloudinary upload failed: %w", err) // Wrap the error for better debugging
    }
    // Add checks for nil response or empty URL
    if resp == nil {
        return "", fmt.Errorf("cloudinary upload response was nil for key %s", key)
    }
    if resp.SecureURL == "" {
        log.Printf("Cloudinary upload successful for key %s, but SecureURL is empty", key)
        return "", fmt.Errorf("cloudinary upload successful, but SecureURL is empty for key %s", key)
    }
    return resp.SecureURL, nil
}