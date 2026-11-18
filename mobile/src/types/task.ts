export type TaskStatus = 'todo' | 'in_progress' | 'done';

export interface Task {
  id: string;
  title: string;
  description: string;
  status: TaskStatus;
  assignee: string;
  created_at: string;
  updated_at: string;
}

export interface PaginationMetaData {
  page: number;
  limit: number;
  total_items: number;
  total_pages: number;
}

export interface TaskListResponse {
  tasks: Task[];
  pagination: PaginationMetaData;
}

export interface TaskFilter {
  status?: TaskStatus | '';
  keyword?: string;
  assignee?: string;
  page?: number;
  limit?: number;
  sort?: string;
}

export interface TaskPayload {
  title: string;
  description: string;
  status: TaskStatus;
  assignee: string;
}

export interface APIResponse<T> {
  success: boolean;
  message?: string;
  data?: T;
  error?: {
    code: string;
    message: string;
  };
}
