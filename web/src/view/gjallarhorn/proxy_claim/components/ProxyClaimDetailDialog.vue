<template>
  <el-dialog :visible.sync="dialogVisible" title="代理声明详情" width="1200px" @closed="handleDialogClosed">
    <el-descriptions :column="2" border class="claim-descriptions">
      <el-descriptions-item label="OA ID">
        <CopyableText :text="detailData.oa_id" />
      </el-descriptions-item>
      <el-descriptions-item label="申请人">
        <CopyableText :text="detailData.applicant" />
      </el-descriptions-item>
      <el-descriptions-item label="状态">
        <el-tag :type="getStateTagType(detailData.claim_state)">
          {{ getStateLabel(detailData.claim_state) }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="ID" label-class-name="detail-id-label">
        <span class="detail-id-value">{{ detailData.id }}</span>
      </el-descriptions-item>
      <el-descriptions-item label="源项目描述" :span="2">{{ detailData.source_project_desc }}</el-descriptions-item>
      <el-descriptions-item label="目标项目描述" :span="2">{{ detailData.target_project_desc }}</el-descriptions-item>
      <el-descriptions-item label="环境名" :span="2">
        <div v-if="detailData.claim_env_ids && detailData.claim_env_ids.length > 0" class="env-tags-container">
          <el-tag
            v-for="id in detailData.claim_env_ids"
            :key="id"
            size="small"
            class="env-tag"
          >
            {{ getEnvLabel(id) }}
          </el-tag>
        </div>
        <span v-else>-</span>
      </el-descriptions-item>
    </el-descriptions>

    <div v-if="hasMigrationRequests === true" class="migration-collapse-panel">
      <div class="migration-collapse-header" @click="migrationExpanded = !migrationExpanded">
        <div class="header-left">
          <i :class="migrationExpanded ? 'el-icon-arrow-down' : 'el-icon-arrow-right'" class="collapse-icon" />
          <span class="header-title">迁移请求</span>
          <el-badge :value="migrationRequests.length" :max="99" class="migration-badge" />
        </div>
        <div class="header-right">
          <span class="expand-hint">{{ migrationExpanded ? '点击折叠' : '点击展开' }}</span>
        </div>
      </div>
      <MigrationRequestList
        v-show="migrationExpanded"
        :key="detailData.oa_id"
        :migration-requests="migrationRequests"
        @view-source-claim="handleViewSourceClaim"
        @retry-migration="handleRetryMigration"
      />
    </div>

    <el-divider content-position="left">代理详情</el-divider>
    <el-tabs v-model="activeEnvTab" type="border-card">
      <el-tab-pane
        v-for="envDetail in detailData.env_details_list"
        :key="envDetail.env_id"
        :label="getEnvLabel(envDetail.env_id)"
        :name="String(envDetail.env_id)"
      >
        <div class="config-tip">
          <el-tag size="mini" type="info" effect="plain">说明</el-tag>
          <span class="tip-text">每个代理服务对应一个对端接口服务，可配置多个后端地址实现负载均衡</span>
        </div>

        <div v-for="(serverDetail, serverIdx) in envDetail.server_details" :key="serverIdx" class="server-group-flat">
          <div class="server-header-flat">
            <i class="el-icon-s-platform server-icon" />
            <span class="server-label-flat">Bifrost 服务器</span>
            <span class="server-path-flat">{{ serverDetail.server_path || '-' }}</span>
            <span class="server-ip-flat">{{ serverDetail.server_ip || '-' }}</span>
          </div>

          <div class="proxy-list">
            <div v-for="(proxy, proxyIndex) in serverDetail.proxies" :key="proxyIndex" class="proxy-item">
              <div class="proxy-header-row">
                <el-tag :type="getProxyTypeTagType(proxy.target_protocol)" size="medium" effect="dark">
                  {{ getProxyTypeLabel(proxy.target_protocol) }}
                </el-tag>
                <div class="proxy-endpoint">
                  <template v-if="proxy.target_protocol !== 'tcp_or_udp'">
                    <div class="endpoint-primary">
                      <div
                        v-for="(vhostName, vhostIdx) in getVirtualHostNames(serverDetail)"
                        :key="vhostIdx"
                        class="endpoint-address-item"
                      >
                        <el-tag v-if="serverDetail.server" size="small" type="success" class="vserver-tag">
                          {{ vhostName }}
                          <span class="vserver-ip">({{ serverDetail.server_ip || '-' }})</span>
                        </el-tag>
                        <span class="endpoint-port">:{{ proxy.listen_port || '-' }}</span>
                        <el-tooltip :content="`复制代理地址: ${getProxyCopyAddress(proxy, serverDetail, vhostName)}`" placement="top">
                          <el-button
                            size="mini"
                            type="text"
                            icon="el-icon-document-copy"
                            class="copy-proxy-btn"
                            @click="copyToClipboard(getProxyCopyAddress(proxy, serverDetail, vhostName))"
                          />
                        </el-tooltip>
                      </div>
                    </div>
                  </template>
                  <template v-else>
                    <span class="endpoint-addr tcp">
                      {{ serverDetail.server_ip || '-' }}:{{ proxy.listen_port || '-' }}
                    </span>
                    <el-tooltip :content="`复制代理地址: ${getProxyCopyAddress(proxy, serverDetail)}`" placement="top">
                      <el-button
                        size="mini"
                        type="text"
                        icon="el-icon-document-copy"
                        class="copy-proxy-btn"
                        @click="copyToClipboard(getProxyCopyAddress(proxy, serverDetail))"
                      />
                    </el-tooltip>
                  </template>
                </div>
                <el-tag :type="getPortStateTagType(proxy.state)" size="small" class="state-tag">
                  {{ getPortStateLabel(proxy.state) }}
                </el-tag>
              </div>
              <div class="proxy-backend-row">
                <span class="backend-label">后端：</span>
                <span class="backend-addresses">
                  <span v-for="(addr, addrIndex) in proxy.backend_addresses" :key="addrIndex" class="backend-addr">
                    <template v-if="proxy.target_protocol !== 'tcp_or_udp'">
                      {{ getBackendProtocol(proxy.target_protocol) }}://{{ addr.host }}:{{ addr.port }}
                    </template>
                    <template v-else>
                      {{ addr.host }}:{{ addr.port }}
                    </template>
                    <span v-if="addrIndex < proxy.backend_addresses.length - 1" class="addr-separator">, </span>
                  </span>
                </span>
                <el-tag v-if="proxy.backend_addresses && proxy.backend_addresses.length > 1" size="mini" type="success" effect="plain" class="lb-tag">
                  负载均衡
                </el-tag>
              </div>
            </div>

            <div v-if="!serverDetail.proxies || serverDetail.proxies.length === 0" class="empty-proxies">
              <el-empty description="暂无代理服务" :image-size="60" />
            </div>
          </div>
        </div>

        <div v-if="!envDetail.server_details || envDetail.server_details.length === 0" class="empty-proxies">
          <el-empty description="暂无服务器信息" :image-size="80" />
        </div>
      </el-tab-pane>
    </el-tabs>

    <div slot="footer" class="dialog-footer">
      <el-button type="primary" @click="dialogVisible = false">关闭</el-button>
    </div>
  </el-dialog>
</template>

<script>
import CopyableText from './CopyableText.vue'
import MigrationRequestList from './MigrationRequestList.vue'
import proxyClaimStateMixin from '../mixins/proxyClaimState'

export default {
  name: 'ProxyClaimDetailDialog',
  components: {
    CopyableText,
    MigrationRequestList
  },
  mixins: [proxyClaimStateMixin],
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    detailData: {
      type: Object,
      default: () => ({})
    },
    migrationRequests: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      activeEnvTab: '',
      migrationExpanded: false
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
    hasMigrationRequests() {
      return this.migrationRequests && this.migrationRequests.length > 0
    }
  },
  watch: {
    detailData: {
      handler(newVal) {
        if (newVal && newVal.env_details_list && newVal.env_details_list.length > 0) {
          this.activeEnvTab = String(newVal.env_details_list[0].env_id)
        }
      },
      immediate: true
    },
    visible(newVal) {
      if (newVal) {
        this.migrationExpanded = false
      }
    }
  },
  methods: {
    handleDialogClosed() {
      this.activeEnvTab = ''
      this.migrationExpanded = false
    },
    getProxyCopyAddress(proxy, serverDetail, vhostName) {
      const port = proxy.listen_port || '-'
      if (proxy.target_protocol === 'tcp_or_udp') {
        const serverIP = serverDetail.server_ip || '-'
        return `${serverIP}:${port}`
      } else {
        const vserverName = vhostName || (serverDetail.server ? serverDetail.server.http_virtual_server_name : '-')
        return `http://${vserverName}:${port}`
      }
    },
    copyToClipboard(text) {
      if (!text) return
      const input = document.createElement('input')
      input.value = text
      document.body.appendChild(input)
      input.select()
      document.execCommand('copy')
      document.body.removeChild(input)
      this.$message({
        type: 'success',
        message: '已复制到剪贴板',
        duration: 1500
      })
    },
    handleViewSourceClaim(oaId) {
      this.$emit('view-source-claim', oaId)
    },
    handleRetryMigration(migrationRequest) {
      this.$emit('retry-migration', migrationRequest)
    }
  }
}
</script>

