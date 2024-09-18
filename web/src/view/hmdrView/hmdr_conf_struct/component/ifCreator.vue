<script>

import { _ContextCreator, NewContextGenerator } from './contextCreator'
import ContextCreator from './contextCreator.vue'

export class NewIfGenerator extends NewContextGenerator {
  constructor(data) {
    super(data)
    var formatData = {
      'context-type': 'if',
      'context-value': ''
    }
    this.data = undefined
    if (data === undefined || data.ifValue === undefined || data.ifValue === '') return
    formatData['context-value'] = '(' + data.ifValue + ')'
    this.data = formatData
  }
  generateNewContextOptions() {
    return this.data
  }
}

export class IfCreator extends _ContextCreator {
  get camelCaseCtxName() { return 'If' }
  defaultFormData() {
    return {
      ifValue: ''
    }
  }
  formRules() {
    return {
      ifValue: [
        { required: true, message: '请输入匹配内容', type: 'string', trigger: 'blur' }
      ]
    }
  }
  initFormDataWithTargetNode(targetNode, position) {
    // TODO: 插入时上下文符合校验
  }
  generateSubmitOptions() {
    return new NewIfGenerator(this.formData).generateNewContextOptions()
  }
}

export default {
  name: 'IfCreator',
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
      creator: new IfCreator()
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
      <el-form-item prop="ifValue" label="匹配参数">
        <el-input v-model="creator.formData.ifValue" placeholder="请输入内容" />
      </el-form-item>
    </context-creator>
  </div>
</template>

<style scoped lang="scss">

</style>
