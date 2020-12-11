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
        <el-col :span="8">
          <el-button size="medium" type="primary" icon="el-icon-video-play" :disabled="!isActive" round @click="handleStartLogOn">开启日志监听</el-button>
          <el-button size="medium" type="danger" icon="el-icon-video-pause" :disabled="isActive" round @click="handleStartLogOff">关闭日志监听</el-button>
        </el-col>
      </el-row>
    </el-card>
    <!--/**  显示日志的窗口  **/ -->
    <el-card>
      <el-input
        v-model="logs"
        type="textarea"
        :rows="26"
        :readonly="true"
        resize="none"
      />
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
      logs: '',
      path: process.env.VUE_APP_WS,
      socket: '',
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
    this.socket.onclose = this.close
  },
  methods: {
    async initOptions() {
      const res = await getOptions()
      if (res.code === 0) {
        this.Options = res.data
        console.log(res.data)
      }
    },
    initWebSocket() {
      if (typeof (WebSocket) === 'undefined') {
        alert('您的浏览器不支持socket')
      } else {
        // 实例化socket
        this.socket = new WebSocket(this.path)
        // 监听socket连接
        this.socket.onopen = this.open
        // 监听socket错误信息
        this.socket.onerror = this.error
        // 监听socket消息
        this.socket.onmessage = this.getMessage
        // 关闭socket时发生消息
        this.socket.onclose = this.close
      }
    },
    handleStartLogOn() {
      this.$refs['elForm'].validate(valid => {
        if (!valid) return
        // TODO 提交表单
        this.isActive = !this.isActive
        this.initWebSocket()
        // this.send()
      })
    },
    handleStartLogOff() {
      this.isActive = !this.isActive
      this.logs += '\n' + '关闭socket连接'
      this.sendOff()
      if (this.socket.bufferedAmount === 0) {
        this.socket.close()
      }
    },
    open: function() {
      this.logs = 'socket连接成功'
      console.log('socket连接成功')
      this.sendOn()
    },
    error: function() {
      console.log('连接错误')
    },
    getMessage: function(msg) {
      this.logs += '\n' + msg.data
      console.log(msg.data)
    },
    sendOn: function() {
      console.log(this.formData.logName)
      const sf = {
        group_id: this.formData.value[0],
        host_id: this.formData.value[1],
        srv_name: this.formData.value[2],
        log_name: this.formData.logName,
        status: true
      }
      const jsonsf = JSON.stringify(sf)
      this.logs += '\n' + '发送认证信息'
      // this.socket.send(sf)
      this.socket.send(jsonsf)
    },
    sendOff() {
      console.log(this.formData.logName)
      const sf = {
        group_id: this.formData.value[0],
        host_id: this.formData.value[1],
        srv_name: this.formData.value[2],
        log_name: this.formData.logName,
        status: false
      }
      const jsonsf = JSON.stringify(sf)
      this.logs += '\n' + '发送认证信息'
      // this.socket.send(sf)
      this.socket.send(jsonsf)
    },
    close: function() {
      console.log('socket已经关闭')
    }
  }
}

</script>

<style scoped>
  /*.hljscss {*/
  /*  height: 70%;*/
  /*  position: relative;*/
  /*  font-size: 18px;*/
  /*  overflow-y: scroll;*/
  /*}*/
  .app {
    height: 100%;
    overflow: hidden;
  }
  .el-scrollbar__wrap {
    overflow: visible;
    overflow-x: hidden;
  }
  .btnClass {
    padding-top: 0;
  }
</style>
