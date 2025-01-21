package server

// // HandlePublicMW 处理公共中间件
// func HandlePublicMW(e *gin.Engine) {
// 	// 跨域
// 	e.Use(corsMW.CorsMiddleware())
//
// 	// e.Use(cors.New(cors.Config{
// 	// 	AllowOrigins:     []string{"http://example.com", "http://localhost:3000"}, // 允许的前端域
// 	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},     // 允许的方法
// 	// 	AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},     // 允许的请求头
// 	// 	ExposeHeaders:    []string{"Content-Length", "Authorization"},             // 允许前端访问的响应头
// 	// 	AllowCredentials: true,                                                    // 允许 Cookie
// 	// 	MaxAge:           12 * 60 * 60,                                            // 预检请求缓存时间
// 	// }))
//
// }
