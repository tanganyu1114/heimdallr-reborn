<template>
  <el-dialog
    :visible.sync="dialogVisible"
    :title="applyType === 'create' ? '提交代理申请' : '更新代理申请'"
    width="800px"
    @closed="handleClosed"
  >
    <el-form ref="applyForm" :model="formData" :rules="rules" label-width="140px" label-position="right">
      <el-form-item v-if="applyType === 'update'" label="ID" prop="id">
        <el-input v-model="formData.id" disabled />
      </el-form-item>
      <el-form-item label="OA ID" prop="oa_id">
        <el-input v-model="formData.oa_id" :disabled="applyType === 'update'" placeholder="请输入唯一OA标识（字母、数字、下划线、破折号）" />
      </el-form-item>
      <el-form-item v-if="applyType === 'create'" label="申请人" prop="applicant">
        <el-input v-model="formData.applicant" placeholder="请输入申请人姓名" />
      </el-form-item>
      <el-form-item label="源项目描述" prop="source_project_desc">
        <el-input v-model="formData.source_project_desc" type="textarea" :rows="2" placeholder="请输入源项目描述" />
      </el-form-item>
      <el-form-item label="目标项目描述" prop="target_project_desc">
        <el-input v-model="formData.target_project_desc" type="textarea" :rows="2" placeholder="请输入目标项目描述" />
      </el-form-item>

      <el-divider content-position="left">代理转发地址配置</el-divider>
      <div class="config-tip">
        <el-tag size="mini" type="info" effect="plain">说明</el-tag>
        <span class="tip-text">每个代理服务对应一个对端接口服务，可配置多个后端地址实现负载均衡</span>
      </div>

      <div v-if="applyType === 'create'" class="import-section">
        <el-button size="small" type="primary" icon="el-icon-download" @click="showImportDialog">导入已有代理地址</el-button>
        <span class="import-tip">从现有OA中导入代理地址配置，避免重复填写</span>
      </div>

      <AddressGroupEditor
        :address-groups="formData.proxy_pass_addresses"
        :show-deprecate-switch="applyType === 'update'"
        @add-group="handleAddGroup"
        @remove-group="handleRemoveGroup"
        @add-address="handleAddAddress"
        @remove-address="handleRemoveAddress"
      />

      <el-form-item label="环境名" prop="claim_env_ids">
        <el-select v-model="formData.claim_env_ids" multiple placeholder="选择环境名" style="width: 100%;">
          <el-option v-for="env in envOptions" :key="env.value" :label="env.label" :value="env.value" />
        </el-select>
      </el-form-item>
    </el-form>

    <div slot="footer" class="dialog-footer">
      <el-button @click="dialogVisible = false">取 消</el-button>
      <el-button type="primary" @click="handleSubmit">确 定</el-button>
    </div>

    <ImportProxyAddressesDialog
      :visible.sync="importDialogVisible"
      :initial-selection="importDialogInitialSelection"
      @import="handleImportData"
    />
  </el-dialog>
</template>

<script>
import AddressGroupEditor from './AddressGroupEditor.vue'
import ImportProxyAddressesDialog from './ImportProxyAddressesDialog.vue'
import proxyClaimStateMixin from '../mixins/proxyClaimState'
import { CLAIM_STATES, ENV_OPTIONS } from '../constants/proxyClaim'
import { getProxyClaim, applyProxyClaim, migrateProxyClaim } from '@/api/gjallarhorn_proxy_claim'

