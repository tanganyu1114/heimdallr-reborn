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
                ref="cascader"
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
            <el-button size="medium" type="success" icon="el-icon-download" round @click="exportToExcel">导出Excel</el-button>
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
        :default-sort="{prop: 'server-port'}"
      >
        <el-table-column sortable label="服务描述" prop="proxy-service-comment" />
        <el-table-column sortable label="服务名" prop="server-name" />
        <el-table-column sortable label="服务侦听端口" prop="server-port" width="140" />
        <el-table-column sortable label="反向代理协议" prop="proxy-protocol" width="140" />
        <el-table-column label="反向代理路由路径" prop="location" />
        <el-table-column label="特殊判断条件" prop="if-condition" />
        <el-table-column sortable label="反向代理原始地址" prop="proxy-original-url" />
        <el-table-column label="反向代理地址集">
          <template slot-scope="props">
            <div v-for="address in props.row['proxy-address']" :key="address['domain-name'] + ':' + address['port']" class="tag-group">
              <span class="tag-group__title">{{ parseUpstreamServerAddresses(address) }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="反向代理网络状态">
          <template slot-scope="props">
            <div v-for="address in props.row['proxy-address']" :key="address['domain-name'] + ':' + address['port']" class="tag-group">
              <el-tooltip
                v-for="item in parseSocketBrief(address, props.row['proxy-protocol'])"
                :key="item.label"
                :content="item.state"
                placement="top"
                size="medium"
                effect="dark"
              >
                <el-tag
                  :type="item.type"
                  size="medium"
                  effect="dark"
                >
                  {{ item.label }}
                </el-tag>
              </el-tooltip>
            </div>
            <div><el-button
              icon="el-icon-refresh"
              @click="connectivityCheck(props.row['context-pos'])"
            >
              代理连通性检查
            </el-button></div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script>
import { getOptions } from '@/api/hmdr_conf.js'
import { getProxyServiceInfo, connectivityCheckOfProxyService, exportProxyServiceInfoToExcel } from '@/api/hmdr_statistics.js'

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
    },
    parseConnectivityLabel(state) {
      switch (state) {
        case 0: return '网络连通性未知'
        case 1: return '网络可达'
        case 2: return '网络不可达'
      }
    },
    parseConnectivityType(state) {
      switch (state) {
        case 0: return 'info'
        case 1: return 'success'
        case 2: return 'danger'
      }
    },
    parseSocketBrief(address, protocol) {
      const ret = []
      const isStream = (protocol === 'TCP/UDP')
      for (let i = 0; i < address.sockets.length; i++) {
        const socket = address.sockets[i]
        if (isStream) {
          ret.push({
            'label': socket['ipv4'] + ':' + socket.port + '(TCP)',
            'state': this.parseConnectivityLabel(socket['tcp-connectivity']),
            'type': this.parseConnectivityType(socket['tcp-connectivity'])
          })
          ret.push({
            'label': socket['ipv4'] + ':' + socket.port + '(UDP)',
            'state': this.parseConnectivityLabel(socket['udp-connectivity']),
            'type': this.parseConnectivityType(socket['udp-connectivity'])
          })
        } else {
          ret.push({
            'label': socket['ipv4'] + ':' + socket.port,
            'state': this.parseConnectivityLabel(socket['tcp-connectivity']),
            'type': this.parseConnectivityType(socket['tcp-connectivity'])
          })
        }
      }
      return ret
    },
    connectivityCheck(ctxPos) {
      this.$refs['elForm'].validate(async(valid) => {
        if (!valid) return
        const reqOpts = {
          group_id: this.formData.value[0],
          host_id: this.formData.value[1],
          srv_name: this.formData.value[2],
          config: ctxPos.config,
          'context-pos-path': ctxPos['context-pos-path']
        }
        const res = await connectivityCheckOfProxyService(reqOpts)
        if (res.code === 0) {
          await this.getProxySvcInfo()
        }
      })
    },
    parseUpstreamServerAddresses(address) {
      const addr = address['domain-name'] + ':' + address.port
      let upsrv = ''
      for (const socket of address.sockets) {
        upsrv += socket['ipv4'] + ':' + socket.port + ', '
      }
      if (upsrv !== '') {
        return addr + '(' + upsrv.slice(0, -2) + ')'
      }
      return addr
    },
    // 导出 Excel 文件
    async exportToExcel() {
      this.$refs['elForm'].validate(async(valid) => {
        if (!valid) {
          this.$message({
            type: 'warning',
            message: '请先选择应用服务器'
          })
          return
        }

        const reqOpts = {
          group_id: this.formData.value[0],
          host_id: this.formData.value[1],
          srv_name: this.formData.value[2]
        }

        try {
          const res = await exportProxyServiceInfoToExcel(reqOpts)

          // 如果 res 是 response 对象（拦截器返回的），取 data
          const data = res.data || res

          // 尝试将 arraybuffer 转换为文本，判断是否为 JSON 错误响应
          if (data instanceof ArrayBuffer) {
            const decoder = new TextDecoder('utf-8')
            const text = decoder.decode(data)

            try {
              // 尝试解析为 JSON
              const jsonData = JSON.parse(text)
              // 如果解析成功且有 code 字段，说明是错误响应
              if (jsonData.code !== undefined && jsonData.code !== 0) {
                this.$message({
                  type: 'error',
                  message: jsonData.msg || '导出失败'
                })
                return
              }
            } catch (e) {
              // 解析失败，说明是真正的 Excel 二进制数据
            }

            // 使用级联选择器的 getCheckedNodes() 方法获取显示名称（O(1) 复杂度）
            let groupName = 'unknown'
            let hostName = 'unknown'

            const nodes = this.$refs.cascader.getCheckedNodes()
            if (nodes.length > 0) {
              const checkedNode = nodes[0]
              // pathNodes 包含从根到叶子的所有节点
              if (checkedNode.pathNodes && checkedNode.pathNodes.length > 0) {
                groupName = checkedNode.pathNodes[0].label || 'unknown'
              }
              if (checkedNode.pathNodes && checkedNode.pathNodes.length > 1) {
                hostName = checkedNode.pathNodes[1].label || 'unknown'
              }
            }

            const serviceName = this.formData.value.length > 2 ? this.formData.value[2] : 'unknown'
            const fileName = `proxy_service_${groupName}_${hostName}_${serviceName}_${new Date().getTime()}.xlsx`

            // 这是真正的 Excel 文件，创建 Blob 并下载
            const blob = new Blob([data], {
              type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
            })
            this.downloadFile(blob, fileName)
            this.$message({
              type: 'success',
              message: '导出成功'
            })
          } else {
            this.$message({
              type: 'error',
              message: '导出失败：响应格式异常'
            })
          }
        } catch (error) {
          console.error('导出失败:', error)
          this.$message({
            type: 'error',
            message: '导出失败：' + (error.message || '未知错误')
          })
        }
      })
    },
    // 下载文件的辅助方法
    downloadFile(blob, filename) {
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = filename
      document.body.appendChild(link)
      link.click()
      // 延迟释放 URL 对象，确保下载完成
      setTimeout(() => {
        document.body.removeChild(link)
        window.URL.revokeObjectURL(url)
      }, 1000)
    }
  }
}
</script>

<style scoped>

.searchClass {
  padding-bottom: 0;
}
</style>
