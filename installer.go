package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/ioprogress"
)

const (
	// Release name for each OS/Arch
	MacOSX64  = "macosx64-binary"
	Linux32   = "linux32-binary"
	Linux64   = "linux64-binary"
	Windows32 = "windows32-exe"
	Windows64 = "windows64-exe"

	// File extension used for packaging cf
	Tgz = "tgz"
	Zip = "zip"

	CfBinary = "cf"
	CfExe    = "cf.exe"
)

type Installer struct {
	URL string
	Ext string

	OutStream io.Writer
}

func NewInstaller(version string) (*Installer, error) {
	var release, ext string

	switch runtime.GOOS {
	case "darwin":
		ext = Tgz
		release = MacOSX64
	case "linux":
		ext = Tgz
		release = Linux64
		if runtime.GOARCH == "386" {
			release = Linux32
		}
	case "windows":
		ext = Zip
		release = Windows64
		if runtime.GOARCH == "386" {
			release = Windows32
		}
	default:
		return nil, fmt.Errorf("no binary is available for your OS: %s (%s)",
			runtime.GOOS, runtime.GOARCH)
	}

	// Construct url from OS/Arch information
	url := fmt.Sprintf(
		"https://cli.run.pivotal.io/stable?release=%s&version=%s",
		release, version)

	return &Installer{
		Ext: ext,
		URL: url,

		OutStream: ioutil.Discard,
	}, nil
}

// Install installs URL target to path
func (i *Installer) Install(savePath string) error {

	tempDir, err := ioutil.TempDir("", "cf-plugin")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	tempPath := filepath.Join(tempDir, "cf."+i.Ext)
	fw, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	// After downloading, tempPath will be opened again
	// to uncompress and extract binary from there.
	// Because of this, `defer` is not used here and Close()
	// is explicitly called when needed.

	Debugf("Start downloading from %s", i.URL)
	client := http.DefaultClient
	res, err := client.Get(i.URL)
	if err != nil {
		fw.Close()
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fw.Close()
		return fmt.Errorf("Invalid status code: %d", res.StatusCode)
	}

	progressR := &ioprogress.Reader{
		Reader: res.Body,
		Size:   res.ContentLength,
		DrawFunc: ioprogress.DrawTerminalf(i.OutStream,
			func(progress, total int64) string {
				return fmt.Sprintf("Downloading: %d/%d", progress, total)
			}),
	}

	if _, err := io.Copy(fw, progressR); err != nil {
		fw.Close()
		return err
	}
	fw.Close()

	Debugf("Start extracting from %s", tempPath)
	switch i.Ext {
	case Tgz:
		if err := extractTgz(tempPath, savePath); err != nil {
			return err
		}
	case Zip:
		if err := extractZip(tempPath, savePath); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unexpected file extention")
	}

	return nil
}

// extractTgz extract cf binary from tgz file and save it
// to the given path.
func extractTgz(tgzPath, savePath string) error {
	f, err := os.Open(tgzPath)
	if err != nil {
		return err
	}
	defer f.Close()

	gzipR, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzipR.Close()

	tarR := tar.NewReader(gzipR)
	if err != nil {
		return err
	}

	for {
		hdr, err := tarR.Next()
		if err == io.EOF {
			break
		}

		if hdr.Name != CfBinary {
			continue
		}

		Debugf("Save %s to %s", CfBinary, savePath)
		dst, err := os.OpenFile(
			savePath,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
			hdr.FileInfo().Mode(),
		)

		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, tarR); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("%s is not found on %s", CfBinary, tgzPath)
}

func extractZip(zipPath, savePath string) error {
	zipR, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}

	for _, f := range zipR.File {

		if f.Name != CfExe {
			continue
		}

		Debugf("Save %s to %s", CfExe, savePath)
		dst, err := os.OpenFile(
			savePath,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
			f.Mode(),
		)

		if err != nil {
			return err
		}
		defer dst.Close()

		fr, err := f.Open()
		if err != nil {
			return err
		}
		defer fr.Close()

		if _, err := io.Copy(dst, fr); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("%s is not found on %s", CfExe, zipPath)
}
