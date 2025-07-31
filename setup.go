package main

import (
	// "crypto/md5"
	"database/sql"
	"fmt"
	"image"
	_ "image/jpeg"

	// "io"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type ImageData struct {
	Name        string
	Path        string
	Http        string
	Idx         int
	Orientation string
}

func img_orient(imgPath string) (string, error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return "", err
	}

	if config.Width > config.Height {
		fmt.Println("Landscape")
		return "landscape", nil
	} else if config.Width < config.Height {
		fmt.Println("Portrait")
		return "portrait", nil
	} else {
		fmt.Println("Square")
		return "square", nil
	}
}

func create_img_db_table(dbpath string) {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		fmt.Println("Failed to open database:", err)
		return
	}
	defer db.Close()

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS images (
		Name TEXT,
		Path TEXT,
		Http TEXT,
		Idx INTEGER,
		Orientation TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		fmt.Println("Failed to create table:", err)
		return
	}
}

func create_http_path(fpath string) string {
	return strings.Replace(fpath, "/home/pimedia/Pictures/", "/static/", 1)
}

func Walk_Img_Dir(dbpath string, dir string) error {
	idx := 0

	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return err
	}
	defer db.Close()

	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		idx += 1
		ext := filepath.Ext(strings.ToLower(d.Name()))
		if ext == ".jpg" {
			orientation, orientErr := img_orient(path)
			if orientErr != nil {
				return orientErr
			}

			imageData := ImageData{
				Name:        d.Name(),
				Path:        path,
				Http:        create_http_path(path),
				Idx:         idx,
				Orientation: orientation,
			}
			fmt.Println(imageData)

			insertSQL := `INSERT INTO images (Name, Path, Http, Idx, Orientation) VALUES (?, ?, ?, ?, ?)`
			_, err = db.Exec(insertSQL, imageData.Name, imageData.Path, imageData.Http, imageData.Idx, imageData.Orientation)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}

func main() {
	dbpath := "/home/pimedia/imagesDB"
	imagedir := "/home/pimedia/Pictures/test"
	create_img_db_table(dbpath)
	Walk_Img_Dir(dbpath, imagedir)
}
