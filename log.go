package tee

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger = *logrus.Entry

// Fields fields
type Fields = logrus.Fields

const STARTMICROTIME = "StartMircoTime"

var Loger Logger

func init() {
	Loger = logrus.WithFields(logrus.Fields{

		"server_ip": ServerIP,
	})

	Loger.Logger.SetReportCaller(true)
}

// Get 获取日志实例
func LogGet(c *TeeContext) (log Logger) {
	//非router请求会nil pointer，这里做一下失败处理
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("日志系统发生panic,错误信息:", err)
			log = Loger
		}
	}()

	//获取request信息
	if c == nil || c.Req == nil {
		return Loger
	}
	path := c.Req.URL.Path
	method := c.Req.Method
	raw := c.Req.URL.RawQuery

	//参数获取，兼容各种格式的请求
	request_param := ""
	if method == "GET" {
		request_param = raw
	} else if method == "POST" {
		if ct, ok := c.Req.Header["Content-Type"]; ok && ct[0] == "application/json" {
			raw, _ := c.GetRawData()
			request_param = string(raw)
		} else {
			//form格式
			_ = c.Req.ParseMultipartForm(128)
			data := c.Req.Form
			for k, v := range data {
				if len(v) == 1 {
					request_param += fmt.Sprintf("%s=%v&", k, v[0])
				} else {
					request_param += fmt.Sprintf("%s=%v&", k, v)
				}
			}
		}
	}

	//debug参数获取
	debugKey := "debug"

	//请求时间戳（毫秒数）
	var microTime int64 = 0
	microTimeInterface, exists := c.Get(STARTMICROTIME)
	if exists {
		microTime = microTimeInterface.(int64)
	}

	return Loger.WithFields(Fields{
		"is_gin":        true,
		"debug_key":     debugKey,
		"trace_id":      GetTraceID(c),
		"method":        c.Req.Method,
		"path":          fmt.Sprintf("【%s】%s", method, path),
		"request_param": request_param,
		"microtime":     microTime,
	})
}

type bodyWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LogFormatterParams is the structure any formatter will be handed when time to log comes
type LogFormatterParams struct {
	TraceID string
	// Request Body
	RequestBody string
	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ClientIP equals teeContext's ClientIP method.
	ClientIP string
	// ServerIP equals teeContext's ServerIP method.
	ServerIP string
	// Method is the HTTP method given to the request.
	Method string
	// UserAgent
	UserAgent string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// BodySize is the size of the Response Body
	BodySize int
	// Response Body
	Body string
	// isError
	isError bool
	// Keys are the keys set on the request's teeContext.
	Keys map[string]interface{}
}

func RecoveryWithWriter(c *TeeContext, err interface{}) {
	// Check for a broken connection, as it is not really a
	// condition that warrants a panic stack trace.
	var brokenPipe bool
	if ne, ok := err.(*net.OpError); ok {
		if se, ok := ne.Err.(*os.SyscallError); ok {
			if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
				brokenPipe = true
			}
		}
	}

	httpRequest, _ := httputil.DumpRequest(c.Req, false)
	headers := strings.Split(string(httpRequest), "\r\n")
	for idx, header := range headers {
		current := strings.Split(header, ":")
		if current[0] == "Authorization" {
			headers[idx] = current[0] + ": *"
		}
	}

	//panic转Error级别日志
	//	LogGet(c).WithField("is_panic", true).Error(err)

	// If the connection is dead, we can't write a status to it.
	if brokenPipe {
		c.Errors = err.(error) // nolint: errcheck
		c.Abort()
		return
	} else {
		c.JSON(http.StatusInternalServerError, "抱歉服务错误，请稍后重试")
		c.Abort()
		return
	}
}

// 处理返回数据,TODO:放到上下文中
func (param *LogFormatterParams) getBody(c *TeeContext, b []byte) {
	if param.StatusCode < http.StatusBadRequest {
		return
	}

	param.RequestBody = getRequestBody(c)

	if len(b) == 0 { // 没有数据返回
		return
	}

	param.Body = string(b)

	return
}

func getRequestBody(c *TeeContext) string {
	method := c.Req.Method

	if method == "GET" {
		return ""
	}

	if data, err := c.GetRawData(); err != nil {
		return string(data)
	}

	return ""
}

// Logger is the logrus logger handler
func TeeLogger() Handler {
	return func(c *TeeContext) {

		c.Set(STARTMICROTIME, time.Now().UnixNano()/1000000)

		// defer func() {
		// 	if err := recover(); err != nil {
		// 		RecoveryWithWriter(c, err)
		// 	}
		// }()

		path := c.Req.URL.Path
		method := c.Req.Method
		raw := c.Req.URL.RawQuery
		start := time.Now()

		if raw != "" {
			path = path + "?" + raw
		}

		//Trace
		traceID := WithTraceID(c)
		WithStartTime(c, start)

		c.Next()

		param := &LogFormatterParams{
			TraceID:  traceID,
			Keys:     c.Keys,
			ServerIP: ServerIP,
		}

		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = c.ClientIp()
		param.Method = method
		param.Path = path
		param.StatusCode = c.StatusCode
		if c.Errors != nil {
			param.ErrorMessage = c.Errors.Error()
		}
		param.BodySize = c.Size
		param.UserAgent = c.Req.UserAgent()

		param.getBody(c, []byte(""))

		param.LogFormatter()
	}
}

func (param *LogFormatterParams) LogFormatter() {
	entry := logrus.WithFields(logrus.Fields{
		"trace_id":     param.TraceID,
		"status":       param.StatusCode,
		"latency":      fmt.Sprintf("%13v", param.Latency),
		"client_ip":    param.ClientIP,
		"server_ip":    ServerIP,
		"timestamp":    param.TimeStamp,
		"method":       param.Method,
		"path":         param.Path,
		"size":         param.BodySize,
		"user_agent":   param.UserAgent,
		"request_body": param.RequestBody,
	})

	if len(param.ErrorMessage) > 0 {
		entry.Error(param.ErrorMessage)
		return
	}

	if param.StatusCode >= http.StatusOK && param.StatusCode < http.StatusBadRequest {
		entry.Info(param.Body)
	} else if param.StatusCode >= http.StatusBadRequest && param.StatusCode < http.StatusInternalServerError {
		entry.Warn(param.Body)
	} else {
		entry.Error(param.Body)
	}
}
