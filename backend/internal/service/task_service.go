package service

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ibrahimhaykal/otomedia-task-manager/internal/cache"
	"github.com/ibrahimhaykal/otomedia-task-manager/internal/model"
	"github.com/ibrahimhaykal/otomedia-task-manager/internal/repository"
)

const taskListCacheTTL = 60 * time.Second

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type TaskService interface {
	ListTasks(ctx context.Context, filter model.TaskFilter) (*model.TaskListResponse, bool, error)
	GetTaskByID(ctx context.Context, id string) (*model.Task, error)
	CreateTask(ctx context.Context, req model.CreateTaskRequest) (*model.Task, error)
	UpdateTask(ctx context.Context, id string, req model.UpdateTaskRequest) (*model.Task, error)
	DeleteTask(ctx context.Context, id string) error
}

type taskService struct {
	repo  repository.TaskRepository
	cache cache.TaskCache
}

func NewTaskService(repo repository.TaskRepository, taskCache cache.TaskCache) *taskService {
	return &taskService{
		repo:  repo,
		cache: taskCache,
	}
}

func (s *taskService) ListTasks(ctx context.Context, filter model.TaskFilter) (*model.TaskListResponse, bool, error) {
	filter.Status = strings.ToLower(strings.TrimSpace(filter.Status))
	filter.Keyword = strings.TrimSpace(filter.Keyword)
	filter.Assignee = strings.TrimSpace(filter.Assignee)
	filter.Sort = strings.ToLower(strings.TrimSpace(filter.Sort))

	if filter.Status != "" && !model.IsValidStatus(filter.Status) {
		return nil, false, &ValidationError{Message: fmt.Sprintf("invalid status: %s", filter.Status)}
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Sort == "" {
		filter.Sort = "created_at:desc"
	}

	cacheKey := cache.BuildTaskListKey(filter)
	if s.cache != nil {
		if cached, err := s.cache.GetTaskList(ctx, cacheKey); err == nil && cached != nil {
			return cached, true, nil
		}
	}

	tasks, totalItems, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, false, err
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(filter.Limit)))
	}

	result := &model.TaskListResponse{
		Tasks: tasks,
		Pagination: model.PaginationMeta{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}

	if s.cache != nil {
		_ = s.cache.SetTaskList(ctx, cacheKey, result, taskListCacheTTL)
	}

	return result, false, nil
}

func (s *taskService) GetTaskByID(ctx context.Context, id string) (*model.Task, error) {
	if strings.TrimSpace(id) == "" {
		return nil, &ValidationError{Message: "invalid task id"}
	}
	return s.repo.FindByID(ctx, id)
}

func (s *taskService) CreateTask(ctx context.Context, req model.CreateTaskRequest) (*model.Task, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.Status = strings.ToLower(strings.TrimSpace(req.Status))
	req.Assignee = strings.TrimSpace(req.Assignee)

	if req.Title == "" {
		return nil, &ValidationError{Message: "title is required"}
	}
	if req.Status == "" {
		req.Status = model.StatusTodo
	}
	if !model.IsValidStatus(req.Status) {
		return nil, &ValidationError{Message: fmt.Sprintf("invalid status: %s", req.Status)}
	}

	created, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		if err := s.cache.InvalidateTaskLists(ctx); err != nil {
			return nil, fmt.Errorf("invalidate cache after create: %w", err)
		}
	}

	return created, nil
}

func (s *taskService) UpdateTask(ctx context.Context, id string, req model.UpdateTaskRequest) (*model.Task, error) {
	if strings.TrimSpace(id) == "" {
		return nil, &ValidationError{Message: "invalid task id"}
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.Status = strings.ToLower(strings.TrimSpace(req.Status))
	req.Assignee = strings.TrimSpace(req.Assignee)

	if req.Title == "" {
		return nil, &ValidationError{Message: "title is required"}
	}
	if req.Status == "" {
		req.Status = model.StatusTodo
	}
	if !model.IsValidStatus(req.Status) {
		return nil, &ValidationError{Message: fmt.Sprintf("invalid status: %s", req.Status)}
	}

	updated, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		if err := s.cache.InvalidateTaskLists(ctx); err != nil {
			return nil, fmt.Errorf("invalidate cache after update: %w", err)
		}
	}

	return updated, nil
}

func (s *taskService) DeleteTask(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return &ValidationError{Message: "invalid task id"}
	}

	if err := s.repo.SoftDelete(ctx, id); err != nil {
		return err
	}

	if s.cache != nil {
		if err := s.cache.InvalidateTaskLists(ctx); err != nil {
			return fmt.Errorf("invalidate cache after delete: %w", err)
		}
	}

	return nil
}
