package helpers

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func SetupCloudinary() (*cloudinary.Cloudinary, error) {
	env := NewEnv()
	cldName := env.Cloudinary.CloudName
	cldKey := env.Cloudinary.ApiKey
	cldSecret := env.Cloudinary.ApiSecret

	cld, err := cloudinary.NewFromParams(cldName, cldKey, cldSecret)
	if err != nil {
		return nil, err
	}

	return cld, nil
}

func UploadToCloudinary(file interface{}, filePath string) (string, error) {
	env := NewEnv()
	assetFolder := env.Cloudinary.AssetFolder
	ctx := context.Background()
	cloudinary, err := SetupCloudinary()
	if err != nil {
		return "", err
	}

	uploadParams := uploader.UploadParams{
		PublicID:       filePath,
		Overwrite:      api.Bool(true),
		Folder:         assetFolder,
		AllowedFormats: []string{"jpg", "png", "webp"},
	}

	result, err := cloudinary.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", err
	}

	imageURL := result.SecureURL

	return imageURL, nil
}

func DeleteFromCloudinary(publicID string) error {
	ctx := context.Background()
	cloudinary, err := SetupCloudinary()
	if err != nil {
		return err
	}

	_, err = cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return err
	}

	return nil
}
