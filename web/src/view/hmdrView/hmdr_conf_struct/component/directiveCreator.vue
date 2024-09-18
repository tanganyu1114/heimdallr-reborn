<script>

import { _ContextCreator, NewContextGenerator } from './contextCreator'
import ContextCreator from './contextCreator.vue'

export class NewDirectiveGenerator extends NewContextGenerator {
  constructor(data) {
    super(data)
    var formatData = {
      'context-type': 'directive',
      'context-value': ''
    }
    this.data = undefined
    if (data === undefined || data.directiveKey === undefined || data.directiveKey === '') return
    formatData['context-value'] += data.directiveKey
    if (data.directiveValue !== undefined || data.directiveValue !== '') formatData['context-value'] += ' ' + data.directiveValue
    this.data = formatData
  }
  generateNewContextOptions() {
    return this.data
  }
}

export class DirectiveCreator extends _ContextCreator {
  get camelCaseCtxName() { return 'Directive' }
  defaultFormData() {
    return {
      directiveKey: '',
      directiveValue: ''
    }
  }
  formRules() {
    return {
      directiveKey: [
        { required: true, message: '请输入指令名', type: 'string', trigger: 'blur' }
      ],
      directiveValue: [
        { required: false, message: '请输入指令参数', type: 'string', trigger: 'blur' }
      ]
    }
  }
  initFormDataWithTargetNode(targetNode, position) {
    // TODO: 插入时上下文符合校验
  }
  generateSubmitOptions() {
    return new NewDirectiveGenerator(this.formData).generateNewContextOptions()
  }
}

export default {
  name: 'DirectiveCreator',
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
      creator: new DirectiveCreator()
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
      <el-form-item prop="directiveKey" label="指令名称">
        <el-input v-model="creator.formData.directiveKey" placeholder="请输入内容" />
      </el-form-item>
      <el-form-item prop="directiveValue" label="指令参数">
        <el-input v-model="creator.formData.directiveValue" placeholder="请输入内容" />
      </el-form-item>
    </context-creator>
  </div>
</template>

<style scoped lang="scss">

</style>
