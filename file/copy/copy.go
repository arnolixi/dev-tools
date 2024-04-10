package copy

import (
	"github.com/arnolixi/dev-tools/file"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CP(src, dst string) error {
	src = strings.TrimRight(src, "/")
	dst = strings.TrimRight(dst, "/")
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if src == path {
			if !info.IsDir() {
				srcSplit := strings.Split(path, "/")
				var newPath = dst
				if file.IsDir(dst) {
					newPath = filepath.Join(newPath, srcSplit[len(srcSplit)-1])
				}
				log.Println(path, newPath)
				err = copyFile(path, newPath)
				if err != nil {
					return err
				}
			} else {
				if err = os.MkdirAll(dst, info.Mode()); err != nil {
					return err
				}
			}
			return nil
		}
		relPath := path[len(src)+1:]
		newPath := filepath.Join(dst, relPath)
		if info.IsDir() {
			os.MkdirAll(newPath, info.Mode())
		} else {
			err = copyFile(path, newPath)
			if err != nil {
				return err
			}
			err = os.Chmod(newPath, info.Mode())
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func copyFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}
	err = destFile.Sync()
	if err != nil {
		return err
	}
	return nil
}

func copyDir(src, dest string) error {
	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return err
	}
	for _, file := range files {
		srcPath := filepath.Join(src, file.Name())
		destPath := filepath.Join(dest, file.Name())
		if file.IsDir() {
			err := copyDir(srcPath, destPath)
			if err != nil {
				return err
			}
		} else {

		}
	}
	return nil

}