export default {
  name: 'ProxyClaimApplyDialog',
  components: {
    AddressGroupEditor,
    ImportProxyAddressesDialog
  },
  mixins: [proxyClaimStateMixin],
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    type: {
      type: String,
      default: 'create',
      validator: (val) => ['create', 'update'].includes(val)
    },
    oaId: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      applyType: 'create',
      formData: {
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
        claim_state: CLAIM_STATES.ClaimStateSubmitted
      },
      envOptions: ENV_OPTIONS,
      importDialogVisible: false,
      importDialogInitialSelection: [],
      rules: {
        oa_id: [
          { required: true, message: '请输入OA ID', trigger: 'blur' },
          { pattern: /^[a-zA-Z0-9_-]+$/, message: 'OA ID仅允许字母、数字、下划线和破折号', trigger: 'blur' }
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
  watch: {
    type: {
      handler(val) {
        this.applyType = val
      },
      immediate: true
    },
    visible: {
      handler(val) {
        if (val && this.applyType === 'update' && this.oaId) {
          this.loadClaimData()
        }
      }
    }
  },
  methods: {
    getInitialFormData() {
      return {
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
        claim_state: CLAIM_STATES.ClaimStateSubmitted
      }
    },
    async loadClaimData() {
      try {
        const res = await getProxyClaim(this.oaId)
        const proxyPassAddresses = res.proxy_pass_addresses || {}
        const addressesArray = Object.values(proxyPassAddresses).map(item => ({
          ...item,
          _isOriginal: true
        }))
        this.formData = {
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
      } catch (error) {
        this.$message({
          type: 'error',
          message: '获取详情失败'
        })
      }
    },
    handleAddGroup() {
      this.formData.proxy_pass_addresses.push({
        protocol: 'http',
        addresses: [
          { host: '', port: 8080 }
        ],
        is_deprecated: false,
        _isOriginal: false
      })
    },
    handleRemoveGroup(index) {
      this.formData.proxy_pass_addresses.splice(index, 1)
    },
    handleAddAddress(groupIndex) {
      this.formData.proxy_pass_addresses[groupIndex].addresses.push({
        host: '',
        port: 8080
      })
    },
    handleRemoveAddress(groupIndex, addressIndex) {
      this.formData.proxy_pass_addresses[groupIndex].addresses.splice(addressIndex, 1)
    },
    async handleSubmit() {
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

        const submitAddressesMap = {}
        let hasMigratedAddresses = false
        for (let index = 0; index < this.formData.proxy_pass_addresses.length; index++) {
          const item = this.formData.proxy_pass_addresses[index]
          if (item._sourceOaId) {
            hasMigratedAddresses = true
          } else {
            const { _isOriginal, ...rest } = item
            submitAddressesMap[index] = rest
          }
        }

        const nextState = this.applyType === 'update'
          ? this.getNextState(this.formData.claim_state)
          : this.formData.claim_state

        const submitData = {
          ...this.formData,
          claim_state: nextState,
          proxy_pass_addresses: submitAddressesMap
        }

        if (this.applyType === 'create') {
          delete submitData.id
        } else {
          submitData.id = Number(submitData.id) || 0
        }

        try {
          await applyProxyClaim(submitData)

          if (this.applyType === 'create' && hasMigratedAddresses) {
            await this.submitMigrateRequests(this.formData.proxy_pass_addresses, this.formData.oa_id)
          }

          this.$message({
            type: 'success',
            message: this.applyType === 'create' ? '申请提交成功' : '更新成功'
          })
          this.dialogVisible = false
          this.$emit('success')
        } catch (error) {
          const errorMsg = error?.message || error?.msg || '未知错误'
          this.$message({
            type: 'error',
            message: `${this.applyType === 'create' ? '申请提交' : '更新'}失败：${errorMsg}，请检查输入参数或联系系统管理员`
          })
        }
      })
    },
    simpleHash(str) {
      let h = 5381
      for (let i = 0; i < str.length; i++) {
        h = ((h << 5) + h + str.charCodeAt(i)) & 0x7fffffff
      }
      return h.toString(16)
    },
    generateMigrationName(sourceOaId, targetOaId, indexesMap) {
      const sortedEntries = Object.entries(indexesMap)
        .sort((a, b) => Number(a[0]) - Number(b[0]))
      const hashInput = sourceOaId + '|' + targetOaId + '|' +
        sortedEntries.map(([k, v]) => `${k}:${v}`).join(',')
      const hash = this.simpleHash(hashInput)
      return `${sourceOaId}-to-${targetOaId}-${hash}`
    },
    async submitMigrateRequests(proxyPassAddresses, targetOaId) {
      const sourceGroups = new Map()

      for (let targetIndex = 0; targetIndex < proxyPassAddresses.length; targetIndex++) {
        const addr = proxyPassAddresses[targetIndex]
        const sourceOaId = addr._sourceOaId
        if (!sourceOaId) continue
        const sourceProxyIndex = addr._sourceProxyIndex

        if (!sourceGroups.has(sourceOaId)) {
          sourceGroups.set(sourceOaId, {
            migratingAddresses: {},
            indexesMap: {},
            sourceClaimState: addr._sourceClaimState || 0
          })
        }

        const group = sourceGroups.get(sourceOaId)

        if (!group.indexesMap[sourceProxyIndex]) {
          group.indexesMap[sourceProxyIndex] = targetIndex
        }

        const { _sourceOaId, _sourceProxyIndex, _sourceClaimState, ...cleanAddr } = addr
        group.migratingAddresses[targetIndex] = cleanAddr
      }

      const migratePromises = []

      for (const [sourceOaId, group] of sourceGroups) {
        const name = this.generateMigrationName(
          sourceOaId,
          targetOaId,
          group.indexesMap
        )

        const migrateData = {
          name: name,
          source_oa_id: sourceOaId,
          target_oa_id: targetOaId,
          source_claim_state: group.sourceClaimState,
          target_claim_state: CLAIM_STATES.ClaimStateSubmitted,
          migrating_addresses: group.migratingAddresses,
          migrating_indexes_map: group.indexesMap,
          migrate_state: 1
        }

        migratePromises.push(
          migrateProxyClaim(migrateData).catch(err => {
            console.error(`Migration failed for source OA ${sourceOaId}:`, err)
            return { sourceOaId, error: err }
          })
        )
      }

      const results = await Promise.all(migratePromises)
      const failedMigrations = results.filter(r => r && r.error)

      if (failedMigrations.length > 0) {
        const failedOaIds = failedMigrations.map(r => r.sourceOaId).join(', ')
        this.$message({
          type: 'warning',
          message: `代理申请已提交，但以下源OA的迁移请求失败：${failedOaIds}，请稍后重试或联系系统管理员`,
          duration: 5000
        })
      } else if (sourceGroups.size > 0) {
        this.$message({
          type: 'success',
          message: `已成功提交 ${sourceGroups.size} 个源OA的代理迁移请求`,
          duration: 3000
        })
      }
    },
    getNextState(currentState) {
      const stateTransitions = {
        [CLAIM_STATES.ClaimStateProvisioned]: CLAIM_STATES.ClaimStateInstantiating,
        [CLAIM_STATES.ClaimStateActive]: CLAIM_STATES.ClaimStateEnabling,
        [CLAIM_STATES.ClaimStateDeprecated]: CLAIM_STATES.ClaimStateDeprecating,
        [CLAIM_STATES.ClaimStateFailed]: CLAIM_STATES.ClaimStateSubmitted,
        [CLAIM_STATES.ClaimStateExpired]: CLAIM_STATES.ClaimStateSubmitted
      }
      return stateTransitions[currentState] || CLAIM_STATES.ClaimStateSubmitted
    },
    handleClosed() {
      this.resetForm()
    },
    resetForm() {
      this.formData = this.getInitialFormData()
      this.$nextTick(() => {
        if (this.$refs.applyForm) {
          this.$refs.applyForm.resetFields()
        }
      })
    },
    showImportDialog() {
      const seen = new Set()
      const initialSelection = []
      for (const item of this.formData.proxy_pass_addresses) {
        if (item._sourceOaId && item._sourceProxyIndex !== undefined) {
          const key = item._sourceOaId + ':' + item._sourceProxyIndex
          if (!seen.has(key)) {
            seen.add(key)
            initialSelection.push({
              oa_id: item._sourceOaId,
              proxy_index: item._sourceProxyIndex
            })
          }
        }
      }
      this.importDialogInitialSelection = initialSelection
      this.importDialogVisible = true
    },
    handleImportData(addressesArray) {
      if (!addressesArray || addressesArray.length === 0) return

      const existingImportKeys = new Set()
      for (const item of this.formData.proxy_pass_addresses) {
        if (item._sourceOaId && item._sourceProxyIndex !== undefined) {
          existingImportKeys.add(item._sourceOaId + ':' + item._sourceProxyIndex)
        }
      }

      const dedupedNew = addressesArray.filter(item => {
        const key = item._sourceOaId + ':' + item._sourceProxyIndex
        if (existingImportKeys.has(key)) return false
        existingImportKeys.add(key)
        return true
      })

      if (dedupedNew.length === 0) {
        this.$message({
          type: 'info',
          message: '所选代理服务均已存在，无需重复导入',
          duration: 2000
        })
        return
      }

      this.formData.proxy_pass_addresses = [
        ...this.formData.proxy_pass_addresses,
        ...dedupedNew
      ]
    }
  }
}
</script>

<style scoped>
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
.import-section {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  padding: 12px;
  background-color: #ecf5ff;
  border-radius: 4px;
  border: 1px solid #d9ecff;
}
.import-tip {
  margin-left: 12px;
  font-size: 13px;
  color: #409eff;
}
</style>
