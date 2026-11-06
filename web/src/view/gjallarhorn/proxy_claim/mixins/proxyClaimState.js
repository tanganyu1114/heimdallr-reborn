import {
  CLAIM_STATES,
  CLAIM_STATE_LABELS,
  CLAIM_STATE_TAG_TYPES,
  TERMINAL_STATES,
  STATE_TRANSITIONS,
  ENV_OPTIONS,
  PORT_STATE_LABELS,
  PORT_STATE_TAG_TYPES,
  PROXY_TYPE_LABELS,
  PROXY_TYPE_TAG_TYPES,
  BACKEND_PROTOCOL_MAP,
  MIGRATION_STATE_LABELS,
  MIGRATION_STATE_TAG_TYPES,
  MIGRATION_STATE_ICONS,
  MIGRATION_STATE_NUMBER_TO_STRING
} from '../constants/proxyClaim'

export default {
  data() {
    return {
      claimStates: CLAIM_STATES,
      envOptions: ENV_OPTIONS
    }
  },
  computed: {
    statusOptions() {
      return Object.entries(CLAIM_STATE_LABELS).map(([value, label]) => ({
        value: Number(value),
        label
      }))
    }
  },
  methods: {
    getStateLabel(state) {
      return CLAIM_STATE_LABELS[state] || '未知'
    },
    getStateTagType(state) {
      return CLAIM_STATE_TAG_TYPES[state] || 'info'
    },
    isTerminalState(claimState) {
      return TERMINAL_STATES.includes(claimState)
    },
    getNextState(currentState) {
      return STATE_TRANSITIONS[currentState] || CLAIM_STATES.ClaimStateSubmitted
    },
    getPortStateLabel(state) {
      return PORT_STATE_LABELS[state] || state || '未知'
    },
    getPortStateTagType(state) {
      return PORT_STATE_TAG_TYPES[state] || 'info'
    },
    getProxyTypeLabel(protocol) {
      return PROXY_TYPE_LABELS[protocol] || protocol || '未知'
    },
    getProxyTypeTagType(protocol) {
      return PROXY_TYPE_TAG_TYPES[protocol] || ''
    },
    getBackendProtocol(proxyProtocol) {
      return BACKEND_PROTOCOL_MAP[proxyProtocol] || 'tcp'
    },
    getEnvLabel(envId) {
      const env = ENV_OPTIONS.find(e => e.value === envId)
      return env ? env.label : envId
    },
    getVirtualHostNames(serverDetail) {
      const server = serverDetail ? serverDetail.server : null
      const rawName = server ? (server.http_virtual_server_name || '') : ''
      const names = String(rawName).split(/\s+/).filter(name => name)
      return names.length > 0 ? names : ['-']
    },
    getMigrationStateLabel(state) {
      return MIGRATION_STATE_LABELS[state] || '未知'
    },
    getMigrationStateTagType(state) {
      return MIGRATION_STATE_TAG_TYPES[state] || 'info'
    },
    getMigrationStateIcon(state) {
      return MIGRATION_STATE_ICONS[state] || 'el-icon-question'
    },
    normalizeMigrationState(state) {
      if (typeof state === 'number') {
        return MIGRATION_STATE_NUMBER_TO_STRING[state] || 'unknown'
      }
      return state || 'unknown'
    }
  }
}
