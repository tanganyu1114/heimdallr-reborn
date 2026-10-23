<template>
  <div>
    <div class="search-term">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="OA ID">
          <el-input v-model="searchInfo.oa_id" placeholder="请输入OA ID" clearable />
        </el-form-item>
        <el-form-item label="申请人">
          <el-input v-model="searchInfo.applicant" placeholder="请输入申请人" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" placeholder="请选择状态" clearable class="status-select">
            <el-option
              v-for="option in statusOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit">查询</el-button>
          <el-button @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
      <div class="quick-filters">
        <span class="filter-label">快捷筛选：</span>
        <el-button-group>
          <el-button :size="'small'" :type="!searchInfo.status ? 'primary' : ''" @click="quickFilter('')">全部</el-button>
          <el-button :size="'small'" :type="searchInfo.status === 14 ? 'success' : ''" @click="quickFilter(14)">已就绪</el-button>
          <el-button :size="'small'" :type="searchInfo.status === 15 ? 'success' : ''" @click="quickFilter(15)">活跃</el-button>
          <el-button :size="'small'" :type="searchInfo.status === 16 ? 'danger' : ''" @click="quickFilter(16)">失败</el-button>
          <el-button :size="'small'" :type="searchInfo.status === 1 ? 'info' : ''" @click="quickFilter(1)">已提交</el-button>
        </el-button-group>
      </div>
      <div class="action-buttons">
        <el-button type="primary" icon="el-icon-plus" @click="openApplyDialog">代理申请</el-button>
        <el-button type="warning" icon="el-icon-sort" @click="openMigrateDialog">代理迁移</el-button>
      </div>
    </div>

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
          <el-button type="primary" @click="openApplyDialog">创建代理申请</el-button>
        </el-empty>
      </template>
      <el-table-column label="OA ID" prop="oa_id" width="150" align="center" header-align="center">
        <template slot-scope="scope">
          <el-tooltip content="点击复制" placement="top">
            <span class="copyable-text" @click="copyToClipboard(scope.row.oa_id)">
              {{ scope.row.oa_id }}
            </span>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column label="申请人" prop="applicant" width="100" header-align="center">
        <template slot-scope="scope">
          <el-tooltip content="点击复制" placement="top">
            <span class="copyable-text" @click="copyToClipboard(scope.row.applicant)">
              {{ scope.row.applicant }}
            </span>
          </el-tooltip>
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
            <el-button size="small" type="primary" icon="el-icon-view" @click="viewDetail(scope.row)">详情</el-button>
            <el-button
              v-if="scope.row.claim_state === claimStates.ClaimStateProvisioned"
              size="small"
              type="success"
              icon="el-icon-check"
              @click="handleApprove(scope.row)"
            >审批</el-button>
            <el-button
              v-if="isTerminalState(scope.row.claim_state)"
              size="small"
              type="info"
              icon="el-icon-edit"
              @click="openUpdateDialog(scope.row)"
            >更新</el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

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

    <!-- 详情对话框 -->
    <el-dialog :visible.sync="detailDialogVisible" title="代理声明详情" width="1200px">
      <el-descriptions :column="2" border class="claim-descriptions">
        <el-descriptions-item label="OA ID">
          <el-tooltip content="点击复制" placement="top">
            <span class="copyable-text" @click="copyToClipboard(detailData.oa_id)">
              {{ detailData.oa_id }}
            </span>
          </el-tooltip>
        </el-descriptions-item>
        <el-descriptions-item label="申请人">
          <el-tooltip content="点击复制" placement="top">
            <span class="copyable-text" @click="copyToClipboard(detailData.applicant)">
              {{ detailData.applicant }}
            </span>
          </el-tooltip>
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
                      <span class="endpoint-primary">
                        <el-tag v-if="serverDetail.server" size="small" type="success" class="vserver-tag">
                          {{ serverDetail.server.http_virtual_server_name }}
                          <span class="vserver-ip">({{ serverDetail.server_ip || '-' }})</span>
                        </el-tag>
                        <span class="endpoint-port">:{{ proxy.listen_port || '-' }}</span>
                      </span>
                    </template>
                    <template v-else>
                      <span class="endpoint-addr tcp">
                        {{ serverDetail.server_ip || '-' }}:{{ proxy.listen_port || '-' }}
                      </span>
                    </template>
                    <el-tooltip :content="`复制代理地址: ${getProxyCopyAddress(proxy, serverDetail)}`" placement="top">
                      <el-button
                        size="mini"
                        type="text"
                        icon="el-icon-document-copy"
                        class="copy-proxy-btn"
                        @click="copyToClipboard(getProxyCopyAddress(proxy, serverDetail))"
                      />
                    </el-tooltip>
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
        <el-button type="primary" @click="detailDialogVisible = false">关闭</el-button>
      </div>
    </el-dialog>

    <!-- 申请/更新对话框 -->
    <el-dialog
      :visible.sync="applyDialogVisible"
      :title="applyType === 'create' ? '提交代理申请' : '更新代理申请'"
      width="800px"
      @closed="resetApplyForm"
    >
      <el-form ref="applyForm" :model="applyFormData" :rules="applyRules" label-width="140px" label-position="right">
        <el-form-item v-if="applyType === 'update'" label="ID" prop="id">
          <el-input v-model="applyFormData.id" disabled />
        </el-form-item>
        <el-form-item label="OA ID" prop="oa_id">
          <el-input v-model="applyFormData.oa_id" :disabled="applyType === 'update'" placeholder="请输入唯一OA标识（字母、数字、下划线）" />
        </el-form-item>
        <el-form-item v-if="applyType === 'create'" label="申请人" prop="applicant">
          <el-input v-model="applyFormData.applicant" placeholder="请输入申请人姓名" />
        </el-form-item>
        <el-form-item label="源项目描述" prop="source_project_desc">
          <el-input v-model="applyFormData.source_project_desc" type="textarea" :rows="2" placeholder="请输入源项目描述" />
        </el-form-item>
        <el-form-item label="目标项目描述" prop="target_project_desc">
          <el-input v-model="applyFormData.target_project_desc" type="textarea" :rows="2" placeholder="请输入目标项目描述" />
        </el-form-item>

        <el-divider content-position="left">代理转发地址配置</el-divider>
        <div class="config-tip">
          <el-tag size="mini" type="info" effect="plain">说明</el-tag>
          <span class="tip-text">每个代理服务对应一个对端接口服务，可配置多个后端地址实现负载均衡</span>
        </div>

        <div v-for="(addrGroup, groupIndex) in applyFormData.proxy_pass_addresses" :key="groupIndex" class="address-group">
          <div class="address-group-header">
            <div class="header-left">
              <span class="group-label">代理服务 {{ groupIndex + 1 }}</span>
              <el-select v-model="addrGroup.protocol" placeholder="选择协议" size="small" style="width: 140px; margin-left: 12px;">
                <el-option label="HTTP" value="http" />
                <el-option label="HTTPS" value="https" />
                <el-option label="TCP/UDP" value="tcp_or_udp" />
              </el-select>
            </div>
            <el-button v-if="applyFormData.proxy_pass_addresses.length > 1 && !addrGroup._isOriginal" size="small" type="danger" plain icon="el-icon-delete" @click="removeAddressGroup(groupIndex)">删除代理服务</el-button>
          </div>

          <div v-if="applyType === 'update' && addrGroup._isOriginal" class="deprecate-row">
            <el-switch
              v-model="addrGroup.is_deprecated"
              active-color="#F56C6C"
              inactive-color="#67C23A"
              size="mini"
            />
            <span class="switch-label">
              {{ addrGroup.is_deprecated ? '已废弃' : '正常' }}
            </span>
          </div>

          <div class="address-list">
            <div v-if="addrGroup.addresses.length > 1" class="address-list-tip">
              <el-tag size="mini" type="success" effect="plain">负载均衡</el-tag>
              <span class="tip-text">多个地址将按顺序轮询分发请求</span>
            </div>
            <div v-for="(addr, addrIndex) in addrGroup.addresses" :key="addrIndex" class="address-item">
              <span class="address-label"><span class="required-star">*</span>主机地址</span>
              <el-input
                v-model="addr.host"
                placeholder="请输入主机地址，如 127.0.0.1"
                size="small"
                :class="{ 'is-error': !addr.host || !addr.host.trim() }"
                style="width: 220px;"
              />
              <el-input-number v-model="addr.port" :min="0" :max="65535" size="small" style="width: 130px; margin-left: 10px;" />
              <el-button v-if="addrGroup.addresses.length > 1" size="small" type="danger" plain icon="el-icon-delete" style="margin-left: 10px;" @click="removeAddress(groupIndex, addrIndex)">删除</el-button>
            </div>
            <el-button class="add-address-btn" size="small" type="text" icon="el-icon-plus" @click="addAddress(groupIndex)">添加地址</el-button>
          </div>
        </div>

        <el-button class="add-group-btn" type="dashed" icon="el-icon-plus" @click="addAddressGroup">添加代理服务</el-button>

        <el-form-item label="环境名" prop="claim_env_ids">
          <el-select v-model="applyFormData.claim_env_ids" multiple placeholder="选择环境名" style="width: 100%;">
            <el-option v-for="env in envOptions" :key="env.value" :label="env.label" :value="env.value" />
          </el-select>
        </el-form-item>
      </el-form>

      <div slot="footer" class="dialog-footer">
        <el-button @click="closeApplyDialog">取 消</el-button>
        <el-button type="primary" @click="submitApply">确 定</el-button>
      </div>
    </el-dialog>

    <!-- 迁移对话框 -->
    <el-dialog :visible.sync="migrateDialogVisible" title="迁移代理申请" width="800px" @closed="resetMigrateForm">
      <el-form ref="migrateForm" :model="migrateFormData" :rules="migrateRules" label-width="160px" label-position="right">
        <el-form-item label="源OA ID" prop="source_oa_id">
          <el-input v-model="migrateFormData.source_oa_id" placeholder="请输入源OA ID" />
        </el-form-item>
        <el-form-item label="目标OA ID" prop="target_oa_id">
          <el-input v-model="migrateFormData.target_oa_id" placeholder="请输入目标OA ID" />
        </el-form-item>

        <el-divider content-position="left">迁移地址配置</el-divider>
        <div class="config-tip">
          <el-tag size="mini" type="info" effect="plain">说明</el-tag>
          <span class="tip-text">每个代理服务对应一个对端接口服务，可配置多个后端地址实现负载均衡</span>
        </div>

        <div v-for="(addrGroup, groupIndex) in migrateFormData.migrating_addresses" :key="groupIndex" class="address-group">
          <div class="address-group-header">
            <div class="header-left">
              <span class="group-label">代理服务 {{ groupIndex + 1 }}</span>
              <el-select v-model="addrGroup.protocol" placeholder="选择协议" size="small" style="width: 140px; margin-left: 12px;">
                <el-option label="HTTP" value="http" />
                <el-option label="HTTPS" value="https" />
                <el-option label="TCP/UDP" value="tcp_or_udp" />
              </el-select>
            </div>
            <el-button v-if="migrateFormData.migrating_addresses.length > 1" size="small" type="danger" plain icon="el-icon-delete" @click="removeMigrateAddressGroup(groupIndex)">删除代理服务</el-button>
          </div>

          <div class="address-list">
            <div v-if="addrGroup.addresses.length > 1" class="address-list-tip">
              <el-tag size="mini" type="success" effect="plain">负载均衡</el-tag>
              <span class="tip-text">多个地址将按顺序轮询分发请求</span>
            </div>
            <div v-for="(addr, addrIndex) in addrGroup.addresses" :key="addrIndex" class="address-item">
              <span class="address-label"><span class="required-star">*</span>主机地址</span>
              <el-input
                v-model="addr.host"
                placeholder="请输入主机地址，如 127.0.0.1"
                size="small"
                :class="{ 'is-error': !addr.host || !addr.host.trim() }"
                style="width: 220px;"
              />
              <el-input-number v-model="addr.port" :min="0" :max="65535" size="small" style="width: 130px; margin-left: 10px;" />
              <el-button v-if="addrGroup.addresses.length > 1" size="small" type="danger" plain icon="el-icon-delete" style="margin-left: 10px;" @click="removeMigrateAddress(groupIndex, addrIndex)">删除</el-button>
            </div>
            <el-button class="add-address-btn" size="small" type="text" icon="el-icon-plus" @click="addMigrateAddress(groupIndex)">添加地址</el-button>
          </div>
        </div>

        <el-button class="add-group-btn" type="dashed" icon="el-icon-plus" @click="addMigrateAddressGroup">添加代理服务</el-button>

        <el-divider content-position="left">迁移索引映射</el-divider>
        <el-form-item label="索引映射" prop="migrating_indexes_map">
          <div v-for="(target, source) in migrateFormData.migrating_indexes_map" :key="source" class="index-mapping-item">
            <el-input-number v-model="migrateFormData.migrating_indexes_map[source]" :min="0" size="small" style="width: 100px;" />
            <el-button size="mini" type="danger" icon="el-icon-delete" style="margin-left: 10px;" @click="removeIndexMapping(source)" />
          </div>
          <el-button size="small" type="primary" icon="el-icon-plus" style="margin-top: 10px;" @click="addIndexMapping">添加映射</el-button>
        </el-form-item>
      </el-form>

      <div slot="footer" class="dialog-footer">
        <el-button @click="closeMigrateDialog">取 消</el-button>
        <el-button type="primary" @click="submitMigrate">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { listProxyClaims, getProxyClaim, getProxyClaimDetail, applyProxyClaim, approveProxyClaim, migrateProxyClaim } from '@/api/gjallarhorn_proxy_claim'
