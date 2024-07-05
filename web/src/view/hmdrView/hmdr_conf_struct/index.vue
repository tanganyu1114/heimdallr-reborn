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
            <el-button size="medium" type="primary" icon="el-icon-search" round @click="searchConfStruct">查询</el-button>
          </el-col>
        </el-form>
      </el-row>
    </el-card>
    <el-row>
      <el-col :span="6">
        <el-card class="configListClass">
          <el-radio-group v-model="currentConfig" vertical @change="changeCurrentConfStruct">
            <el-radio v-if="configsData.mainConfig" :key="0" :style="{ fontWeight: 'bold'}" :label="configsData.mainConfig">
              主配置：{{ configsData.mainConfig }}
            </el-radio>
            <template v-for="(item,index) in Object.keys(configsData.configs)">
              <el-radio v-if="item !== configsData.mainConfig" :key="index+1" :label="item">
                {{ item }}
              </el-radio>
            </template>
          </el-radio-group>
        </el-card>
      </el-col>
      <el-col :span="18">
        <el-card class="configStructureClass">

          <!--<div style="height:600px;">
            <el-scrollbar style="height:100%">
              <highlightjs language="nginx" :code="code" />
            </el-scrollbar>
          </div>-->
          <!--          <highlightjs language="nginx" :code="code" class="hljs" />-->
          <!--            @node-drag-start="backupTreeNodeDrag"-->
          <el-tree
            :data="currentConfStruct"
            :props="defaultProps"
            :default-expanded-keys="currentTreeExpandedKeysMap[currentConfig]"
            node-key="id"
            draggable
            highlight-current
            :allow-drop="allowTreeDrop"
            :allow-drag="allowTreeDrag"
            :render-content="treeRenderContent"
            @node-expand="(data)=>handleTreeNodeExpand(data)"
            @node-collapse="(data)=>handleTreeNodeCollapse(data)"
            @node-drop="handleTreeNodeDrop"
          />
          <el-dialog title="拖拽确认" width="25%" :visible.sync="dragDialogVisible">
            <p>请选择拖拽操作定义选项:</p>
            <el-radio-group v-model="dragRadio" size="mini">
              <el-radio label="mv">移动</el-radio>
              <el-radio label="cp">复制</el-radio>
            </el-radio-group>
            <span slot="footer" class="dialog-footer">
              <el-button @click="cancelDragDialog">取 消</el-button>
              <el-button type="primary" :loading="isConfirmingDrag" @click.stop="confirmDragDialog">确 定</el-button>
            </span>
          </el-dialog>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>

class ContextStructBuilder {
  constructor() {
    this.builders = {
      'server': function(data) {
        return new ServerCtx(data)
      },
      'location': function(data) {
        return new LocationCtx(data)
      }
    }
  }
  build(data) {
    if (!(data.ctxType in this.builders)) return new ContextStruct(data)
    return this.builders[data.ctxType](data)
  }
}

class ContextStruct {
  constructor(data) {
    this.data = data
  }
  genExtendLabel() {
    return ''
  }
}

class ServerCtx extends ContextStruct {
  constructor(data) {
    super(data)
    this.ports = []
    this.server_names = []
    const listenRE = /^listen\s+(.*)/
    const serverNameRE = /^server_name\s+(.*)/
    for (let i = 0; i < this.data.children.length; i++) {
      if (this.data.children[i].ctxType === 'directive') {
        switch (true) {
          case listenRE.test(this.data.children[i].value): {
            this.ports.push(listenRE.exec(this.data.children[i].value)[1].split(/\s+/))
            break
          }
          case serverNameRE.test(this.data.children[i].value): {
            this.server_names.push(serverNameRE.exec(this.data.children[i].value)[1].split(/\s+/))
            break
          }
        }
      }
    }
  }
  genExtendLabel() {
    var el = '<=='
    if (this.ports.length > 0) {
      el += ' 侦听端口：“' + this.ports.join('”，“') + '”'
    } else {
      el += '侦听默认端口：“80/8000”？'
    }
    if (this.server_names.length > 0) {
      el += '，服务名：“' + this.server_names.join('”，“') + '”'
    }
    return el
  }
}

class LocationCtx extends ContextStruct {
  constructor(data) {
    super(data)
    // TODO: 支持if上下文内的代理或静态资源配置
    this.proxy = ''
    this.root = ''
    const proxyPassRE = /^proxy_pass\s+(.*)/
    const rootRE = /^root\s+(.*)/
    for (let i = 0; i < this.data.children.length; i++) {
      if (this.data.children[i].ctxType === 'directive') {
        switch (true) {
          case proxyPassRE.test(this.data.children[i].value): {
            this.proxy = proxyPassRE.exec(this.data.children[i].value)[1]
            break
          }
          case rootRE.test(this.data.children[i].value): {
            this.root = rootRE.exec(this.data.children[i].value)[1]
            break
          }
        }
      }
    }
  }
  genExtendLabel() {
    if (this.proxy !== '') return '<== 代理至：“' + this.proxy + '”'
    if (this.root !== '') return '<== 静态资源路径：“' + this.root + '”'; else return '<== 默认静态资源路径：“html”？'
  }
}

