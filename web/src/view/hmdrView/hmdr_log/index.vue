<template>
  <div>
    <el-card>
      <el-row :gutter="15">
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
          <el-col :span="4">
            <el-form-item label-width="80px" label="日志名字" prop="logName">
              <el-select v-model="formData.logName" placeholder="请选择日志名字" :style="{width: '100%'}">
                <el-option
                  v-for="(item, index) in logNameOptions"
                  :key="index"
                  :label="item.label"
                  :value="item.value"
                  :disabled="item.disabled"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-form>
      </el-row>
      <el-row class="btnClass">
        <el-col :span="24">
          <el-button size="medium" type="primary" icon="el-icon-video-play" :disabled="!isActive" round @click="handleStartLogOn">开启日志监听</el-button>
          <el-button size="medium" type="danger" icon="el-icon-video-pause" :disabled="isActive" round @click="handleStartLogOff">关闭日志监听</el-button>
          <el-button size="medium" type="warning" icon="el-icon-delete" round style="float:right;" @click="handleCleanLog">清空日志</el-button>
        </el-col>
      </el-row>
    </el-card>
    <!--/**  显示日志的窗口  **/ -->
    <!-- TODO: 显示框切换为ul li span 遍历data生成 -->
    <el-card>
      <ul id="log-box" class="logBox">
        <li v-for="(item,index) in logs" :key="index"><span class="logLine">{{ item }}</span></li>
      </ul>
    </el-card>
  </div>
</template>

<script>
import { getOptions } from '@/api/hmdr_conf.js'

export default {
  name: 'HmdrLog',
  data() {
    return {
      key: '',
      logs: [],
      path: process.env.VUE_APP_WS,
      socket: '',
      healthId: Number(),
      isActive: true,
      formData: {
        logName: 'access.log',
        value: []
      },
      rules: {
        logName: [{
          required: true,
          message: '请选择日志名字',
          trigger: 'change'
        }],
        value: [{
          required: true,
          type: 'array',
          message: '请至少选择一个应用服务器选择',
          trigger: 'change'
        }]
      },
      logNameOptions: [{
        'label': 'access.log',
        'value': 'access.log'
      }, {
        'label': 'error.log',
        'value': 'error.log'
      }],
      Options: []
    }
  },
  created() {
    this.initOptions()
  },
  mounted() {
    // this.initWebSocket()
  },
  destroyed() {
    // 销毁监听
    if (this.socket.readyState === 1) {
      this.socket.close()
      clearInterval(this.healthId)
    }
  },
  methods: {
    async initOptions() {
      const res = await getOptions()
      if (res.code === 0) {
        this.Options = res.data
      }
    },
    healthCheck() {
      this.healthId = setInterval(() => {
        this.socket.send('ping')
      }, 20000)
    },
    initWebSocket() {
      if (typeof (WebSocket) === 'undefined') {
        alert('您的浏览器不支持socket')
      } else {
        // 实例化socket
        this.socket = new WebSocket(this.path)
        // 监听socket连接
        this.socket.onopen = this.onOpen
        // 监听socket错误信息
        this.socket.onerror = this.onError
        // 监听socket消息
        this.socket.onmessage = this.onMessage
        // 关闭socket时发生消息
        this.socket.onclose = this.onClose
      }
    },
    handleStartLogOn() {
      this.$refs['elForm'].validate(valid => {
        if (!valid) return
        // 按钮取反
        this.isActive = !this.isActive
        // 初始化websocket
        this.initWebSocket()
      })
    },
    handleStartLogOff() {
      // 按钮取反
      this.isActive = !this.isActive
      // 关闭websocket
      this.socket.close()
    },
    handleCleanLog() {
      this.logs = []
    },
    onOpen: function() {
      this.addLogs('socket连接成功')
      this.sendOnSignal()
      // 心跳检查
      // this.healthCheck()
    },
    onError: function(e) {
      this.addLogs('socket error :' + e)
    },
    onMessage(e) {
      this.addLogs(e.data)
      this.freshSrollBar()
    },
    addLogs(msg) {
      const arr = msg.split('\n')
      // 去除空元素
      const r = arr.filter(s => {
        return s && s.trim()
      })
      // 添加到日志信息
      r.forEach(item => {
        this.logs.push(item)
      })
      // this.logs = this.logs.concat(arr)
    },
    onClose(e) {
      this.addLogs('websocket已关闭')
      this.addLogs('websocket 断开: ' + e.code + ' ' + e.reason + ' ' + e.wasClean)
      clearInterval(this.healthId)
    },
    sendOnSignal: function() {
      const sf = {
        group_id: this.formData.value[0],
        host_id: this.formData.value[1],
        srv_name: this.formData.value[2],
        log_name: this.formData.logName
      }
      const jsonsf = JSON.stringify(sf)
      this.addLogs('开始获取日志信息')
      this.socket.send(jsonsf)
    },
    freshSrollBar() {
      this.$nextTick(() => {
        setTimeout(() => {
          const textarea = document.getElementById('log-box')
          textarea.scrollTop = textarea.scrollHeight
        }, 13)
      })
    }
  }
}

</script>

<style scoped>
  .btnClass {
    padding-top: 0;
  }
  .logBox {
    height: 500px;
    background: #E4E7ED;
    border: #DCDFE6 2px solid;
    overflow-y: scroll;
  }
  .logLine {
    margin-left: 5px;
    margin-right: 5px;
    font-size: 16px;
    color: #157a0c;
  }
</style>
