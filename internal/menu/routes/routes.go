package routes

/*import (
	"forum/casbin"
	"forum/internal/menu/controllers"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
)

// 菜单管理 casbin鉴权

func Menu(e *gin.Engine) {
	casbinService, err := casbin.NewCasbinService(globals.DB)
	if err != nil {
		globals.Log.Errorf("casbin启动错误")
	}
	r := e.Group("/menu")

	r.Use(token.AuthMiddleware(), casbin.NewCasbinAuth(casbinService))
	{
		//菜单初始化
		r.GET("/init", controllers.MenuInitCtrl)
	}
}*/
