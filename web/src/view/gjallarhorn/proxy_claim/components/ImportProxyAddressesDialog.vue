<template>
  <el-dialog
    :visible.sync="dialogVisible"
    title="导入已有代理地址"
    width="1000px"
    append-to-body
    @closed="handleClosed"
  >
    <el-form :inline="true" :model="searchForm" class="import-search-form">
      <el-form-item label="OA ID">
        <el-input v-model="searchForm.oa_id" placeholder="精确匹配OA ID" clearable style="width: 200px;" />
      </el-form-item>
      <el-form-item label="代理地址">
        <el-input v-model="searchForm.proxy_address" placeholder="模糊匹配主机地址" clearable style="width: 200px;" />
      </el-form-item>
      <el-form-item label="端口">
        <el-input v-model="searchForm.port" placeholder="精确匹配端口号" clearable style="width: 150px;" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" icon="el-icon-search" @click="handleSearch">检索</el-button>
        <el-button @click="resetSearch">重置</el-button>
      </el-form-item>
    </el-form>

    <div v-loading="searchLoading" class="import-tree-container">
      <!-- TODO: Implement lazy loading for tree nodes to improve performance with large datasets -->
      <el-tree
        ref="importTree"
        :data="treeData"
        show-checkbox
        node-key="nodeKey"
        :props="treeProps"
        :default-expand-all="true"
        @check="handleTreeCheck"
        @node-click="handleTreeNodeClick"
      >
        <span slot-scope="{ data }" class="tree-node-content">
          <span v-if="data.type === 'oa'" class="oa-node">
            <el-tag size="mini" type="info">{{ data.oa_id }}</el-tag>
            <span class="node-meta">{{ data.applicant }} | {{ getStateLabel(data.claim_state) }}</span>
          </span>
          <span v-else-if="data.type === 'address-group'" class="proxy-group-node">
            <span class="proxy-index-badge">{{ data.oa_id }}:{{ data.proxy_index }}</span>
            <el-tag size="mini" :type="getProtocolTagType(data.protocol)" effect="plain">{{ data.protocol }}</el-tag>
            <span class="separator">|</span>
            <span class="compact-info">
              <span class="backend-part">
                <span class="backend-label">后端：</span>
                <span class="backend-list">{{ data.addressesDisplay }}</span>
              </span>
              <span v-if="data.servicesByEnv && data.servicesByEnv.length > 0" class="separator">|</span>
              <span v-if="data.servicesByEnv && data.servicesByEnv.length > 0" class="services-part">
                <span class="services-label">代理服务：</span>
                <span v-for="(svcEnv, idx) in data.servicesByEnv" :key="idx" class="svc-env-group">
                  <span class="env-prefix">({{ svcEnv.envs.join(',') }})</span>
                  <span class="svc-endpoint">{{ svcEnv.domain }}:{{ svcEnv.port }}</span>
                </span>
              </span>
            </span>
          </span>
        </span>
      </el-tree>

      <div v-if="treeData.length === 0 && !searchLoading" class="empty-tree-tip">
        <el-empty description="未找到匹配的代理声明，请调整检索条件" />
      </div>
    </div>

    <div v-if="selectedNodes.length > 0" class="selected-import-info">
      <el-alert
        :title="'已选择 ' + selectedNodes.length + ' 个代理服务'"
        type="info"
        :closable="false"
        show-icon
      >
        <div class="selected-nodes-list">
          <el-tag
            v-for="node in selectedNodes"
            :key="node.nodeKey"
            size="small"
            style="margin-right: 8px; margin-bottom: 8px;"
            closable
            @close="removeSelectedNode(node)"
          >
            {{ node.oa_id }}:{{ node.proxy_index }}
          </el-tag>
        </div>
      </el-alert>
    </div>

    <div slot="footer" class="import-dialog-footer">
      <el-button @click="dialogVisible = false">取 消</el-button>
      <el-button type="primary" :disabled="selectedNodes.length === 0" @click="confirmImport">确认导入</el-button>
    </div>
  </el-dialog>
</template>

