<template>
  <el-tooltip :content="tooltipContent" :placement="placement">
    <span class="copyable-text" @click="handleCopy">
      <slot>{{ text }}</slot>
    </span>
  </el-tooltip>
</template>

<script>
export default {
  name: 'CopyableText',
  props: {
    text: {
      type: String,
      default: ''
    },
    tooltipContent: {
      type: String,
      default: '点击复制'
    },
    placement: {
      type: String,
      default: 'top'
    }
  },
  methods: {
    handleCopy() {
      const copyText = this.text || this.$slots.default?.[0]?.text
      if (!copyText) return

      const input = document.createElement('input')
      input.value = copyText
      document.body.appendChild(input)
      input.select()
      document.execCommand('copy')
      document.body.removeChild(input)

      this.$message({
        type: 'success',
        message: '已复制到剪贴板',
        duration: 1500
      })
    }
  }
}
</script>

<style scoped>
.copyable-text {
  cursor: pointer;
  color: #409eff;
  transition: color 0.2s;
}
.copyable-text:hover {
  color: #66b1ff;
  text-decoration: underline;
}
</style>
