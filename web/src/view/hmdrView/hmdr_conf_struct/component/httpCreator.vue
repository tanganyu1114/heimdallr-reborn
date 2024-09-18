<script>

import { _ContextCreator, NewContextGenerator } from './contextCreator'
import ContextCreator from './contextCreator.vue'

export class NewHTTPGenerator extends NewContextGenerator {
  constructor(data) {
    super(data)
    // this.data = undefined
    // if (data !== undefined) {
    // }
    this.data = {
      'context-type': 'http',
      'context-value': ''
    }
  }
  generateNewContextOptions() {
    return this.data
  }
}

export class HTTPCreator extends _ContextCreator {
  get camelCaseCtxName() { return 'HTTP' }
  generateSubmitOptions() {
    return new NewHTTPGenerator(this.formData).generateNewContextOptions()
  }
}

export default {
  name: 'HTTPCreator',
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
      creator: new HTTPCreator()
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
    />
  </div>
</template>

<style scoped lang="scss">

</style>
