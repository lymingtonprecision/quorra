package ova

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type Ova struct {
	SourcePath string
	LocalPath  string
	TempFile   bool
}

type OvaEntry struct {
	io.Reader
	f    *os.File
	Size int64
}

func (e *OvaEntry) Close() error {
	return e.f.Close()
}

func Open(path string) (*Ova, error) {
	urlre := regexp.MustCompile("(?i)^(https?|ftp)://")

	if urlre.MatchString(path) {
		return fromURL(path)
	} else {
		return fromLocalFile(path)
	}
}

func fromURL(path string) (*Ova, error) {
	/*ova := Ova{
		SourcePath: path,
		LocalPath:  "",
		TempFile:   true,
	}*/

	return nil, os.ErrNotExist
}

func fromLocalFile(path string) (*Ova, error) {
	ova := Ova{
		SourcePath: path,
		LocalPath:  path,
		TempFile:   false,
	}

	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	return &ova, nil
}

func (ova *Ova) Close() error {
	if !ova.TempFile {
		return nil
	}

	return os.Remove(ova.LocalPath)
}

func (ova *Ova) GetFile(name string) (io.ReadCloser, error) {
	f, err := os.Open(ova.LocalPath)
	if err != nil {
		return nil, err
	}

	r := tar.NewReader(f)

	for {
		h, err := r.Next()

		if err == io.EOF {
			break
		}

		match, err := filepath.Match(name, h.Name)
		if err != nil {
			return nil, err
		}

		if match {
			return &OvaEntry{r, f, h.Size}, nil
		}
	}

	_ = f.Close()

	return nil, os.ErrNotExist
}

func (ova *Ova) descriptorFileName() string {
	p := strings.Replace(ova.LocalPath, "\\", "/", -1)
	return strings.TrimSuffix(path.Base(p), path.Ext(p)) + "*.ovf"
}

func (ova *Ova) ReadDescriptor() (string, error) {
	f, err := ova.GetFile(ova.descriptorFileName())
	if err != nil {
		return "", fmt.Errorf(
			"failed reading descriptor %s from %s: %s",
			ova.descriptorFileName(),
			ova.LocalPath,
			err,
		)
	}

	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
