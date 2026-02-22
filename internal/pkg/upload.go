package pkg

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type File struct {
	Reader   io.Reader
	Filename string
	Src      string
}

func UploadFile(input *File) error {

	if !pathCheck(input.Src) {
		return NewError(KindInternal, "source is required", nil)
	}

	allowedExt := map[string]bool{
		".jpg": true,
		".png": true,
	}

	ext := filepath.Ext(input.Filename)
	if !allowedExt[ext] {
		return NewError(KindBadRequest, "extension must be .jpg or .png", nil)
	}

	input.Filename = fmt.Sprintf("PRODUCT_%s%s", time.Now().Format("20060102_150405"), ext)

	file, err := os.Create(input.Src + input.Filename)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, input.Reader)
	if err != nil {
		return err
	}

	return nil
}

func pathCheck(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Delete(src string) error {

	if !pathCheck(src) {
		log.Println("not found image: ", src)
		return nil
	}

	if err := os.Remove(src); err != nil {
		return err
	}

	return nil
}
