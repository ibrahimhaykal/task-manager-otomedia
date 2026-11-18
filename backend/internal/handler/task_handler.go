package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ibrahimhaykal/otomedia-task-manager/internal/model"
	"github.com/ibrahimhaykal/otomedia-task-manager/internal/repository"
	"github.com/ibrahimhaykal/otomedia-task-manager/internal/response"
	"github.com/ibrahimhaykal/otomedia-task-manager/internal/service"
)

type TaskHandler struct {
	service service.TaskService
}

func NewTaskHandler(svc service.TaskService) *TaskHandler {
	return &TaskHandler{service: svc}
}

func (h *TaskHandler) RegisterRoutes(router *gin.RouterGroup) {
	tasks := router.Group("/tasks")
	{
		tasks.GET("", h.ListTasks)
		tasks.GET("/:id", h.GetTaskByID)
		tasks.POST("", h.CreateTask)
		tasks.PUT("/:id", h.UpdateTask)
		tasks.DELETE("/:id", h.DeleteTask)
	}
}

func (h *TaskHandler) ListTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filter := model.TaskFilter{
		Status:   c.Query("status"),
		Keyword:  c.Query("keyword"),
		Assignee: c.Query("assignee"),
		Page:     page,
		Limit:    limit,
		Sort:     c.DefaultQuery("sort", "created_at:desc"),
	}

	result, cacheHit, err := h.service.ListTasks(c.Request.Context(), filter)
	if err != nil {
		h.handleError(c, err)
		return
	}

	if cacheHit {
		c.Header("X-Cache", "HIT")
	} else {
		c.Header("X-Cache", "MISS")
	}

	response.Success(c, http.StatusOK, "Tasks retrieved successfully", result)
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	task, err := h.service.GetTaskByID(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Task retrieved successfully", task)
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req model.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	task, err := h.service.CreateTask(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Task created successfully", task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	var req model.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	task, err := h.service.UpdateTask(c.Request.Context(), id, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Task updated successfully", task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	if err := h.service.DeleteTask(c.Request.Context(), id); err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Task deleted successfully", nil)
}

func parseIDParam(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_ID", "Task ID must be a valid UUID")
		return "", false
	}
	return id, true
}

func (h *TaskHandler) handleError(c *gin.Context, err error) {
	var valErr *service.ValidationError
	switch {
	case errors.As(err, &valErr):
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", valErr.Message)
	case errors.Is(err, repository.ErrTaskNotFound):
		response.Error(c, http.StatusNotFound, "TASK_NOT_FOUND", err.Error())
	case errors.Is(err, repository.ErrDuplicateTask):
		response.Error(c, http.StatusConflict, "DUPLICATE_TITLE", err.Error())
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An unexpected error occurred")
	}
}
