import Constants from 'expo-constants';
import { Platform } from 'react-native';
import { APIResponse, Task, TaskFilter, TaskListResponse, TaskPayload } from '../types/task';

export class ApiError extends Error {
  status: number;
  code: string;

  constructor(status: number, code: string, message: string) {
    super(message);
    this.status = status;
    this.code = code;
  }
}

function getBaseUrl(): string {
  const hostUri = Constants.expoConfig?.hostUri;
  if (hostUri) {
    const hostIp = hostUri.split(':')[0];
    return `http://${hostIp}:8080/api`;
  }

  if (Platform.OS === 'android') {
    return 'http://10.0.2.2:8080/api';
  }
  return 'http://localhost:8080/api';
}

const BASE_URL = getBaseUrl();

async function request<T>(
  path: string,
  options?: RequestInit,
): Promise<{ data: T; cacheHeader: string | null }> {
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 8000);

  let response: Response;
  try {
    response = await fetch(`${BASE_URL}${path}`, {
      headers: {
        'Content-Type': 'application/json',
        ...(options?.headers || {}),
      },
      signal: controller.signal,
      ...options,
    });
  } catch (err) {
    if (err instanceof Error && err.name === 'AbortError') {
      throw new ApiError(
        0,
        'TIMEOUT',
        'Koneksi terasa lambat. Silakan periksa jaringan internet Anda dan coba lagi.',
      );
    }
    throw new ApiError(
      0,
      'NETWORK_ERROR',
      'Tidak dapat terhubung saat ini. Pastikan perangkat terhubung ke internet lalu coba lagi.',
    );
  } finally {
    clearTimeout(timeoutId);
  }

  const cacheHeader = response.headers.get('X-Cache');
  const json: APIResponse<T> = await response.json();

  if (!response.ok || !json.success) {
    throw new ApiError(
      response.status,
      json.error?.code || 'UNKNOWN_ERROR',
      json.error?.message || 'Terjadi kendala saat memproses permintaan Anda.',
    );
  }

  return { data: json.data as T, cacheHeader };
}

export const taskApi = {
  async listTasks(
    filter: TaskFilter = {},
  ): Promise<{ result: TaskListResponse; cacheStatus: string | null }> {
    const params = new URLSearchParams();
    if (filter.status) params.append('status', filter.status);
    if (filter.keyword) params.append('keyword', filter.keyword);
    if (filter.assignee) params.append('assignee', filter.assignee);
    if (filter.page) params.append('page', String(filter.page));
    if (filter.limit) params.append('limit', String(filter.limit));
    if (filter.sort) params.append('sort', filter.sort);

    const queryString = params.toString();
    const path = queryString ? `/tasks?${queryString}` : '/tasks';
    const { data, cacheHeader } = await request<TaskListResponse>(path);

    return { result: data, cacheStatus: cacheHeader };
  },

  async createTask(payload: TaskPayload): Promise<Task> {
    const { data } = await request<Task>('/tasks', {
      method: 'POST',
      body: JSON.stringify(payload),
    });
    return data;
  },

  async updateTask(id: string, payload: TaskPayload): Promise<Task> {
    const { data } = await request<Task>(`/tasks/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    });
    return data;
  },

  async deleteTask(id: string): Promise<void> {
    await request<null>(`/tasks/${id}`, {
      method: 'DELETE',
    });
  },
};
