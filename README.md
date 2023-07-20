# TeeGo

teeGo是类似gin的一个极简框架，路由分发性能是gin的3倍，是iris的1.6倍


teeGo支持路径参数
## 背景
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;路由匹配算法一般使用前缀树进行匹配，如何优化匹配算法
## 优化点：
### 第一点优化： 
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;假设有三个路由 /task/:type/service/url/list, /task/:id/service/url/info, /task/:id/service/url/tag,
浏览器输入请求路径为/task/1/service/url/tag，会匹配路由/task/:id/service/url/tag，如果前缀树如下图所示，那么路由查找的时候，在匹配了task之后，:type和：id都会被匹配，之后还会分别匹配后面的service，会分两条路径进行匹配。我们发现这样的匹配会有多余匹配。因为本来我只会匹配/task/:id/service/url/tag，结果是我即匹配/task/:type/service/url/，也匹配/task/:id/service/url，只有匹配到最后的叶子节点才发现不匹配。怎么解决产生的多余匹配问题，teego框架已经给出方案。

​
​
![image](https://github.com/teemoyangyyq/TeeGo/assets/33918440/ee6bee1c-9e6d-4360-ad98-2dddf3f93441)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;那就是把:id,:type合并成/，这样就好了，如下图。但是这样优化好了之后，我们要获取路径参数，得知道路径参数名是id而不是type。如果这样处理，会遇到新的问题，那就是怎么知道路径参数：
比如浏览器输入请求路径为/task/1/service/url/tag， 我们在控制器里获取路径参数id是为1，获取路径参数type就是空，因为路由上对应的是:id，而不是:type。解决思路如下，每个路由的插入的前缀树叶子节点肯定不同，如果叶子节点相同，代表输入url会匹配两个路由，就会有问题。所以我们可以个给每个路由叶子节点一个索引，理论上通过这个索引，我们是可以知道这个完整路由，从而拿到注册路由的路径参数名，实现见后文


![image](https://github.com/teemoyangyyq/TeeGo/assets/33918440/20c633fe-c18e-4356-843b-4608bbbbaf2a)

### 第二点优化：
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;对于/task/service/url/list, /task//service/url/info/:id, /task/service/url/tag,这三个路由，当浏览器输入 请求路径为/task/service/url/tag，那么需要四次匹配，分别匹配task，service，url，tag，这个时候对于没有路径参数的路由其实我们可以存个全局路由map，key为/task/service/url/tag，value为对应执行方法，这样匹配全局路由map可以一次匹配到位。对于有路径参数的，那就只能在前缀树上一一匹配了


​
![image](https://github.com/teemoyangyyq/TeeGo/assets/33918440/0e272d8a-cef7-48c8-92d2-695db7c6530e)



### 第三点优化：
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;其实我们匹配到路由后，会获取上面我们所说的路由索引。在路径匹配前，我们遍历这个前缀树，初始化一下一个新的全局map，这个map以路由索引为key，value直接存储该路由的所有执行方法和路径参数。这样匹配后，我们拿到路由索引，直接在一个全局map中获取执行方法和路径参数返回和执行

### 路由匹配算法：
``` go
         e := tee.NewEngine()
	
	v1 := e.Group("/task/:type")
	{
		v1.AddRoute("/service/url/list", UserController)
	}
        v2 := ve.Group("/task/:id")
	{       v2.AddRoute("/service/url/tag", UserController)
		v2.AddRoute("/service/url/info/:id", UserController)
        }
```
### 思路图解：
    1.矩形代表Engine结构体节点，每一个Group和AddRoute都会创建一个Engine，Engine的当前节点（取名CurNode）指向前缀树TreeNode节点；
    2.椭圆代表TreeNode结构体节点（TreeNode是前缀树节点），当TreeNode为叶子节点会记录它对应的addRoute的Engine节点（取名PreEngine）
![切片 1](https://github.com/teemoyangyyq/TeeGo/assets/33918440/23df5862-0d88-4acf-8bfe-67d1f496f25d)

    3.每一个Engine会存储路径参数名称和父亲Engine（取名PreEngine）

![切片 4 (1)](https://github.com/teemoyangyyq/TeeGo/assets/33918440/c532d61c-9baf-4c75-86fc-f202728f6eaa)


   4.在路由前缀树建立成功后，遍历前缀树，做两件事：<br><br>
    &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  4.1  把不带路径参数的路由放进一个全局map，key为路由index,value为路由索引<br><br>
    &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  4.2 存两个路由索引map，key为路由索引。第一个value存路由对应所有方法，包括group分组方法，这些分组方法可以通过addroute叶子节点反查引擎;第二个value存路由对应路由路径参数数组，按先后顺序存放<br><br>
   5.浏览器输入路径，匹配时，先匹配不带参数的全局map路由索引，如果不存在再查询前缀树，获取路由索引，然后一一匹配路径方法，匹配到路径参数过程中记录路径值，可以用递归回溯解决


## 性能测试
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;测试文件在test/目录下，三个一模一样的路由注册，分别是teego, iris,  gin,拥有相同控制器方法，中间件，为了测试性能，这些方法内什么操作都没有，仅仅测试框架路由分发性能

### 操作-teeGo性能测试
 &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;1.在/test/testTee.go中，取消main函数注释，运行代码监听8082端口
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 2.打开新窗口，切换到/test/hi_test目录下，执行 go test -bench Hi -benchmem  命令

teeGo性能测试：![image](https://github.com/teemoyangyyq/TeeGo/assets/33918440/56692b2a-70ae-4266-99d3-2d724a54a8a3)


### 操作-iris性能测试
 &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 1.在/test/testIris.go中，取消main函数注释，运行代码监听8082端口
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;2.打开新窗口，切换到/test/hi_test目录下，执行 go test -bench Hi -benchmem  命令

iris性能测试：![image](https://github.com/teemoyangyyq/TeeGo/assets/33918440/3a19e17c-0468-47b5-bb40-75c618e32508)

### 操作-gin性能测试
 &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 1.在/test/testGin.go中，取消main函数注释，运行代码监听8082端口
 &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 2.打开新窗口，切换到/test/hi_test目录下，执行 go test -bench Hi -benchmem  命令

gin性能测试：![image](https://github.com/teemoyangyyq/TeeGo/assets/33918440/0ee57c26-10cb-457c-b0c0-5a2ea3773551)




### 性能单元测试代码：

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


## 路由用法：
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

## 获取参数：
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


## 使用日志：
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








