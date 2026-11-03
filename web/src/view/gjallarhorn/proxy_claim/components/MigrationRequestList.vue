<template>
  <div class="migration-request-list">
    <el-table
      :data="migrationRequests"
      border
      stripe
      style="width: 100%"
      :empty-text="emptyText"
    >
      <el-table-column
        prop="source_oa_id"
        label="源OA ID"
        min-width="180"
      >
        <template slot-scope="{ row }">
          <el-link
            type="primary"
            :underline="false"
            @click="handleViewSourceClaim(row.source_oa_id)"
          >
            {{ row.source_oa_id }}
          </el-link>
        </template>
      </el-table-column>

      <el-table-column
        prop="source_claim_origin_state"
        label="源声明状态"
        width="120"
        align="center"
      >
        <template slot-scope="{ row }">
          <el-tag
            v-if="row.source_claim_origin_state"
            :type="getStateTagType(row.source_claim_origin_state)"
            size="small"
          >
            {{ getStateLabel(row.source_claim_origin_state) }}
          </el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>

      <el-table-column
        prop="target_claim_origin_state"
        label="目标声明状态"
        width="120"
        align="center"
      >
        <template slot-scope="{ row }">
          <el-tag
            v-if="row.target_claim_origin_state"
            :type="getStateTagType(row.target_claim_origin_state)"
            size="small"
          >
            {{ getStateLabel(row.target_claim_origin_state) }}
          </el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>

      <el-table-column
        prop="migrate_state"
        label="迁移状态"
        width="150"
        align="center"
      >
        <template slot-scope="{ row }">
          <el-tag
            :type="getMigrationStateTagType(row.migrate_state)"
            size="small"
            effect="dark"
          >
            <i :class="getMigrationStateIcon(row.migrate_state)" style="margin-right: 4px;" />
            {{ getMigrationStateLabel(row.migrate_state) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column
        prop="retry_count"
        label="重试次数"
        width="100"
        align="center"
      >
        <template slot-scope="{ row }">
          <span v-if="row.retry_count > 0" class="retry-count-warning">
            <i class="el-icon-warning" style="margin-right: 4px; color: #E6A23C;" />
            {{ row.retry_count }}
          </span>
          <span v-else class="text-muted">0</span>
        </template>
      </el-table-column>

      <el-table-column
        label="迁移地址数"
        width="120"
        align="center"
      >
        <template slot-scope="{ row }">
          <el-tag size="mini" type="info" effect="plain">
            {{ getAddressCount(row.migrating_addresses) }} 个地址组
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column
        label="操作"
        width="120"
        align="center"
        fixed="right"
      >
        <template slot-scope="{ row }">
          <el-button
            size="mini"
            type="text"
            @click="handleViewDetail(row)"
          >
            查看详情
          </el-button>
          <el-button
            v-if="row.migrate_state === 'failed'"
            size="mini"
            type="text"
            class="retry-btn"
            @click="handleRetry(row)"
          >
            重试
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <MigrationRequestDetailDialog
      :visible.sync="detailDialogVisible"
      :migration-request="selectedMigrationRequest"
    />
  </div>
</template>

<script>
import proxyClaimStateMixin from '../mixins/proxyClaimState'
import MigrationRequestDetailDialog from './MigrationRequestDetailDialog.vue'

export default {
  name: 'MigrationRequestList',
  components: {
    MigrationRequestDetailDialog
  },
  mixins: [proxyClaimStateMixin],
  props: {
    migrationRequests: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      detailDialogVisible: false,
      selectedMigrationRequest: null
    }
  },
  computed: {
    emptyText() {
      return '暂无迁移请求'
    }
  },
  methods: {
    getAddressCount(addresses) {
      if (!addresses) {
        return 0
      }
      if (Array.isArray(addresses)) {
        return addresses.length
      }
      return Object.keys(addresses).length
    },
    handleViewSourceClaim(oaId) {
      this.$emit('view-source-claim', oaId)
    },
    handleViewDetail(row) {
      this.selectedMigrationRequest = row
      this.detailDialogVisible = true
    },
    handleRetry(row) {
      this.$confirm('确定要重试此迁移请求吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.$emit('retry-migration', row)
      }).catch(() => {})
    }
  }
}
</script>

<style scoped>
.migration-request-list {
  padding: 16px 0;
}
.text-muted {
  color: #909399;
  font-size: 13px;
}
.retry-count-warning {
  color: #E6A23C;
  font-weight: 600;
}
.retry-btn {
  color: #e6a23c;
}
.retry-btn:hover {
  color: #ebb563;
}
</style>
