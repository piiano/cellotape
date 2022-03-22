package restcontroller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (fn ControllerFn[B, P, Q, R]) GinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params Params[B, P, Q]
		if err := c.Bind(&params.Body); err != nil {
			c.JSON(400, err.Error())
			return
		}
		if err := c.BindUri(&params.Path); err != nil {
			c.JSON(400, err.Error())
			return
		}
		if err := c.BindQuery(&params.Query); err != nil {
			c.JSON(400, err.Error())
			return
		}
		params.Headers = c.Request.Header
		resp, err := fn(params)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
