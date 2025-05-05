package test

import (
	"forum/pkg/globals"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 目标处理函数
func BenchmarkGinHandler(b *testing.B) {
	r := globals.Router
	/*r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})*/

	// 创建测试请求和响应记录器
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/xxxxxx", nil)

	b.ResetTimer() // 计时开始
	for i := 0; i < b.N; i++ {
		r.ServeHTTP(w, req)
	}
	b.StopTimer() // 计时结束
}