<script>
import proxyClaimStateMixin from '../mixins/proxyClaimState'
import proxyClaimDetailMixin from '../mixins/proxyClaimDetail'
import { listProxyClaims, getProxyClaimDetail } from '@/api/gjallarhorn_proxy_claim'

export default {
  name: 'ImportProxyAddressesDialog',
  mixins: [proxyClaimStateMixin, proxyClaimDetailMixin],
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    initialSelection: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      searchLoading: false,
      searchForm: {
        oa_id: '',
        proxy_address: '',
        port: ''
      },
      treeData: [],
      selectedNodes: [],
      treeProps: {
        children: 'children',
        label: 'label'
      }
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
    }
  },
  methods: {
    async handleSearch() {
      this.searchLoading = true
      try {
        const params = {
          page: 1,
          page_size: 100
        }
        if (this.searchForm.oa_id) {
          params.oa_id = this.searchForm.oa_id
        }
        if (this.searchForm.proxy_address) {
          params.proxy_address = this.searchForm.proxy_address
        }
        if (this.searchForm.port) {
          params.port = parseInt(this.searchForm.port)
        }

        const res = await listProxyClaims(params)
        const results = res.items || []

        const detailPromises = results.map(claim =>
          getProxyClaimDetail(claim.oa_id).catch(err => {
            console.error('Failed to get detail for ' + claim.oa_id + ':', err)
            return null
          })
        )

        const detailResults = await Promise.all(detailPromises)

        const processedResults = await Promise.all(
          detailResults.map(detailRes => this.buildEnvDetailsList(detailRes))
        )

        this.treeData = this.buildTreeData(results, processedResults)
        this.$nextTick(() => this.applyInitialSelection())
      } catch (error) {
        const errorMsg = error?.message || error?.msg || '未知错误'
        this.$message({
          type: 'error',
          message: '检索失败：' + errorMsg + '，请检查网络连接或联系系统管理员'
        })
      } finally {
        this.searchLoading = false
      }
    },
    resetSearch() {
      this.searchForm = {
        oa_id: '',
        proxy_address: '',
        port: ''
      }
      this.treeData = []
      this.selectedNodes = []
    },
    handleTreeCheck(data, checkedInfo) {
      this.selectedNodes = checkedInfo.checkedNodes.filter(node => node.type === 'address-group')
    },
    applyInitialSelection() {
      if (!this.initialSelection || this.initialSelection.length === 0) return

      const keySet = new Set(
        this.initialSelection.map(s => s.oa_id + ':' + s.proxy_index)
      )

      const keysToCheck = []
      for (const oaNode of this.treeData) {
        for (const group of oaNode.children || []) {
          if (keySet.has(group.oa_id + ':' + group.proxy_index)) {
            keysToCheck.push(group.nodeKey)
          }
        }
      }

      if (keysToCheck.length > 0) {
        this.$refs.importTree.setCheckedKeys(keysToCheck, true)
      }
      this.selectedNodes = this.$refs.importTree.getCheckedNodes().filter(n => n.type === 'address-group')
    },
    handleTreeNodeClick(data) {
      const current = this.$refs.importTree.isChecked(data.nodeKey)
      this.$refs.importTree.setChecked(data.nodeKey, !current)
    },
    removeSelectedNode(node) {
      this.$refs.importTree.setChecked(node.nodeKey, false)
      this.selectedNodes = this.$refs.importTree.getCheckedNodes().filter(n => n.type === 'address-group')
    },
    getProtocolTagType(protocol) {
      const typeMap = {
        'http': 'success',
        'https': 'warning',
        'tcp_or_udp': 'info'
      }
      return typeMap[protocol] || 'info'
    },
    buildTreeData(results, processedResults) {
      return results.map((claim, claimIdx) => {
        const envDetailsList = processedResults[claimIdx]
        if (!envDetailsList || envDetailsList.length === 0) {
          return {
            nodeKey: 'claim-' + claimIdx,
            type: 'oa',
            oa_id: claim.oa_id,
            applicant: claim.applicant,
            claim_state: claim.claim_state,
            label: claim.oa_id,
            children: []
          }
        }

        const proxyGroupMap = new Map()
        let groupCounter = 0

        for (const envDetail of envDetailsList) {
          const envLabel = this.getEnvLabel(envDetail.env_id)
          const serverDetails = envDetail.server_details || []

          for (const serverDetail of serverDetails) {
            const proxiesArray = serverDetail.proxies || []

            for (const proxy of proxiesArray) {
              const isTcpUdp = proxy.target_protocol === 'tcp_or_udp'
              const serverIP = serverDetail.server_ip || '-'
              const rawVserverName = isTcpUdp
                ? serverIP
                : (serverDetail.server ? (serverDetail.server.http_virtual_server_name || '') : '')
              const vserverNames = String(rawVserverName).split(/\s+/).filter(name => name)
              if (vserverNames.length === 0) {
                vserverNames.push('-')
              }
              const listenPort = proxy.listen_port || '-'
              const protocol = proxy.target_protocol || 'http'
              const backendAddresses = proxy.backend_addresses || []
              const proxyIndex = proxy.proxy_index || 0

              const groupKey = claim.oa_id + ':' + proxyIndex

              if (!proxyGroupMap.has(groupKey)) {
                proxyGroupMap.set(groupKey, {
                  nodeKey: 'claim-' + claimIdx + '-group-' + groupCounter,
                  type: 'address-group',
                  oa_id: claim.oa_id,
                  proxy_index: proxyIndex,
                  env_ids: new Set([envDetail.env_id]),
                  protocol: protocol,
                  addresses: backendAddresses.map(addr => ({
                    key: (addr.host || '') + ':' + (addr.port || 0),
                    host: addr.host || '',
                    port: addr.port || 0
                  })),
                  services: [],
                  label: groupKey
                })
                groupCounter++
              }

              const group = proxyGroupMap.get(groupKey)
              group.env_ids.add(envDetail.env_id)

              for (const vserverName of vserverNames) {
                const serviceKey = vserverName + ':' + listenPort
                const existingService = group.services.find(s => s.key === serviceKey)

                if (!existingService) {
                  group.services.push({
                    key: serviceKey,
                    domain: vserverName,
                    port: listenPort,
                    envs: new Set([envLabel])
                  })
                } else {
                  existingService.envs.add(envLabel)
                }
              }
            }
          }
        }

        const groupNodes = Array.from(proxyGroupMap.values()).map(group => {
          const serviceMap = new Map()

          for (const svc of group.services) {
            const svcKey = svc.domain + ':' + svc.port
            if (!serviceMap.has(svcKey)) {
              serviceMap.set(svcKey, {
                domain: svc.domain,
                port: svc.port,
                envs: new Set()
              })
            }
            for (const env of svc.envs) {
              serviceMap.get(svcKey).envs.add(env)
            }
          }

          const servicesByEnv = Array.from(serviceMap.values()).map(svc => ({
            domain: svc.domain,
            port: svc.port,
            envs: Array.from(svc.envs).sort()
          }))

          return {
            ...group,
            addressesDisplay: group.addresses.map(a => a.host + ':' + a.port).join(', '),
            servicesByEnv: servicesByEnv
          }
        })

        return {
          nodeKey: 'claim-' + claimIdx,
          type: 'oa',
          oa_id: claim.oa_id,
          applicant: claim.applicant,
          claim_state: claim.claim_state,
          label: claim.oa_id,
          children: groupNodes
        }
      }).filter(claim => claim.children.length > 0)
    },
    confirmImport() {
      if (this.selectedNodes.length === 0) {
        this.$message({
          type: 'warning',
          message: '请至少选择一个代理服务',
          duration: 2000
        })
        return
      }

      const oaClaimStateMap = new Map()
      for (const node of this.selectedNodes) {
        if (!oaClaimStateMap.has(node.oa_id)) {
          const oaNode = this.treeData.find(t => t.type === 'oa' && t.oa_id === node.oa_id)
          oaClaimStateMap.set(node.oa_id, oaNode ? oaNode.claim_state : 0)
        }
      }

      const proxyIndexMap = new Map()

      for (const node of this.selectedNodes) {
        const key = node.oa_id + ':' + node.proxy_index

        if (!proxyIndexMap.has(key)) {
          proxyIndexMap.set(key, {
            protocol: node.protocol,
            addresses: [],
            is_deprecated: false,
            _isOriginal: false,
            _sourceOaId: node.oa_id,
            _sourceProxyIndex: node.proxy_index,
            _sourceClaimState: oaClaimStateMap.get(node.oa_id)
          })
        }

        const group = proxyIndexMap.get(key)
        for (const addr of node.addresses) {
          const addrKey = addr.host + ':' + addr.port
          const existingAddr = group.addresses.find(a => a.host + ':' + a.port === addrKey)
          if (!existingAddr) {
            group.addresses.push({
              host: addr.host,
              port: addr.port
            })
          }
        }
      }

      const addressesArray = Array.from(proxyIndexMap.values())

      if (addressesArray.length > 0) {
        this.$emit('import', addressesArray)
        const oaIds = [...new Set(this.selectedNodes.map(n => n.oa_id))]
        this.$message({
          type: 'success',
          message: '已成功导入 ' + addressesArray.length + ' 个代理服务（来自 ' + oaIds.join(', ') + '）',
          duration: 2000
        })
      }

      this.dialogVisible = false
    },
    handleClosed() {
      this.treeData = []
      this.selectedNodes = []
      this.searchForm = {
        oa_id: '',
        proxy_address: '',
        port: ''
      }
    }
  }
}
</script>

