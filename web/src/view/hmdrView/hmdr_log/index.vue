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
            <el-form-item label-width="120px" label="应用服务器选择">
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
            <el-button size="medium" type="primary" icon="el-icon-video-play" round @click="GetLogsInfo">开启日志监听</el-button>
            <el-button size="medium" type="danger" icon="el-icon-video-pause" round @click="close">关闭日志监听</el-button>
          </el-col>
        </el-form>
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
import { getLogsInfo } from '@/api/hmdr_websocket.js'

export default {
  name: 'HmdrLog',
  data() {
    return {
      key: '',
      logs: '',
      path: 'ws://127.0.0.1:8888/hmdrWebSocket/ws',
      socket: '',
      formData: {
        value: []
      },
      rules: {
        hmdr_groupOptions: [{
          required: true,
          type: 'array',
          message: '请至少选择一个应用服务器选择',
          trigger: 'change'
        }]
      },
      Options: []
    }
  },
  created() {
    this.initOptions()
  },
  mounted() {
    this.initWebSocket()
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
    async GetLogsInfo() {
      const res = await getLogsInfo()
      if (res.code === 0) {
        if (res.data.status === 'ok') {
          this.key = res.data.key
        }
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
      }
    },
    open: function() {
      this.logs = 'socket连接成功'
      console.log('socket连接成功')
    },
    error: function() {
      console.log('连接错误')
    },
    getMessage: function(msg) {
      this.logs += '\n' + msg.data
      console.log(msg.data)
    },
    send: function() {
      const sf = {
        group_id: this.formData.value[0],
        host_id: this.formData.value[1],
        srv_name: this.formData.value[2]
      }
      const jsonsf = JSON.stringify(sf)
      this.logs += '\n' + '发送认证信息'
      this.socket.send(sf)
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
</style>
