package main

import (
	// Modules in GOROOT
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	// External modules
)

func calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	md5Hash := md5.New()
	if _, err := io.Copy(md5Hash, file); err != nil {
		return "", err
	}

	md5Checksum := hex.EncodeToString(md5Hash.Sum(nil))
	return md5Checksum, nil
}

func calculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	sha256Hash := sha256.New()
	if _, err := io.Copy(sha256Hash, file); err != nil {
		return "", err
	}

	sha256Checksum := hex.EncodeToString(sha256Hash.Sum(nil))
	return sha256Checksum, nil
}
