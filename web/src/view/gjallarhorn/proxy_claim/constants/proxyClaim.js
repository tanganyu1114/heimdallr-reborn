export const CLAIM_STATES = {
  ClaimStateSubmitted: 1,
  ClaimStatePendingEnable: 2,
  ClaimStatePendingMigrate: 3,
  ClaimStatePendingDeprecate: 4,
  ClaimStateInstantiating: 10,
  ClaimStateEnabling: 11,
  ClaimStateMigrating: 12,
  ClaimStateDeprecating: 13,
  ClaimStateProvisioned: 14,
  ClaimStateActive: 15,
  ClaimStateFailed: 16,
  ClaimStateExpired: 17,
  ClaimStateDeprecated: 18
}

export const CLAIM_STATE_LABELS = {
  1: '已提交',
  2: '待启用',
  3: '待迁移',
  4: '待废弃',
  10: '实例化中',
  11: '启用中',
  12: '迁移中',
  13: '废弃中',
  14: '已就绪',
  15: '活跃',
  16: '失败',
  17: '已过期',
  18: '已废弃'
}

export const CLAIM_STATE_TAG_TYPES = {
  1: 'info',
  2: 'warning',
  3: 'warning',
  4: 'danger',
  10: 'warning',
  11: 'warning',
  12: 'warning',
  13: 'warning',
  14: 'success',
  15: 'success',
  16: 'danger',
  17: 'info',
  18: 'danger'
}

export const TERMINAL_STATES = [
  CLAIM_STATES.ClaimStateProvisioned,
  CLAIM_STATES.ClaimStateActive,
  CLAIM_STATES.ClaimStateFailed,
  CLAIM_STATES.ClaimStateExpired,
  CLAIM_STATES.ClaimStateDeprecated
]

export const STATE_TRANSITIONS = {
  [CLAIM_STATES.ClaimStateProvisioned]: CLAIM_STATES.ClaimStateInstantiating,
  [CLAIM_STATES.ClaimStateActive]: CLAIM_STATES.ClaimStateEnabling,
  [CLAIM_STATES.ClaimStateDeprecated]: CLAIM_STATES.ClaimStateDeprecating,
  [CLAIM_STATES.ClaimStateFailed]: CLAIM_STATES.ClaimStateSubmitted,
  [CLAIM_STATES.ClaimStateExpired]: CLAIM_STATES.ClaimStateSubmitted
}

export const ENV_OPTIONS = [
  { value: 0, label: 'DEV' },
  { value: 1, label: 'SIT' },
  { value: 2, label: 'UAT' },
  { value: 3, label: 'SANDBOX' }
]

export const PROTOCOL_OPTIONS = [
  { label: 'HTTP', value: 'http' },
  { label: 'HTTPS', value: 'https' },
  { label: 'TCP/UDP', value: 'tcp_or_udp' }
]

export const PORT_STATE_LABELS = {
  1: '已禁用',
  2: '已启用',
  3: '已废弃',
  'disabled': '已禁用',
  'enabled': '已启用',
  'deprecated': '已废弃'
}

export const PORT_STATE_TAG_TYPES = {
  1: 'danger',
  2: 'success',
  3: 'info',
  'disabled': 'danger',
  'enabled': 'success',
  'deprecated': 'info'
}

export const PROXY_TYPE_LABELS = {
  'http': 'HTTP代理',
  'https': 'HTTP代理',
  'tcp_or_udp': 'TCP/UDP代理'
}

export const MIGRATION_STATES = {
  Pending: 'pending',
  ClaimsTransitioned: 'claims_transitioned',
  ProxiesReconciled: 'proxies_reconciled',
  Completed: 'completed',
  Failed: 'failed'
}

// 数字状态码到字符串状态码的映射（后端返回数字，前端使用字符串）
export const MIGRATION_STATE_NUMBER_TO_STRING = {
  1: 'pending',
  2: 'claims_transitioned',
  3: 'proxies_reconciled',
  4: 'completed',
  5: 'failed'
}

export const MIGRATION_STATE_LABELS = {
  'pending': '待处理',
  'claims_transitioned': '声明已转换',
  'proxies_reconciled': '代理已协调',
  'completed': '已完成',
  'failed': '失败',
  1: '待处理',
  2: '声明已转换',
  3: '代理已协调',
  4: '已完成',
  5: '失败'
}

export const MIGRATION_STATE_TAG_TYPES = {
  'pending': 'info',
  'claims_transitioned': 'warning',
  'proxies_reconciled': 'warning',
  'completed': 'success',
  'failed': 'danger',
  1: 'info',
  2: 'warning',
  3: 'warning',
  4: 'success',
  5: 'danger'
}

export const MIGRATION_STATE_ICONS = {
  'pending': 'el-icon-time',
  'claims_transitioned': 'el-icon-check',
  'proxies_reconciled': 'el-icon-check',
  'completed': 'el-icon-circle-check',
  'failed': 'el-icon-circle-close',
  1: 'el-icon-time',
  2: 'el-icon-check',
  3: 'el-icon-check',
  4: 'el-icon-circle-check',
  5: 'el-icon-circle-close'
}

export const PROXY_TYPE_TAG_TYPES = {
  'http': '',
  'https': '',
  'tcp_or_udp': 'warning'
}

export const BACKEND_PROTOCOL_MAP = {
  'http': 'http',
  'https': 'https',
  'tcp_or_udp': 'tcp'
}
