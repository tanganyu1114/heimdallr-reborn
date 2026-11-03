<template>
  <el-dialog
    :visible.sync="dialogVisible"
    title="迁移请求详情"
    width="900px"
    :append-to-body="true"
    @closed="handleClosed"
  >
    <div v-if="migrationRequest" class="migration-detail">
      <el-descriptions :column="2" border class="detail-descriptions">
        <el-descriptions-item label="迁移ID">
          <span class="id-value">{{ migrationRequest.id }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="迁移名称">
          <CopyableText :text="migrationRequest.name" />
        </el-descriptions-item>
        <el-descriptions-item label="源OA ID">
          <el-link
            type="primary"
            :underline="false"
            @click="handleViewSourceClaim(migrationRequest.source_oa_id)"
          >
            {{ migrationRequest.source_oa_id }}
          </el-link>
        </el-descriptions-item>
        <el-descriptions-item label="目标OA ID">
          <CopyableText :text="migrationRequest.target_oa_id" />
        </el-descriptions-item>
        <el-descriptions-item label="源声明状态">
          <el-tag
            v-if="migrationRequest.source_claim_origin_state"
            :type="getStateTagType(migrationRequest.source_claim_origin_state)"
            size="small"
          >
            {{ getStateLabel(migrationRequest.source_claim_origin_state) }}
          </el-tag>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="目标声明状态">
          <el-tag
            v-if="migrationRequest.target_claim_origin_state"
            :type="getStateTagType(migrationRequest.target_claim_origin_state)"
            size="small"
          >
            {{ getStateLabel(migrationRequest.target_claim_origin_state) }}
          </el-tag>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="迁移状态" :span="2">
          <el-tag
            :type="getMigrationStateTagType(migrationRequest.migrate_state)"
            size="medium"
            effect="dark"
          >
            <i :class="getMigrationStateIcon(migrationRequest.migrate_state)" style="margin-right: 6px;" />
            {{ getMigrationStateLabel(migrationRequest.migrate_state) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="重试次数">
          <span v-if="migrationRequest.retry_count > 0" class="retry-count-warning">
            <i class="el-icon-warning" style="margin-right: 4px; color: #E6A23C;" />
            {{ migrationRequest.retry_count }} 次
          </span>
          <span v-else class="text-muted">0 次</span>
        </el-descriptions-item>
      </el-descriptions>

      <el-divider content-position="left">迁移进度</el-divider>
      <div class="migration-progress">
        <el-steps :active="currentStepIndex" finish-status="success" align-center>
          <el-step
            v-for="step in migrationSteps"
            :key="step.key"
            :title="step.label"
            :status="getStepStatus(step.key, migrationRequest.migrate_state)"
          />
        </el-steps>
      </div>

      <el-divider content-position="left">迁移地址与索引映射</el-divider>
      <div v-if="hasMigrationAddresses" class="migration-address-table">
        <el-table :data="mergedMigrationData" border>
          <el-table-column prop="sourceIndex" label="源代理索引" align="center" width="120">
            <template slot-scope="{ row }">
              <span class="index-label">索引 {{ row.sourceIndex }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="targetIndex" label="目标代理索引" align="center" width="120">
            <template slot-scope="{ row }">
              <span class="index-label">索引 {{ row.targetIndex }}</span>
            </template>
          </el-table-column>
          <el-table-column label="目标后端地址">
            <template slot-scope="{ row }">
              <div class="address-group-cell">
                <div class="address-group-header">
                  <el-tag :type="getProtocolTagType(row.protocol)" size="mini" effect="dark">
                    {{ getProtocolLabel(row.protocol) }}
                  </el-tag>
                  <span class="address-count">{{ row.addresses.length }} 个地址</span>
                </div>
                <div class="address-list-cell">
                  <div
                    v-for="(addr, idx) in row.addresses"
                    :key="idx"
                    class="address-item-cell"
                  >
                    <span class="addr-text">
                      <template v-if="row.protocol !== 'tcp_or_udp'">
                        {{ row.protocol }}://{{ addr.host }}:{{ addr.port }}
                      </template>
                      <template v-else>
                        {{ addr.host }}:{{ addr.port }}
                      </template>
                    </span>
                    <el-tag v-if="row.is_deprecated" size="mini" type="info" effect="plain">
                      已废弃
                    </el-tag>
                  </div>
                </div>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <el-empty v-else description="暂无迁移地址" :image-size="60" />
    </div>

    <div slot="footer" class="dialog-footer">
      <!-- TODO: 后续按需实现重试迁移功能 -->
      <el-button type="primary" @click="dialogVisible = false">关闭</el-button>
    </div>
  </el-dialog>
</template>

<script>
import CopyableText from './CopyableText.vue'
import { PROTOCOL_OPTIONS } from '../constants/proxyClaim'
import proxyClaimStateMixin from '../mixins/proxyClaimState'

export default {
  name: 'MigrationRequestDetailDialog',
  components: {
    CopyableText
  },
  mixins: [proxyClaimStateMixin],
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    migrationRequest: {
      type: Object,
      default: null
    }
  },
  data() {
    return {
      migrationSteps: [
        { key: 'pending', label: '待处理' },
        { key: 'claims_transitioned', label: '声明已转换' },
        { key: 'proxies_reconciled', label: '代理已协调' },
        { key: 'completed', label: '已完成' }
      ]
    }
  },
  computed: {
    dialogVisible: {
      get() {
        return this.visible
      },
      set(val) {
        this.$emit('update:visible', val)
      }
    },
    currentStepIndex() {
      if (!this.migrationRequest || !this.migrationRequest.migrate_state) {
        return 0
      }
      const currentState = this.normalizeMigrationState(this.migrationRequest.migrate_state)
      const stepKeys = this.migrationSteps.map(s => s.key)
      const currentIndex = stepKeys.indexOf(currentState)
      if (currentState === 'failed') {
        return currentIndex
      }
      return currentIndex + 1
    },
    hasMigrationAddresses() {
      if (!this.migrationRequest || !this.migrationRequest.migrating_addresses) {
        return false
      }
      const addresses = this.migrationRequest.migrating_addresses
      if (Array.isArray(addresses)) {
        return addresses.length > 0
      }
      return Object.keys(addresses).length > 0
    },
    migratingAddressesArray() {
      if (!this.migrationRequest || !this.migrationRequest.migrating_addresses) {
        return []
      }
      const addresses = this.migrationRequest.migrating_addresses
      if (Array.isArray(addresses)) {
        return addresses
      }
      return Object.values(addresses)
    },
    mergedMigrationData() {
      if (!this.migrationRequest) {
        return []
      }
      const indexesMap = this.migrationRequest.migrating_indexes_map || {}
      const addresses = this.migrationRequest.migrating_addresses || {}

      // migrating_indexes_map: {"0": 5, "1": 6, "2": 7}
      // migrating_addresses: {"5": {protocol, addresses, is_deprecated}, ...}
      return Object.entries(indexesMap).map(([sourceIdx, targetIdx]) => {
        const addrGroup = addresses[targetIdx] || {}
        return {
          sourceIndex: sourceIdx,
          targetIndex: targetIdx,
          protocol: addrGroup.protocol || '-',
          addresses: addrGroup.addresses || [],
          is_deprecated: addrGroup.is_deprecated || false
        }
      })
    },
    isMigrationFailed() {
      if (!this.migrationRequest || !this.migrationRequest.migrate_state) {
        return false
      }
      const state = this.normalizeMigrationState(this.migrationRequest.migrate_state)
      return state === 'failed'
    },
    indexMappingData() {
      if (!this.migrationRequest || !this.migrationRequest.migrating_indexes_map) {
        return []
      }
      return Object.entries(this.migrationRequest.migrating_indexes_map).map(([source, target]) => ({
        source: '索引 ' + source,
        target: '索引 ' + target
      }))
    }
  },
  methods: {
    getStepStatus(stepKey, currentState) {
      const stepOrder = ['pending', 'claims_transitioned', 'proxies_reconciled', 'completed']
      const normalizedState = this.normalizeMigrationState(currentState)
      const currentIdx = stepOrder.indexOf(normalizedState)
      const stepIdx = stepOrder.indexOf(stepKey)

      if (normalizedState === 'failed') {
        return stepIdx < currentIdx ? 'success' : (stepIdx === currentIdx ? 'error' : 'wait')
      }
      return stepIdx < currentIdx ? 'success' : (stepIdx === currentIdx ? 'process' : 'wait')
    },
    getProtocolLabel(protocol) {
      const option = PROTOCOL_OPTIONS.find(o => o.value === protocol)
      return option ? option.label : protocol
    },
    getProtocolTagType(protocol) {
      const typeMap = {
        'http': '',
        'https': 'success',
        'tcp_or_udp': 'warning'
      }
      return typeMap[protocol] || 'info'
    },
    handleViewSourceClaim(oaId) {
      this.$emit('view-source-claim', oaId)
    },
    handleRetry() {
      this.$confirm('确定要重试此迁移请求吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.$emit('retry-migration', this.migrationRequest)
      }).catch(() => {})
    },
    handleClosed() {
      this.$emit('closed')
    }
  }
}
</script>

<style scoped>
.migration-detail {
  padding: 0 8px;
}
.detail-descriptions {
  margin-bottom: 20px;
}
::v-deep .detail-descriptions .el-descriptions-item__label {
  background-color: #fafbfc;
  font-weight: 500;
  color: #606266;
  width: 140px;
}
.id-value {
  font-family: 'Courier New', monospace;
  font-size: 13px;
  color: #909399;
}
.text-muted {
  color: #909399;
  font-size: 13px;
}
.retry-count-warning {
  color: #E6A23C;
  font-weight: 600;
}
.migration-progress {
  padding: 20px 0;
}
.migration-address-table {
  margin-top: 16px;
}
.index-label {
  font-family: 'Courier New', monospace;
  font-size: 13px;
  color: #606266;
}
.address-group-cell {
  padding: 4px 0;
}
.address-group-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}
.address-count {
  font-size: 12px;
  color: #909399;
}
.address-list-cell {
  padding-left: 4px;
}
.address-item-cell {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 0;
  border-bottom: 1px solid #f0f0f0;
}
.address-item-cell:last-child {
  border-bottom: none;
}
.addr-text {
  font-size: 13px;
  color: #606266;
  font-family: 'Courier New', monospace;
}
</style>
