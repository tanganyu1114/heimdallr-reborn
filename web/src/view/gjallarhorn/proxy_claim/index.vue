<template>
  <div>
    <ProxyClaimSearch
      :status-options="statusOptions"
      :search-info="searchInfo"
      @search="handleSearch"
      @reset="handleReset"
      @quick-filter="handleQuickFilter"
      @apply="openApplyDialog"
    />

    <ProxyClaimTable
      :table-data="tableData"
      :loading="loading"
      @view="viewDetail"
      @approve="handleApprove"
      @update="openUpdateDialog"
      @apply="openApplyDialog"
    />

    <el-pagination
      class="pagination-wrapper"
      :current-page="page"
      :page-size="pageSize"
      :page-sizes="[10, 30, 50, 100]"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
    />

    <ProxyClaimDetailDialog
      ref="detailDialog"
      :visible.sync="detailDialogVisible"
      :detail-data="detailData"
      :migration-requests="migrationRequests"
    />

    <ProxyClaimApplyDialog
      :visible.sync="applyDialogVisible"
      :type="applyType"
      :oa-id="applyOaId"
      @success="refreshTable"
    />
  </div>
</template>

<script>
import ProxyClaimSearch from './components/ProxyClaimSearch.vue'
import ProxyClaimTable from './components/ProxyClaimTable.vue'
import ProxyClaimDetailDialog from './components/ProxyClaimDetailDialog.vue'
import ProxyClaimApplyDialog from './components/ProxyClaimApplyDialog.vue'
import proxyClaimStateMixin from './mixins/proxyClaimState'
import proxyClaimDetailMixin from './mixins/proxyClaimDetail'
import { listProxyClaims, getProxyClaimDetail, approveProxyClaim, listMigrationRequests } from '@/api/gjallarhorn_proxy_claim'

export default {
  name: 'ProxyClaim',
  components: {
    ProxyClaimSearch,
    ProxyClaimTable,
    ProxyClaimDetailDialog,
    ProxyClaimApplyDialog
  },
  mixins: [proxyClaimStateMixin, proxyClaimDetailMixin],
  data() {
    return {
      page: 1,
      total: 0,
      pageSize: 10,
      tableData: [],
      loading: false,
      searchInfo: {
        oa_id: '',
        applicant: '',
        proxy_address: '',
        port: '',
        status: ''
      },
      detailDialogVisible: false,
      detailData: {},
      migrationRequests: [],
      applyDialogVisible: false,
      applyType: 'create',
      applyOaId: ''
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    async getTableData() {
      this.loading = true
      try {
        const res = await listProxyClaims(this.buildQueryParams())
        this.tableData = res.items || []
        this.total = res.total_count || 0
      } catch (error) {
        const errorMsg = error?.message || error?.msg || '未知错误'
        this.$message({
          type: 'error',
          message: `获取列表失败：${errorMsg}，请检查网络连接或刷新页面重试`
        })
      } finally {
        this.loading = false
      }
    },
    buildQueryParams() {
      const params = {
        page: this.page,
        page_size: this.pageSize
      }
      if (this.searchInfo.oa_id) {
        params.oa_id = this.searchInfo.oa_id
      }
      if (this.searchInfo.applicant) {
        params.applicant = this.searchInfo.applicant
      }
      if (this.searchInfo.proxy_address) {
        params.proxy_address = this.searchInfo.proxy_address
      }
      if (this.searchInfo.port) {
        params.port = parseInt(this.searchInfo.port)
      }
      if (this.searchInfo.status !== '' && this.searchInfo.status !== null) {
        params.status = this.searchInfo.status
      }
      return params
    },
    handleSearch(searchInfo) {
      this.searchInfo = { ...searchInfo }
      this.page = 1
      this.getTableData()
    },
    handleReset() {
      this.searchInfo = {
        oa_id: '',
        applicant: '',
        proxy_address: '',
        port: '',
        status: ''
      }
      this.page = 1
      this.getTableData()
    },
    handleQuickFilter(status) {
      this.searchInfo.status = status
      this.page = 1
      this.getTableData()
    },
    handleSizeChange(val) {
      this.pageSize = val
      this.getTableData()
    },
    handleCurrentChange(val) {
      this.page = val
      this.getTableData()
    },
    refreshTable() {
      this.getTableData()
    },
    openApplyDialog() {
      this.applyType = 'create'
      this.applyOaId = ''
      this.applyDialogVisible = true
    },
    openUpdateDialog(row) {
      this.applyType = 'update'
      this.applyOaId = row.oa_id
      this.applyDialogVisible = true
    },
    async viewDetail(row) {
      try {
        const res = await getProxyClaimDetail(row.oa_id)

        const claimData = res.claim || {}
        const proxyPassAddresses = claimData.proxy_pass_addresses || {}

        const envDetailsList = await this.buildEnvDetailsList(res)

        this.detailData = {
          ...claimData,
          proxy_pass_addresses: Object.values(proxyPassAddresses),
          env_details_list: envDetailsList
        }

        const migrationRes = await listMigrationRequests({ target_oa_id: row.oa_id })
        this.migrationRequests = Array.isArray(migrationRes) ? migrationRes : []

        this.detailDialogVisible = true
      } catch (error) {
        const errorMsg = error?.message || error?.msg || '未知错误'
        this.$message({
          type: 'error',
          message: `获取详情失败：${errorMsg}，请检查网络连接或联系系统管理员`
        })
      }
    },
    async handleApprove(row) {
      this.$confirm(
        `确认审批以下代理申请？\n\nOA ID: ${row.oa_id}\n申请人: ${row.applicant}\n状态: ${this.getStateLabel(row.claim_state)}`,
        '审批确认',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
          dangerouslyUseHTMLString: false
        }
      ).then(async() => {
        try {
          await approveProxyClaim(row.oa_id)
          this.$message({
            type: 'success',
            message: '审批成功，代理将进入实例化流程'
          })
          this.getTableData()
        } catch (error) {
          const errorMsg = error?.message || error?.msg || '未知错误'
          this.$message({
            type: 'error',
            message: `审批失败：${errorMsg}，请检查代理配置或联系系统管理员`
          })
        }
      }).catch(() => {})
    }
  }
}
</script>

<style scoped>
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: 20px;
}
</style>
