<script>

// import { searchCtxPoses } from '@/api/hmdr_conf.js'

export default {
  name: 'ContextPosesSearcher',
  data() {
    return {
      total: -1,
      index: 0,
      posList: [],
      onlyInCurrentConfig: true,
      isRegExp: false,
      searchKeywords: '',
      searchTypeRadio: 0
    }
  },
  methods: {
    async handleSearch() {
      this.$nextTick(async() => {
        await this.$emit('search',
          {
            onlyInCurrentConfig: this.onlyInCurrentConfig,
            isRegExp: this.isRegExp,
            searchKeywords: this.searchKeywords,
            searchTypeRadio: this.searchTypeRadio
          },
          this.posList,
          this.total,
          this.index)
        // if (!this.searchRequestParams.serverOptions || !this.searchRequestParams.originalFingerprints || !this.searchRequestParams.currentConfig) return
        // const reqData = {
        //   'web-server-options': this.searchRequestParams.serverOptions,
        //   'original-fingerprints': this.searchRequestParams.originalFingerprints,
        //   keywords: this.searchKeywords,
        //   'is-regexp-rule': this.isRegExp,
        //   'is-only-in-current': this.onlyInCurrentConfig
        // }
        // if (this.searchTypeRadio === 1) {
        //   reqData['starting-position-list'] = this.posList
        // } else {
        //   reqData['starting-position-list'] = {
        //     config: this.searchRequestParams.currentConfig,
        //     'context-pos-path': []
        //   }
        // }
        // const resp = await searchCtxPoses(reqData)
        // if (resp.code === 0) {
        //   this.total = resp.data.length
        //   this.index = 0
        //   this.posList = resp.data
        //   this.handleNext()
        // } else {
        //   this.total = -1
        //   this.index = 0
        //   this.posList = []
        // }
      })
    },
    handleIndexChangeEvent(posList, total, index) {
      this.$nextTick(() => {
        if (this.index < 0) return
        this.$emit('index-change-event', posList, total, index)
      })
    },
    handlePrev(posList, total, index) {
      if (index > 0) {
        index--
      } else {
        index = total
      }
      this.handleIndexChangeEvent(posList, total, index)
    },
    handleNext(posList, total, index) {
      if (total < 0) return
      if (total > index) {
        index++
      } else {
        index = 1
      }
      this.handleIndexChangeEvent(posList, total, index)
    }
  }
}
</script>

<template>
  <div>
    <span>
      <el-input
        v-model="searchKeywords"
        placeholder="请输入搜索关键字"
        @keyup.enter.native="handleSearch"
      >
        <el-checkbox slot="prepend" v-model="onlyInCurrentConfig">仅在当前配置搜索</el-checkbox>
        <el-checkbox slot="prepend" v-model="isRegExp">使用正则表达式搜索关键字</el-checkbox>
        <i v-if="total >= 0" slot="suffix">{{ index }} / {{ total }}</i>
        <el-radio-group v-if="total >= 0" slot="suffix" v-model="searchTypeRadio">
          <el-radio :label="0">重新搜索</el-radio>
          <el-radio :label="1">在当前搜索结果内搜索</el-radio>
        </el-radio-group>
        <el-button v-if="total >= 0" slot="append" icon="el-icon-arrow-left" @click="handlePrev(posList, total, index)">上一个</el-button>
        <el-button v-if="total >= 0" slot="append" icon="el-icon-arrow-right" @click="handleNext(posList, total, index)">下一个</el-button>
        <el-button slot="append" icon="el-icon-search" @click="handleSearch">搜索</el-button>
      </el-input>
    </span>
  </div>
</template>

<style scoped lang="scss">

</style>
