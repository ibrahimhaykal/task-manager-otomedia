package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/ibrahimhaykal/otomedia-task-manager/internal/model"
)

var (
	ErrTaskNotFound  = errors.New("task not found")
	ErrDuplicateTask = errors.New("task with the same title already exists")
)

type TaskRepository interface {
	FindAll(ctx context.Context, filter model.TaskFilter) ([]model.Task, int64, error)
	FindByID(ctx context.Context, id string) (*model.Task, error)
	Create(ctx context.Context, req model.CreateTaskRequest) (*model.Task, error)
	Update(ctx context.Context, id string, req model.UpdateTaskRequest) (*model.Task, error)
	SoftDelete(ctx context.Context, id string) error
}

type mysqlTaskRepository struct {
	db *sql.DB
}

func NewMySQLTaskRepository(db *sql.DB) *mysqlTaskRepository {
	return &mysqlTaskRepository{db: db}
}

var allowedSortColumns = map[string]string{
	"created_at:asc":  "created_at ASC",
	"created_at:desc": "created_at DESC",
	"updated_at:asc":  "updated_at ASC",
	"updated_at:desc": "updated_at DESC",
	"title:asc":       "title ASC",
	"title:desc":      "title DESC",
	"status:asc":      "status ASC",
	"status:desc":     "status DESC",
	"assignee:asc":    "assignee ASC",
	"assignee:desc":   "assignee DESC",
}

func ResolveOrderBy(sort string) (string, error) {
	cleanedSort := strings.ToLower(strings.TrimSpace(sort))
	if cleanedSort == "" {
		return "created_at DESC", nil
	}
	if orderBy, ok := allowedSortColumns[cleanedSort]; ok {
		return orderBy, nil
	}
	return "", fmt.Errorf("invalid sort parameter: %s", sort)
}

func (r *mysqlTaskRepository) FindAll(ctx context.Context, filter model.TaskFilter) ([]model.Task, int64, error) {
	conditions := []string{"deleted_at IS NULL"}
	args := make([]interface{}, 0)

	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}

	if filter.Keyword != "" {
		conditions = append(conditions, "(title LIKE ? OR description LIKE ?)")
		searchPattern := "%" + filter.Keyword + "%"
		args = append(args, searchPattern, searchPattern)
	}

	if filter.Assignee != "" {
		conditions = append(conditions, "assignee LIKE ?")
		args = append(args, "%"+filter.Assignee+"%")
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	var totalItems int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tasks %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return nil, 0, fmt.Errorf("count tasks: %w", err)
	}

	orderBy, err := ResolveOrderBy(filter.Sort)
	if err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit

	selectQuery := fmt.Sprintf(`
		SELECT id, title, COALESCE(description, ''), status, COALESCE(assignee, ''), created_at, updated_at
		FROM tasks
		%s
		ORDER BY %s
		LIMIT ? OFFSET ?
	`, whereClause, orderBy)

	queryArgs := append(args, filter.Limit, offset)
	rows, err := r.db.QueryContext(ctx, selectQuery, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Assignee,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate tasks: %w", err)
	}

	return tasks, totalItems, nil
}

func (r *mysqlTaskRepository) FindByID(ctx context.Context, id string) (*model.Task, error) {
	query := `
		SELECT id, title, COALESCE(description, ''), status, COALESCE(assignee, ''), created_at, updated_at
		FROM tasks
		WHERE id = ? AND deleted_at IS NULL
	`

	var task model.Task
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Assignee,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTaskNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find task by id: %w", err)
	}

	return &task, nil
}

func (r *mysqlTaskRepository) Create(ctx context.Context, req model.CreateTaskRequest) (*model.Task, error) {
	taskID := uuid.NewString()

	query := `
		INSERT INTO tasks (id, title, description, status, assignee)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, taskID, req.Title, req.Description, req.Status, req.Assignee)
	if err != nil {
		if IsDuplicateEntryError(err) {
			return nil, ErrDuplicateTask
		}
		return nil, fmt.Errorf("insert task: %w", err)
	}

	return r.FindByID(ctx, taskID)
}

func (r *mysqlTaskRepository) Update(ctx context.Context, id string, req model.UpdateTaskRequest) (*model.Task, error) {
	if _, err := r.FindByID(ctx, id); err != nil {
		return nil, err
	}

	query := `
		UPDATE tasks
		SET title = ?, description = ?, status = ?, assignee = ?
		WHERE id = ? AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, req.Title, req.Description, req.Status, req.Assignee, id)
	if err != nil {
		if IsDuplicateEntryError(err) {
			return nil, ErrDuplicateTask
		}
		return nil, fmt.Errorf("update task: %w", err)
	}

	return r.FindByID(ctx, id)
}

func (r *mysqlTaskRepository) SoftDelete(ctx context.Context, id string) error {
	query := `
		UPDATE tasks
		SET deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("soft delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func IsDuplicateEntryError(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
