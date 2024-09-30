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
          <el-radio-group v-model="currentConfig" style="width: 100%" @change="(c) => changeConfStructTo(c)">
            <el-tree
              ref="configPathTree"
              :data="configsData.pathsStruct"
              :props="configPathTreeProps"
              node-key="fullPath"
              accordion
              highlight-current
            >
              <span slot-scope="{ node, data }" style="display: flex">
                <i v-if="data.isDir">
                  <i v-if="!node.expanded" class="el-icon-folder" />
                  <i v-else class="el-icon-folder-opened" />
                  {{ node.label }}
                </i>
                <i v-else>
                  <el-radio :label="data.configPath">
                    <i v-if="currentConfig === data.configPath" class="el-icon-document-checked" style="text-indent: -0.75em" />
                    <i v-else class="el-icon-document" style="text-indent: -0.75em" />
                    {{ node.label }}
                  </el-radio>
                </i>
              </span>
            </el-tree>
          </el-radio-group>
        </el-card>
      </el-col>
      <el-col :span="18" style="height: 610px; max-height: 610px; overflow: hidden">
        <el-card-collapse v-show="configStructEditable" ref="ctxCreatorsCard" :is-collapse="ctxCreatorsCardIsCollapse" @click-event="handleCtxCreatorsCardClickEvent">
          <div slot="header" class="flex-between">
            <transition name="el-fade-in">
              <el-button v-show="ctxCreatorsCardIsCollapse" style="position: absolute;top: 50%;left: 50%;transform: translate(-50%, -50%);padding: 3px 0; margin-right: 10px;" type="text">上下文新建模块</el-button>
            </transition>
          </div>
          <el-row>
            <template v-for="creator in ctxCreatorCompMetaList">
              <el-col :key="creator.key" :span="2.5">
                <component
                  :is="creator.comp"
                  :ref="creator.refName"
                  :editable="configStructEditable"
                  @form-commit-event="handleNewCtxCommitEvent"
                  @node-drag-start="handleCtxCreatorDragStart"
                  @node-drag-end="handleCtxCreatorDragEnd"
                  @dialog-before-close-event="resetUpdateRequest"
                />
              </el-col>
            </template>
          </el-row>
        </el-card-collapse>
        <el-card ref="configStructCard" class="configStructureClass">
          <el-tree
            ref="configTree"
            :data="currentConfStruct"
            :props="configStructTreeProps"
            :default-expanded-keys="currentTreeExpandedKeysMap[currentConfig]"
            node-key="id"
            :draggable="configStructEditable"
            highlight-current
            :allow-drop="allowTreeDrop"
            :allow-drag="allowTreeDrag"
            @node-expand="(data)=>handleTreeNodeExpand(data)"
            @node-collapse="(data)=>handleTreeNodeCollapse(data)"
            @node-drop="handleTreeNodeDrop"
          >
            <span
              slot-scope="{ node, data }"
              class="custom-tree-node"
              @mouseleave="() => handleTreeNodeMouseLeave(node)"
              @mouseenter="() => handleTreeNodeMouseEnter(node)"
            >
              <span>{{ node.label }}</span> <span v-if="extendTreeNodeLabel(node)" style="color: darkseagreen">{{ node.data.extendLabel }}</span> <span v-if="node.parent.data.ctxType !== undefined && node.parent.data.ctxType === 'include'">
                <el-button size="mini" type="info" icon="el-icon-position" @click="() => changeConfStructTo(node.label)">跳转至该配置</el-button>
              </span><span v-else v-show="node.data.delButtonShow && configStructEditable">
                <el-button v-show="isCtxWithValue(data)" size="mini" type="primary" icon="el-icon-edit" circle @click="() => handleTreeNodeModify(node, data)" />
                <el-button size="mini" type="danger" icon="el-icon-delete" circle @click="() => handleTreeNodeDelete(node, data)" />
              </span>
            </span>
          </el-tree>
          <div class="floating-button-container">
            <el-button
              v-show="configsData.stackCursor > 0"
              :class="{ 'floating-button': true }"
              :style="{ 'right': '70px' }"
              type="primary"
              icon="el-icon-back"
              @click="changeBackConfStruct"
            />
            <el-button
              v-show="configsData.stackCursor >= 0 && configsData.cachedConfigStack.length - 1 > configsData.stackCursor"
              :class="{ 'floating-button': true }"
              :style="{ 'right': '20px' }"
              type="primary"
              icon="el-icon-right"
              @click="changeForwardConfStruct"
            />
          </div>
          <el-dialog title="拖拽确认" width="25%" :visible.sync="dragTreeNodeDialogVisible">
            <p>请选择拖拽操作定义选项:</p>
            <el-radio-group v-model="dragTreeNodeRadio" size="mini">
              <el-radio label="mv">移动</el-radio>
              <el-radio label="cp">复制</el-radio>
            </el-radio-group>
            <span slot="footer" class="dialog-footer">
              <el-button @click="cancelDragDialog">取 消</el-button>
              <el-button type="primary" :loading="isConfirmingDrag" @click.stop="confirmDragDialog">确 定</el-button>
            </span>
          </el-dialog>
          <el-dialog title="修改上下文" width="25%" :visible.sync="modifyTreeNodeDialogVisible">
            <el-form
              v-if="updateRequestData['target-config-context-options'] !== undefined && updateRequestData['target-config-context-options']['target-context'] !== undefined"
              :model="updateRequestData['target-config-context-options']['target-context']"
              label-width="100px"
              label-position="left"
            >
              <el-form-item prop="context-value" label="上下文参数值">
                <el-input v-model="updateRequestData['target-config-context-options']['target-context']['context-value']" placeholder="请输入内容" />
              </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
              <el-button @click="cancelModifyDialog">取 消</el-button>
              <el-button type="primary" :loading="isConfirmingModify" @click.stop="confirmModifyDialog">确 定</el-button>
            </span>
          </el-dialog>
          <el-dialog title="删除确认" width="25%" :visible.sync="delTreeNodeDialogVisible">
            <p>请确认是否要删除该上下文[{{ delTreeNodeLabel }}]</p>
            <span slot="footer" class="dialog-footer">
              <el-button @click="cancelDelDialog">取 消</el-button>
              <el-button type="primary" :loading="isConfirmingDel" @click.stop="confirmDelDialog">确 定</el-button>
            </span>
          </el-dialog>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>

