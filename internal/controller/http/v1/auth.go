package v1

import (
	"net/http"

	"github.com/arunima10a/task-manager/internal/usecase"
	"github.com/gin-gonic/gin"
)

type authRoutes struct {
	a usecase.AuthUseCase
}

type signUpRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func newAuthRoutes(handler *gin.RouterGroup, a usecase.AuthUseCase) {
	r := &authRoutes{a}
	h := handler.Group("/auth")
	{
		h.POST("/sign-up", r.signUp)
		h.POST("/login", r.login)
	}
}

func (r *authRoutes) signUp(c *gin.Context) {
	var req signUpRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	_ = r.a.SignUp(c.Request.Context(), req.Email, req.Password)
	c.JSON(201, gin.H{"message": "user created"})
}

func (r *authRoutes) login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	token, err := r.a.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}
