<script>

class NewLocationGenerator {
  constructor(data) {
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

export default {
  name: 'LocationCreator',
  data() {
    return {
      matchTypes: [
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
      ],
      defaultFormData: {
        ctxType: 'location',
        createType: 'empty',
        matchType: ''
      },
      formData: {
        ctxType: 'location',
        createType: 'empty',
        matchType: ''
      },
      dialogVisible: false,
      defaultTagNodeProps: {
        label: 'label',
        isLeaf: 'isLeaf'
      },
      defaultTagNodeData: [
        {
          id: 'new-location',
          label: 'location',
          ctxType: 'location',
          dragType: 'create',
          isLeaf: false
        }
      ],
      tagNodeData: [
        {
          id: 'new-location',
          label: 'location',
          ctxType: 'location',
          dragType: 'create',
          isLeaf: false
        }
      ]
    }
  },
  methods: {
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
    },
    openDialog() {
      this.dialogVisible = true
    },
    closeDialog() {
      this.dialogVisible = false
      this.$nextTick(() => {
        this.resetForm()
      })
    },
    submitForm(formName) {
      // this.resetTagNodeData()
      this.$refs[formName].validate((valid) => {
        if (valid) {
          var options = new NewLocationGenerator(this.formData).generateNewContextOptions()
          if (options === undefined) {
            // alert('提交的表单数据异常：' + this.formData.toString())
            this.$message({
              showClose: true,
              message: '提交的表单数据异常，表单内容：' + JSON.stringify(this.formData),
              type: 'warning'
            })
            return
          }
          this.$emit('form-commit-event', options, (res) => {
            this.$nextTick(() => {
              if (!res) {
                this.$message({
                  showClose: true,
                  message: '新建Location提交失败！',
                  type: 'error'
                })
                return
              }
              this.$message({
                showClose: true,
                message: '新建Location提交成功',
                type: 'success'
              })
              this.closeDialog()
            })
          })
        } else {
          // console.log(this.$refs[formName].toString())
          this.$message({
            showClose: true,
            message: '表单校验异常！',
            type: 'warning'
          })
        }
      })
    },
    resetForm() {
      this.formData = JSON.parse(JSON.stringify(this.defaultFormData))
      this.$emit('dialog-before-close-event')
    },
    resetTagNodeData() {
      this.tagNodeData = JSON.parse(JSON.stringify(this.defaultTagNodeData))
    },
    handleDragStart(node, event) {
      // console.log('handle location tag drag start')
      this.$emit('node-drag-start', { node: node }, event)
    },
    handleDragEnd(draggingNode, endNode, position, event) {
      this.$emit('node-drag-end', draggingNode, endNode, position, event)
      this.$nextTick(() => {
        this.resetTagNodeData()
      })
    }
  }
}
</script>

<template>
  <div>
    <el-tree
      ref="locTagTree"
      :data="tagNodeData"
      :props="defaultTagNodeProps"
      draggable
      @node-drag-start="handleDragStart"
      @node-drag-end="handleDragEnd"
    >
      <el-tag effect="plain">Location</el-tag>
    </el-tree>
    <el-dialog title="新建Location" width="25%" :visible.sync="dialogVisible" :before-close="closeDialog">
      <el-form
        ref="locationCreatorForm"
        :model="formData"
        :rules="formRules()"
        :validate-on-rule-change="false"
        label-width="100px"
        label-position="left"
      >
        <el-form-item prop="matchType" label="匹配类型">
          <el-select v-model="formData.matchType" placeholder="请选择">
            <el-option
              v-for="item in matchTypes"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item prop="matchValue" label="匹配路径">
          <el-input v-model="formData.matchValue" placeholder="请输入内容" />
        </el-form-item>
        <el-form-item prop="creatorType" label="Location类型">
          <el-radio-group v-model="formData.createType" size="mini">
            <el-radio label="empty">空配置</el-radio>
            <el-radio label="static">静态资源</el-radio>
            <el-radio label="proxy">反向代理</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="formData.createType !== undefined && formData.createType === 'proxy'" prop="proxyURL" label="被反向代理的URL" label-width="140px">
          <el-input v-model="formData.proxyURL" placeholder="请输入内容" />
        </el-form-item>
        <el-form-item v-if="formData.createType !== undefined && formData.createType === 'static'" prop="rootPath" label="静态资源根路径" label-width="110px">
          <el-input v-model="formData.rootPath" placeholder="请输入内容（缺省值：html）" />
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="submitForm('locationCreatorForm')">创建</el-button>
        <el-button @click="closeDialog">取消</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">

</style>
