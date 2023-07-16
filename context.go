package tee

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

var DefaultMultipartMemory int64
var ServerIP = localIP()
func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
// 32 MB
// 上下文
type Context struct {
	Res           http.ResponseWriter    // 请求信息
	Req           *http.Request          // 返回信息
	RouteParamMap map[string]interface{} // 路径参数
	HandlerSlice  []Handler              // 中间件，控制器方法数组
	Index         int                    // 指定路由的当前执行的方法索引
	StatusCode    int                    // 错误码
	Errors        error
	Size          int
	sameSite      http.SameSite
	Keys         map[string]any
	mu           sync.RWMutex
}
var (
	TraceIDKey = "x-trace-id"
	StartTimeKey = "x-start-time"
)
func WithTraceID(ctx *Context) string {
	traceID := GetTraceID(ctx)
	if len(traceID) > 0 {
		return traceID
	}
	traceID = NewTraceID()
	
	ctx.Res.Header().Set(TraceIDKey, traceID)
	ctx.Set(TraceIDKey, traceID)
	return traceID
}

// WithTraceID 注入 trace_id
func WithStartTime(ctx *Context, startTime time.Time) {
	ctx.Set(StartTimeKey, startTime)
	return 
}
func (c *Context) Set(key string, value any) {
	c.mu.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]any)
	}

	c.Keys[key] = value
	c.mu.Unlock()
}

// 进入对应路由的下一个方法
func (c *Context) Abort() {
	c.Index = 100
}

// 进入对应路由的下一个方法
func (c *Context) Next() {
	c.Index++
	for c.Index < len(c.HandlerSlice) {
		c.HandlerSlice[c.Index](c)
		c.Index++
	}

}

func (c *Context) ClientIp() string{
	ip := ClientPublicIP(c.Req)
	if ip == ""{
	ip = ClientIP(c.Req)
	}
	return ip

}
func ClientPublicIP(r *http.Request) string {
	var ip string
	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		ip = strings.TrimSpace(ip)
		if ip != ""  {
			return ip
		}
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != ""  {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
	
		return ip

	}
	return ""
	}


func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
	return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
	return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
	return ip
	}
	return ""
	}


func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// GetTraceID 获取用户请求标识
func GetTraceID(ctx *Context) string {
	id, _ := ctx.Get(TraceIDKey)
	return id.(string)
}

func (c *Context) Get(key string) (value any, exists bool) {
	c.mu.RLock()
	value, exists = c.Keys[key]
	c.mu.RUnlock()
	return
}


func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) PathParam(key string) interface{} {
	if v, ok := c.RouteParamMap[key]; ok {
		return v
	}
	return nil
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
	n, err := c.Res.Write([]byte(fmt.Sprintf(format, values...)))
	if err != nil {
		c.Errors = err
	}
	c.Size += n			
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	resbyte , err := json.Marshal(obj)
	if err != nil {
		http.Error(c.Res, err.Error(), 500)
	} 
	c.Size += len(resbyte)
	c.Res.Write(resbyte)
    
	
		
}

// GetRawData returns stream data.
func (c *Context) GetRawData() ([]byte, error) {
	return ioutil.ReadAll(c.Req.Body)
}

// SetSameSite with cookie
func (c *Context) SetSameSite(samesite http.SameSite) {
	c.sameSite = samesite
}

func (c *Context) GetKeys(key string) (value any, exists bool) {
	c.mu.RLock()
	value, exists = c.Keys[key]
	c.mu.RUnlock()
	return
}
// SetCookie adds a Set-Cookie header to the ResponseWriter's headers.
// The provided cookie must have a valid Name. Invalid cookies may be
// silently dropped.
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Res, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: c.sameSite,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

// Cookie returns the named cookie provided in the request or
// ErrNoCookie if not found. And return the named cookie is unescaped.
// If multiple cookies match the given name, only one cookie will
// be returned.
func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.Req.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

func NewTraceID() string {
	return uuid.New().String()
}

func (c *Context) PathParamInt(key string) int {
	v :=  c.PathParam(key)
	if v == nil {
		return 0
	}
	return v.(int)
}
func (c *Context) PathParamInt64(key string) int64 {
	v :=  c.PathParam(key)
	if v == nil {
		return 0
	}
	return v.(int64)
}
func (c *Context) PathParamString(key string) string {
	v :=  c.PathParam(key)
	if v == nil {
		return ""
	}
	return v.(string)
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Size += len(data)
	c.Res.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	resbyte := []byte(html)
	c.Size += len(resbyte)
	c.Res.Write(resbyte)
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
