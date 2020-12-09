<template>
  <div>
    <div class="search-term">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="组名">
          <el-input v-model="searchInfo.groupId" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="openDialog">新增主机信息</el-button>
        </el-form-item>
        <el-form-item>
          <el-popover v-model="deleteVisible" placement="top" width="160">
            <p>确定要删除吗？</p>
            <div style="text-align: right; margin: 0">
              <el-button size="mini" type="text" @click="deleteVisible = false">取消</el-button>
              <el-button size="mini" type="primary" @click="onDelete">确定</el-button>
            </div>
            <el-button slot="reference" icon="el-icon-delete" size="mini" type="danger">批量删除</el-button>
          </el-popover>
        </el-form-item>
      </el-form>
    </div>
    <el-table
      ref="multipleTable"
      :data="tableData"
      border
      stripe
      style="width: 100%"
      tooltip-effect="dark"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column label="日期" width="180">
        <template slot-scope="scope">{{ scope.row.CreatedAt|formatDate }}</template>
      </el-table-column>

      <el-table-column label="组名" prop="groupId" width="120">
        <template slot-scope="scope">
          {{ filterDict(scope.row.groupId,"hmdr_group") }}
        </template>
      </el-table-column>

      <el-table-column label="主机名" prop="name" width="120" />

      <el-table-column label="描述信息" prop="description" width="120" />

      <el-table-column label="状态" prop="status" width="120">
        <template slot-scope="scope">{{ scope.row.status|formatBoolean }}</template>
      </el-table-column>

      <el-table-column label="IP地址" prop="ipaddr" width="120" />

      <el-table-column label="端口" prop="port" width="120" />

      <el-table-column label="token认证" prop="token" width="400" />

      <el-table-column label="操作">
        <template slot-scope="scope">
          <el-button class="table-button" size="small" type="primary" icon="el-icon-edit" @click="updateHmdrHost(scope.row)">变更</el-button>
          <el-popover v-model="scope.row.visible" placement="top" width="160">
            <p>确定要删除吗？</p>
            <div style="text-align: right; margin: 0">
              <el-button size="mini" type="text" @click="scope.row.visible = false">取消</el-button>
              <el-button type="primary" size="mini" @click="deleteHmdrHost(scope.row)">确定</el-button>
            </div>
            <el-button slot="reference" type="danger" icon="el-icon-delete" size="mini">删除</el-button>
          </el-popover>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      :current-page="page"
      :page-size="pageSize"
      :page-sizes="[10, 30, 50, 100]"
      :style="{float:'right',padding:'20px'}"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
    />

    <el-dialog :before-close="closeDialog" :visible.sync="dialogFormVisible" title="弹窗操作">
      <el-form :model="formData" label-position="right" label-width="80px">
        <el-form-item label="分组信息:">
          <el-select v-model="formData.groupId" placeholder="请选择" clearable>
            <el-option v-for="(item,key) in hmdr_groupOptions" :key="key" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="主机名:">
          <el-input v-model="formData.name" clearable placeholder="请输入" />
        </el-form-item>

        <el-form-item label="描述信息:">
          <el-input v-model="formData.description" clearable placeholder="请输入" />
        </el-form-item>

        <el-form-item label="状态:">
          <el-switch v-model="formData.status" active-color="#13ce66" inactive-color="#ff4949" active-text="启用" inactive-text="禁用" clearable />
        </el-form-item>

        <el-form-item label="IP地址:">
          <el-input v-model="formData.ipaddr" clearable placeholder="请输入" />
        </el-form-item>

        <el-form-item label="端口:">
          <el-input v-model="formData.port" clearable placeholder="请输入" />
        </el-form-item>

        <el-form-item label="token认证:">
          <el-input v-model="formData.token" clearable placeholder="请输入" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="closeDialog">取 消</el-button>
        <el-button type="primary" @click="enterDialog">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {
  createHmdrHost,
  deleteHmdrHost,
  deleteHmdrHostByIds,
  updateHmdrHost,
  findHmdrHost,
  getHmdrHostList
} from '@/api/hmdr_host' //  此处请自行替换地址
import { formatTimeToStr } from '@/utils/date'
import infoList from '@/mixins/infoList'
import { getHmdrGroupList } from '@/api/hmdr_group.js'
export default {
  name: 'HmdrHost',
  filters: {
    formatDate: function(time) {
      if (time != null && time !== '') {
        var date = new Date(time)
        return formatTimeToStr(date, 'yyyy-MM-dd hh:mm:ss')
      } else {
        return ''
      }
    },
    formatBoolean: function(bool) {
      if (bool != null) {
        return bool ? '启用' : '禁用'
      } else {
        return ''
      }
      // if (typeof val === 'number' && !isNaN(val)) {
      //   return val === 1 ? '启用' : '禁用'
      // } else {
      //   return ''
      // }
    }
  },
  mixins: [infoList],
  data() {
    return {
      listApi: getHmdrHostList,
      dialogFormVisible: false,
      visible: false,
      type: '',
      deleteVisible: false,
      multipleSelection: [],
      hmdr_groupOptions: [],
      formData: {
        groupId: 1,
        name: '',
        description: '',
        status: Number(),
        ipaddr: '',
        port: '',
        token: ''

      }
    }
  },
  async created() {
    await this.getTableData()
    const res = await getHmdrGroupList()
    if (res.code === 0) {
      res.data.list.forEach((item) => {
        const obj = {}
        obj['label'] = item.name
        obj['value'] = item.ID
        this.hmdr_groupOptions.push(obj)
      })
    }

    // await this.getDict('hmdr_group')
  },
  methods: {
    // 条件搜索前端看此方法
    onSubmit() {
      this.page = 1
      this.pageSize = 10
      if (this.searchInfo.status === '') {
        this.searchInfo.status = null
      }
      this.getTableData()
    },
    handleSelectionChange(val) {
      this.multipleSelection = val
    },
    async onDelete() {
      const ids = []
      if (this.multipleSelection.length === 0) {
        this.$message({
          type: 'warning',
          message: '请选择要删除的数据'
        })
        return
      }
      this.multipleSelection &&
          this.multipleSelection.map(item => {
            ids.push(item.ID)
          })
      const res = await deleteHmdrHostByIds({ ids })
      if (res.code === 0) {
        this.$message({
          type: 'success',
          message: '删除成功'
        })
        this.deleteVisible = false
        this.getTableData()
      }
    },
    async updateHmdrHost(row) {
      const res = await findHmdrHost({ ID: row.ID })
      this.type = 'update'
      if (res.code === 0) {
        this.formData = res.data.rehmdrHost
        this.dialogFormVisible = true
      }
    },
    closeDialog() {
      this.dialogFormVisible = false
      this.formData = {
        groupId: 1,
        name: '',
        description: '',
        status: 1,
        ipaddr: '',
        port: '',
        token: ''

      }
    },
    async deleteHmdrHost(row) {
      this.visible = false
      const res = await deleteHmdrHost({ ID: row.ID })
      if (res.code === 0) {
        this.$message({
          type: 'success',
          message: '删除成功'
        })
        this.getTableData()
      }
    },
    async enterDialog() {
      let res
      switch (this.type) {
        case 'create':
          res = await createHmdrHost(this.formData)
          break
        case 'update':
          console.log(this.formData)
          res = await updateHmdrHost(this.formData)
          break
        default:
          res = await createHmdrHost(this.formData)
          break
      }
      if (res.code === 0) {
        this.$message({
          type: 'success',
          message: '创建/更改成功'
        })
        this.closeDialog()
        this.getTableData()
      }
    },
    openDialog() {
      this.type = 'create'
      this.dialogFormVisible = true
    }
  }
}
</script>

<style>
</style>
