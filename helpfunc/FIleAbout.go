package helpfunc

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
)

/**
用于读取文件的md5值
*/
func ComputeMD5(file io.Reader) (string, error) {
	hash := md5.New()
	chunk := make([]byte, 4096) // 4KB chunk size
	for {
		bytesRead, err := file.Read(chunk)
		if err != nil {
			if err != io.EOF {
				return "", fmt.Errorf("读取错误: %v", err)
			}
			break
		}
		hash.Write(chunk[:bytesRead])
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
