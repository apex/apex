package archive

import (
	"archive/zip"
	"compress/flate"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// NewZip creates compressed (deflate) zip archive.
func NewZip(dest io.Writer) *Zip {
	writer := zip.NewWriter(dest)

	writer.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.DefaultCompression)
	})

	return &Zip{writer: writer}
}

// Zip represents zip archive.
type Zip struct {
	writer *zip.Writer
	lock   sync.Mutex
}

// AddBytes add bytes to archive.
func (z *Zip) AddBytes(path string, contents []byte) error {
	z.lock.Lock()
	defer z.lock.Unlock()

	header := &zip.FileHeader{
		Name:   path,
		Method: zip.Deflate,
	}

	header.SetModTime(time.Unix(0, 0))

	zippedFile, err := z.writer.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = zippedFile.Write(contents)
	return err
}

// AddFile adds a file to archive.
// AddFile resets mtime.
func (z *Zip) AddFile(path string, file *os.File) error {
	path = strings.Replace(path, "\\", "/", -1)

	z.lock.Lock()
	defer z.lock.Unlock()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	if !info.Mode().IsRegular() {
		return errors.New("Only regular files supported: " + path)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = path
	header.Method = zip.Deflate
	header.SetModTime(time.Unix(0, 0))

	zippedFile, err := z.writer.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(zippedFile, file)
	return err
}

// AddDir to target path in archive. This function doesn't follow symlinks.
func (z *Zip) AddDir(root, target string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		archivePath := filepath.Join(target, rel)
		return z.AddFile(archivePath, file)
	})
}

// Close Zip writer.
func (z *Zip) Close() error {
	return z.writer.Close()
}
