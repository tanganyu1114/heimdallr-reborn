<template>
  <div>
    <el-card>
      <el-row :gutter="15" class="searchClass">
        <el-form
          ref="elForm"
          :model="formData"
          :rules="rules"
          size="medium"
          label-width="100px"
          label-position="left"
        >
          <el-col :span="12">
            <el-form-item label-width="120px" prop="value" label="应用服务器选择">
              <el-cascader
                v-model="formData.value"
                :options="Options"
                :props="{ expandTrigger: 'hover' }"
                :style="{width: '100%'}"
                placeholder="请选择环境以及主机信息应用服务器选择"
                clearable
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-button size="medium" type="primary" icon="el-icon-search" round @click="getProxySvcInfo">查询</el-button>
          </el-col>
        </el-form>
      </el-row>
    </el-card>
    <el-card>
      <el-table
        :data="proxySvcInfos"
        border
        stripe
        height="600"
        tooltip-effect="dark"
        :default-sort="{prop: 'port'}"
      >
        <el-table-column sortable label="代理类型" prop="proxy-type" width="120" />
        <el-table-column sortable label="服务名" prop="server-name" />
        <el-table-column sortable label="服务侦听端口" prop="port" width="140" />
        <el-table-column label="服务路径" prop="location" />
        <el-table-column sortable label="代理地址" prop="proxy-address" />
      </el-table>
    </el-card>
  </div>
</template>

<script>
import { getOptions } from '@/api/hmdr_conf.js'
import { getProxyServiceInfo } from '@/api/hmdr_statistics.js'

export default {
  name: 'HmdrStatistics',
  data() {
    return {
      code: '',
      formData: {
        value: []
      },
      rules: {
        value: [{
          required: true,
          type: 'array',
          message: '请至少选择一个应用服务器选择',
          trigger: 'change'
        }]
      },
      Options: [],
      proxySvcInfos: []
    }
  },
  created() {
    this.initOptions()
  },
  methods: {
    async initOptions() {
      const res = await getOptions()
      if (res.code === 0) {
        this.Options = res.data
      }
    },
    async getProxySvcInfo() {
      this.$refs['elForm'].validate(async(valid) => {
        if (!valid) return
        const reqOpts = {
          group_id: this.formData.value[0],
          host_id: this.formData.value[1],
          srv_name: this.formData.value[2]
        }
        const res = await getProxyServiceInfo(reqOpts)
        if (res.code === 0) {
          this.proxySvcInfos = res.data
        }
      })
    }
  }
}
</script>

<style scoped>

.searchClass {
  padding-bottom: 0;
}
</style>
