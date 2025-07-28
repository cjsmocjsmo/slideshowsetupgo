package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"
	_ "github.com/mattn/go-sqlite3"
)

type ImageData struct {
	Name string
	Path string
	Idx int
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
		Idx INTEGER,
		Orientation TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		fmt.Println("Failed to create table:", err)
		return
	}
}

func Walk_Img_Dir(dbpath string, dir string) error {
	idx := 0

	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return err
	}
	defer db.Close()

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		idx += 1
		ext := filepath.Ext(strings.ToLower(info.Name()))
		if ext == ".jpg" {
			orientation, orientErr := img_orient(path)
			if orientErr != nil {
				return orientErr
			}
		
			imageData := ImageData{
				Name:        info.Name(),
				Path:        path,
				Idx:         idx,
				Orientation: orientation,
			}
			fmt.Println(imageData)
		
			insertSQL := `INSERT INTO images (Name, Path, Idx, Orientation) VALUES (?, ?, ?, ?)`
			_, err = db.Exec(insertSQL, imageData.Name, imageData.Path, imageData.Idx, imageData.Orientation)
			if err != nil {
				return err
			}
		}
		

		return nil
	})

	return nil
}

func main() {
	dbpath := "/home/whitepi/go/slideshowsetupgo/imagesDB"
	imagedir := "/home/whitepi/Pictures/"
	create_img_db_table(dbpath)
	Walk_Img_Dir(dbpath, imagedir)
}

	

