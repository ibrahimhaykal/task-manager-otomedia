import { StyleSheet } from 'react-native';

export const appStyles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#F8FAFC',
    paddingTop: 32,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 16,
    paddingVertical: 12,
    backgroundColor: '#FFFFFF',
    borderBottomWidth: 1,
    borderBottomColor: '#E2E8F0',
  },
  headerTitle: {
    fontSize: 22,
    fontWeight: '800',
    color: '#0F172A',
  },
  subHeaderRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    marginTop: 2,
  },
  headerSubtitle: {
    fontSize: 13,
    color: '#64748B',
  },
  cacheBadge: {
    paddingHorizontal: 8,
    paddingVertical: 2,
    borderRadius: 6,
  },
  cacheHit: {
    backgroundColor: '#DCFCE7',
  },
  cacheMiss: {
    backgroundColor: '#E0F2FE',
  },
  cacheBadgeText: {
    fontSize: 11,
    fontWeight: '700',
    color: '#0F172A',
  },
  addBtn: {
    backgroundColor: '#2563EB',
    paddingHorizontal: 16,
    paddingVertical: 10,
    borderRadius: 10,
  },
  addBtnText: {
    color: '#FFFFFF',
    fontWeight: '700',
    fontSize: 14,
  },
  searchSection: {
    flexDirection: 'row',
    gap: 8,
    paddingHorizontal: 16,
    paddingTop: 12,
  },
  searchInput: {
    backgroundColor: '#FFFFFF',
    borderWidth: 1,
    borderColor: '#CBD5E1',
    borderRadius: 10,
    paddingHorizontal: 12,
    paddingVertical: 8,
    fontSize: 13,
    color: '#0F172A',
  },
  filterBarWrapper: {
    paddingVertical: 10,
  },
  filterBar: {
    paddingHorizontal: 16,
    gap: 8,
  },
  filterChip: {
    paddingHorizontal: 14,
    paddingVertical: 7,
    borderRadius: 999,
    backgroundColor: '#FFFFFF',
    borderWidth: 1,
    borderColor: '#CBD5E1',
  },
  filterChipActive: {
    backgroundColor: '#0F172A',
    borderColor: '#0F172A',
  },
  filterChipText: {
    fontSize: 13,
    fontWeight: '600',
    color: '#475569',
  },
  filterChipTextActive: {
    color: '#FFFFFF',
  },
  sortRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: 16,
    paddingBottom: 8,
  },
  sortLabel: {
    fontSize: 12,
    color: '#64748B',
    fontWeight: '600',
  },
  sortOptions: {
    flexDirection: 'row',
    gap: 12,
  },
  sortOptionText: {
    fontSize: 12,
    color: '#64748B',
    fontWeight: '500',
  },
  sortOptionTextActive: {
    color: '#2563EB',
    fontWeight: '700',
  },
  listContent: {
    paddingHorizontal: 16,
    paddingTop: 4,
    paddingBottom: 24,
  },
  centerBox: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 32,
    marginTop: 40,
  },
  emptyTitle: {
    fontSize: 16,
    fontWeight: '700',
    color: '#334155',
    marginBottom: 4,
  },
  errorTitle: {
    fontSize: 16,
    fontWeight: '700',
    color: '#DC2626',
    marginBottom: 4,
  },
  helperText: {
    fontSize: 13,
    color: '#64748B',
    textAlign: 'center',
    marginTop: 6,
  },
  retryBtn: {
    marginTop: 14,
    backgroundColor: '#2563EB',
    paddingHorizontal: 18,
    paddingVertical: 9,
    borderRadius: 8,
  },
  retryBtnText: {
    color: '#FFFFFF',
    fontWeight: '600',
    fontSize: 13,
  },
  paginationRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginTop: 8,
    paddingVertical: 8,
  },
  pageBtn: {
    backgroundColor: '#FFFFFF',
    borderWidth: 1,
    borderColor: '#CBD5E1',
    paddingHorizontal: 14,
    paddingVertical: 8,
    borderRadius: 8,
  },
  pageBtnDisabled: {
    opacity: 0.4,
  },
  pageBtnText: {
    fontSize: 13,
    fontWeight: '600',
    color: '#334155',
  },
  pageIndicator: {
    fontSize: 13,
    fontWeight: '600',
    color: '#475569',
  },
});