class ConfigPathTreeStruct {
  constructor(mainConfigPath, configPaths) {
    if (Array.isArray(configPaths) && typeof mainConfigPath === 'string' && mainConfigPath.charAt(0) === '/') {
      this.data = {
        main: mainConfigPath,
        paths: {
          children: {}
        },
        configPathTreeNodeKeyMap: {},
        treeStruct: []
      }

      for (const configPath of configPaths) {
        this.push(configPath)
      }
      this.data.treeStruct = this.genTreeStruct()
    }
    return undefined
  }
  getPathNode(pos) {
    if (pos === []) return this.data.paths
    if (!Array.isArray(pos)) return undefined
    var node = this.data.paths
    for (const idx of pos) {
      node = node.children[idx]
      if (node === undefined) return undefined
    }
    return node
  }
  push(configPath) {
    if (typeof configPath === 'string' && configPath !== '') {
      var pos = []
      if (configPath.charAt(0) === '/') {
        this.pushWithParsedPath(pos, configPath.substring(1))
      } else {
        pos = [...this.data.main.substring(1).split('/').slice(0, -1)]
        this.pushWithParsedPath(pos, configPath)
      }
      var node = this.getPathNode(pos)
      if (node !== undefined) {
        node.configPath = configPath
        this.data.configPathTreeNodeKeyMap[configPath] = node.fullPath
      }
    }
  }
  pushWithParsedPath(pos, path) {
    if (path === undefined || typeof path !== 'string' || path === '' || path.charAt(0) === '/' || path.charAt(-1) === '/') return
    var dsIdx = path.indexOf('/')
    var name, subPath
    if (dsIdx >= 0) {
      name = path.substring(0, dsIdx)
      subPath = path.substring(dsIdx + 1)
    } else {
      name = path
    }
    if (this.getPathNode(pos).children[name] === undefined) {
      this.getPathNode(pos).children[name] = {
        filename: name,
        fullPath: [...pos, name].join('/'),
        isDir: subPath !== undefined && subPath !== '',
        children: {}
      }
    }
    pos.push(name)
    this.pushWithParsedPath(pos, subPath)
  }
  formatToTreeStruct(paths) {
    var struct = []
    for (const k in paths.children) {
      struct.push({
        filename: paths.children[k].filename,
        fullPath: paths.children[k].fullPath,
        isDir: paths.children[k].isDir,
        configPath: paths.children[k].configPath,
        children: this.formatToTreeStruct(paths.children[k])
      })
    }
    return struct
  }
  genTreeStruct() {
    if (this.data === undefined || this.data.paths === undefined || this.data.main === undefined) return undefined
    return this.formatToTreeStruct(this.data.paths)
  }
}

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

