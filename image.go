package main

import (
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	imageIDLength = 10
)

// A map of accepted mime types and their file extension
var mimeExtension = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/gif":  ".gif",
}

type Image struct {
	ID          string `bson:"_id" json:"id"`
	UserID      string `bson:"user_id"`
	Name        string
	Location    string
	Size        int64
	CreatedAt   time.Time `bson:"created_at"`
	Description string
}

func NewImage(user *User) *Image {
	return &Image{
		ID:        GenerateID("img", imageIDLength),
		UserID:    user.ID,
		CreatedAt: time.Now(),
	}
}

func (image *Image) CreatedFromFile(file multipart.File, headers *multipart.FileHeader) error {
	// Move our file to an appropriate place, with an appropriate name
	image.Name = headers.Filename
	image.Location = image.ID + filepath.Ext(image.Name)

	// Open a file at the target location
	savedFile, err := os.Create("./data/images/" + image.Location)
	if err != nil {
		return err
	}
	defer savedFile.Close()

	// Copy the uploaded file to the target location
	size, err := io.Copy(savedFile, file)
	if err != nil {
		return err
	}
	image.Size = size
	// Save the image to the database
	db := NewDBImageStore()
	defer db.Close()
	return db.Save(image)
}

func (image *Image) CreatedFromURL(imageUrl string) error {
	// Get response from url
	response, err := http.Get(imageUrl)
	if err != nil {
		return err
	}

	// Make sure we have a response
	if response.StatusCode != http.StatusOK {
		return errImageURLInvalid
	}
	defer response.Body.Close()

	// Ascertain the type of file we downloaded
	mimeType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		return errInvalidImageType
	}

	// Get an extension for the file
	ext, valid := mimeExtension[mimeType]
	if !valid {
		return errInvalidImageType
	}

	// Get a name from the URL
	image.Name = filepath.Base(imageUrl)
	image.Location = image.ID + ext

	// Open a file at target location
	savedFile, err := os.Create("./data/images/" + image.Location)
	if err != nil {
		return err
	}
	defer savedFile.Close()

	// Copy the entire response to the output file
	size, err := io.Copy(savedFile, response.Body)
	if err != nil {
		return err
	}

	// The returned value from io.Copy is the number of bytes copied
	image.Size = size

	// Save our image to the store
	db := NewDBImageStore()
	defer db.Close()
	return db.Save(image)
}

// This method automatically called by html template
func (image *Image) StaticRoute() string {
	return "/im/" + image.Location
}

// This method automatically called by html template
func (image *Image) ShowRoute() string {
	return "/image/" + image.ID
}
