<script>

import { _ContextCreator, NewContextGenerator } from './contextCreator'
import ContextCreator from './contextCreator.vue'

export class NewLocationGenerator extends NewContextGenerator {
  constructor(data) {
    super(data)
    var formatData = {
      'context-type': 'location',
      'context-value': '',
      'children-context-meta': []
    }
    this.data = undefined
    if (data === undefined || data.matchType === undefined || data.matchValue === undefined || data.matchValue === '') {
      return
    }

    if (data.matchType === '') {
      formatData['context-value'] = data.matchValue
    } else {
      formatData['context-value'] = data.matchType + ' ' + data.matchValue
    }

    switch (data.createType) {
      case 'empty': {
        break
      }
      case 'static': {
        if (data.rootPath === undefined || data.rootPath === '') {
          // 静态资源 root 缺省值 html
          data.rootPath = 'html'
        }
        formatData['children-context-meta'].push({
          'context-type': 'directive',
          'context-value': 'root ' + data.rootPath
        })
        break
      }
      case 'proxy': {
        if (data.proxyURL === undefined || data.proxyURL === '') {
          return
        }
        formatData['children-context-meta'].push({
          'context-type': 'directive',
          'context-value': 'proxy_pass ' + data.proxyURL
        })
        break
      }
      default: {
        return
      }
    }
    this.data = formatData
  }
  generateNewContextOptions() {
    return this.data
  }
}

export class LocationCreator extends _ContextCreator {
  matchTypes = [
    {
      value: '',
      label: '前缀位置匹配（空)'
    },
    {
      value: '=',
      label: '精确完全匹配（=）'
    },
    {
      value: '~',
      label: '区分大小写的正则匹配（~）'
    },
    {
      value: '~*',
      label: '不区分大小写的正则匹配（~*）'
    },
    {
      value: '^~',
      label: '无正则前缀位置匹配（^~）'
    }
  ]
  get camelCaseCtxName() { return 'Location' }
  defaultFormData() {
    return {
      createType: 'empty',
      matchType: ''
    }
  }
  formRules() {
    return {
      matchValue: [
        { required: true, message: '请输入需匹配路径的规则', type: 'string', trigger: 'blur' }
      ],
      proxyURL: [
        { required: this.formData.createType === 'proxy', message: '请输入被代理的URL', type: 'string', trigger: 'blur' }
      ],
      rootPath: [
        { required: false, message: '请输入静态资源根目录路径', type: 'string', trigger: 'blur' }
      ]
    }
  }
  initFormDataWithTargetNode(targetNode, position) {
    // TODO: 插入时上下文符合校验
  }
  generateSubmitOptions() {
    return new NewLocationGenerator(this.formData).generateNewContextOptions()
  }
}

export default {
  name: 'LocationCreator',
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
      creator: new LocationCreator()
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
      <el-form-item prop="matchType" label="匹配类型">
        <el-select v-model="creator.formData.matchType" placeholder="请选择">
          <el-option
            v-for="item in creator.matchTypes"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
      </el-form-item>
      <el-form-item prop="matchValue" label="匹配路径">
        <el-input v-model="creator.formData.matchValue" placeholder="请输入内容" />
      </el-form-item>
      <el-form-item prop="creatorType" label="Location类型">
        <el-radio-group v-model="creator.formData.createType" size="mini">
          <el-radio label="empty">空配置</el-radio>
          <el-radio label="static">静态资源</el-radio>
          <el-radio label="proxy">反向代理</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item v-if="creator.formData.createType !== undefined && creator.formData.createType === 'proxy'" prop="proxyURL" label="被反向代理的URL" label-width="140px">
        <el-input v-model="creator.formData.proxyURL" placeholder="请输入内容" />
      </el-form-item>
      <el-form-item v-if="creator.formData.createType !== undefined && creator.formData.createType === 'static'" prop="rootPath" label="静态资源根路径" label-width="110px">
        <el-input v-model="creator.formData.rootPath" placeholder="请输入内容（缺省值：html）" />
      </el-form-item>
    </context-creator>
  </div>
</template>

<style scoped lang="scss">

</style>
