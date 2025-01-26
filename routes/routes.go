package routes

import "github.com/gin-gonic/gin"

func Routes(r *gin.RouterGroup) {
	v1 := r.Group("/v1")

	UserRoutes(v1)

}
