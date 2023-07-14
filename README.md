# TeeGo

teeGo是类似gin的一个极简框架，性能是gin的3倍，是iris的1.07倍


teeGo支持路径参数


teeGo性能测试：
![a7da04c8ce648f4301077c6bf92b339](https://github.com/teemoyangyyq/TeeGo/assets/33918440/ec019825-2efa-4fb7-a704-3269cfaa957a)



iris性能测试：
![539c8dc6e4f84ae91b4d883ecdd132d](https://github.com/teemoyangyyq/TeeGo/assets/33918440/09eebac4-8933-45a5-94ae-585265eb3f26)



gin性能测试：
![036bea6e7ae7ea0ee792dc59569fd50](https://github.com/teemoyangyyq/TeeGo/assets/33918440/2ad6c913-c16c-4f39-bb67-9d8f7de15371)


测试代码：

``` go
func BenchmarkHi(b *testing.B) {
	var validTests = []struct {
		data string
		ok   bool
	}{
		{`http://127.0.0.1:8082/tee/api/api/1/qq/2`, false},
		{`http://127.0.0.1:8082/tee/api/api/qq`, false},
		{`http://127.0.0.1:8082/tee/api/api/1/qqq/2`, false},
		{`http://127.0.0.1:8082/tee/api/api/qqq`, false},
		{`http://127.0.0.1:8082/tee/api/api/1/qqqq/2`, false},
		{`http://127.0.0.1:8082/tee/api/api/qqqq`, false},

		{`http://127.0.0.1:8082/tee/api/hh/api/qq`, false},
		{`http://127.0.0.1:8082/tee/api/hh/api/2/qq`, false},
		{`http://127.0.0.1:8082/tee/api/hh/api/qqq`, false},
		{`http://127.0.0.1:8082/tee/api/hh/api/2/qqq`, false},
		{`http://127.0.0.1:8082/tee/api/hh/api/qqqq`, false},
		{`http://127.0.0.1:8082/tee/api/hh/api/2/qqqq`, false},

		{`http://127.0.0.1:8082/tee/service/api/qq`, false},
		{`http://127.0.0.1:8082/tee/service/api/1/qq`, true},
		{`http://127.0.0.1:8082/tee/service/api/qqq`, false},
		{`http://127.0.0.1:8082/tee/service/api/1/qqq`, true},
		{`http://127.0.0.1:8082/tee/service/api/qqqq`, false},
		{`http://127.0.0.1:8082/tee/service/api/1/qqqq`, true},

		{`http://127.0.0.1:8082/yq/yy1`, true},
		{`http://127.0.0.1:8082/yq/yy2`, true},
		{`http://127.0.0.1:8082/yq/yy3`, true},
		{`http://127.0.0.1:8082/yq/yy4`, true},
		{`http://127.0.0.1:8082/yq/yy5`, true},
		{`http://127.0.0.1:8082/yq/yy6`, true},
		{`http://127.0.0.1:8082/yyq/yy7`, true},

		{`http://127.0.0.1:8082/yq/yy1/1`, true},
		{`http://127.0.0.1:8082/yq/yy2/2`, true},
		{`http://127.0.0.1:8082/yq/yy3/3`, true},
		{`http://127.0.0.1:8082/yq/yy4/4`, true},
		{`http://127.0.0.1:8082/yq/yy5/5`, true},
		{`http://127.0.0.1:8082/yq/yy6/6`, true},
		{`http://127.0.0.1:8082/yq/yy7/7`, true},

		{`http://127.0.0.1:8082/tee/yyq/yy1`, true},
		{`http://127.0.0.1:8082/tee/yyq/yy2`, true},
		{`http://127.0.0.1:8082/tee/yyq/yy3`, true},
		{`http://127.0.0.1:8082/tee/yyq/yy4`, true},
		{`http://127.0.0.1:8082/tee/yyq/yy5`, true},
		{`http://127.0.0.1:8082/tee/yyq/yy6`, true},
		{`http://127.0.0.1:8082/tee/yyq/yy7`, true},

		{`http://127.0.0.1:8082/yyq/yy3`, true},
		{`http://127.0.0.1:8082/yyq/yy4`, true},
		{`http://127.0.0.1:8082/yyq/yy5`, true},
		{`http://127.0.0.1:8082/yyq/yy6`, true},
		{`http://127.0.0.1:8082/yyq/yy7`, true},
	}
	
		b.Run("", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					for _, v := range validTests {
						http.Get(v.data)
				    }
				}
		})
	
}
```


路由用法：
``` go
package main

import (
	tee "teego/TeeGo"
)

//  测试代码

func main() {
	//  获取引擎
	e := tee.NewEngine()
	e.Use(HiMiddle)
	v0 := e.Group("/tee")
	{
		v0.Use(HelloMiddle)
		v1 := v0.Group("/api").Use(HiMiddle) //使用中间件
		{
			v1.Use(HelloMiddle)
			
			v1.POST("/api/:type/qq/:id", UserController)
			v1.GET("/api/:id/qq", UserMiddle, UserController)
			v5 := v1.Group("/hh")
			{
				v5.Use(HelloMiddle)
				v5.GET("/api/qq", UserController)
				
			}
		}
		v2 := v0.Group("/service")
		{
			v2.GET("/api/:type/qqqqq", UserController)
		}

	}
	v3 := e.Group("/yq")
	{
		
		v3.GET("/yy7/:id", UserController)
	}
	e.GET("/tee/yyq/yy1", UserController)

	tee.Start("127.0.0.1:8083")
}

func HelloMiddle(c *tee.Context) {
	//fmt.Println("before hello")
	c.Next()
	//fmt.Println("after hello")
}

func HiMiddle(c *tee.Context) {
	//fmt.Println("before hi")
	c.Next()
	//fmt.Println("after hi")
}

func UserMiddle(c *tee.Context) {
	//fmt.Println("before UserMiddle")
	c.Next()
	//fmt.Println("after UserMiddle")

}

```

获取参数：
``` go
// 浏览器输入  http://127.0.0.1:8083/tee/api/api/1/qq?name=yyq

type ResData struct {
	Name string
	Id   int
}


func UserController(c *tee.Context) {
        // 获取路径参数
	idint := c.PathParamInt("id")
        idstring := c.PathParamString("id")
        idint64 :=  c.PathParamString("id")
        // 获取请求参数
        name ：= c.Query("name")
        // 获取文件
        file,_ ：= c.FormFile("fileName")
        // 保存文件
        c.SaveUploadedFile(file, "./tmp/"+file.Filename)
	c.JSON( 200, map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": &ResData{
			Id:   1,
			Name: "杨云强",
		},
	})
}
```










