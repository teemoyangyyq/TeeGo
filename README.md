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
