import { StatusBar } from 'expo-status-bar';
import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { SafeAreaProvider, SafeAreaView } from 'react-native-safe-area-context';
import {
  ActivityIndicator,
  Alert,
  FlatList,
  RefreshControl,
  ScrollView,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';
import { taskApi } from './src/api/taskApi';
import { TaskCard } from './src/components/TaskCard';
import { TaskFormModal } from './src/components/TaskFormModal';
import { appStyles as styles } from './src/styles/app.styles';
import { PaginationMetaData, Task, TaskPayload, TaskStatus } from './src/types/task';

const STATUS_FILTERS: { value: TaskStatus | ''; label: string }[] = [
  { value: '', label: 'Semua' },
  { value: 'todo', label: 'To Do' },
  { value: 'in_progress', label: 'In Progress' },
  { value: 'done', label: 'Done' },
];

const SORT_OPTIONS: { value: string; label: string }[] = [
  { value: 'created_at:desc', label: 'Terbaru' },
  { value: 'created_at:asc', label: 'Terlama' },
  { value: 'title:asc', label: 'Judul A-Z' },
];

export default function App() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [pagination, setPagination] = useState<PaginationMetaData>({
    page: 1,
    limit: 10,
    total_items: 0,
    total_pages: 0,
  });

  const [statusFilter, setStatusFilter] = useState<TaskStatus | ''>('');
  const [keywordInput, setKeywordInput] = useState('');
  const [debouncedKeyword, setDebouncedKeyword] = useState('');
  const [assigneeInput, setAssigneeInput] = useState('');
  const [debouncedAssignee, setDebouncedAssignee] = useState('');
  const [sort, setSort] = useState('created_at:desc');
  const [page, setPage] = useState(1);

  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [cacheStatus, setCacheStatus] = useState<string | null>(null);
  const [fetchError, setFetchError] = useState<string | null>(null);

  const [modalVisible, setModalVisible] = useState(false);
  const [editingTask, setEditingTask] = useState<Task | null>(null);
  const existingAssignees = useMemo(() => {
    const unique = new Set<string>();
    tasks.forEach((t) => {
      if (t.assignee && t.assignee.trim() !== '') {
        unique.add(t.assignee.trim());
      }
    });
    return Array.from(unique).sort();
  }, [tasks]);
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedKeyword(keywordInput);
      setDebouncedAssignee(assigneeInput);
      setPage(1);
    }, 400);
    return () => clearTimeout(timer);
  }, [keywordInput, assigneeInput]);

  const fetchTasks = useCallback(
    async (isRefresh = false) => {
      if (isRefresh) {
        setRefreshing(true);
      } else {
        setLoading(true);
      }
      setFetchError(null);

      try {
        const { result, cacheStatus: xCache } = await taskApi.listTasks({
          status: statusFilter,
          keyword: debouncedKeyword,
          assignee: debouncedAssignee,
          page,
          limit: 10,
          sort,
        });
        setTasks(result.tasks || []);
        setPagination(result.pagination);
        setCacheStatus(xCache);
      } catch (err) {
        setFetchError(
          err instanceof Error ? err.message : 'Gagal memuat data dari server backend.',
        );
      } finally {
        setLoading(false);
        setRefreshing(false);
      }
    },
    [statusFilter, debouncedKeyword, debouncedAssignee, page, sort],
  );

  useEffect(() => {
    fetchTasks();
  }, [fetchTasks]);

  const handleOpenCreate = () => {
    setEditingTask(null);
    setModalVisible(true);
  };

  const handleOpenEdit = (task: Task) => {
    setEditingTask(task);
    setModalVisible(true);
  };

  const handleSubmitForm = async (payload: TaskPayload, taskId?: string) => {
    if (taskId) {
      await taskApi.updateTask(taskId, payload);
    } else {
      await taskApi.createTask(payload);
    }
    await fetchTasks();
  };

  const handleConfirmDelete = (task: Task) => {
    Alert.alert('Hapus Task', `Yakin ingin menghapus "${task.title}"?`, [
      { text: 'Batal', style: 'cancel' },
      {
        text: 'Hapus',
        style: 'destructive',
        onPress: async () => {
          try {
            await taskApi.deleteTask(task.id);
            await fetchTasks();
          } catch (err) {
            Alert.alert(
              'Gagal Menghapus',
              err instanceof Error ? err.message : 'Terjadi kesalahan.',
            );
          }
        },
      },
    ]);
  };

  return (
    <SafeAreaProvider>
      <SafeAreaView style={styles.container}>
        <StatusBar style="dark" />

        {/* header */}
        <View style={styles.header}>
          <View>
            <Text style={styles.headerTitle}>Task Manager</Text>
            <View style={styles.subHeaderRow}>
              <Text style={styles.headerSubtitle}>Total: {pagination.total_items} task</Text>
              {cacheStatus ? (
                <View
                  style={[
                    styles.cacheBadge,
                    cacheStatus === 'HIT' ? styles.cacheHit : styles.cacheMiss,
                  ]}
                >
                  <Text style={styles.cacheBadgeText}>Cache: {cacheStatus}</Text>
                </View>
              ) : null}
            </View>
          </View>

          <TouchableOpacity
            style={[styles.addBtn, Boolean(fetchError) && { opacity: 0.5 }]}
            onPress={handleOpenCreate}
            disabled={Boolean(fetchError)}
            activeOpacity={0.8}
          >
            <Text style={styles.addBtnText}>+ Task</Text>
          </TouchableOpacity>
        </View>

        {/* filter search and assignee */}
        <View style={styles.searchSection}>
          <TextInput
            style={[styles.searchInput, { flex: 1.4 }]}
            placeholder="Cari judul / deskripsi..."
            placeholderTextColor="#94A3B8"
            value={keywordInput}
            onChangeText={setKeywordInput}
          />
          <TextInput
            style={[styles.searchInput, { flex: 1 }]}
            placeholder="Assignee..."
            placeholderTextColor="#94A3B8"
            value={assigneeInput}
            onChangeText={setAssigneeInput}
          />
        </View>

        {/* Status Filter Chips */}
        <View style={styles.filterBarWrapper}>
          <ScrollView
            horizontal
            showsHorizontalScrollIndicator={false}
            contentContainerStyle={styles.filterBar}
          >
            {STATUS_FILTERS.map((item) => {
              const active = statusFilter === item.value;
              return (
                <TouchableOpacity
                  key={item.label}
                  style={[styles.filterChip, active && styles.filterChipActive]}
                  onPress={() => {
                    setStatusFilter(item.value);
                    setPage(1);
                  }}
                >
                  <Text style={[styles.filterChipText, active && styles.filterChipTextActive]}>
                    {item.label}
                  </Text>
                </TouchableOpacity>
              );
            })}
          </ScrollView>
        </View>

        {/* Sort Row */}
        <View style={styles.sortRow}>
          <Text style={styles.sortLabel}>Urutkan:</Text>
          <View style={styles.sortOptions}>
            {SORT_OPTIONS.map((opt) => {
              const active = sort === opt.value;
              return (
                <TouchableOpacity
                  key={opt.value}
                  onPress={() => {
                    setSort(opt.value);
                    setPage(1);
                  }}
                >
                  <Text style={[styles.sortOptionText, active && styles.sortOptionTextActive]}>
                    {opt.label}
                  </Text>
                </TouchableOpacity>
              );
            })}
          </View>
        </View>

        {/* Task List */}
        {loading && !refreshing ? (
          <View style={styles.centerBox}>
            <ActivityIndicator size="large" color="#2563EB" />
            <Text style={styles.helperText}>Memuat daftar task...</Text>
          </View>
        ) : fetchError ? (
          <View style={styles.centerBox}>
            <Text style={styles.errorTitle}>Koneksi Bermasalah</Text>
            <Text style={styles.helperText}>{fetchError}</Text>
            <TouchableOpacity style={styles.retryBtn} onPress={() => fetchTasks()}>
              <Text style={styles.retryBtnText}>Coba Lagi</Text>
            </TouchableOpacity>
          </View>
        ) : (
          <FlatList
            data={tasks}
            keyExtractor={(item) => item.id}
            contentContainerStyle={styles.listContent}
            refreshControl={
              <RefreshControl
                refreshing={refreshing}
                onRefresh={() => fetchTasks(true)}
                colors={['#2563EB']}
              />
            }
            renderItem={({ item }) => (
              <TaskCard task={item} onEdit={handleOpenEdit} onDelete={handleConfirmDelete} />
            )}
            ListEmptyComponent={
              <View style={styles.centerBox}>
                <Text style={styles.emptyTitle}>Belum ada task</Text>
                <Text style={styles.helperText}>
                  Tekan tombol "+ Task" di pojok kanan atas untuk menambahkan.
                </Text>
              </View>
            }
            ListFooterComponent={
              pagination.total_pages > 1 ? (
                <View style={styles.paginationRow}>
                  <TouchableOpacity
                    style={[styles.pageBtn, page <= 1 && styles.pageBtnDisabled]}
                    disabled={page <= 1}
                    onPress={() => setPage((p) => Math.max(1, p - 1))}
                  >
                    <Text style={styles.pageBtnText}>Sebelumnya</Text>
                  </TouchableOpacity>

                  <Text style={styles.pageIndicator}>
                    Halaman {pagination.page} / {pagination.total_pages}
                  </Text>

                  <TouchableOpacity
                    style={[
                      styles.pageBtn,
                      page >= pagination.total_pages && styles.pageBtnDisabled,
                    ]}
                    disabled={page >= pagination.total_pages}
                    onPress={() => setPage((p) => p + 1)}
                  >
                    <Text style={styles.pageBtnText}>Berikutnya</Text>
                  </TouchableOpacity>
                </View>
              ) : undefined
            }
          />
        )}

        {/* task modal create and update */}
        <TaskFormModal
          visible={modalVisible}
          initialTask={editingTask}
          existingAssignees={existingAssignees}
          onClose={() => setModalVisible(false)}
          onSubmit={handleSubmitForm}
        />
      </SafeAreaView>
    </SafeAreaProvider>
  );
}
