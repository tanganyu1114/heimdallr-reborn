<template>
  <el-table
    ref="multipleTable"
    :data="tableData"
    :loading="loading"
    border
    stripe
    style="width: 100%"
    tooltip-effect="dark"
    :header-cell-style="{ background: 'linear-gradient(180deg, #f8f9fb 0%, #f0f2f5 100%)', color: '#303133', fontWeight: '600', fontSize: '13px', borderBottom: '2px solid #e4e7ed' }"
    :cell-style="{ padding: '6px 0', fontSize: '13px' }"
  >
    <template slot="empty">
      <el-empty description="暂无代理声明数据">
        <el-button type="primary" @click="$emit('apply')">创建代理申请</el-button>
      </el-empty>
    </template>
    <el-table-column label="OA ID" prop="oa_id" width="150" align="center" header-align="center">
      <template slot-scope="scope">
        <CopyableText :text="scope.row.oa_id" />
      </template>
    </el-table-column>
    <el-table-column label="申请人" prop="applicant" width="100" header-align="center">
      <template slot-scope="scope">
        <CopyableText :text="scope.row.applicant" />
      </template>
    </el-table-column>
    <el-table-column label="源项目描述" prop="source_project_desc" min-width="180" show-overflow-tooltip header-align="center" />
    <el-table-column label="目标项目描述" prop="target_project_desc" min-width="180" show-overflow-tooltip header-align="center" />
    <el-table-column label="状态" prop="claim_state" width="90" align="center">
      <template slot-scope="scope">
        <el-tag :type="getStateTagType(scope.row.claim_state)">
          {{ getStateLabel(scope.row.claim_state) }}
        </el-tag>
      </template>
    </el-table-column>
    <el-table-column label="环境" prop="claim_env_ids" min-width="100" max-width="220" header-align="center">
      <template slot-scope="scope">
        <div class="env-tags-container">
          <el-tag
            v-for="id in scope.row.claim_env_ids"
            :key="id"
            size="mini"
            class="env-tag"
          >
            {{ getEnvLabel(id) }}
          </el-tag>
          <span v-if="!scope.row.claim_env_ids || scope.row.claim_env_ids.length === 0">-</span>
        </div>
      </template>
    </el-table-column>
    <el-table-column label="操作" width="240" align="center">
      <template slot-scope="scope">
        <div class="operation-cell">
          <el-button size="small" type="primary" icon="el-icon-view" @click="$emit('view', scope.row)">详情</el-button>
          <el-button
            v-if="scope.row.claim_state === claimStates.ClaimStateProvisioned"
            size="small"
            type="success"
            icon="el-icon-check"
            @click="$emit('approve', scope.row)"
          >审批</el-button>
          <el-button
            v-if="isTerminalState(scope.row.claim_state)"
            size="small"
            type="info"
            icon="el-icon-edit"
            @click="$emit('update', scope.row)"
          >更新</el-button>
        </div>
      </template>
    </el-table-column>
  </el-table>
</template>

<script>
import CopyableText from './CopyableText.vue'
import proxyClaimStateMixin from '../mixins/proxyClaimState'

export default {
  name: 'ProxyClaimTable',
  components: {
    CopyableText
  },
  mixins: [proxyClaimStateMixin],
  props: {
    tableData: {
      type: Array,
      required: true
    },
    loading: {
      type: Boolean,
      default: false
    }
  }
}
</script>

<style scoped>
.copyable-text {
  cursor: pointer;
  color: #409eff;
  transition: color 0.2s;
}
.copyable-text:hover {
  color: #66b1ff;
  text-decoration: underline;
}
.env-tag {
  margin-right: 4px;
  margin-bottom: 4px;
}
.env-tag:last-child {
  margin-right: 0;
}
.env-tags-container {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 4px;
  min-height: 24px;
}
::v-deep .el-table .cell .el-tag:not(.env-tag) {
  min-width: 60px;
  text-align: center;
  justify-content: center;
}
.env-tag {
  min-width: auto !important;
  padding: 0 8px;
}
.operation-cell {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 6px;
  flex-wrap: nowrap;
}
::v-deep .el-table__body tr:hover > td {
  background-color: #e8f4ff !important;
}
::v-deep .el-table--striped .el-table__body tr.el-table__row--striped:hover td.el-table__cell {
  background-color: #e8f4ff !important;
}
</style>
