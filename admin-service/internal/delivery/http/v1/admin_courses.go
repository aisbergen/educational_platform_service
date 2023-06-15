package v1

import (
	"admin/internal/domain"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	admin "admin/pkg/admin/api/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

var (
	apiCourses = "http://localhost:8000/api/v1/courses/"
)

// @Summary		Admin Get Course By ID
// @Security AdminAuth
// @Tags			Courses
// @Description	Admin Create New Courses
// @ModuleID adminGetCourseByID
// @Accept			json
// @Produce		json
// @Param			id path string	true	"course id"
// @Success		201		{object}	domain.ResponseCourse
// @Failure		400		{object}	Resposne
// @Failure		500		{object}	Resposne
// @Router			/admin/courses/{id} [get]
func (h *Handler) adminGetCourseByID(ctx *gin.Context) {
	param := ctx.Param("id")

	url := apiCourses + param

	resp, err := http.Get(url)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Failed to make the request")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Failed to read the response body")
		return
	}

	if resp.StatusCode != http.StatusOK {
		newResponse(ctx, resp.StatusCode, string(body))
		return
	}

	fmt.Println(string(body))

	var course domain.Courses

	if err = json.Unmarshal(body, &course); err != nil {
		newResponse(ctx, resp.StatusCode, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, course)
}

// @Summary		Admin Create New Courses
// @Security AdminAuth
// @Tags			Courses
// @Description	Admin Create New Courses
// @Accept			json
// @Produce		json
// @Param			account	body		createCourses	true	"Courses"
// @Success		201		{object}	Resposne
// @Failure		400		{object}	Resposne
// @Failure		500		{object}	Resposne
// @Router			/admin/courses/create [post]
func (h *Handler) adminCreateCourses(ctx *gin.Context) {
	var inp admin.Courses

	if err := ctx.BindJSON(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect input data format")
		return
	}

	url := apiCourses + "create"

	serializedData, err := proto.Marshal(&inp)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Failed to marshal input data")
		return
	}

	resp, err := http.Post(url, "application/x-protobuf", bytes.NewBuffer(serializedData))
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Failed to make the request")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		newResponse(ctx, resp.StatusCode, "Failed to create courses")
		return
	}

	ctx.JSON(http.StatusOK, Resposne{"Courses created successfully"})
}

// @Summary		Admin Update Course
// @Security AdminAuth
// @Tags			Courses
// @Description	Admin Update Course
// @Accept			json
// @Produce		json
// @Param			account	body		inputCourse	true	"course update info"
// @Param			id path string		true	"course id"
// @Success		200		{object}	Resposne
// @Failure		400		{object}	Resposne
// @Failure		500		{object}	Resposne
// @Router			/admin/courses/{id}/update [put]
func (h *Handler) adminUpdateCourses(ctx *gin.Context) {
	var inp admin.Courses
	param := ctx.Param("id")

	if err := ctx.BindJSON(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect input data format")
		return
	}

	url := apiCourses + param + "/" + "update"

	serializedData, err := proto.Marshal(&inp)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Failed to marshal input data")
		return
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(serializedData))
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Failed to create the request")
		return
	}
	req.Header.Set("Content-Type", "application/x-protobuf")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Failed to make the request")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		newResponse(ctx, resp.StatusCode, "Failed to update the course")
		return
	}

	ctx.JSON(http.StatusOK, Resposne{"Course updated successfully"})
}

// @Summary		Admin Delete Course
// @Security AdminAuth
// @Tags			Courses
// @Description	Admin Delete Course
// @Accept			json
// @Produce		json
// @Param			id path string		true	"course id"
// @Success		200		{object}	Resposne
// @Failure		400		{object}	Resposne
// @Failure		500		{object}	Resposne
// @Router			/admin/courses/{id}/delete [delete]
func (h *Handler) adminDeleteCourse(ctx *gin.Context) {
	param := ctx.Param("id")

	url := apiCourses + param + "/" + "delete"

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Failed to create the request")
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Failed to make the request")
		return
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode != http.StatusOK {
		newResponse(ctx, resp.StatusCode, "Failed to delete the course")
		return
	}

	ctx.JSON(http.StatusOK, Resposne{"Course deleted successfully"})
}

// @Summary		Admin Get Students By CoursId
// @Security AdminAuth
// @Tags			Courses
// @Description	Admin Get Students By CoursId
// @Accept			json
// @Produce		json
// @Param			id path string		true	"course id"
// @Success		200		{object}	Resposne
// @Failure		500		{object}	Resposne
// @Router			/admin/courses/{id}/students [get]
func (h *Handler) adminGetStudentsByCoursId(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.services.Kafka.SendMessages("students-request", id)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Failed to get information about courses")
		return
	}
	responseData := <-h.responseCh

	ctx.Data(http.StatusOK, "application/json", responseData)
}

/*
var (
	api = "http://localhost:8000/api/v1/students/"
)

func (h *Handler) getStudentsByCoursId(ctx *gin.Context) {
	param := ctx.Param("id")

	url := api + param + "/students"
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get information about students" + err.Error(),
		})
		return
	}

	if resp.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get information about students",
		})
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read the body of the answer",
		})
		return
	}

	var students domain.Response
	if err := json.Unmarshal(body, &students); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to decode response body",
		})
		return
	}
	if len(students.Students) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Students not found.",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"students": students,
	})
}
*/