<style scoped>
.import-search-form {
  margin-bottom: 16px;
}
.import-tree-container {
  max-height: 400px;
  overflow-y: auto;
  overflow-x: hidden;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 12px;
  background-color: #fafafa;
  width: 100%;
  box-sizing: border-box;
}

.import-tree-container /deep/ .el-tree {
  width: 100%;
  overflow: hidden;
}
.import-tree-container /deep/ .el-tree-node {
  max-width: 100%;
  overflow: hidden;
}
.import-tree-container /deep/ .el-tree-node__content {
  height: auto !important;
  min-height: 32px;
  padding: 4px 0;
  line-height: 1.6;
  display: flex;
  align-items: flex-start;
  max-width: 100%;
  overflow: hidden;
}
.import-tree-container /deep/ .el-tree-node__children {
  overflow: visible;
}

.tree-node-content {
  flex: 1;
  display: block;
  font-size: 12px;
  line-height: 1.8;
  min-width: 0;
  max-width: 100%;
  overflow-wrap: break-word;
  word-break: break-word;
  white-space: normal;
}
.oa-node {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}
.oa-node .node-meta {
  font-size: 12px;
  color: #909399;
}
.proxy-group-node {
  display: flex;
  align-items: baseline;
  gap: 6px;
  flex-wrap: wrap;
  max-width: 100%;
}
.proxy-index-badge {
  display: inline-block;
  padding: 1px 6px;
  font-size: 11px;
  font-family: 'Consolas', 'Monaco', monospace;
  color: #409eff;
  background: #ecf5ff;
  border: 1px solid #d9ecff;
  border-radius: 3px;
  font-weight: 600;
  flex-shrink: 0;
}
.compact-info {
  display: block;
  font-size: 12px;
  min-width: 0;
  flex: 1 1 auto;
  width: 0;
  overflow-wrap: break-word;
}
.separator {
  color: #c0c4cc;
  font-weight: 300;
  margin: 0 2px;
  white-space: nowrap;
}
.backend-part {
  display: inline;
}
.backend-label {
  color: #909399;
  white-space: nowrap;
}
.backend-list {
  color: #606266;
  font-family: 'Courier New', monospace;
  word-break: break-all;
  overflow-wrap: break-word;
}
.services-part {
  display: inline;
}
.services-label {
  color: #909399;
  white-space: nowrap;
}
.svc-env-group {
  display: inline-flex;
  align-items: center;
  gap: 2px;
}
.env-prefix {
  color: #67c23a;
  font-size: 11px;
}
.svc-endpoint {
  color: #409eff;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  word-break: break-all;
  overflow-wrap: break-word;
}
.empty-tree-tip {
  padding: 40px 0;
}
.selected-import-info {
  margin-top: 16px;
}
.selected-nodes-list {
  margin-top: 8px;
}
.import-dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
