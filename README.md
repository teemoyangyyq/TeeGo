# TeeGo

teeGo是类似gin的一个极简框架，路由分发性能是gin的3倍，是iris的1.6倍


teeGo支持路径参数


teeGo性能测试：![image](https://github.com/teemoyangyyq/TeeGo/assets/33918440/56692b2a-70ae-4266-99d3-2d724a54a8a3)





iris性能测试：![image](https://github.com/teemoyangyyq/TeeGo/assets/33918440/3a19e17c-0468-47b5-bb40-75c618e32508)




gin性能测试：![image](https://github.com/teemoyangyyq/TeeGo/assets/33918440/0ee57c26-10cb-457c-b0c0-5a2ea3773551)


测试文件在test/目录下，三个一模一样的路由注册，分别是teego, iris,  gin,拥有相同控制器方法，中间件，为了测试性能，这些方法内什么操作都没有，仅仅测试框架路由分发性能

性能测试代码：

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
        name := c.Query("name")
        // 获取文件
        file, _ := c.FormFile("fileName")
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


使用日志：
``` go
package main

import (
	"fmt"
	tee "teego"
)

//  测试代码

func Testtee() {
	//  获取引擎
	e := tee.NewEngine()
	e.Use(tee.TeeLogger())  // 使用日志
	e.GET("/yyq/yy7", UserController)

	tee.Start("127.0.0.1:8082")
}

func HelloMiddle(c *tee.Context) {
	//fmt.Println("before HelloMiddle")
	c.Next()
	//fmt.Println("after HelloMiddle")
}

func HiMiddle(c *tee.Context) {
	//fmt.Println("before HiMiddle")
	c.Next()
//	fmt.Println("after HiMiddle")
}

func UserMiddle(c *tee.Context) {
	
	//fmt.Println("before UserMiddle")
	c.Next()
	//fmt.Println("after UserMiddle")

}
func UserController(c *tee.Context) {
	 fmt.Println("controller UserController")
         tee.LogGet(c).Info("1234567=====")
	c.Res.Write([]byte(c.Req.RequestURI))
	
}

func main() {
	Testtee()
}
```

![image](https://github.com/teemoyangyyq/TeeGo/assets/33918440/daeddffe-e07a-4d0b-92b1-f91aa717c8a9)








