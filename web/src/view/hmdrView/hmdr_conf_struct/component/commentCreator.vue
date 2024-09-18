<script>

import { _ContextCreator, NewContextGenerator } from './contextCreator'
import ContextCreator from './contextCreator.vue'

export class NewCommentGenerator extends NewContextGenerator {
  constructor(data) {
    super(data)
    var formatData = {
      'context-type': 'comment',
      'context-value': ''
    }
    this.data = undefined
    if (data !== undefined) {
      if (data.commentValue !== undefined && data.commentValue !== '') {
        formatData['context-value'] = data.commentValue
      }
      if (data.isInline) formatData['context-type'] = 'inline_comment'
    }
    this.data = formatData
  }
  generateNewContextOptions() {
    return this.data
  }
}

export class CommentCreator extends _ContextCreator {
  get camelCaseCtxName() { return 'Comment' }
  defaultFormData() {
    return {
      commentValue: '',
      isInline: false
    }
  }
  formRules() {
    return {
      isInline: [
        { required: false, type: 'boolean', trigger: 'blur' }
      ],
      commentValue: [
        { required: false, message: '请输入注释内容', type: 'string', trigger: 'blur' }
      ]
    }
  }
  initFormDataWithTargetNode(targetNode, position) {
    // TODO: 插入时上下文符合校验
  }
  generateSubmitOptions() {
    return new NewCommentGenerator(this.formData).generateNewContextOptions()
  }
}

export default {
  name: 'CommentCreator',
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
      creator: new CommentCreator()
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
      <el-form-item prop="isInline" label="是否为（前一非注解上下文）末尾同行注释" label-width="280px">
        <el-checkbox v-model="creator.formData.isInline" />
      </el-form-item>
      <el-form-item prop="commentValue" label="注释内容">
        <el-input v-model="creator.formData.commentValue" placeholder="请输入内容" />
      </el-form-item>
    </context-creator>
  </div>
</template>

<style scoped lang="scss">

</style>
