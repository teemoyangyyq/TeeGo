package main

import (
	

	"github.com/kataras/iris/v12"
)

func testIris() {
	e := iris.New()


	e.Use(HiMiddleIris)
	v0 := e.Party("/tee")
	{
		v0.Use(HelloMiddleIris)
		v1 := v0.Party("/api") //使用中间件
		{
			v1.Use(HelloMiddleIris)
			v1.Get("/api/qq", UserMiddleIris, UserControllerIris)
			v1.Get("/api/:type/qq", UserControllerIris)
			v1.Get("/api/qqq", UserMiddleIris, UserControllerIris)
			v1.Get("/api/:type/qqq/:id", UserControllerIris)
			v1.Get("/api/qqqq", UserMiddleIris, UserControllerIris)
			v1.Get("/api/:type/qqqq/:id", UserControllerIris)
			v5 := v1.Party("/hh")
			{
				v5.Use(HelloMiddleIris)
				v5.Get("/api/qq", UserControllerIris)
				v5.Get("/api/:type/qq", UserControllerIris)
				v5.Get("/api/qqq", UserControllerIris)
				v5.Get("/api/:type/qqq", UserControllerIris)
				v5.Get("/api/qqqq", UserControllerIris)
			}
		}
		v2 := v0.Party("/service") 
		{
			v2.Get("/api/qq", UserControllerIris)
		    v2.Get("/api/:type/qq", UserControllerIris)
			v2.Get("/api/qqq", UserControllerIris)
			v2.Get("/api/:type/qqq", UserControllerIris)
			v2.Get("/api/qqqq", UserControllerIris)
			v2.Get("/api/:type/qqqq", UserControllerIris)
			v2.Get("/api/qqqqq", UserControllerIris)
			v2.Get("/api/:type/qqqqq", UserControllerIris)
		}
		
	}
	v3 := e.Party("/yq")
	{
		v3.Get("/yy1", UserControllerIris)
		v3.Get("/yy2", UserControllerIris)
		v3.Get("/yy3", UserControllerIris)
		v3.Get("/yy4", UserControllerIris)
		v3.Get("/yy5", UserControllerIris)
		v3.Get("/yy6", UserControllerIris)
		v3.Get("/yy7", UserControllerIris)
		v3.Get("/yy1/:id", UserControllerIris)
		v3.Get("/yy2/:id", UserControllerIris)
		v3.Get("/yy3/:id", UserControllerIris)
		v3.Get("/yy4/:id", UserControllerIris)
		v3.Get("/yy5/:id", UserControllerIris)
		v3.Get("/yy6/:id", UserControllerIris)
		v3.Get("/yy7/:id", UserControllerIris)
	}
	e.Get("/tee/yyq/yy1", UserControllerIris)
	e.Get("/tee/yyq/yy2", UserControllerIris)
	e.Get("/tee/yyq/yy3", UserControllerIris)
	e.Get("/tee/yyq/yy4", UserControllerIris)
	e.Get("/tee/yyq/yy5", UserControllerIris)
	e.Get("/tee/yyq/yy6", UserControllerIris)
	e.Get("/tee/yyq/yy7", UserControllerIris)
	e.Get("/yyq/yy1", UserControllerIris)
	e.Get("/yyq/yy2", UserControllerIris)
	e.Get("/yyq/yy3", UserControllerIris)
	e.Get("/yyq/yy4", UserControllerIris)
	e.Get("/yyq/yy5", UserControllerIris)
	e.Get("/yyq/yy6", UserControllerIris)
	e.Get("/yyq/yy7", UserControllerIris)
	e.Run(iris.Addr("localhost:8082"))
}
	
	
func HelloMiddleIris(c iris.Context) {
	c.Next()
}

func HiMiddleIris(c iris.Context) {
	c.Next()
}

func UserMiddleIris(c iris.Context) {
	c.Next()
}
func UserControllerIris(c iris.Context) {
	
	
	c.Write([]byte("hh"))
}

// func main() {
// 	testIris()
// }