import { getOptions, getConfStruct, removeCtx, moveCtx, modifyCtxValue, insertCloneCtx, insertNewCtx } from '@/api/hmdr_conf.js'
import CommentCreator from '@/view/hmdrView/hmdr_conf_struct/component/commentCreator.vue'
import DirectiveCreator from '@/view/hmdrView/hmdr_conf_struct/component/directiveCreator.vue'
import EventsCreator from '@/view/hmdrView/hmdr_conf_struct/component/eventsCreator.vue'
import HttpCreator from '@/view/hmdrView/hmdr_conf_struct/component/httpCreator.vue'
import IfCreator from '@/view/hmdrView/hmdr_conf_struct/component/ifCreator.vue'
import LocationCreator from '@/view/hmdrView/hmdr_conf_struct/component/locationCreator.vue'
import ServerCreator from '@/view/hmdrView/hmdr_conf_struct/component/serverCreator.vue'
import StreamCreator from '@/view/hmdrView/hmdr_conf_struct/component/streamCreator.vue'
import UpstreamCreator from '@/view/hmdrView/hmdr_conf_struct/component/upstreamCreator.vue'
import ElCardCollapse from '@/components/ElCardCollapse.vue'

export default {
  name: 'HmdrConfStruct',
  components: {
    CommentCreator,
    DirectiveCreator,
    EventsCreator,
    HttpCreator,
    IfCreator,
    LocationCreator,
    ServerCreator,
    StreamCreator,
    UpstreamCreator,
    ElCardCollapse
  },
  props: {
    configStructEditable: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      currentConfig: '',
      currentConfStruct: [],
      currentTreeExpandedKeysMap: {},
      deleteRequestData: {},
      updateRequestData: {},
      configsData: {
        mainConfig: '',
        pathsStruct: [],
        configStructs: {},
        cachedConfigStack: [],
        stackCursor: -1
      },
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
      configStructTreeProps: {
        label: this.toTreeLabel,
        isLeaf: this.isTreeLeaf
      },
      configPathTreeProps: {
        label: (data, node) => {
          return data.filename
        },
        isLeaf: (data, node) => {
          return !data.isDir
        }
      },
      modifyTreeNodeLabel: '未知上下文',
      modifyTreeNodeDialogVisible: false,
      isConfirmingModify: false,
      delTreeNodeLabel: '未知上下文',
      delTreeNodeDialogVisible: false,
      isConfirmingDel: false,
      dragTreeNodeRadio: '',
      dragTreeNodeDialogVisible: false,
      isConfirmingDrag: false,
      ctxCreatorsCardIsCollapse: true,
      ctxCreatorCompMetaList: this.ctxCreatorComponentsMeta()
    }
  },
  created() {
    this.initOptions()
  },
  mounted() {
    this.$nextTick(() => {
      this.setConfigStructCardHeight()
    })
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
        this.configsData.configs = Object.keys(res.data.configs)
        var pathStruct = new ConfigPathTreeStruct(res.data['main-config'], Object.keys(res.data.configs))
        this.configsData.pathsStruct = pathStruct.data.treeStruct
        this.configsData.configPathTreeNodeKeyMap = pathStruct.data.configPathTreeNodeKeyMap
        this.configsData.cachedConfigStack = []
        this.configsData.stackCursor = -1
        for (const c of Object.keys(res.data.configs)) {
          this.configsData.configStructs[c] = []
          for (let i = 0; i < res.data.configs[c].params.length; i++) {
            var childCtx = this.formatConfStruct([i], res.data.configs[c].params[i])
            if (childCtx !== {}) {
              this.configsData.configStructs[c].push(childCtx)
            }
          }
        }
        return true
      }
      return false
    },
    async searchConfStruct() {
      this.$refs['elForm'].validate(async(valid) => {
        if (!valid) return
        if (await this.refreshConfStruct()) {
          this.currentTreeExpandedKeysMap = {}
          await this.changeConfStructTo(this.configsData.mainConfig)
        }
      })
    },
    async changeCurrentConfStruct() {
      this.currentConfStruct = this.configsData.configStructs[this.currentConfig]
      this.$refs.configPathTree.setCurrentKey(this.configsData.configPathTreeNodeKeyMap[this.currentConfig])
      this.accordionExpandTreeNode(this.$refs.configPathTree, this.configsData.configPathTreeNodeKeyMap[this.currentConfig])
    },
    async changeConfStructTo(configName) {
      this.currentConfig = configName
      await this.changeCurrentConfStruct()
      this.configsData.stackCursor++
      this.configsData.cachedConfigStack.splice(this.configsData.stackCursor)
      this.configsData.cachedConfigStack.push(configName)
    },
    async changeBackConfStruct() {
      if (this.configsData.stackCursor > 0 && this.configsData.cachedConfigStack.length > 1) {
        this.currentConfig = this.configsData.cachedConfigStack[this.configsData.stackCursor - 1]
        await this.changeCurrentConfStruct()
        this.configsData.stackCursor--
      }
    },
    async changeForwardConfStruct() {
      if (this.configsData.stackCursor >= 0 && this.configsData.cachedConfigStack.length - 1 > this.configsData.stackCursor) {
        this.currentConfig = this.configsData.cachedConfigStack[this.configsData.stackCursor + 1]
        await this.changeCurrentConfStruct()
        this.configsData.stackCursor++
      }
    },
    formatConfStruct(pos, contextNode) {
      if (Array.isArray(pos) && typeof contextNode !== 'string') {
        var formattedContext = {
          ctxType: contextNode['context-type'],
          value: contextNode['value'],
          pos: pos,
          id: pos.toString(),
          children: [],
          delButtonShow: false,
          extendLabel: ''
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
    extendTreeNodeLabel(node) {
      var el = new ContextStructBuilder().build(node.data).genExtendLabel()
      if (el !== undefined && el !== '') {
        node.data.extendLabel = el
        return true
      } else {
        return false
      }
    },
    accordionExpandTreeNode(tree, key) {
      if (key) {
        this.accordionExpandParentNode(tree.getNode(key))
      }
    },
    accordionExpandParentNode(node) {
      if (node) {
        if (node.parent) {
          for (const childNode of node.parent.childNodes) {
            childNode.expanded = false
          }
          this.accordionExpandParentNode(node.parent)
        }
        node.expanded = true
      }
    },
    handleTreeNodeMouseEnter(node) {
      // this.$set(node.data, 'delButtonShow', true)
      node.data.delButtonShow = true
    },
    handleTreeNodeMouseLeave(node) {
      // this.$set(node.data, 'delButtonShow', false)
      node.data.delButtonShow = false
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
    handleTreeNodeModify(node, data) {
      if (data === undefined) return
      this.updateRequestData = {
        'web-server-options': {
          group_id: this.formData.value[0],
          host_id: this.formData.value[1],
          srv_name: this.formData.value[2]
        },
        'target-config-context-options': {
          position: {
            config: this.currentConfig,
            'context-pos-path': data.pos
          },
          'target-context': {
            'context-type': data.ctxType,
            'context-value': data.value
          }
        }
      }
      this.modifyTreeNodeLabel = data.label
      this.modifyTreeNodeDialogVisible = true
    },
    handleTreeNodeDelete(node, data) {
      if (data === undefined) return
      this.deleteRequestData = {
        group_id: this.formData.value[0],
        host_id: this.formData.value[1],
        srv_name: this.formData.value[2],
        config: this.currentConfig,
        'context-pos-path': data.pos
      }
      this.delTreeNodeLabel = data.label
      this.delTreeNodeDialogVisible = true
    },
    handleTreeNodeDrop(draggingNode, dropNode, dropType, ev) {
      // console.log('draggingNode.data: ' + JSON.stringify(draggingNode.data))
      if (draggingNode.data === undefined) return
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
          }
        }
      }

      switch (draggingNode.data.dragType) {
        case 'create': { // 新建上下文
          // this.$refs[draggingNode.data.ctxType + 'Creator'].openDialog()
          // component组件动态注入的组件，动态生成的ref是list
          this.$refs[draggingNode.data.ctxType + 'Creator'][0].$refs.creator.initFormDataWithTargetNode(dropNode, dropType) // 初始化注入目标节点信息
          this.$refs[draggingNode.data.ctxType + 'Creator'][0].$refs.creator.openDialog() // 打开对话框
          break
        }
        default: {
          var draggingPos = draggingNode.data.pos
          this.updateRequestData['target-config-context-options']['target-context'] = {
            'clone-context-pos': {
              config: this.currentConfig,
              'context-pos-path': draggingPos
            }
          }
          this.dragTreeNodeRadio = 'mv'
          this.dragTreeNodeDialogVisible = true
        }
      }
    },
    isCtxWithValue(data) {
      const reg = /^(events|http|include|server|stream|types)$/
      return !reg.test(data.ctxType.toLowerCase())
    },
    cancelModifyDialog() {
      this.changeConfStructTo(this.currentConfig)
      this.modifyTreeNodeDialogVisible = false
      this.isConfirmingModify = false
    },
    async confirmModifyDialog() {
      if (this.isConfirmingModify) return
      this.isConfirmingModify = true
      var res = {
        code: 7,
        msg: '放弃删除操作并离开页面'
      }
      var currentConfName = this.currentConfig
      res = await modifyCtxValue(this.updateRequestData)
      if (res.code !== 0) {
        // TODO: 错误提示框
        // console.log('cancel dialog')
        this.cancelModifyDialog()
        return
      }
      // console.log('search config struct')
      if (await this.refreshConfStruct()) {
        // 置空当前配置树节点展开状态
        this.currentTreeExpandedKeysMap[currentConfName] = []
        // console.log('change config struct to current config')
        await this.changeConfStructTo(currentConfName)
      }
      this.modifyTreeNodeDialogVisible = false
      this.isConfirmingModify = false
    },
    cancelDelDialog() {
      this.changeConfStructTo(this.currentConfig)
      this.delTreeNodeDialogVisible = false
      this.isConfirmingDel = false
    },
    async confirmDelDialog() {
      if (this.isConfirmingDel) return
      this.isConfirmingDel = true
      var res = {
        code: 7,
        msg: '放弃删除操作并离开页面'
      }
      var currentConfName = this.currentConfig
      res = await removeCtx(this.deleteRequestData)
      if (res.code !== 0) {
        // TODO: 错误提示框
        // console.log('cancel dialog')
        this.cancelDelDialog()
        return
      }
      // console.log('search config struct')
      if (await this.refreshConfStruct()) {
        // 置空当前配置树节点展开状态
        this.currentTreeExpandedKeysMap[currentConfName] = []
        // console.log('change config struct to current config')
        await this.changeConfStructTo(currentConfName)
      }
      this.delTreeNodeDialogVisible = false
      this.isConfirmingDel = false
    },
    cancelDragDialog() {
      this.changeConfStructTo(this.currentConfig)
      this.dragTreeNodeDialogVisible = false
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
      switch (this.dragTreeNodeRadio) {
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
      this.dragTreeNodeDialogVisible = false
      this.isConfirmingDrag = false
    },
    ctxCreatorComponentsMeta() {
      var meta = []
      var regExp = /^(\S+)Creator$/
      for (const name of Object.keys(this.$options.components).filter(comp => regExp.test(comp))) {
        var ctxName = regExp.exec(name)[1]
        ctxName = ctxName.charAt(0).toLowerCase() + ctxName.slice(1)
        meta.push({
          comp: name,
          refName: ctxName + 'Creator',
          key: ctxName + '-creator'
        })
      }
      return meta
    },
    getCtxBuildersCardHeight() {
      return this.$refs.ctxCreatorsCard.$el.scrollHeight
    },
    setConfigStructCardHeight() {
      var configStructCardHeight = 600 - this.getCtxBuildersCardHeight()
      this.$refs.configStructCard.$el.style.height = `${configStructCardHeight}px`
      // console.log("set config struct card's height to " + `${configStructCardHeight}px`)
    },
    handleCtxCreatorsCardClickEvent(eventCb) {
      this.ctxCreatorsCardIsCollapse = eventCb()
      this.$nextTick(() => {
        this.setConfigStructCardHeight()
      })
    },
    handleCtxCreatorDragStart(draggingNode, event) {
      // console.log('handle context creator drag start')
      this.$refs.configTree.$emit('tree-node-drag-start', event, draggingNode)
    },
    handleCtxCreatorDragEnd(draggingNode, endNode, position, event) {
      // console.log('handle context creator drag end')
      this.$refs.configTree.$emit('tree-node-drag-end', event)
    },
    async handleNewCtxCommitEvent(newCtxOpts, cb) {
      this.updateRequestData['target-config-context-options']['target-context'] = JSON.parse(JSON.stringify(newCtxOpts))
      // console.log(this.updateRequestData)
      var res = await insertNewCtx(this.updateRequestData)
      if (res.code !== 0) {
        await this.$nextTick(() => {
          this.$message({
            showClose: true,
            message: '提交新建上下文异常：' + JSON.stringify(res),
            type: 'error'
          })
        })
        await cb(false)
        return
      }
      var currentConfName = this.currentConfig
      if (await this.refreshConfStruct()) {
        // 置空当前配置树节点展开状态
        this.currentTreeExpandedKeysMap[currentConfName] = []
        this.changeConfStructTo(currentConfName)
      }
      cb(true)
    },
    resetUpdateRequest() {
      this.updateRequestData = {}
      this.changeCurrentConfStruct()
    }
  }
}

</script>

<style scoped>
.searchClass {
  padding-bottom: 0;
}
.configListClass {
  height: 600px;
  width: 100%;
  overflow-y: auto;
  overflow-x: auto;
  font-size: 16px;
}
.configStructureClass {
  margin-bottom: 20px;
  width: 100%;
  overflow-y: auto;
  overflow-x: auto;
  font-size: 16px;
}
.radio-messagebox .el-message-box__content {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}
.floating-button-container {
  position: absolute;
  top: 30px;
  right: 20px;
  bottom: 20px;
  z-index: 9999;
}
.floating-button {
  position: absolute;
  top: 0;
  height: 40px;
  bottom: 20px;
  transition: all 0.3s;
}
.floating-button:hover {
  /* 悬浮时的样式变化 */
  position: absolute;
  top: 1px;
  right: 15px;
  bottom: 15px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}
/* 隐藏单选按钮的圆点 */
::v-deep .el-radio .el-radio__inner {
  display: none;
}
</style>
