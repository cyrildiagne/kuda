package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	docker "github.com/docker/docker/builder/dockerignore"
	fileutils "github.com/docker/docker/pkg/fileutils"
)

// Tar takes a source and variable writers and walks 'source' writing each file
// found to the tar writer.
// Skip patterns from dockerignore file if found.
// From https://gist.github.com/sdomino/e6bc0c98f87843bc26bb
func Tar(src string, tarName string, writer io.Writer, dockerignore io.Reader) error {

	// Ensure the src actually exists before trying to tar it.
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("Unable to tar files - %v", err.Error())
	}

	gzw := gzip.NewWriter(writer)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	excludes := []string{}
	if dockerignore != nil {
		ex, err := docker.ReadAll(dockerignore)
		if err != nil {
			return err
		}
		excludes = ex
	}

	// Go through each file.
	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Return if file matches any pattern in the .dockerignore file.
		if rm, _ := fileutils.Matches(file, excludes); rm {
			// fmt.Println("Ignore", file)
			return nil
		}
		// Return on non-regular files.
		if !fi.Mode().IsRegular() {
			return nil
		}
		// Return if file matches the output tar file name.
		if fi.Name() == tarName {
			return nil
		}
		// Create a new dir/file header.
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}
		// Update the name to correctly reflect the desired destination when untaring.
		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))
		// Write the header.
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// Open files for taring.
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		// Copy file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}
		// Manually close here after each file operation; defering would cause each file close
		// to wait until all operations have completed.
		f.Close()

		return nil
	})
}

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func Untar(dst string, r io.Reader) error {

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}
