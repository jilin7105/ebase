package helpfunc

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

//大文件上传
func UploadBigFileWithParams(url, filePath string, data, headers map[string]string) (string, error) {

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()
	//videoSignature, _ := ComputeMD5(file)
	//
	//file.Seek(0, 0)
	stat, _ := file.Stat()
	byteBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(byteBuf)

	// 添加额外的参数到multipart.Writer的请求体中
	for key, value := range data {
		sendkey := key
		sendValue := value
		//if key == "md5" {
		//	sendkey = value
		//	sendValue = videoSignature
		//}
		if err := writer.WriteField(sendkey, sendValue); err != nil {
			return "", fmt.Errorf("无法写入参数: %v", err)
		}

	}

	// part: file
	_, err = writer.CreateFormFile("video_file", stat.Name())
	if err != nil {
		return "", err
	}

	//contentType := writer.FormDataContentType()

	nheader := byteBuf.Len()
	header := make([]byte, nheader)
	_, _ = byteBuf.Read(header)

	// part: latest boundary
	// when multipart closed, latest boundary is added
	writer.Close()
	nboundary := byteBuf.Len()
	boundary := make([]byte, nboundary)
	_, _ = byteBuf.Read(boundary)

	// calculate content length
	totalSize := int64(nheader) + stat.Size() + int64(nboundary)

	//use pipe to pass request
	rd, wr := io.Pipe()
	defer rd.Close()

	go func() {
		defer wr.Close()
		// write multipart
		_, _ = wr.Write(header)

		// write file
		buf := make([]byte, 16*1024)
		for {
			n, err := file.Read(buf)
			if err != nil {
				break
			}
			_, _ = wr.Write(buf[:n])
		}

		// write boundary
		_, _ = wr.Write(boundary)
	}()

	// 完成multipart.Writer的写入
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("无法关闭multipart.Writer: %v", err)
	}

	// 创建HTTP请求
	request, err := http.NewRequest("POST", url, rd)
	if err != nil {
		return "", fmt.Errorf("无法创建HTTP请求: %v", err)
	}

	// 设置请求头
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	// 设置Content-Type头为multipart/form-data
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Content-Length", fmt.Sprintf("%v", totalSize))
	// 发送请求并获取响应
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("无法发送HTTP请求: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", errors.New("HTTP请求返回了非200状态码: " + response.Status)
	}

	// 读取响应数据
	var responseBody bytes.Buffer
	_, err = io.Copy(&responseBody, response.Body)
	if err != nil {
		return "", fmt.Errorf("无法读取响应数据: %v", err)
	}

	return responseBody.String(), nil
}