import { getOptions, getConfStruct, moveCtx, insertCloneCtx } from '@/api/hmdr_conf.js'

export default {
  name: 'HmdrConfStruct',
  data() {
    return {
      currentConfig: '',
      currentConfStruct: [],
      currentTreeExpandedKeysMap: {},
      updateRequestData: {},
      configsData: {
        mainConfig: '',
        configs: {}
      },
      // code: [],
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
      Options: [],
      defaultProps: {
        label: this.toTreeLabel,
        isLeaf: this.isTreeLeaf
      },
      dragRadio: '',
      dragDialogVisible: false,
      isConfirmingDrag: false
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
      }
    },
    async refreshConfStruct() {
      const sf = {
        group_id: this.formData.value[0],
        host_id: this.formData.value[1],
        srv_name: this.formData.value[2]
      }
      const res = await getConfStruct(sf)
      // console.log(res)
      if (res.code === 0) {
        // console.log('配置赋值')
        this.configsData.mainConfig = res.data['main-config']
        this.configsData.configs = res.data.configs
        return true
      }
      return false
    },
    async searchConfStruct() {
      this.$refs['elForm'].validate(async(valid) => {
        if (!valid) return
        if (await this.refreshConfStruct()) {
          await this.changeConfStructTo(this.configsData.mainConfig)
        }
      })
    },
    async changeCurrentConfStruct() {
      this.currentConfStruct = []
      for (let i = 0; i < this.configsData.configs[this.currentConfig].params.length; i++) {
        var childCtx = this.formatConfStruct([i], this.configsData.configs[this.currentConfig].params[i])
        if (childCtx !== {}) {
          this.currentConfStruct.push(childCtx)
        }
      }
    },
    async changeConfStructTo(configName) {
      this.currentConfig = configName
      await this.changeCurrentConfStruct()
    },
    formatConfStruct(pos, contextNode) {
      if (Array.isArray(pos) && typeof contextNode !== 'string') {
        var formattedContext = {
          ctxType: contextNode['context-type'],
          value: contextNode['value'],
          pos: pos,
          id: pos.toString(),
          children: []
        }
        formattedContext.label = this.toTreeLabel(formattedContext)
        formattedContext.isLeaf = this.isTreeLeaf(formattedContext)
        if ('params' in contextNode) {
          for (let i = 0; i < contextNode.params.length; i++) {
            const childPos = pos.concat(i)
            const childCtx = this.formatConfStruct(childPos, contextNode.params[i])
            if (childCtx !== {}) {
              formattedContext.children.push(childCtx)
            }
          }
        }
        return formattedContext
      } else if (typeof contextNode === 'string') {
        return {
          label: contextNode,
          isLeaf: true,
          id: pos.toString()
        }
      } else {
        return {}
      }
    },
    allowTreeDrop(draggingNode, dropNode, type) {
      if (dropNode.data.isLeaf) {
        return type !== 'inner'
      } else return dropNode.data.pos !== undefined
    },
    allowTreeDrag(node) {
      return !(node.data.isLeaf && node.parent.data.ctxType === 'include')
    },
    isTreeLeaf(data, node) {
      return !(data !== {} && typeof data !== 'string' && data.ctxType !== undefined && !['directive', 'inline_comment', 'comment'].includes(data.ctxType))
    },
    toTreeLabel(data, node) {
      if (data !== {}) {
        if (data.ctxType !== undefined) {
          switch (data.ctxType) {
            case 'directive': {
              return data.value + ';'
            }
            case 'inline_comment': {
              return '└─ # ' + data.value
            }
            case 'comment': {
              return '# ' + data.value
            }
            default: {
              if (data.value !== undefined) {
                return data.ctxType + ': ' + data.value
              }
              return data.ctxType
            }
          }
        } else if ('label' in data && typeof data.label === 'string') return data.label
      }
      return 'Error Context!'
    },
    treeRenderContent(h, { data, node, store }) {
      var el = new ContextStructBuilder().build(data).genExtendLabel()
      if (node.parent.data.ctxType !== undefined && node.parent.data.ctxType === 'include') {
        // var includeConfName = node.parent.label
        // console.log(node.label)
        var includeConfCheckOut = <span> <el-button size='mini' type='info' icon='el-icon-position' on-click = { () => this.changeConfStructTo(node.label) }>跳转至该配置< /el-button></span>
      }
      // return (
      //   <span class='custom-tree-node'>
      //     <span>{node.label}</span><span v-if="el !== ''" style='color: darkseagreen'> {el}</span>
      //   </span>
      // )
      if (el === undefined || el === '') {
        return (
          <span class='custom-tree-node'>
            <span>{node.label}</span>{includeConfCheckOut}
          </span>)
      } else {
        return (
          <span class='custom-tree-node'>
            <span>{node.label}</span>{includeConfCheckOut} <span style='color: darkseagreen'>{el}</span>
          </span>
        )
      }
    },
    handleTreeNodeExpand(data) {
      if (this.currentTreeExpandedKeysMap[this.currentConfig] === undefined) this.currentTreeExpandedKeysMap[this.currentConfig] = []
      if (this.currentTreeExpandedKeysMap[this.currentConfig].indexOf(data.id) === -1) {
        // console.log('expand: ' + data.id)
        this.currentTreeExpandedKeysMap[this.currentConfig].push(data.id)
      }
    },
    handleTreeNodeCollapse(data) {
      if (this.currentTreeExpandedKeysMap[this.currentConfig] !== undefined) {
        var index = this.currentTreeExpandedKeysMap[this.currentConfig].indexOf(data.id)
        if (index > -1) {
          // console.log('collapse: ' + data.id)
          this.currentTreeExpandedKeysMap[this.currentConfig].splice(index, 1)
        }
      }
    },
    handleTreeNodeDrop(draggingNode, dropNode, dropType, ev) {
      var draggingPos = draggingNode.data.pos
      var targetPos = dropNode.data.pos
      switch (dropType) {
        case 'before': {
          // console.log('dragging to before')
          break
        }
        case 'after': {
          // console.log('dragging to after')
          // console.log(targetPos.toString())
          targetPos[targetPos.length - 1] += 1
          // console.log(targetPos.toString())
          break
        }
        case 'inner': {
          // console.log('dragging into latest')
          targetPos.push(dropNode.data.children.length)
          break
        }
        default: {
          return
        }
      }
      this.updateRequestData = {
        'web-server-options': {
          group_id: this.formData.value[0],
          host_id: this.formData.value[1],
          srv_name: this.formData.value[2]
        },
        'target-config-context-options': {
          position: {
            config: this.currentConfig,
            'context-pos-path': targetPos
          },
          'target-context': {
            'clone-context-pos': {
              config: this.currentConfig,
              'context-pos-path': draggingPos
            }
          }
        }
      }
      this.dragRadio = 'mv'
      this.dragDialogVisible = true
    },
    cancelDragDialog() {
      this.changeConfStructTo(this.currentConfig)
      this.dragDialogVisible = false
      this.isConfirmingDrag = false
    },
    async confirmDragDialog() {
      if (this.isConfirmingDrag) return
      this.isConfirmingDrag = true
      var res = {
        code: 7,
        msg: '放弃拖拽操作并离开页面'
      }
      var currentConfName = this.currentConfig
      switch (this.dragRadio) {
        case 'mv': {
          res = await moveCtx(this.updateRequestData)
          break
        }
        case 'cp': {
          res = await insertCloneCtx(this.updateRequestData)
          break
        }
        default: {
          res.msg = '未知的拖拽操作'
          break
        }
      }
      // console.log(res)
      if (res.code !== 0) {
        // TODO: 错误提示框
        // console.log('cancel dialog')
        this.cancelDragDialog()
        return
      }
      // console.log('search config struct')
      if (await this.refreshConfStruct()) {
        // 置空当前配置树节点展开状态
        this.currentTreeExpandedKeysMap[currentConfName] = []
        // console.log('change config struct to current config')
        await this.changeConfStructTo(currentConfName)
      }
      this.dragDialogVisible = false
      this.isConfirmingDrag = false
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
/*.app {
  height: 400px;
  overflow: hidden;
}
.el-scrollbar__wrap {
  overflow: visible;
  overflow-x: hidden;
}*/
.searchClass {
  padding-bottom: 0;
}
.configListClass {
  height: 600px;
  width: 100%;
  overflow-y: scroll;
  overflow-x: hidden!important;
  font-size: 16px;
}
.configStructureClass {
  height: 600px;
  width: 100%;
  overflow-y: scroll;
  overflow-x: hidden!important;
  font-size: 16px;
  //background-color: #DCDFE6;
}
.radio-messagebox .el-message-box__content {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}
/*.grayColor {
  color: darkseagreen;
}*/
/*.auto-wrap {
  width: 100%;         !* 设置容器宽度 *!
  word-wrap: break-word; !* 允许在长单词或URL中间进行换行 *!
  white-space: normal;   !* 允许自动换行 *!
  overflow-wrap: break-word; !* 确保旧版浏览器中的换行行为 *!
}*/
</style>
