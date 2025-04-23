package main

import (
	"net/http"
	tee "github.com/teemoyangyyq/TeeGo"
)

//  测试代码

func Testtee() {
	//  获取引擎
	e := tee.NewEngine()
	e.Use(HiMiddle)
	v0 := e.Group("/tee")
	{
		v0.Use(HelloMiddle)
		v1 := v0.Group("/api") //使用中间件
		{   
			v1.Use(HelloMiddle)
			v1.GET("/api/qq", UserMiddle, UserController)
			v1.GET("/api/:type/qq", UserController)
			v1.GET("/api/qqq", UserMiddle, UserController)
			v1.GET("/api/:type/qqq/:id", UserController)
			v1.GET("/api/qqqq", UserMiddle, UserController)
			v1.GET("/api/:type/qqqq/:id", UserController)
			v5 := v1.Group("/hh")
			{
				v5.Use(HelloMiddle)
				v5.GET("/api/qq", UserController)
				v5.GET("/api/:type/qq", UserController)
				v5.GET("/api/qqq", UserController)
				v5.GET("/api/:type/qqq", UserController)
				v5.GET("/api/qqqq", UserController)
			}
		}
		v2 := v0.Group("/service") 
		{
			v2.GET("/api/qq", UserController)
		    v2.GET("/api/:type/qq", UserController)
			v2.GET("/api/qqq", UserController)
			v2.GET("/api/:type/qqq", UserController)
			v2.GET("/api/qqqq", UserController)
			v2.GET("/api/:type/qqqq", UserController)
			v2.GET("/api/qqqqq", UserController)
			v2.GET("/api/:type/qqqqq", UserController)
		}
		
	}
	v3 := e.Group("/yq")
	{
		v3.GET("/yy1", UserController)
		v3.GET("/yy2", UserController)
		v3.GET("/yy3", UserController)
		v3.GET("/yy4", UserController)
		v3.GET("/yy5", UserController)
		v3.GET("/yy6", UserController)
		v3.GET("/yy7", UserController)
		v3.GET("/yy1/:id", UserController)
		v3.GET("/yy2/:id", UserController)
		v3.GET("/yy3/:id", UserController)
		v3.GET("/yy4/:id", UserController)
		v3.GET("/yy5/:id", UserController)
		v3.GET("/yy6/:id", UserController)
		v3.GET("/yy7/:id", UserController)
	}
	e.GET("/tee/yyq/yy1", UserController)
	e.GET("/tee/yyq/yy2", UserController)
	e.GET("/tee/yyq/yy3", UserController)
	e.GET("/tee/yyq/yy4", UserController)
	e.GET("/tee/yyq/yy5", UserController)
	e.GET("/tee/yyq/yy6", UserController)
	e.GET("/tee/yyq/yy7", UserController)
	e.GET("/yyq/yy1", UserController)
	e.GET("/yyq/yy2", UserController)
	e.GET("/yyq/yy3", UserController)
	e.GET("/yyq/yy4", UserController)
	e.GET("/yyq/yy5", UserController)
	e.GET("/yyq/yy6", UserController)
	e.GET("/yyq/yy7", UserController)

	tee.Start("127.0.0.1:8082")
}

func HelloMiddle(c *tee.Context) {
	c.Next()
}

func HiMiddle(c *tee.Context) {
	c.Next()
}

func UserMiddle(c *tee.Context) {
	c.Next()
	
}
func UserController(c *tee.Context) {	
	c.JSON(http.StatusOK,"杨云强")
}

 func main(){
 	Testtee()
 }
