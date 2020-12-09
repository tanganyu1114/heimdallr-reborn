<template>
  <div>
    <template v-for="(item,i1) in group">
      <el-divider :key="i1" content-position="left">{{ item.name }}</el-divider>
      <el-row :key="i1" :gutter="20">
        <template v-for="(val,i2) in item.hosts">
          <el-col :key="i2" :span="6">
            <el-card :key="String(i1)+String(i2)" class="card-box">
              <div slot="header">
                <el-tooltip placement="top">
                  <div slot="content">描述: {{ val.descrip }}</div>
                  <span class="card-header">{{ val.name }}</span>
                </el-tooltip>
                <el-tag class="card-right" :type="val.status === false ? 'info' : val.status === true ? 'success' : 'danger'">{{ val.status === false ? '禁用' : val.status === true ? '正常' : '异常' }}</el-tag>
              </div>
              <div>
                <div class="card-line">IP地址: <span class="card-right">{{ val.ipaddr }}</span></div>
                <div class="card-line">操作系统: <span class="card-right">{{ val.agent.system }}</span></div>
                <div class="card-line">系统时间: <span class="card-right">{{ val.agent.time }}</span></div>
                <div class="card-line">CPU使用率(%): <span class="card-right" :style="formatColor(val.agent.cpu)">{{ val.agent.cpu }}</span></div>
                <div class="card-line">内存使用率(%): <span class="card-right" :style="formatColor(val.agent.mem)">{{ val.agent.mem }}</span></div>
                <div class="card-line">磁盘使用率(%): <span class="card-right" :style="formatColor(val.agent.disk)">{{ val.agent.disk }}</span></div>
                <div class="card-line">应用状态:
                  <template v-for="(ng,i3) in val.agent.status_list">
                    <el-tooltip :key="String(i1)+String(i2)+String(i3)" placement="top">
                      <div slot="content">应用名称: {{ ng.name }}<br>{{ ng.version }}</div>
                      <!-- /**  应用状态代码：0 未知;  1  禁用; 2 异常; 3  正常;  **/ -->
                      <el-tag class="card-right" :type="ng.status === 0 ? 'info' : ng.status === 3 ? 'success' : 'danger'">{{ ng.status === 0 ? '未知' : ng.status === 3 ? '正常' : '异常' }}</el-tag>
                    </el-tooltip>
                  </template>
                </div>
              </div>
            </el-card>
          </el-col>
        </template>
      </el-row>
    </template>
  </div>
</template>

<script>
import { getAgentInfo } from '@/api/agent'

export default {
  name: 'Dashboard',
  data() {
    return {
      group: []
    }
  },
  computed: {
    formatColor() {
      return (val) => {
        if (val < 70) {
          return { 'color': '#67C23A' }
        } else if (val < 85) {
          return { 'color': '#E6A23C' }
        } else {
          return { 'color': '#F56C6C' }
        }
      }
    }
  },
  created() {
    this.initCard()
  },
  methods: {
    async initCard() {
      const res = await getAgentInfo()
      console.log(res)
      if (res.code === 0) {
        this.group = res.data
      }
    }
  }
}
</script>

<style scoped>
.card-box {
  width: 320px;
}
.card-header {
  font-weight: bold;
  font-size: 16px;
  color: #409EFF;
}
div.el-card__body {
  padding-top: 10px;
}
.card-line {
  font-style:normal;
  font-size: 14px;
  line-height: 1.7;
  color: #409EFF;
}
.card-right {
  float: right;
  margin-right: 5px;
}
</style>