<style scoped>
.claim-descriptions {
  margin-bottom: 20px;
}
::v-deep .claim-descriptions .el-descriptions-item__label {
  background-color: #fafbfc;
  font-weight: 500;
  color: #606266;
  width: 120px;
}
::v-deep .claim-descriptions .el-descriptions-item__content {
  color: #303133;
}
.detail-id-label {
  color: #909399 !important;
  font-weight: 400 !important;
}
.detail-id-value {
  color: #909399;
  font-family: 'Courier New', monospace;
  font-size: 12px;
}
.env-tag {
  margin-right: 4px;
  margin-bottom: 4px;
}
.env-tags-container {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 4px;
  min-height: 24px;
}
.migration-collapse-panel {
  margin-bottom: 20px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;
  background: #fff;
}
.migration-collapse-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: linear-gradient(135deg, #fff7e6 0%, #fff0d6 100%);
  border-bottom: 1px solid #e4e7ed;
  cursor: pointer;
  user-select: none;
  transition: background-color 0.2s;
}
.migration-collapse-header:hover {
  background: linear-gradient(135deg, #fff0d6 0%, #ffe8c2 100%);
}
.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}
.migration-badge {
  margin-left: 4px;
}
::v-deep .migration-badge .el-badge__content {
  background-color: #e6a23c;
  border: none;
}
.collapse-icon {
  font-size: 14px;
  color: #e6a23c;
  transition: transform 0.3s;
}
.header-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}
.failed-badge {
  margin-left: 4px;
}
.header-right {
  display: flex;
  align-items: center;
}
.expand-hint {
  font-size: 12px;
  color: #909399;
}
.migration-collapse-body {
  padding: 16px;
  background: #fafbfc;
}
.config-tip {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  padding: 8px 12px;
  background-color: #f4f4f5;
  border-radius: 4px;
  border: 1px solid #e9e9eb;
}
.config-tip .tip-text {
  margin-left: 8px;
  font-size: 13px;
  color: #909399;
}
.server-group-flat {
  margin-bottom: 24px;
  background: #fafbfc;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e4e7ed;
}
.server-header-flat {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  background: linear-gradient(135deg, #f0f5ff 0%, #e6f0ff 100%);
  border-bottom: 1px solid #d9e2ec;
}
.server-icon {
  font-size: 18px;
  color: #409eff;
}
.server-label-flat {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  white-space: nowrap;
}
.server-path-flat {
  font-size: 13px;
  color: #606266;
  font-family: 'Courier New', monospace;
  flex: 1;
}
.server-ip-flat {
  font-size: 13px;
  color: #909399;
  white-space: nowrap;
}
.proxy-list {
  padding: 8px 0;
}
.proxy-item {
  padding: 12px 16px;
  border-bottom: 1px solid #f0f0f0;
  transition: background-color 0.2s;
}
.proxy-item:last-child {
  border-bottom: none;
}
.proxy-item:hover {
  background-color: #f5f7fa;
}
.proxy-header-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 10px;
}
.proxy-endpoint {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}
.endpoint-primary {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 6px;
}
.endpoint-address-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}
.endpoint-port {
  font-size: 13px;
  font-weight: 600;
  color: #303133;
  font-family: 'Courier New', monospace;
  letter-spacing: 0.5px;
  vertical-align: middle;
}
.endpoint-addr {
  font-size: 13px;
  font-weight: 600;
  color: #303133;
  font-family: 'Courier New', monospace;
  letter-spacing: 0.5px;
  vertical-align: middle;
  margin-right: 2px;
}
.vserver-tag {
  font-size: 12px;
  color: #909399;
  font-family: 'Courier New', monospace;
  margin-left: 4px;
  vertical-align: middle;
}
.state-tag {
  margin-left: auto;
}
.proxy-backend-row {
  display: flex;
  align-items: center;
  gap: 6px;
  padding-left: 4px;
}
.backend-label {
  font-size: 13px;
  color: #909399;
  white-space: nowrap;
}
.backend-addresses {
  flex: 1;
  font-size: 13px;
  color: #606266;
  font-family: 'Courier New', monospace;
}
.backend-addr {
  color: #409eff;
}
.addr-separator {
  color: #c0c4cc;
  margin: 0 4px;
}
.lb-tag {
  margin-left: auto;
}
.copy-proxy-btn {
  color: #409eff;
  padding: 4px 6px;
  margin: 0 4px;
  transition: all 0.2s;
}
.copy-proxy-btn:hover {
  color: #66b1ff;
  background-color: #ecf5ff;
}
.empty-proxies {
  padding: 20px 0;
}
</style>
