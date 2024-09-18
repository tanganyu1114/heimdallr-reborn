<script>

import { _ContextCreator, NewContextGenerator } from './contextCreator'
import ContextCreator from './contextCreator.vue'

export class NewStreamGenerator extends NewContextGenerator {
  constructor(data) {
    super(data)
    // this.data = undefined
    // if (data !== undefined) {
    // }
    this.data = {
      'context-type': 'stream',
      'context-value': ''
    }
  }
  generateNewContextOptions() {
    return this.data
  }
}

export class StreamCreator extends _ContextCreator {
  get camelCaseCtxName() { return 'Stream' }
  generateSubmitOptions() {
    return new NewStreamGenerator(this.formData).generateNewContextOptions()
  }
}

export default {
  name: 'StreamCreator',
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
      creator: new StreamCreator()
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
