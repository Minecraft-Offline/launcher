package main

import (
	//std necessities
	"archive/zip"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

func hideConsole() {
	if runtime.GOOS != "windows" {
		return
	}

	getConsoleWindow := syscall.NewLazyDLL("kernel32.dll").NewProc("GetConsoleWindow")
	if getConsoleWindow.Find() != nil {
		return
	}

	showWindow := syscall.NewLazyDLL("user32.dll").NewProc("ShowWindow")
	if showWindow.Find() != nil {
		return
	}

	hwnd, _, _ := getConsoleWindow.Call()
	if hwnd == 0 {
		return
	}

	showWindow.Call(hwnd, 0)
}

func extract(zipPath, extractDir string, exclusions ...string) error {
	archive, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, file := range archive.File {
		if file.FileInfo().IsDir() {
			continue
		}

		if len(exclusions) > 0 {
			exclude := false
			for _, exclusion := range exclusions {
				log.Trace("extract: checking exclusion ", exclusion, " against file ", file.Name)
				if file.Name == exclusion {
					log.Trace("extract: excluding ", exclusion)
					exclude = true
					break
				}
			}
			if exclude {
				continue
			}
		}
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		if err := fileWriteRC(rc, extractDir+"/"+file.Name); err != nil {
			return err
		}
	}

	return nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func mkdir(dir string) error {
	return os.MkdirAll(dir, 0770)
}

func rm(dir string) {
	dirRead, _ := os.Open(dir)
	dirFiles, _ := dirRead.Readdir(0)

	if len(dirFiles) > 0 {
		for _, file := range dirFiles {
			filePath := dir + "/" + file.Name()
			if file.IsDir() {
				rm(filePath)
			}
			os.Remove(filePath)
		}
	}

	os.Remove(dir)
}

func fileMove(src, dst string) error {
	return os.Rename(src, dst)
}

func fileCreate(path string) (*os.File, error) {
	err := mkdir(filepath.Dir(path))
	if err != nil {
		return nil, err
	}

	file, err := os.Create(path)
	return file, err
}

func fileRead(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	return data, err
}

func fileWrite(data []byte, dst string) error {
	file, err := fileCreate(dst)
	if err != nil {
		return err
	}
	file.Close()

	err = ioutil.WriteFile(dst, data, 0644)
	return err
}

func fileWriteRC(src io.ReadCloser, dst string) error {
	dstFile, err := fileCreate(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, src)
	if err != nil {
		return err
	}

	return nil
}

func fileCopy(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := fileCreate(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	srcHash, err := hash_file_sha1(src)
	if err != nil {
		return err
	}

	dstHash, err := hash_file_sha1(dst)
	if err != nil {
		return err
	}

	if srcHash != dstHash {
		return fmt.Errorf("copy: hash for (\"%s\") doesn't match hash for (\"%s\")", src, dst)
	}

	return nil
}

func download(url, filepath, sha1 string) error {
	//log.Trace("download(", url, ", ", filepath, ", ", sha1, ")")
	if fileExists(filepath) {
		hash, err := hash_file_sha1(filepath)
		if err == nil {
			if hash == sha1 {
				return nil
			}
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := fileCreate(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	if sha1 != "" {
		hash, err := hash_file_sha1(filepath)
		if err != nil {
			return err
		}

		if hash != sha1 {
			return fmt.Errorf("download: hash for (\"%s\") doesn't match hash for (\"%s\")", url, hash)
		}
	}

	return nil
}

func hash_file_sha1(filePath string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnSHA1String string

	//Open the filepath passed by the argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnSHA1String, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new SHA1 hash interface to write to
	hash := sha1.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnSHA1String, err
	}

	//Get the 20 bytes hash
	hashInBytes := hash.Sum(nil)[:20]

	//Convert the bytes to a string
	returnSHA1String = hex.EncodeToString(hashInBytes)

	return returnSHA1String, nil

}
