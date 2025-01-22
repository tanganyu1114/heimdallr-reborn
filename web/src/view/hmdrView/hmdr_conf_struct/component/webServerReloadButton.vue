<script>

import { binCMDExec } from '@/api/hmdr_bin_cmd.js'

export default {
  name: 'WebServerReloadButton',
  props: {
    size: {
      type: String,
      required: false,
      default: 'mini'
    },
    type: {
      type: String,
      required: false,
      default: 'primary'
    },
    icon: {
      type: String,
      required: false,
      default: 'el-icon-refresh'
    },
    disabled: {
      type: Boolean,
      required: false,
      default: false
    },
    serverOptions: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      loading: false,
      dialogVisible: false
    }
  },
  methods: {
    cancel() {
      this.dialogVisible = false
      this.loading = false
    },
    async submit() {
      this.loading = true
      var request = {
        'web-server-options': {
          group_id: this.serverOptions.group_id,
          host_id: this.serverOptions.host_id,
          srv_name: this.serverOptions.srv_name
        },
        args: [
          '-s',
          'reload'
        ]
      }
      var resp = await binCMDExec(request)
      if (resp.code !== 0) {
        this.loading = false
        return
      }
      var msgType = 'error'
      var msgTitle = '服务器重载失败'
      var msgOutput = resp.data.stderr
      var msgColor = 'color: #F56C6C'
      if (resp.data.successful === true) {
        msgType = 'success'
        msgTitle = '服务器重载成功'
        msgOutput = resp.data.stdout
        msgColor = 'color: #67C23A'
      }
      const h = this.$createElement
      this.$msgbox({
        title: msgTitle,
        type: msgType,
        message: h('p', null, [
          h('span', null, '命令执行回执：'),
          h('br'),
          h('i', { style: msgColor }, msgOutput)
        ])
      })
      this.loading = false
      this.dialogVisible = false
    }
  }
}
</script>

<template>
  <div>
    <el-button
      :size="size"
      :type="type"
      :icon="icon"
      :disabled="disabled"
      @click="()=>{dialogVisible = true}"
    />
    <el-dialog
      title="服务端配置重载确认"
      width="35%"
      :visible.sync="dialogVisible"
      :before-close="cancel"
    >
      <h1>请确认是否要重载“{{ serverOptions.srv_name }}”服务端配置</h1>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" :loading.sync="loading" @click="submit">确认重载</el-button>
        <el-button @click="cancel">取消</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">

</style>
