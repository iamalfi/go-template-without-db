package routes

import (
	"adhar-verification/controller/aadhaar"
	"adhar-verification/middleware"

	"github.com/gin-gonic/gin"
)

func AdharRoutes(r *gin.RouterGroup) {
	adharRoutes := r.Group("/aadhaar")
	adharRoutes.Use(middleware.Authenticate())

	adharRoutes.POST("/otp", aadhaar.GenerateOtp)
	adharRoutes.POST("/otp/verify/:aadhaar_no", aadhaar.Verify)

}
