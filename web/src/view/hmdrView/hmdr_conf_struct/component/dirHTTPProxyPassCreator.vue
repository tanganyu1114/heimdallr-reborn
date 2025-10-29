<script>

import { _ContextCreator, NewContextGenerator } from './contextCreator'
import ContextCreator from './contextCreator.vue'

export class NewDirHTTPProxyPassGenerator extends NewContextGenerator {
  constructor(data) {
    super(data)
    var formatData = {
      'enabled': true,
      'context-type': 'dir_http_proxy_pass',
      'context-value': ''
    }
    this.data = undefined
    if (data === undefined || data.proxyPassValue === undefined || data.proxyPassValue === '') return
    formatData['context-value'] += data.proxyPassValue
    this.data = formatData
  }
  generateNewContextOptions() {
    return this.data
  }
}

export class DirHTTPProxyPassCreator extends _ContextCreator {
  get camelCaseCtxName() { return 'ProxyPass(HTTP)' }
  className() {
    return 'DirHTTPProxyPassCreator'
  }
  defaultFormData() {
    return {
      proxyPassValue: ''
    }
  }
  formRules() {
    return {
      proxyPassValue: [
        { required: true, message: '请输入反向代理URL地址', type: 'string', trigger: 'blur' }
      ]
    }
  }
  initFormDataWithTargetNode(targetNode, position) {
    // TODO: 插入时上下文符合校验
  }
  generateSubmitOptions() {
    return new NewDirHTTPProxyPassGenerator(this.formData).generateNewContextOptions()
  }
}

export default {
  name: 'DirHTTPProxyPassCreator',
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
      creator: new DirHTTPProxyPassCreator()
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
      <el-form-item prop="proxyPassValue" label="反向代理URL地址" label-width="140px">
        <el-input v-model="creator.formData.proxyPassValue" placeholder="请输入内容" />
      </el-form-item>
    </context-creator>
  </div>
</template>

<style scoped lang="scss">

</style>
