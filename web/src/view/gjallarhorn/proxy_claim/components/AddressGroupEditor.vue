<template>
  <div class="address-group-editor">
    <div v-for="(addrGroup, groupIndex) in addressGroups" :key="groupIndex" class="address-group">
      <div class="address-group-header">
        <div class="header-left">
          <span class="group-label">{{ groupLabel }} {{ groupIndex + 1 }}</span>
          <el-select v-model="addrGroup.protocol" placeholder="选择协议" size="small" style="width: 140px; margin-left: 12px;">
            <el-option label="HTTP" value="http" />
            <el-option label="HTTPS" value="https" />
            <el-option label="TCP/UDP" value="tcp_or_udp" />
          </el-select>
          <el-tag v-if="addrGroup._sourceOaId" size="mini" type="info" effect="plain" style="margin-left: 8px;">
            迁移自 {{ addrGroup._sourceOaId }}:{{ addrGroup._sourceProxyIndex }}
          </el-tag>
        </div>
        <el-button
          v-if="canDeleteGroup && addressGroups.length > 1 && !addrGroup._isOriginal"
          size="small"
          type="danger"
          plain
          icon="el-icon-delete"
          @click="handleRemoveGroup(groupIndex)"
        >删除代理服务</el-button>
      </div>

      <div v-if="showDeprecateSwitch && addrGroup._isOriginal" class="deprecate-row">
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
          <el-button
            v-if="addrGroup.addresses.length > 1"
            size="small"
            type="danger"
            plain
            icon="el-icon-delete"
            style="margin-left: 10px;"
            @click="handleRemoveAddress(groupIndex, addrIndex)"
          >删除</el-button>
        </div>
        <el-button class="add-address-btn" size="small" type="text" icon="el-icon-plus" @click="handleAddAddress(groupIndex)">添加地址</el-button>
      </div>
    </div>

    <el-button class="add-group-btn" type="dashed" icon="el-icon-plus" @click="handleAddGroup">添加代理服务</el-button>
  </div>
</template>

<script>
export default {
  name: 'AddressGroupEditor',
  props: {
    addressGroups: {
      type: Array,
      required: true
    },
    groupLabel: {
      type: String,
      default: '代理服务'
    },
    canDeleteGroup: {
      type: Boolean,
      default: true
    },
    showDeprecateSwitch: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    handleAddGroup() {
      this.$emit('add-group')
    },
    handleRemoveGroup(index) {
      this.$emit('remove-group', index)
    },
    handleAddAddress(groupIndex) {
      this.$emit('add-address', groupIndex)
    },
    handleRemoveAddress(groupIndex, addressIndex) {
      this.$emit('remove-address', groupIndex, addressIndex)
    }
  }
}
</script>

<style scoped>
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
.group-label {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
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
</style>
