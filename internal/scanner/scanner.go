package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/yazdun/bitc-go/internal/compress"
)

var imageExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".bmp":  true,
	".tiff": true,
	".tif":  true,
	".webp": true,
	".avif": true,
	".heic": true,
	".heif": true,
}

func ScanDirectory(dir string) ([]compress.ImageFile, error) {
	var images []compress.ImageFile

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if imageExtensions[ext] {
			images = append(images, compress.ImageFile{
				Path:     path,
				Name:     filepath.Base(path),
				OrigSize: info.Size(),
			})
		}
		return nil
	})

	return images, err
}
