<template>
  <div class="bottom-info">
    <div>
      <span>Powered by 测试环境组</span>
      <el-divider direction="vertical" />
      <span>Copyright @ Heimdallr-Reborn</span>
    </div>
  </div>
</template>
<script>
export default {
  name: 'BottomInfo'
}
</script>
<style lang="scss">
.bottom-info {
  color: #888;
  height: 30px;
  line-height: 12px;
  a {
    color: #888;
  }
  div {
    display: flex;
    justify-content: center;
    span{
      margin: 0 3px;
    }
  }
}
</style>
