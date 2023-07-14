package tee

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"mime/multipart"
	"net/http"
)

var DefaultMultipartMemory int64

// 32 MB
// 上下文
type Context struct {
	Res           http.ResponseWriter    // 请求信息
	Req           *http.Request          // 返回信息
	RouteParamMap map[string]interface{} // 路径参数
	HandlerSlice  []Handler              // 中间件，控制器方法数组
	Index         int                    // 指定路由的当前执行的方法索引
	StatusCode    int                    // 错误码
}

// 进入对应路由的下一个方法
func (c *Context) Next() {
	c.Index++
	for c.Index < len(c.HandlerSlice) {
		c.HandlerSlice[c.Index](c)
		c.Index++
	}

}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) PathParam(key string) string {
	if v, ok := c.RouteParamMap[key]; ok {
		return v.(string)
	}
	return ""
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Res.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Res.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Res.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Res)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Res, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Res.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Res.Write([]byte(html))
}

func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	if c.Req.MultipartForm == nil {
		DefaultMultipartMemory = 32 << 20
		if err := c.Req.ParseMultipartForm(DefaultMultipartMemory); err != nil {
			return nil, err
		}
	}
	f, fh, err := c.Req.FormFile(name)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, err
}

//保存文件到指定目录
func (c *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
