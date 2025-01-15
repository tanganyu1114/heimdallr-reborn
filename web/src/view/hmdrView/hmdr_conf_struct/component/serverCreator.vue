<script>

import { _ContextCreator, NewContextGenerator } from './contextCreator'
import ContextCreator from './contextCreator.vue'

export class NewServerGenerator extends NewContextGenerator {
  constructor(data) {
    super(data)
    var formatData = {
      'enabled': true,
      'context-type': 'server',
      'children-context-meta': []
    }
    this.data = undefined
    if (data === undefined || data.listenPort === undefined || data.listenPort === '') {
      return
    }

    formatData['children-context-meta'].push({
      'enabled': true,
      'context-type': 'directive',
      'context-value': 'listen ' + data.listenPort
    })

    if (data.fatherCtxType === 'http') {
      if (data.isSSL) {
        formatData['children-context-meta'].push({
          'enabled': true,
          'context-type': 'directive',
          'context-value': 'ssl on'
        })
      }

      if (data.certCrtPath !== undefined && data.certCrtPath !== '') {
        formatData['children-context-meta'].push({
          'enabled': true,
          'context-type': 'directive',
          'context-value': 'ssl_certificate ' + data.certCrtPath
        })
      }

      if (data.certKeyPath !== undefined && data.certKeyPath !== '') {
        formatData['children-context-meta'].push({
          'enabled': true,
          'context-type': 'directive',
          'context-value': 'ssl_certificate_key ' + data.certKeyPath
        })
      }

      var serverName = 'localhost'
      if (data.serverName !== undefined && data.serverName !== '') {
        serverName = data.serverName
      }
      formatData['children-context-meta'].push({
        'enabled': true,
        'context-type': 'directive',
        'context-value': 'server_name ' + serverName
      })
    }
    this.data = formatData
  }
  generateNewContextOptions() {
    return this.data
  }
}

export class ServerCreator extends _ContextCreator {
  get camelCaseCtxName() { return 'Server' }
  defaultFormData() {
    return {
      fatherCtxType: '',
      isSSL: false,
      certCrtPath: '',
      certKeyPath: '',
      listenPort: 0
    }
  }
  formRules() {
    return {
      listenPort: [
        { required: true, message: '请输入服务侦听的端口', type: 'number', trigger: 'blur' },
        {
          validator: function(rule, value, cb) {
            const intValue = parseInt(value, 10)
            if (isNaN(intValue) || value !== '' && (intValue < 80 || intValue > 65535)) {
              cb(new Error('请输入80~65535之间的整数'))
            } else cb()
          },
          trigger: 'blur'
        }
      ],
      isSSL: [
        { required: false, type: 'boolean', trigger: 'change' }
      ],
      certCrtPath: [
        { required: this.formData.fatherCtxType !== undefined && this.formData.fatherCtxType === 'http' && this.formData.isSSL, message: '请输入正确的证书文件路径', trigger: 'blur' }
      ],
      certKeyPath: [
        { required: this.formData.fatherCtxType !== undefined && this.formData.fatherCtxType === 'http' && this.formData.isSSL, message: '请输入正确的密钥文件路径', trigger: 'blur' }
      ],
      serverName: [
        { required: false, message: '请输入正确格式的虚拟服务名', type: 'string', trigger: 'blur' }
      ]
    }
  }
  initFormDataWithTargetNode(targetNode, position) {
    // TODO: 插入时上下文符合校验
    if (position === 'inner') {
      this.formData.fatherCtxType = targetNode.data.ctxType
    } else {
      this.formData.fatherCtxType = targetNode.parent.data.ctxType
    }
  }
  generateSubmitOptions() {
    return new NewServerGenerator(this.formData).generateNewContextOptions()
  }
}

export default {
  name: 'ServerCreator',
  components: { ContextCreator },
  props: {
    editable: {
      type: Boolean,
      required: true,
      default: false
    }
  },
  data() {
    return {
      creator: new ServerCreator()
    }
  },
  methods: {
    handleDragStart(node, event) {
      this.$emit('node-drag-start', node, event)
    },
    handleDragEnd(draggingNode, endNode, position, event) {
      this.$emit('node-drag-end', draggingNode, endNode, position, event)
    },
    handleFormCommit(newCtxOpts, cb) {
      this.$emit('form-commit-event', newCtxOpts, cb)
    },
    handleDialogBeforeClose() {
      this.$emit('dialog-before-close-event')
    }
  }
}
</script>

<template>
  <div>
    <context-creator
      ref="creator"
      :creator="creator"
      :editable="editable"
      @node-drag-start="handleDragStart"
      @node-drag-end="handleDragEnd"
      @form-commit-event="handleFormCommit"
      @dialog-before-close-event="handleDialogBeforeClose"
    >
      <el-form-item prop="listenPort" label="侦听端口">
        <el-input v-model.number="creator.formData.listenPort" placeholder="请输入数值" />
      </el-form-item>
      <el-form-item v-if="creator.formData.fatherCtxType !== undefined && creator.formData.fatherCtxType === 'http'" prop="isSSL" label="是否启用SSL" label-width="140px">
        <el-checkbox v-model="creator.formData.isSSL" />
      </el-form-item>
      <el-form-item v-if="creator.formData.fatherCtxType !== undefined && creator.formData.fatherCtxType === 'http' && creator.formData.isSSL" prop="certCrtPath" label="证书文件路径" label-width="140px">
        <el-input v-model="creator.formData.certCrtPath" placeholder="请输入内容" />
      </el-form-item>
      <el-form-item v-if="creator.formData.fatherCtxType !== undefined && creator.formData.fatherCtxType === 'http' && creator.formData.isSSL" prop="certKeyPath" label="密钥文件路径" label-width="140px">
        <el-input v-model="creator.formData.certKeyPath" placeholder="请输入内容" />
      </el-form-item>
      <el-form-item v-if="creator.formData.fatherCtxType !== undefined && creator.formData.fatherCtxType === 'http'" prop="serverName" label="虚拟服务名" label-width="140px">
        <el-input v-model="creator.formData.serverName" placeholder="请输入内容（缺省值：localhost）" />
      </el-form-item>
    </context-creator>
  </div>
</template>

<style scoped lang="scss">

</style>
