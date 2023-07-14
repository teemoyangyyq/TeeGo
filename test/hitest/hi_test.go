package test

import (
	"net/http"
	"testing"
)

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
