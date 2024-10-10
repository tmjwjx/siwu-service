package routes

/*import (
	"forum/casbin_r_m_a"
	"forum/internal/menu/controllers"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
)

// 菜单管理 casbin鉴权

func Menu(e *gin.Engine) {
	casbinService, err := casbin_r_m_a.NewCasbinService(globals.DB)
	if err != nil {
		globals.Log.Errorf("casbin启动错误")
	}
	r := e.Group("/menu")

	r.Use(token.AuthMiddleware(), casbin_r_m_a.NewCasbinAuth(casbinService))
	{
		//菜单初始化
		r.GET("/init", controllers.MenuInitCtrl)
	}
}*/
