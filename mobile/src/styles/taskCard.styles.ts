import { StyleSheet } from 'react-native';

export const taskCardStyles = StyleSheet.create({
  card: {
    backgroundColor: '#FFFFFF',
    borderRadius: 12,
    padding: 16,
    marginBottom: 12,
    borderWidth: 1,
    borderColor: '#E2E8F0',
    shadowColor: '#0F172A',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.05,
    shadowRadius: 3,
    elevation: 2,
  },
  headerRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    gap: 10,
    marginBottom: 6,
  },
  title: {
    flex: 1,
    fontSize: 16,
    fontWeight: '700',
    color: '#0F172A',
  },
  badge: {
    paddingHorizontal: 10,
    paddingVertical: 4,
    borderRadius: 999,
  },
  badgeText: {
    fontSize: 12,
    fontWeight: '600',
  },
  description: {
    fontSize: 14,
    color: '#475569',
    lineHeight: 20,
    marginBottom: 12,
  },
  metaRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingTop: 10,
    borderTopWidth: 1,
    borderTopColor: '#F1F5F9',
    marginBottom: 12,
  },
  metaText: {
    fontSize: 12,
    color: '#64748B',
  },
  metaBold: {
    fontWeight: '600',
    color: '#334155',
  },
  actionRow: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
    gap: 8,
  },
  actionBtn: {
    paddingHorizontal: 14,
    paddingVertical: 7,
    borderRadius: 8,
  },
  editBtn: {
    backgroundColor: '#EFF6FF',
  },
  editBtnText: {
    color: '#1D4ED8',
    fontSize: 13,
    fontWeight: '600',
  },
  deleteBtn: {
    backgroundColor: '#FEF2F2',
  },
  deleteBtnText: {
    color: '#DC2626',
    fontSize: 13,
    fontWeight: '600',
  },
});
