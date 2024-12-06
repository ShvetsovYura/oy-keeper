package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// вычисляет hash-строку
func Hash(value []byte) string {
	h := md5.New()
	h.Write(value)
	res := h.Sum(nil)

	return hex.EncodeToString(res)
}

// получает хэш md5 от файла
func MD5Sum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
