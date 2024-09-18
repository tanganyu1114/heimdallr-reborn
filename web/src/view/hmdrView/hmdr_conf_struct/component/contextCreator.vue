<script>

export default {
  name: 'ContextCreator',
  props: {
    submitForm: {
      type: Function,
      required: false,
      default: undefined
    },
    editable: {
      type: Boolean,
      required: true,
      default: false
    },
    creator: {
      type: Object,
      required: true
    }
  },
  methods: {
    defaultSubmitForm() {
      if (this.submitForm !== undefined && typeof this.submitForm === 'function') {
        this.submitForm()
      } else {
        this.$refs[this.creator.formRefName()].validate((valid) => {
          if (valid) {
            var options = this.creator.generateSubmitOptions()
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
                    message: this.creator.dialogTitle() + '表单提交失败！',
                    type: 'error'
                  })
                  return
                }
                this.$message({
                  showClose: true,
                  message: this.creator.dialogTitle() + '表单提交成功',
                  type: 'success'
                })
                this.closeDialog()
              })
            })
          } else {
            // console.log(this.$refs[formRefName].toString())
            this.$message({
              showClose: true,
              message: '表单校验异常！',
              type: 'warning'
            })
          }
        })
      }
    },
    initFormDataWithTargetNode(targetNode, position) { this.creator.initFormDataWithTargetNode(targetNode, position) },
    openDialog() { this.creator.openDialog() },
    closeDialog() {
      this.creator.closeDialog()
      this.$nextTick(() => {
        this.$emit('dialog-before-close-event')
        this.creator.resetFormData()
      })
    },
    handleDragStart(node, event) {
      // console.log('handle location tag drag start')
      this.$emit('node-drag-start', { node: node }, event)
    },
    handleDragEnd(draggingNode, endNode, position, event) {
      this.$emit('node-drag-end', draggingNode, endNode, position, event)
      this.$nextTick(() => {
        this.creator.resetTagNodeData()
      })
    }
  }
}
</script>

<template>
  <div>
    <el-tree
      :ref="creator.tagTreeRefName()"
      :data="creator.tagNodeData"
      :props="creator.defaultTagNodeProps"
      :draggable="editable"
      @node-drag-start="handleDragStart"
      @node-drag-end="handleDragEnd"
    >
      <el-tag effect="plain">{{ creator.camelCaseCtxName }}</el-tag>
    </el-tree>
    <el-dialog :title="creator.dialogTitle()" width="35%" :visible.sync="creator.dialogVisible" :before-close="closeDialog">
      <el-form
        :ref="creator.formRefName()"
        :model="creator.formData"
        :rules="creator.formRules()"
        :validate-on-rule-change="false"
        label-width="100px"
        label-position="left"
      >
        <slot />
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="defaultSubmitForm">创建</el-button>
        <el-button @click="closeDialog">取消</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">

</style>
