package middlewares

import (
	"fmt"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// UserRankMW 用户热度排行中间件
func UserRankMW(c *gin.Context) {
	// 绑定数据
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		globals.Log.Error(response.ErrBindDataIsWrong)
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrBindDataIsWrong), nil))
		c.Abort()
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		globals.Log.Error(response.ErrBindDataIsWrong)
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrBindDataIsWrong), nil))
		c.Abort()
		return
	}

	// 检验数据
	if page <= 0 {
		globals.Log.Error(response.ErrPageRangeIsWrong)
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrPageRangeIsWrong), nil))
		c.Abort()
		return
	}
	if limit <= 0 {
		globals.Log.Error(response.ErrLimitRangeIsWrong)
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrLimitRangeIsWrong), nil))
		c.Abort()
		return
	}

	c.Set("page", page)
	c.Set("limit", limit)

	c.Next()
}

// AttentionMW 搜索用户关注的人中间件
func AttentionMW(c *gin.Context) {
	// 绑定数据
	userId, err := strconv.Atoi(c.Query("userId"))
	if err != nil {
		globals.Log.Error(response.ErrDataTypeIsWrong)
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrDataTypeIsWrong), nil))
		c.Abort()
		return
	}

	keyword := c.Query("keyword")
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		globals.Log.Error(response.ErrDataTypeIsWrong)
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrDataTypeIsWrong), nil))
		c.Abort()
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		globals.Log.Error(response.ErrDataTypeIsWrong)
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrDataTypeIsWrong), nil))
		c.Abort()
		return
	}

	// 检验数据
	if page <= 0 {
		globals.Log.Error(response.ErrPageRangeIsWrong)
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrPageRangeIsWrong), nil))
		c.Abort()
		return
	}
	if limit <= 0 {
		globals.Log.Error(response.ErrLimitRangeIsWrong)
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrLimitRangeIsWrong), nil))
		c.Abort()
		return
	}

	attentionReq := requests.AttentionReq{
		UserId:  uint(userId),
		Keyword: keyword,
		Page:    page,
		Limit:   limit,
	}

	c.Set("req", attentionReq)

	c.Next()
}
