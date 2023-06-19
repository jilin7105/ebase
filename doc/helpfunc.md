# 大文件上传
```go
    imports(
        "github.com/jilin7105/ebase/helpfunc"
		) 
	params, err := helpfunc.UploadBigFileWithParams(
		"https://uploadapi_url_path",
		"/your/file_path",
		map[string]string{
			"params1": "22",
			"params2": "1111",
		},
		map[string]string{
			"Access-Token": "token",
		},
	)
	if err != nil {
		return err
	}

```