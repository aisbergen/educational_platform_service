package v1

import (
	"admin/internal/domain"
	"bytes"
	"encoding/json"
	"io"

	admin "admin/pkg/admin/api/proto"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

var (
	api = "http://localhost:8001/api/v1/admin/students/"
)

// @Summary		Admin Create Student
// @Security AdminAuth
// @Tags			Students
// @Description	Admin Create Student
// @Accept			json
// @Produce		json
// @Param			account	body		inputStudent	true	"Admin"
// @Success		200		{object}	Resposne
// @Failure		400		{object}	Resposne
// @Failure		500		{object}	Resposne
// @Router			/admin/students/create [post]
func (h *Handler) adminCreatestudent(ctx *gin.Context) {
	var inp admin.Student

	if err := ctx.BindJSON(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect input data format")
		return
	}

	url := api + "create"

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Failed to read the response body")
		return
	}

	if resp.StatusCode != http.StatusCreated {
		newResponse(ctx, resp.StatusCode, string(body))
		return
	}

	newResponse(ctx, http.StatusOK, "Student created successfully")
}

// @Summary		Admin Get Student By ID
// @Security AdminAuth
// @Tags			Students
// @Description	Admin Get Student By ID
// @Accept			json
// @Produce		json
// @Param			id path string	true	"student id"
// @Success		200		{object}	domain.Student
// @Failure		500		{object}	Resposne
// @Router			/admin/students/{id} [get]
func (h *Handler) adminGetStudentByID(ctx *gin.Context) {
	param := ctx.Param("id")

	url := api + param

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

	var student domain.Student

	if err := json.Unmarshal(body, &student); err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Failed to Unmarshal")
		return
	}

	ctx.JSON(http.StatusOK, student)
}

// @Summary		Admin Update Student
// @Security AdminAuth
// @Tags			Students
// @Description	Admin Update Student
// @Accept			json
// @Produce		json
// @Param			account	body		domain.UpdateStudentInput	true	"student update info"
// @Param			id path string		true	"student id"
// @Success		200		{object}	Resposne
// @Failure		400		{object}	Resposne
// @Failure		500		{object}	Resposne
// @Router			/admin/students/update/{id} [put]
func (h *Handler) adminUpdateStudent(ctx *gin.Context) {
	var inp admin.Admin
	param := ctx.Param("id")

	if err := ctx.BindJSON(&inp); err != nil {
		newResponse(ctx, http.StatusBadRequest, "Incorrect input data format")
		return
	}

	url := api + "update" + "/" + param
	log.Println(url)

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
		newResponse(ctx, resp.StatusCode, "Failed to update student")
		return
	}

	newResponse(ctx, http.StatusOK, "Student updated successfully")
}

// @Summary		Admin Delete Student
// @Security AdminAuth
// @Tags			Students
// @Description	Admin Delete Student
// @Accept			json
// @Produce		json
// @Param			id path string		true	"student id"
// @Success		200		{object}	Resposne
// @Failure		500		{object}	Resposne
// @Router			/admin/courses/{id}/delete [delete]
func (h *Handler) adminDeleteStudent(ctx *gin.Context) {
	param := ctx.Param("id")

	url := api + "delete" + "/" + param

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

	if resp.StatusCode != http.StatusOK {
		newResponse(ctx, resp.StatusCode, "Failed to delete student")
		return
	}

	newResponse(ctx, http.StatusOK, "Student deleted successfully")
}

// @Summary		Admin Get Course By Students
// @Security AdminAuth
// @Tags			Students
// @Description	Admin Get Students By CoursId
// @Accept			json
// @Produce		json
// @Param			id path string		true	"course id"
// @Success		200		{byte}	[]byte
// @Failure		500		{object}	Resposne
// @Router			/admin/students/{id}/courses [get]
func (h *Handler) adminGetCoursesStudents(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.services.Kafka.SendMessages("courses-request", id)
	if err != nil {
		newResponse(ctx, http.StatusInternalServerError, "Failed to get information about courses")
		return
	}
	responseData := <-h.responseCh

	ctx.Data(http.StatusOK, "application/json", responseData)
}
