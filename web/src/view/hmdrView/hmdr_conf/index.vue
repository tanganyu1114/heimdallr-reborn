<template>
  <div>
    <el-card>
      <el-row :gutter="15" class="searchClass">
        <el-form
          ref="elForm"
          :model="formData"
          :rules="rules"
          size="medium"
          label-width="100px"
          label-position="left"
        >
          <el-col :span="12">
            <el-form-item label-width="120px" prop="value" label="应用服务器选择">
              <el-cascader
                v-model="formData.value"
                :options="Options"
                :props="{ expandTrigger: 'hover' }"
                :style="{width: '100%'}"
                placeholder="请选择环境以及主机信息应用服务器选择"
                clearable
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-button size="medium" type="primary" icon="el-icon-search" round @click="searchConfInfo">查询</el-button>
          </el-col>
        </el-form>
      </el-row>
    </el-card>
    <el-card>

      <!--<div style="height:600px;">
        <el-scrollbar style="height:100%">
          <highlightjs language="nginx" :code="code" />
        </el-scrollbar>
      </div>-->
      <highlightjs language="nginx" :code="code" class="hljs" />
    </el-card>
  </div>
</template>

<script>
import { getOptions, getConfInfo } from '@/api/hmdr_conf.js'

export default {
  name: 'HmdrConf',
  data() {
    return {
      code: '',
      formData: {
        value: []
      },
      rules: {
        value: [{
          required: true,
          type: 'array',
          message: '请至少选择一个应用服务器选择',
          trigger: 'change'
        }]
      },
      Options: []
    }
  },
  created() {
    this.initOptions()
  },
  methods: {
    async initOptions() {
      const res = await getOptions()
      if (res.code === 0) {
        this.Options = res.data
        console.log(res.data)
      }
    },
    async searchConfInfo() {
      this.$refs['elForm'].validate(async(valid) => {
        if (!valid) return
        const sf = {
          group_id: this.formData.value[0],
          host_id: this.formData.value[1],
          srv_name: this.formData.value[2]
        }
        const res = await getConfInfo(sf)
        if (res.code === 0) {
          this.code = res.data
          console.log(res.data)
        }
      })
    }
  }
}

</script>

<style scoped>
/*.hljscss {*/
/*  height: 70%;*/
/*  position: relative;*/
/*  font-size: 18px;*/
/*  overflow-y: scroll;*/
/*}*/
.app {
  height: 400px;
  overflow: hidden;
}
.el-scrollbar__wrap {
  overflow: visible;
  overflow-x: hidden;
}
.searchClass {
  padding-bottom: 0;
}
.hljs {
  max-height: 600px;
  width: 100%;
  overflow-y: scroll;
  overflow-x: hidden!important;
  font-size: 16px;
}
</style>
