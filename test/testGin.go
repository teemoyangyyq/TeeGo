package main

import (
	

	"github.com/gin-gonic/gin"
)

//  测试代码

func Testgin() {
	//  获取引擎
	e := gin.Default()
	e.Use(HiMiddlegin)
	v0 := e.Group("/tee")
	{

		v0.Use(HiMiddlegin)
		v1 := v0.Group("/api") //使用中间件
		{
			v1.Use(HelloMiddlegin)
			v1.GET("/api/qq", UserMiddlegin, UserControllergin)
			v1.GET("/api/:type/qq/:id", UserControllergin)
			v1.GET("/api/qqq", UserMiddlegin, UserControllergin)
			v1.GET("/api/:type/qqq/:id", UserControllergin)
			v1.GET("/api/qqqq", UserMiddlegin, UserControllergin)
			v1.GET("/api/:type/qqqq/:id", UserControllergin)
			v5 := v1.Group("/hh")
			{
				v5.Use(HelloMiddlegin)
				v5.GET("/api/qq", UserControllergin)
				v5.GET("/api/:type/qq", UserControllergin)
				v5.GET("/api/qqq", UserControllergin)
				v5.GET("/api/:type/qqq", UserControllergin)
				v5.GET("/api/qqqq", UserControllergin)
				v5.GET("/api/:type/qqqq", UserControllergin)
			}
			
		}
		v2 := v0.Group("/service")
		 {
		     v2.GET("/api/qq", UserControllergin)
		     v2.GET("/api/:type/qq", UserControllergin)
			 v2.GET("/api/qqq", UserControllergin)
		     v2.GET("/api/:type/qqq", UserControllergin)
			 v2.GET("/api/qqqq", UserControllergin)
		     v2.GET("/api/:type/qqqq", UserControllergin)
			 v2.GET("/api/qqqqq", UserControllergin)
		     v2.GET("/api/:type/qqqqq", UserControllergin)
		 }
	}
	v3 := e.Group("/yq")
	{
		v3.GET("/yy1", UserControllergin)
		v3.GET("/yy2", UserControllergin)
		v3.GET("/yy3", UserControllergin)
		v3.GET("/yy4", UserControllergin)
		v3.GET("/yy5", UserControllergin)
		v3.GET("/yy6", UserControllergin)
		v3.GET("/yy7", UserControllergin)
		
		v3.GET("/yy1/:id", UserControllergin)
		v3.GET("/yy2/:id", UserControllergin)
		v3.GET("/yy3/:id", UserControllergin)
		v3.GET("/yy4/:id", UserControllergin)
		v3.GET("/yy5/:id", UserControllergin)
		v3.GET("/yy6/:id", UserControllergin)
		v3.GET("/yy7/:id", UserControllergin)
		

	}

	e.GET("/yyq/yy1", UserControllergin)
	e.GET("/yyq/yy2", UserControllergin)
	e.GET("/yyq/yy3", UserControllergin)
	e.GET("/yyq/yy4", UserControllergin)
	e.GET("/yyq/yy5", UserControllergin)
	e.GET("/yyq/yy6", UserControllergin)
	e.GET("/yyq/yy7", UserControllergin)
	e.GET("/tee/yyq/yy1", UserControllergin)
	e.GET("/tee/yyq/yy2", UserControllergin)
	e.GET("/tee/yyq/yy3", UserControllergin)
	e.GET("/tee/yyq/yy4", UserControllergin)
	e.GET("/tee/yyq/yy5", UserControllergin)
	e.GET("/tee/yyq/yy6", UserControllergin)
	e.GET("/tee/yyq/yy7", UserControllergin)
	e.Run(":8082")
}

func HelloMiddlegin(c*gin.Context) {
	
	c.Next()
	
	//c.Res.Write([]byte("qqqqqqqqqqqq"))
}

func HiMiddlegin(c*gin.Context) {
	// fmt.Println("before HiMiddle")
	c.Next()
	// fmt.Println("after HiMiddle")

}

func UserMiddlegin(c*gin.Context) {
	
	c.Next()
	

}
func UserControllergin(c*gin.Context) {
	
}

// func main(){
// 	Testgin()
// }

