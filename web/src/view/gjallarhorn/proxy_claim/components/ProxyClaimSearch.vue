<template>
  <div class="search-term">
    <el-form :inline="true" :model="localSearchInfo" class="demo-form-inline">
      <el-form-item label="OA ID">
        <el-input v-model="localSearchInfo.oa_id" placeholder="请输入OA ID" clearable />
      </el-form-item>
      <el-form-item label="申请人">
        <el-input v-model="localSearchInfo.applicant" placeholder="请输入申请人" clearable />
      </el-form-item>
      <el-form-item label="代理地址">
        <el-input v-model="localSearchInfo.proxy_address" placeholder="按主机地址检索" clearable />
      </el-form-item>
      <el-form-item label="端口">
        <el-input v-model="localSearchInfo.port" placeholder="按端口号检索" clearable />
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="localSearchInfo.status" placeholder="请选择状态" clearable class="status-select">
          <el-option
            v-for="option in statusOptions"
            :key="option.value"
            :label="option.label"
            :value="option.value"
          />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSearch">查询</el-button>
        <el-button @click="handleReset">重置</el-button>
      </el-form-item>
    </el-form>
    <div class="quick-filters">
      <span class="filter-label">快捷筛选：</span>
      <el-button-group>
        <el-button :size="'small'" :type="!localSearchInfo.status ? 'primary' : ''" @click="quickFilter('')">全部</el-button>
        <el-button :size="'small'" :type="localSearchInfo.status === 14 ? 'success' : ''" @click="quickFilter(14)">已就绪</el-button>
        <el-button :size="'small'" :type="localSearchInfo.status === 15 ? 'success' : ''" @click="quickFilter(15)">活跃</el-button>
        <el-button :size="'small'" :type="localSearchInfo.status === 16 ? 'danger' : ''" @click="quickFilter(16)">失败</el-button>
        <el-button :size="'small'" :type="localSearchInfo.status === 1 ? 'info' : ''" @click="quickFilter(1)">已提交</el-button>
      </el-button-group>
    </div>
    <div class="action-buttons">
      <el-button type="primary" icon="el-icon-plus" @click="handleApply">代理申请</el-button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ProxyClaimSearch',
  props: {
    statusOptions: {
      type: Array,
      required: true
    },
    searchInfo: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      localSearchInfo: { ...this.searchInfo }
    }
  },
  watch: {
    searchInfo: {
      handler(val) {
        this.localSearchInfo = { ...val }
      },
      deep: true
    }
  },
  methods: {
    handleSearch() {
      this.$emit('search', { ...this.localSearchInfo })
    },
    handleReset() {
      this.localSearchInfo = {
        oa_id: '',
        applicant: '',
        proxy_address: '',
        port: '',
        status: ''
      }
      this.$emit('reset')
    },
    quickFilter(status) {
      this.localSearchInfo.status = status
      this.$emit('quick-filter', status)
    },
    handleApply() {
      this.$emit('apply')
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
</style>
