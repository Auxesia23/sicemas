package database

import (
	"context"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

func NewCloudinary(ctx context.Context) (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		log.Println("Error creating cloudinary:", err.Error())
		return nil, err
	}
	return cld, nil
}
