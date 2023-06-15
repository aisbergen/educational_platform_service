package v1

import (
	"admin/internal/domain"
	"context"
	"errors"
	"net/http"

	admin "admin/pkg/admin/api/proto"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initAdminRoutes(api *gin.RouterGroup) {
	admin := api.Group("/admin")
	{
		admin.POST("/sign-up", h.adminSignUp)
		admin.POST("/sign-in", h.adminSignIn)
		admin.POST("/auth/refresh", h.adminRefreshToken)
		authentocated := admin.Group("/", h.adminIdentity)
		{
			students := authentocated.Group("/students")
			{
				students.POST("/create", h.adminCreatestudent)
				students.GET("/:id", h.adminGetStudentByID)
				students.PUT("/update/:id", h.adminUpdateStudent)
				students.DELETE("/delete/:id", h.adminDeleteStudent)
				students.GET("/:id/courses", h.adminGetCoursesStudents)
			}

			courses := authentocated.Group("/courses")

			{
				courses.GET("/:id", h.adminGetCourseByID)
				courses.POST("/create", h.adminCreateCourses)
				courses.PUT("/:id/update", h.adminUpdateCourses)
				courses.DELETE("/:id/delete", h.adminDeleteCourse)
				courses.GET("/:id/students", h.adminGetStudentsByCoursId)
			}
		}
	}
}

// @Summary		Create a new admin
// @Tags			admin
// @Description	Create a new admin with the input payload
// @Accept			json
// @Produce		json
// @Param			account	body		inputAdmin	true	"Admin"
// @Success		201		{object}	Resposne
// @Failure		400		{object}	Resposne
// @Failure		500		{object}	Resposne
// @Router			/admin/sign-up [post]
func (h *Handler) adminSignUp(ctx *gin.Context) {
	var inp admin.Admin

	if err := ctx.Bind(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect input data format")
		return
	}

	if err := h.services.Admin.SignUp(ctx, domain.Admin{
		Email:    inp.Email,
		Name:     inp.Name,
		Password: inp.Password,
	}); err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Error when registering as an administrator")
		return
	}
	ctx.JSON(http.StatusCreated, Resposne{"Administrator successfully registered"})
}

// type signInInput struct {
// 	Email    string `json:"email" binding:"required,email,max=64"`
// 	Password string `json:"password" binding:"required,min=8,max=64"`
// }

// @Summary		Sign-in
// @Tags			admin
// @Description	Sign-in
// @Accept			json
// @Produce		json
// @Param			account	body		signInInput	true	"Admin"
// @Success		200		{object}	Resposne
// @Failure		400		{object}	Resposne
// @Failure		500		{object}	Resposne
// @Router			/admin/sign-in [post]
func (h *Handler) adminSignIn(ctx *gin.Context) {
	var inp admin.Admin

	if err := ctx.Bind(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect input data format")
		return
	}

	token, err := h.services.Admin.SignIn(context.Background(), inp.Email, inp.Password)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			newResponse(ctx, http.StatusUnauthorized, "Incorrect email address or password")
			return
		}

		newResponse(ctx, http.StatusInternalServerError, "A login error occurred")
		return
	}

	ctx.JSON(http.StatusOK, token)
}

// @Summary		Refresh Token
// @Tags			admin
// @Description	Refresh Token
// @Accept			json
// @Produce		json
// @Param			account	body		domain.Session	true	"Admin"
// @Success		200		{object}	domain.Token
// @Failure		400		{object}	Resposne
// @Failure		500		{object}	Resposne
// @Router			/admin/auth/refresh [post]
func (h *Handler) adminRefreshToken(ctx *gin.Context) {
	var inp admin.Session

	if err := ctx.Bind(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect input data format")

		return
	}

	token, err := h.services.Admin.GetByRefreshToken(context.Background(), inp.RefreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			newResponse(ctx, http.StatusBadRequest, "Refresh token not found")
			return
		}
		newResponse(ctx, http.StatusInternalServerError, "Error retrieving token")
		return
	}

	ctx.JSON(http.StatusOK, token)
}
