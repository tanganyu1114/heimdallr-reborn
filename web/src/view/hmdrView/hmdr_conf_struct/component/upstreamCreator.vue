<script>

import { _ContextCreator, NewContextGenerator } from './contextCreator'
import ContextCreator from './contextCreator.vue'

export class NewUpstreamGenerator extends NewContextGenerator {
  constructor(data) {
    super(data)
    var formatData = {
      'enabled': true,
      'context-type': 'upstream',
      'context-value': '',
      'children-context-meta': []
      // TODO: 被接管的服务
    }
    this.data = undefined
    if (data === undefined || data.upstreamValue === undefined || data.upstreamValue === '' || data.backends === undefined || data.backends === []) return
    formatData['context-value'] = data.upstreamValue
    for (let i = 0; i < data.backends.length; i++) {
      var directive = 'server ' + data.backends[i].server
      if (data.backends[i].params !== undefined && data.backends[i].params !== '') directive += ' ' + data.backends[i].params
      formatData['children-context-meta'].push({
        'enabled': true,
        'context-type': 'directive',
        'context-value': directive
      })
    }
    this.data = formatData
  }
  generateNewContextOptions() {
    return this.data
  }
}

export class UpstreamCreator extends _ContextCreator {
  get camelCaseCtxName() { return 'Upstream' }
  defaultFormData() {
    return {
      upstreamValue: '',
      backends: [
        { server: '' }
      ]
    }
  }
  formRules() {
    return {
      upstreamValue: [
        { required: true, message: '请输入上游服务名', type: 'string', trigger: 'blur' }
      ]
    }
  }
  initFormDataWithTargetNode(targetNode, position) {
    // TODO: 插入时上下文符合校验
  }
  generateSubmitOptions() {
    return new NewUpstreamGenerator(this.formData).generateNewContextOptions()
  }
}

export default {
  name: 'UpstreamCreator',
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
      creator: new UpstreamCreator()
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
    },
    removeBackend(index) {
      this.creator.formData.backends.splice(index, 1)
    },
    addBackend(index) {
      this.creator.formData.backends.splice(index + 1, 0, { server: '' })
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
      <el-form-item prop="upstreamValue" label="上游服务名">
        <el-input v-model="creator.formData.upstreamValue" placeholder="请输入内容" />
      </el-form-item>
      <el-row v-for="(item, index) in creator.formData.backends" :key="index">
        <el-col :span="10">
          <el-form-item
            label="后端地址"
            :prop="'backends.' + index + '.server'"
            :rules="[{ required: true, message: '请输入后端地址', type: 'string', trigger: 'blur' }]"
            label-width="80px"
          >
            <el-input v-model="item.server" placeholder="请输入内容" />
          </el-form-item>
        </el-col>
        <el-col :span="10">
          <el-form-item
            label="服务参数"
            :prop="'backends.' + index + '.params'"
            :rules="[{ required: false, message: '请输入正确的服务参数', type: 'string', trigger: 'blur' }]"
            label-width="80px"
          >
            <el-input v-model="item.params" placeholder="请输入内容" />
          </el-form-item>
        </el-col>
        <el-col :span="1.5">
          <el-button v-if="index > 0" size="mini" type="danger" icon="el-icon-minus" circle @click="removeBackend(index)" />
        </el-col>
        <el-col :span="1.5">
          <el-button size="mini" type="primary" icon="el-icon-plus" circle @click="addBackend(index)" />
        </el-col>
      </el-row>
    </context-creator>
  </div>
</template>

<style scoped lang="scss">

</style>