import { findHmdrHost } from '@/api/hmdr_host'
import { findHmdrGroup } from '@/api/hmdr_group'

export default {
  name: 'ProxyClaim',
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
        status: ''
      },
      detailDialogVisible: false,
      detailData: {},
      activeEnvTab: '',
      applyDialogVisible: false,
      applyType: 'create',
      applyFormData: {
        id: 0,
        oa_id: '',
        applicant: '',
        source_project_desc: '',
        target_project_desc: '',
        proxy_pass_addresses: [],
        claim_env_ids: [],
        claim_state: 1
      },
      applyRules: {
        oa_id: [
          { required: true, message: '请输入OA ID', trigger: 'blur' },
          { pattern: /^[a-zA-Z0-9_]+$/, message: 'OA ID仅允许字母、数字和下划线', trigger: 'blur' }
        ],
        applicant: [{ required: true, message: '请输入申请人', trigger: 'blur' }],
        target_project_desc: [{ required: true, message: '请输入目标项目描述', trigger: 'blur' }],
        proxy_pass_addresses: [
          {
            validator: (rule, value, callback) => {
              if (!value || value.length === 0) {
                callback(new Error('请至少添加一个对端地址组'))
              } else {
                const hasEmptyAddress = value.some(group => !group.addresses || group.addresses.length === 0)
                if (hasEmptyAddress) {
                  callback(new Error('每个对端地址组请至少添加一个地址'))
                } else {
                  const hasEmptyHost = value.some(group =>
                    group.addresses.some(addr => !addr.host || !addr.host.trim())
                  )
                  if (hasEmptyHost) {
                    callback(new Error('请填写所有地址的主机地址'))
                  } else {
                    callback()
                  }
                }
              }
            },
            trigger: 'change'
          }
        ]
      },
      migrateDialogVisible: false,
      migrateFormData: {
        source_oa_id: '',
        target_oa_id: '',
        migrating_addresses: [],
        migrating_indexes_map: {}
      },
      migrateRules: {
        source_oa_id: [
          { required: true, message: '请输入源OA ID', trigger: 'blur' },
          { pattern: /^[a-zA-Z0-9_]+$/, message: 'OA ID仅允许字母、数字和下划线', trigger: 'blur' }
        ],
        target_oa_id: [
          { required: true, message: '请输入目标OA ID', trigger: 'blur' },
          { pattern: /^[a-zA-Z0-9_]+$/, message: 'OA ID仅允许字母、数字和下划线', trigger: 'blur' }
        ],
        migrating_addresses: [
          {
            validator: (rule, value, callback) => {
              if (!value || value.length === 0) {
                callback(new Error('请至少添加一个对端地址组'))
              } else {
                const hasEmptyAddress = value.some(group => !group.addresses || group.addresses.length === 0)
                if (hasEmptyAddress) {
                  callback(new Error('每个对端地址组请至少添加一个地址'))
                } else {
                  const hasEmptyHost = value.some(group =>
                    group.addresses.some(addr => !addr.host || !addr.host.trim())
                  )
                  if (hasEmptyHost) {
                    callback(new Error('请填写所有地址的主机地址'))
                  } else {
                    callback()
                  }
                }
              }
            },
            trigger: 'change'
          }
        ]
      },
      claimStates: {
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
      },
      claimStateLabels: {
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
      },
      envOptions: [
        { value: 0, label: 'DEV' },
        { value: 1, label: 'SIT' },
        { value: 2, label: 'UAT' },
        { value: 3, label: 'SANDBOX' }
      ]
    }
  },
  computed: {
    statusOptions() {
      return Object.entries(this.claimStateLabels).map(([value, label]) => ({
        value: Number(value),
        label
      }))
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    async getTableData() {
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
      if (this.searchInfo.status !== '' && this.searchInfo.status !== null) {
        params.status = this.searchInfo.status
      }

      this.loading = true
      try {
        const res = await listProxyClaims(params)
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
    onSubmit() {
      this.page = 1
      this.getTableData()
    },
    onReset() {
      this.searchInfo = {
        oa_id: '',
        applicant: '',
        status: ''
      }
      this.page = 1
      this.getTableData()
    },
    quickFilter(status) {
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
    async viewDetail(row) {
      try {
        const res = await getProxyClaimDetail(row.oa_id)

        const claimData = res.claim || {}
        const proxyPassAddresses = claimData.proxy_pass_addresses || {}

        const envDetailsList = []
        if (res.env_details && typeof res.env_details === 'object') {
          const envDetailsObj = res.env_details
          for (const envId in envDetailsObj) {
            const envDetail = envDetailsObj[envId]
            const serverDetailsArray = []

            if (envDetail.server_details && Array.isArray(envDetail.server_details)) {
              for (const serverDetail of envDetail.server_details) {
                let serverIP = '-'
                let serverPath = '-'

                if (serverDetail.server && serverDetail.server.host_id) {
                  try {
                    const hostRes = await findHmdrHost({ ID: serverDetail.server.host_id })
                    if (hostRes && hostRes.data && hostRes.data.rehmdrHost) {
                      const host = hostRes.data.rehmdrHost
                      serverIP = host.ipaddr || '-'

                      const promises = []
                      if (host.groupId) {
                        promises.push(findHmdrGroup({ ID: host.groupId }).catch(err => {
                          console.error(`Failed to get group info for group_id ${host.groupId}:`, err)
                          return null
                        }))
                      }

                      const results = await Promise.all(promises)
                      const groupRes = results[0]

                      if (groupRes && groupRes.data && groupRes.data.rehmdrGroup) {
                        const groupName = groupRes.data.rehmdrGroup.name || '-'
                        const hostName = host.name || '-'
                        const bifrostName = serverDetail.server.name || '-'
                        serverPath = `${groupName} / ${hostName} / ${bifrostName}`
                      } else {
                        const hostName = host.name || '-'
                        const bifrostName = serverDetail.server.name || '-'
                        serverPath = `${hostName} / ${bifrostName}`
                      }
                    } else {
                      serverIP = `HostID: ${serverDetail.server.host_id}`
                    }
                  } catch (error) {
                    console.error(`Failed to get host info for host_id ${serverDetail.server.host_id}:`, error)
                    serverIP = `HostID: ${serverDetail.server.host_id}`
                  }
                }

                const proxiesObj = serverDetail.proxies || {}
                const proxiesArray = typeof proxiesObj === 'object' && !Array.isArray(proxiesObj)
                  ? Object.values(proxiesObj).sort((a, b) => (a.proxy_index || 0) - (b.proxy_index || 0))
                  : (Array.isArray(proxiesObj) ? proxiesObj : [])

                serverDetailsArray.push({
                  server: serverDetail.server,
                  server_ip: serverIP,
                  server_path: serverPath,
                  proxies: proxiesArray
                })
              }
            }

            envDetailsList.push({
              env_id: parseInt(envId),
              server_details: serverDetailsArray
            })
          }
          envDetailsList.sort((a, b) => a.env_id - b.env_id)
        }

        this.detailData = {
          ...claimData,
          proxy_pass_addresses: Object.values(proxyPassAddresses),
          env_details_list: envDetailsList
        }

        if (envDetailsList.length > 0) {
          this.activeEnvTab = String(envDetailsList[0].env_id)
        }

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
    },
    isTerminalState(claimState) {
      const terminalStates = [
        this.claimStates.ClaimStateProvisioned,
        this.claimStates.ClaimStateActive,
        this.claimStates.ClaimStateFailed,
        this.claimStates.ClaimStateExpired,
        this.claimStates.ClaimStateDeprecated
      ]
      return terminalStates.includes(claimState)
    },
    getNextState(currentState) {
      const stateTransitions = {
        [this.claimStates.ClaimStateProvisioned]: this.claimStates.ClaimStateInstantiating,
        [this.claimStates.ClaimStateActive]: this.claimStates.ClaimStateEnabling,
        [this.claimStates.ClaimStateDeprecated]: this.claimStates.ClaimStateDeprecating,
        [this.claimStates.ClaimStateFailed]: this.claimStates.ClaimStateSubmitted,
        [this.claimStates.ClaimStateExpired]: this.claimStates.ClaimStateSubmitted
      }
      return stateTransitions[currentState] || this.claimStates.ClaimStateSubmitted
    },
    openApplyDialog() {
      this.applyType = 'create'
      this.resetApplyForm()
      this.applyDialogVisible = true
    },
    async openUpdateDialog(row) {
      try {
        const res = await getProxyClaim(row.oa_id)
        this.applyType = 'update'
        // 将后端返回的对象格式转换为数组格式
        const proxyPassAddresses = res.proxy_pass_addresses || {}
        const addressesArray = Object.values(proxyPassAddresses).map(item => ({
          ...item,
          _isOriginal: true
        }))
        this.applyFormData = {
          id: res.id,
          oa_id: res.oa_id,
          applicant: res.applicant,
          source_project_desc: res.source_project_desc,
          target_project_desc: res.target_project_desc,
          proxy_pass_addresses: addressesArray.length > 0 ? addressesArray : [
            {
              protocol: 'http',
              addresses: [{ host: '', port: 8080 }],
              is_deprecated: false,
              _isOriginal: false
            }
          ],
          claim_env_ids: res.claim_env_ids || [],
          claim_state: res.claim_state
        }
        this.applyDialogVisible = true
      } catch (error) {
        this.$message({
          type: 'error',
          message: '获取详情失败'
        })
      }
    },
    closeApplyDialog() {
      this.applyDialogVisible = false
    },
    resetApplyForm() {
      this.applyType = 'create'
      this.applyFormData = {
        id: 0,
        oa_id: '',
        applicant: '',
        source_project_desc: '',
        target_project_desc: '',
        proxy_pass_addresses: [
          {
            protocol: 'http',
            addresses: [
              { host: '', port: 8080 }
            ],
            is_deprecated: false
          }
        ],
        claim_env_ids: [],
        claim_state: this.claimStates.ClaimStateSubmitted
      }
      this.$nextTick(() => {
        if (this.$refs.applyForm) {
          this.$refs.applyForm.resetFields()
        }
      })
    },
    async submitApply() {
      this.$refs.applyForm.validate(async(valid) => {
        if (!valid) {
          this.$message({
            type: 'warning',
            message: '请检查表单填写是否完整和正确',
            duration: 2000
          })
          this.$nextTick(() => {
            const firstError = document.querySelector('.el-form-item__error')
            if (firstError) {
              firstError.scrollIntoView({ behavior: 'smooth', block: 'center' })
            }
          })
          return
        }

        // 将 proxy_pass_addresses 数组转换为对象格式（后端 MapList 期望的格式）
        // 过滤掉内部标记字段 _isOriginal
        const cleanedAddresses = this.applyFormData.proxy_pass_addresses.map(item => {
          // eslint-disable-next-line no-unused-vars
          const { _isOriginal, ...rest } = item
          return rest
        })

        // 根据当前状态计算下一个状态
        const nextState = this.applyType === 'update'
          ? this.getNextState(this.applyFormData.claim_state)
          : this.applyFormData.claim_state

        const submitData = {
          ...this.applyFormData,
          claim_state: nextState,
          proxy_pass_addresses: cleanedAddresses.reduce((obj, item, index) => {
            obj[index] = item
            return obj
          }, {})
        }

        try {
          await applyProxyClaim(submitData)
          this.$message({
            type: 'success',
            message: this.applyType === 'create' ? '申请提交成功' : '更新成功'
          })
          this.closeApplyDialog()
          this.getTableData()
        } catch (error) {
          const errorMsg = error?.message || error?.msg || '未知错误'
          this.$message({
            type: 'error',
            message: `${this.applyType === 'create' ? '申请提交' : '更新'}失败：${errorMsg}，请检查输入参数或联系系统管理员`
          })
        }
      })
    },
    addAddressGroup() {
      this.applyFormData.proxy_pass_addresses.push({
        protocol: 'http',
        addresses: [
          { host: '', port: 8080 }
        ],
        is_deprecated: false,
        _isOriginal: false
      })
    },
    removeAddressGroup(index) {
      this.applyFormData.proxy_pass_addresses.splice(index, 1)
    },
    addAddress(groupIndex) {
      this.applyFormData.proxy_pass_addresses[groupIndex].addresses.push({
        host: '',
        port: 8080
      })
    },
    removeAddress(groupIndex, addressIndex) {
      this.applyFormData.proxy_pass_addresses[groupIndex].addresses.splice(addressIndex, 1)
    },
    openMigrateDialog() {
      this.resetMigrateForm()
      this.migrateDialogVisible = true
    },
    closeMigrateDialog() {
      this.migrateDialogVisible = false
    },
    resetMigrateForm() {
      this.migrateFormData = {
        source_oa_id: '',
        target_oa_id: '',
        migrating_addresses: [],
        migrating_indexes_map: {}
      }
      this.$nextTick(() => {
        if (this.$refs.migrateForm) {
          this.$refs.migrateForm.resetFields()
        }
      })
    },
    async submitMigrate() {
      this.$refs.migrateForm.validate(async(valid) => {
        if (!valid) {
          this.$message({
            type: 'warning',
            message: '请检查表单填写是否完整和正确',
            duration: 2000
          })
          this.$nextTick(() => {
            const firstError = document.querySelector('.el-form-item__error')
            if (firstError) {
              firstError.scrollIntoView({ behavior: 'smooth', block: 'center' })
            }
          })
          return
        }

        try {
          await migrateProxyClaim(this.migrateFormData)
          this.$message({
            type: 'success',
            message: '迁移请求提交成功，请等待处理完成'
          })
          this.closeMigrateDialog()
          this.getTableData()
        } catch (error) {
          const errorMsg = error?.message || error?.msg || '未知错误'
          this.$message({
            type: 'error',
            message: `迁移失败：${errorMsg}，请检查迁移参数或联系系统管理员`
          })
        }
      })
    },
    addMigrateAddressGroup() {
      this.migrateFormData.migrating_addresses.push({
        protocol: 'http',
        addresses: [
          { host: '', port: 8080 }
        ],
        is_deprecated: false
      })
    },
    removeMigrateAddressGroup(index) {
      this.migrateFormData.migrating_addresses.splice(index, 1)
    },
    addMigrateAddress(groupIndex) {
      this.migrateFormData.migrating_addresses[groupIndex].addresses.push({
        host: '',
        port: 8080
      })
    },
    removeMigrateAddress(groupIndex, addressIndex) {
      this.migrateFormData.migrating_addresses[groupIndex].addresses.splice(addressIndex, 1)
    },
    addIndexMapping() {
      const newIndex = Object.keys(this.migrateFormData.migrating_indexes_map).length
      this.$set(this.migrateFormData.migrating_indexes_map, String(newIndex), newIndex)
    },
    removeIndexMapping(source) {
      this.$delete(this.migrateFormData.migrating_indexes_map, source)
    },
    getStateLabel(state) {
      return this.claimStateLabels[state] || '未知'
    },
    getStateTagType(state) {
      const typeMap = {
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
      return typeMap[state] || 'info'
    },
    getPortStateLabel(state) {
      const stateMap = {
        1: '已禁用',
        2: '已启用',
        3: '已废弃',
        'disabled': '已禁用',
        'enabled': '已启用',
        'deprecated': '已废弃'
      }
      return stateMap[state] || state || '未知'
    },
    getPortStateTagType(state) {
      const typeMap = {
        1: 'danger',
        2: 'success',
        3: 'info',
        'disabled': 'danger',
        'enabled': 'success',
        'deprecated': 'info'
      }
      return typeMap[state] || 'info'
    },
    getProxyTypeLabel(protocol) {
      const labelMap = {
        'http': 'HTTP代理',
        'https': 'HTTP代理',
        'tcp_or_udp': 'TCP/UDP代理'
      }
      return labelMap[protocol] || protocol || '未知'
    },
    getProxyTypeTagType(protocol) {
      const typeMap = {
        'http': '',
        'https': '',
        'tcp_or_udp': 'warning'
      }
      return typeMap[protocol] || ''
    },
    getBackendProtocol(proxyProtocol) {
      const protocolMap = {
        'http': 'http',
        'https': 'https',
        'tcp_or_udp': 'tcp'
      }
      return protocolMap[proxyProtocol] || 'tcp'
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
    getEnvLabel(envId) {
      const env = this.envOptions.find(e => e.value === envId)
      return env ? env.label : envId
    },
    getProxyCopyAddress(proxy, serverDetail) {
      const port = proxy.listen_port || '-'
      if (proxy.target_protocol === 'tcp_or_udp') {
        const serverIP = serverDetail.server_ip || '-'
        return `${serverIP}:${port}`
      } else {
        const vserverName = serverDetail.server ? serverDetail.server.http_virtual_server_name : '-'
        return `http://${vserverName}:${port}`
      }
    }
  }
}
</script>

<style scoped>
.search-term {
  margin-bottom: 20px;
}
.quick-filters {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  padding: 12px 16px;
  background-color: #f5f7fa;
  border-radius: 6px;
}
.filter-label {
  font-size: 13px;
  color: #606266;
  font-weight: 500;
  margin-right: 12px;
  white-space: nowrap;
}
.action-buttons {
  display: flex;
  gap: 12px;
  margin-top: 12px;
}
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
/* 状态标签固定尺寸 */
::v-deep .el-table .cell .el-tag:not(.env-tag) {
  min-width: 60px;
  text-align: center;
  justify-content: center;
}
.env-tag {
  min-width: auto !important;
  padding: 0 8px;
}

/* 操作列按钮左对齐 */
.operation-cell {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 6px;
  flex-wrap: nowrap;
}

/* 表格行悬停效果 - 整行统一高亮 */
::v-deep .el-table__body tr:hover > td {
  background-color: #e8f4ff !important;
}
/* 斑马纹行悬停高亮 */
::v-deep .el-table--striped .el-table__body tr.el-table__row--striped:hover td.el-table__cell {
  background-color: #e8f4ff !important;
}

/* 详情描述表格优化 */
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

/* 代理地址复制按钮 */
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
.address-group {
  border: 1px solid #e4e7ed;
  padding: 16px;
  margin-bottom: 16px;
  border-radius: 6px;
  background-color: #fff;
  transition: box-shadow 0.2s;
}
.address-group:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
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
.address-group-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}
.header-left {
  display: flex;
  align-items: center;
  flex-wrap: nowrap;
}
.switch-label {
  margin-left: 8px;
  font-size: 12px;
  color: #606266;
  white-space: nowrap;
}
.deprecate-row {
  display: block;
  margin-bottom: 10px;
  padding-left: 4px;
  width: 100%;
}
.deprecate-row .el-switch {
  vertical-align: middle;
}
.group-label {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}
.address-list {
  padding-left: 0;
}
.address-list-tip {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
  padding: 6px 10px;
  background-color: #f0f9eb;
  border-radius: 4px;
  border: 1px solid #e1f3d8;
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
.address-list-tip .tip-text {
  margin-left: 8px;
  font-size: 12px;
  color: #67c23a;
}
.address-item {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}
.address-item:last-child {
  margin-bottom: 0;
}
.address-label {
  display: inline-block;
  width: 70px;
  text-align: right;
  font-size: 14px;
  color: #606266;
  flex-shrink: 0;
  margin-right: 10px;
}
.is-error /deep/ .el-input__inner {
  border-color: #f56c6c;
}
.required-star {
  color: #f56c6c;
  margin-right: 4px;
}
.add-address-btn {
  margin-top: 8px;
  padding-left: 0;
  color: #409eff;
}
.add-address-btn:hover {
  color: #66b1ff;
}
.add-group-btn {
  width: 100%;
  margin-bottom: 16px;
  color: #409eff;
  border-color: #409eff;
}
.index-mapping-item {
  display: inline-block;
  margin-right: 10px;
  margin-bottom: 10px;
}

/* Fix el-input-number button vertical alignment for small size */
.address-item /deep/ .el-input-number__decrease,
.address-item /deep/ .el-input-number__increase {
  top: 0 !important;
  height: 100% !important;
  line-height: 32px !important;
}

.address-item /deep/ .el-input-number {
  height: 32px !important;
}

.empty-proxies {
  padding: 20px 0;
}
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: 20px;
}
</style>
