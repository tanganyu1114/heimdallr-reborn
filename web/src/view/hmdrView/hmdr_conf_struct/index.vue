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
      <el-row>
        <el-input
          v-if="currentConfStruct.length > 0"
          v-model="searchRequestData.searchKeywords"
          placeholder="请输入搜索关键字"
          @keyup.enter.native="() => handleSearch()"
        >
          <el-checkbox slot="prepend" v-model="searchRequestData.onlyInCurrentConfig">仅在当前配置搜索</el-checkbox>
          <el-checkbox slot="prepend" v-model="searchRequestData.isRegExp">使用正则表达式搜索</el-checkbox>
          <i v-if="searchResponse.total >= 0" slot="suffix">{{ searchResponse.index + 1 }} / {{ searchResponse.total }} </i>
          <el-radio-group v-if="searchResponse.total >= 0" slot="suffix" v-model="searchRequestData.searchTypeRadio">
            <el-radio :label="1">重新搜索</el-radio>
            <el-radio :label="2">在当前结果中搜索</el-radio>
          </el-radio-group>
          <el-button v-if="searchResponse.total >= 0" slot="append" icon="el-icon-arrow-left" @click="handleSearchPrev">上一个</el-button>
          <el-button v-if="searchResponse.total >= 0" slot="append" icon="el-icon-arrow-right" @click="handleSearchNext">下一个</el-button>
          <el-button slot="append" icon="el-icon-search" @click="() => handleSearch()">搜索</el-button>
        </el-input>
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
                  <el-radio class="hidden-el-radio" :label="data.configPath" style="white-space: pre-wrap; font-style: normal">
                    <i v-if="currentConfig === data.configPath || currentConfig === data.fullPath" class="el-icon-document-checked" style="text-indent: -0.75em" />
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
            lazy
            :load="handleTreeNodeLazyLoad"
            @node-expand="(data)=>handleTreeNodeExpand(data)"
            @node-collapse="(data)=>handleTreeNodeCollapse(data)"
            @node-drop="handleTreeNodeDrop"
          >
            <span
              :id="data.id"
              slot-scope="{ node, data }"
              :style="handleSearchHighlight(node, data)"
              @mouseleave="() => handleTreeNodeMouseLeave(node)"
              @mouseenter="() => handleTreeNodeMouseEnter(node)"
            >
              <span :style="data.labelStyle">{{ node.label }}</span> <span v-if="extendTreeNodeLabel(node)" style="color: darkseagreen">{{ node.data.extendLabel }}</span> <span v-if="node.parent.data.ctxType !== undefined && node.parent.data.ctxType === 'include'">
                <el-button size="mini" type="info" icon="el-icon-position" @click="() => changeConfStructTo(node.label)">跳转至该配置</el-button>
              </span><span v-else v-show="node.data.contextHoverButtonsShow">
                <el-switch :value="node.data.enabled" :disabled="!configStructEditable" @change="() => handleTreeNodeEnabledStateChanges(node, data)" /> <span /> <span />
                <el-button size="mini" type="info" icon="el-icon-search" circle @click="event => handleTargetCtxSearch(event, node)" />
                <el-button size="mini" type="info" icon="el-icon-more" circle @click="event => handleTreeNodeDetailedConfigDisplay(event, node, data)" />
                <el-button v-if="configStructEditable && isCtxWithValue(data)" size="mini" type="primary" icon="el-icon-edit" circle @click="event => handleTreeNodeModify(event, node, data)" />
                <el-button v-if="configStructEditable" size="mini" type="danger" icon="el-icon-delete" circle @click="event => handleTreeNodeDelete(event, node, data)" />
              </span>
            </span>
          </el-tree>
          <div ref="floating-button" :style="floatingButtonContainerStyle">
            <el-button
              v-show="configsData.stackCursor > 0"
              :class="{ 'floating-button': true }"
              :style="{ 'right': '120px' }"
              type="primary"
              icon="el-icon-back"
              @click="changeBackConfStruct"
            />
            <el-button
              v-show="configsData.stackCursor >= 0 && configsData.cachedConfigStack.length - 1 > configsData.stackCursor"
              :class="{ 'floating-button': true }"
              :style="{ 'right': '70px' }"
              type="primary"
              icon="el-icon-right"
              @click="changeForwardConfStruct"
            />
            <WebServerReloadButton
              v-if="configStructEditable"
              :server-options="serverOptions"
              :disabled="Object.keys(serverOptions).length === 0"
              :class="{ 'floating-button': true }"
              :style="{ 'right': '20px' }"
            />
          </div>
          <el-dialog title="拖拽确认" width="25%" :visible.sync="dragTreeNodeDialogVisible">
            <p>请选择拖拽操作定义选项:</p>
            <el-radio-group v-model="dragTreeNodeRadio" size="mini">
              <el-radio label="mv">移动</el-radio>
              <el-radio label="cp">复制</el-radio>
            </el-radio-group>
            <span slot="footer" class="dialog-footer">
              <el-checkbox v-model="dragTreeNodeDisableTheTarget">禁用拖放点上下文</el-checkbox>
              <el-button @click="cancelDragDialog">取 消</el-button>
              <el-button type="primary" :loading="isConfirmingDrag" @click.stop="confirmDragDialog">确 定</el-button>
            </span>
          </el-dialog>
          <el-dialog title="上下文启用状态修改确认" width="25%" :visible.sync="enabledStateChangeTreeNodeDialogVisible">
            <p>请确认是否要{{ enabledStateChangeLabel }}该上下文[{{ enabledStateChangeTreeNodeLabel }}]</p>
            <span slot="footer" class="dialog-footer">
              <el-button @click="cancelEnabledStateChangeDialog">取 消</el-button>
              <el-button type="primary" :loading="isConfirmingDrag" @click.stop="confirmEnabledStateChangeDialog">确 定</el-button>
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
          <el-drawer
            title="配置上下文详情"
            :visible.sync="ctxDetailDrawerVisible"
            :wrapper-closable="false"
            direction="rtl"
            :before-close="handleCtxDetailDrawerClose"
          >
            <highlightjs language="nginx" :code="currentCtxText" class="hljs" />
          </el-drawer>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>

import WebServerReloadButton from './component/webServerReloadButton.vue'
import {
  changeCtxEnabledState,
  getConfStruct,
  getContextText,
  getIncludes,
  getOptions,
  insertCloneCtx,
  insertNewCtx,
  modifyCtxValue,
  moveCtx,
  removeCtx,
  searchCtxPoses
} from '@/api/hmdr_conf.js'
import CommentCreator from '@/view/hmdrView/hmdr_conf_struct/component/commentCreator.vue'
import DirectiveCreator from '@/view/hmdrView/hmdr_conf_struct/component/directiveCreator.vue'
import DirHTTPProxyPassCreator from '@/view/hmdrView/hmdr_conf_struct/component/dirHTTPProxyPassCreator.vue'
import DirStreamProxyPassCreator from '@/view/hmdrView/hmdr_conf_struct/component/dirStreamProxyPassCreator.vue'
import EventsCreator from '@/view/hmdrView/hmdr_conf_struct/component/eventsCreator.vue'
import HTTPCreator from '@/view/hmdrView/hmdr_conf_struct/component/httpCreator.vue'
import IfCreator from '@/view/hmdrView/hmdr_conf_struct/component/ifCreator.vue'
import LocationCreator from '@/view/hmdrView/hmdr_conf_struct/component/locationCreator.vue'
import ServerCreator from '@/view/hmdrView/hmdr_conf_struct/component/serverCreator.vue'
import StreamCreator from '@/view/hmdrView/hmdr_conf_struct/component/streamCreator.vue'
import UpstreamCreator from '@/view/hmdrView/hmdr_conf_struct/component/upstreamCreator.vue'
import ElCardCollapse from '@/components/ElCardCollapse.vue'

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
        this.data.configPathTreeNodeKeyMap[node.fullPath] = node.fullPath
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
        fullPath: '/' + [...pos, name].join('/'),
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
      },
      'dir_http_proxy_pass': function(data) {
        return new HTTPProxyPassCtx(data)
      },
      'dir_stream_proxy_pass': function(data) {
        return new StreamProxyPassCtx(data)
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
    this.el = ''
    this.disabledCtxExtendLabel = ''
    if (!this.data.enabled && this.data.ctxType !== undefined && this.data.ctxType !== 'comment' && this.data.ctxType !== 'inline_comment') this.disabledCtxExtendLabel = '(已禁用的上下文)'
  }
  genExtendLabel() {
    return this.disabledCtxExtendLabel + this.el
  }
}

class ServerCtx extends ContextStruct {
  constructor(data) {
    super(data)
    this.ports = []
    this.server_names = []
    this.stream_proxy = ''
    const listenRE = /^listen\s+(.*)/
    const serverNameRE = /^server_name\s+(.*)/
    for (let i = 0; i < this.data.children.length; i++) {
      if (!this.data.children[i].enabled) {
        continue
      }
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
      } else if (this.data.children[i].ctxType === 'dir_stream_proxy_pass') {
        this.stream_proxy = this.data.children[i].value + ' (' + this.data.children[i].extendLabelOriginalData + ')'
      }
    }
    this.el = '<=='
    if (this.ports.length > 0) {
      this.el += ' 侦听端口：“' + this.ports.join('”，“') + '”'
    } else {
      this.el += ' 侦听默认端口：“80/8000”？'
    }
    if (this.server_names.length > 0) {
      this.el += '，服务名：“' + this.server_names.join('”，“') + '”'
    }
    if (this.stream_proxy !== '') {
      this.el += '，代理至：“' + this.stream_proxy + '”'
    }
  }
}

class LocationCtx extends ContextStruct {
  constructor(data) {
    super(data)
    // TODO: 支持if上下文内的代理或静态资源配置
    this.http_proxy = ''
    this.root = ''
    const rootRE = /^root\s+(.*)/
    for (let i = 0; i < this.data.children.length; i++) {
      if (!this.data.children[i].enabled) {
        continue
      }
      if (this.data.children[i].ctxType === 'directive' && rootRE.test(this.data.children[i].value)) {
        this.root = rootRE.exec(this.data.children[i].value)[1]
        break
      } else if (this.data.children[i].ctxType === 'dir_http_proxy_pass') {
        this.http_proxy = this.data.children[i].value + ' (' + this.data.children[i].extendLabelOriginalData + ')'
        break
      }
    }
    this.el = '<=='
    if (this.http_proxy !== '') {
      this.el += ' 代理至：“' + this.http_proxy + '”'
    } else if (this.root !== '') {
      this.el += ' 静态资源路径：“' + this.root + '”'
    } else {
      this.el += ' 默认静态资源路径：“html”？'
    }
  }
}

class HTTPProxyPassCtx extends ContextStruct {
  constructor(data) {
    super(data)
    this.extendLabelOriginalData = data.extendLabelOriginalData
  }
  genExtendLabel() {
    return this.disabledCtxExtendLabel + '<== HTTP(s)反向代理至：' + this.extendLabelOriginalData
  }
}

class StreamProxyPassCtx extends ContextStruct {
  constructor(data) {
    super(data)
    this.extendLabelOriginalData = data.extendLabelOriginalData
  }
  genExtendLabel() {
    return this.disabledCtxExtendLabel + '<== TCP/UDP反向代理至：' + this.extendLabelOriginalData
  }
}

// import ContextPosesSearcher from '@/view/hmdrView/hmdr_conf_struct/component/contextPosesSearcher.vue'

export default {
  name: 'HmdrConfStruct',
  components: {
    WebServerReloadButton,
    CommentCreator,
    DirectiveCreator,
    DirHTTPProxyPassCreator,
    DirStreamProxyPassCreator,
    EventsCreator,
    HTTPCreator,
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
      currentCtxText: '',
      currentConfStruct: [],
      currentTreeExpandedKeysMap: {},
      enabledStateChangeRequestData: {},
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
      serverOptions: {},
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
      enabledStateChangeLabel: '启用/禁用',
      enabledStateChangeTreeNodeLabel: '未知上下文',
      enabledStateChangeTreeNodeDialogVisible: false,
      isConfirmingEnabledStateChange: false,
      modifyTreeNodeLabel: '未知上下文',
      modifyTreeNodeDialogVisible: false,
      isConfirmingModify: false,
      delTreeNodeLabel: '未知上下文',
      delTreeNodeDialogVisible: false,
      isConfirmingDel: false,
      dragTreeNodeRadio: '',
      dragTreeNodeDialogVisible: false,
      dragTreeNodeDisableTheTarget: false,
      ctxDetailDrawerVisible: false,
      isConfirmingDrag: false,
      ctxCreatorsCardIsCollapse: true,
      ctxCreatorCompMetaList: this.ctxCreatorComponentsMeta(),
      searchRequestData: {
        searchKeywords: '',
        onlyInCurrentConfig: true,
        isRegExp: false,
        searchTypeRadio: 1
      },
      searchResponse: {
        total: -1,
        index: 0,
        posList: []
      },
      floatingButtonContainerStyle: {
        position: 'absolute',
        top: '30px',
        right: '20px',
        bottom: '20px'
      }
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
      const res = await getConfStruct(this.serverOptions)
      // console.log(res)
      if (res.code === 0) {
        // console.log('配置赋值')
        this.configsData.mainConfig = res.data.config['main-config']
        this.configsData.configs = Object.keys(res.data.config.configs)
        // refresh config original fingerprints
        this.configsData.originalFingerprints = res.data['original-fingerprints']
        var pathStruct = new ConfigPathTreeStruct(res.data.config['main-config'], Object.keys(res.data.config.configs))
        this.configsData.pathsStruct = pathStruct.data.treeStruct
        this.configsData.configPathTreeNodeKeyMap = pathStruct.data.configPathTreeNodeKeyMap
        this.configsData.cachedConfigStack = []
        this.configsData.stackCursor = -1
        for (const c of Object.keys(res.data.config.configs)) {
          const configFullPath = this.configsData.configPathTreeNodeKeyMap[c]
          this.configsData.configStructs[configFullPath] = []
          for (let i = 0; i < res.data.config.configs[c].params.length; i++) {
            var childCtx = this.formatConfStruct([i], res.data.config.configs[c].params[i], configFullPath)
            if (childCtx !== {}) {
              this.configsData.configStructs[configFullPath].push(childCtx)
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
        this.serverOptions = {
          group_id: JSON.parse(JSON.stringify(this.formData.value[0])),
          host_id: JSON.parse(JSON.stringify(this.formData.value[1])),
          srv_name: JSON.parse(JSON.stringify(this.formData.value[2]))
        }
        if (await this.refreshConfStruct()) {
          this.currentTreeExpandedKeysMap = {}
          await this.changeConfStructTo(this.configsData.mainConfig)
        }
        this.searchRequestData = {
          searchKeywords: '',
          onlyInCurrentConfig: true,
          isRegExp: false,
          searchTypeRadio: 1
        }
        this.searchResponse = {
          total: -1,
          index: 0,
          posList: []
        }
      })
    },
    async changeCurrentConfStruct() {
      const p = await this.currentConfigTreeRerender()
      this.$refs.configPathTree.setCurrentKey(this.configsData.configPathTreeNodeKeyMap[this.currentConfig])
      this.accordionExpandTreeNode(this.$refs.configPathTree, this.configsData.configPathTreeNodeKeyMap[this.currentConfig])
      return p
    },
    async currentConfigTreeRerender() {
      let p
      // loading child include context
      const configFullPath = this.currentConfig
      for (let i = 0; i < this.configsData.configStructs[configFullPath].length; i++) {
        p = await this.handleParseIncludes(this.configsData.configStructs[configFullPath][i], configFullPath)
      }
      this.currentConfStruct = JSON.parse(JSON.stringify(this.configsData.configStructs[configFullPath]))
      return p
    },
    changeCurrentConfig(configName) {
      this.currentConfig = this.configsData.configPathTreeNodeKeyMap[configName]
    },
    async changeConfStructTo(configName) {
      this.changeCurrentConfig(configName)
      const p = await this.changeCurrentConfStruct()
      this.configsData.stackCursor++
      this.configsData.cachedConfigStack.splice(this.configsData.stackCursor)
      this.configsData.cachedConfigStack.push(this.configsData.configPathTreeNodeKeyMap[configName])
      return p
    },
    async changeBackConfStruct() {
      let p
      if (this.configsData.stackCursor > 0 && this.configsData.cachedConfigStack.length > 1) {
        this.changeCurrentConfig(this.configsData.cachedConfigStack[this.configsData.stackCursor - 1])
        p = await this.changeCurrentConfStruct()
        this.configsData.stackCursor--
      }
      return p
    },
    async changeForwardConfStruct() {
      let p
      if (this.configsData.stackCursor >= 0 && this.configsData.cachedConfigStack.length - 1 > this.configsData.stackCursor) {
        this.changeCurrentConfig(this.configsData.cachedConfigStack[this.configsData.stackCursor + 1])
        p = await this.changeCurrentConfStruct()
        this.configsData.stackCursor++
      }
      return p
    },
    formatConfStruct(pos, contextNode, configFullPath) {
      if (Array.isArray(pos) && typeof contextNode !== 'string') {
        var formattedContext = {
          enabled: contextNode['enabled'],
          ctxType: contextNode['context-type'],
          value: contextNode['value'],
          pos: pos,
          id: pos.toString(),
          configFullPath: configFullPath,
          children: [],
          contextHoverButtonsShow: false,
          labelStyle: {
            color: '#000000'
          },
          extendLabel: '',
          extendLabelOriginalData: ''
        }
        formattedContext.extendLabelOriginalData = this.formatExtendLabel(contextNode)
        formattedContext.label = this.toTreeLabel(formattedContext)
        formattedContext.isLeaf = this.isTreeLeaf(formattedContext)
        if ('params' in contextNode) {
          for (let i = 0; i < contextNode.params.length; i++) {
            const childPos = pos.concat(i)
            const childCtx = this.formatConfStruct(childPos, contextNode.params[i], configFullPath)
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
          id: pos.toString(),
          labelStyle: {
            color: '#000000',
            'font-style': 'normal'
          }
        }
      } else {
        return {}
      }
    },
    async handleParseIncludes(data, wbConfig) {
      if (data.ctxType === 'include' && !data.isFormatted && !data.isFormatting) {
        data.isFormatting = true
        const reqOpts = {
          group_id: this.serverOptions.group_id,
          host_id: this.serverOptions.host_id,
          srv_name: this.serverOptions.srv_name,
          config: wbConfig,
          'context-pos-path': data.pos
        }
        this.setOFP2RequestData(reqOpts)
        // console.log(this.currentTreeExpandedKeysMap[this.currentConfig].toString())
        // console.log(wbConfig)
        const res = await getIncludes(reqOpts)
        if (res.code === 0) {
          for (let i = 0; i < res.data.length; i++) {
            data.children.push(this.formatConfStruct(data.pos.concat(i), res.data[i]))
          }
          let children = this.configsData.configStructs[wbConfig]
          for (let i = 0; i < data.pos.length - 1; i++) {
            children = children[data.pos[i]].children
          }
          data.isFormatted = true
          children[data.pos[data.pos.length - 1]] = data
        }
        data.isFormatting = false
        return res.code
      }
    },
    async handleTreeNodeLazyLoad(node, resolve) {
      if (!node.data.children) return
      let p
      // console.log('expanding... ' + node.data.id)
      for (let i = 0; i < node.data.children.length; i++) {
        // p = await this.handleParseIncludes(node.data.children[i], this.currentConfig)
        p = await this.handleParseIncludes(node.data.children[i], node.data.configFullPath)
      }
      resolve(node.data.children)
      return p
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
      return !(data !== {} && typeof data !== 'string' && data.ctxType !== undefined && !['directive', 'inline_comment', 'comment', 'dir_http_proxy_pass', 'dir_stream_proxy_pass'].includes(data.ctxType))
    },
    toTreeLabel(data, node) {
      if (data) {
        if (data.ctxType) {
          switch (data.ctxType) {
            case 'directive': {
              if (data.enabled) {
                data.labelStyle.color = '#606266'
              } else {
                data.labelStyle.color = '#C0C4CC'
              }
              return data.value + ';'
            }
            case 'dir_http_proxy_pass': {
              if (data.enabled) {
                data.labelStyle.color = '#03b16b'
                return 'proxy_pass ' + data.value + ';'
              } else {
                data.labelStyle.color = '#b9d3c9'
                return 'proxy_pass ' + data.value + ';'
              }
            }
            case 'dir_stream_proxy_pass': {
              if (data.enabled) {
                data.labelStyle.color = '#c29a07'
                return 'proxy_pass ' + data.value + ';'
              } else {
                data.labelStyle.color = '#d3cbb2'
                return '# proxy_pass ' + data.value + ';'
              }
            }
            case 'inline_comment': {
              data.labelStyle.color = '#C0C4CC'
              data.labelStyle['font-style'] = 'italic'
              return '└─ # ' + data.value
            }
            case 'comment': {
              data.labelStyle.color = '#C0C4CC'
              data.labelStyle['font-style'] = 'italic'
              return '# ' + data.value
            }
            default: {
              if (data.enabled) {
                data.labelStyle.color = '#409EFF'
              } else {
                data.labelStyle.color = '#C0C4CC'
              }
              if (data.value) {
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
      node.data.contextHoverButtonsShow = true
    },
    handleTreeNodeMouseLeave(node) {
      node.data.contextHoverButtonsShow = false
    },
    handleTreeNodeExpand(data) {
      // console.log('expand: ' + data.id)
      this.configTreeExpandedKeysMapPush(this.currentConfig, data.id)
    },
    configTreeExpandedKeysMapPush(configFullPath, key) {
      if (!this.currentTreeExpandedKeysMap) this.currentTreeExpandedKeysMap = {}
      if (!this.currentTreeExpandedKeysMap[configFullPath]) this.currentTreeExpandedKeysMap[configFullPath] = []
      if (!this.currentTreeExpandedKeysMap[configFullPath].includes(key)) {
        this.currentTreeExpandedKeysMap[configFullPath].push(key)
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
    handleTreeNodeEnabledStateChanges(node, data) {
      if (data === undefined) return
      this.enabledStateChangeRequestData = {
        'web-server-options': {
          group_id: this.serverOptions.group_id,
          host_id: this.serverOptions.host_id,
          srv_name: this.serverOptions.srv_name
        },
        'target-config-context-options': {
          position: {
            config: this.currentConfig,
            'context-pos-path': data.pos
          },
          'target-context': {
            'enabled': !data.enabled
          }
        }
      }
      if (data.enabled) {
        this.enabledStateChangeLabel = '禁用'
      } else {
        this.enabledStateChangeLabel = '启用'
      }
      this.enabledStateChangeTreeNodeLabel = data.label
      this.enabledStateChangeTreeNodeDialogVisible = true
    },
    async handleTreeNodeDetailedConfigDisplay(event, node, data) {
      event.stopPropagation()
      if (data) {
        var reqOpts = {
          group_id: this.serverOptions.group_id,
          host_id: this.serverOptions.host_id,
          srv_name: this.serverOptions.srv_name,
          config: this.currentConfig,
          'context-pos-path': data.pos
        }
        this.setOFP2RequestData(reqOpts)
        var res = await getContextText(reqOpts)
        if (res.code === 0) {
          this.currentCtxText = res.data
          this.ctxDetailDrawerVisible = true
        }
      }
    },
    handleCtxDetailDrawerClose() {
      this.ctxDetailDrawerVisible = false
      this.$nextTick(() => {
        this.currentCtxText = ''
      })
    },
    handleTreeNodeModify(event, node, data) {
      event.stopPropagation()
      if (data === undefined) return
      this.updateRequestData = {
        'web-server-options': {
          group_id: this.serverOptions.group_id,
          host_id: this.serverOptions.host_id,
          srv_name: this.serverOptions.srv_name
        },
        'target-config-context-options': {
          position: {
            config: this.currentConfig,
            'context-pos-path': data.pos
          },
          'target-context': {
            'enabled': data.enabled,
            'context-type': data.ctxType,
            'context-value': data.value
          }
        }
      }
      this.modifyTreeNodeLabel = data.label
      this.modifyTreeNodeDialogVisible = true
    },
    handleTreeNodeDelete(event, node, data) {
      event.stopPropagation()
      if (data === undefined) return
      this.deleteRequestData = {
        group_id: this.serverOptions.group_id,
        host_id: this.serverOptions.host_id,
        srv_name: this.serverOptions.srv_name,
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
          group_id: this.serverOptions.group_id,
          host_id: this.serverOptions.host_id,
          srv_name: this.serverOptions.srv_name
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
          this.$refs[draggingNode.data.className][0].$refs.creator.initFormDataWithTargetNode(dropNode, dropType) // 初始化注入目标节点信息
          this.$refs[draggingNode.data.className][0].$refs.creator.openDialog() // 打开对话框
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
    cancelEnabledStateChangeDialog() {
      this.changeConfStructTo(this.currentConfig)
      this.enabledStateChangeTreeNodeDialogVisible = false
      this.enabledStateChangeTreeNodeLabel = '未知上下文'
      this.enabledStateChangeLabel = '启用/禁用'
      this.enabledStateChangeRequestData = {}
      this.isConfirmingEnabledStateChange = false
    },
    async confirmEnabledStateChangeDialog() {
      if (this.isConfirmingEnabledStateChange) return
      this.isConfirmingEnabledStateChange = true
      var res = {
        code: 7,
        msg: '放弃[' + this.enabledStateChangeTreeNodeLabel + ']上下文' + this.enabledStateChangeLabel + '操作并离开页面'
      }
      var currentConfName = this.currentConfig
      this.setOFP2EnabledStateChangeRequest()
      res = await changeCtxEnabledState(this.enabledStateChangeRequestData)
      if (res.code !== 0) {
        // TODO: 错误提示框
        // console.log('cancel dialog')
        this.cancelEnabledStateChangeDialog()
        return
      }
      // console.log('search config struct')
      if (await this.refreshConfStruct()) {
        // 置空当前配置树节点展开状态
        this.currentTreeExpandedKeysMap[currentConfName] = []
        // console.log('change config struct to current config')
        await this.changeConfStructTo(currentConfName)
      }
      this.enabledStateChangeTreeNodeDialogVisible = false
      this.isConfirmingEnabledStateChange = false
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
      this.setOFP2UpdateRequest()
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
      this.setOFP2DeleteRequest()
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
      this.$nextTick(() => {
        this.dragTreeNodeDisableTheTarget = false
      })
    },
    async confirmDragDialog() {
      if (this.isConfirmingDrag) return
      this.isConfirmingDrag = true
      var res = {
        code: 7,
        msg: '放弃拖拽操作并离开页面'
      }
      var currentConfName = this.currentConfig
      this.setOFP2UpdateRequest()
      // set `DisableTheTarget` to update request data
      this.updateRequestData['disable-the-target'] = this.dragTreeNodeDisableTheTarget
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
      this.$nextTick(() => {
        this.dragTreeNodeDisableTheTarget = false
      })
    },
    ctxCreatorComponentsMeta() {
      var meta = []
      var regExp = /^(\S+)Creator$/
      for (const name of Object.keys(this.$options.components).filter(comp => regExp.test(comp))) {
        var ctxName = regExp.exec(name)[1]
        ctxName = ctxName.charAt(0).toLowerCase() + ctxName.slice(1)
        meta.push({
          comp: name,
          refName: name,
          key: ctxName.replace(/([a-z])([A-Z])/g, '$1-$2').toLowerCase() + '-creator'
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
      var floatingButtonTop = 30 + this.getCtxBuildersCardHeight()
      this.floatingButtonContainerStyle.top = `${floatingButtonTop}px`
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
      this.setOFP2UpdateRequest()
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
        await this.changeConfStructTo(currentConfName)
      }
      cb(true)
    },
    resetUpdateRequest() {
      this.updateRequestData = {}
      this.changeCurrentConfStruct()
    },
    setOFP2RequestData(data) {
      data['original-fingerprints'] = JSON.parse(JSON.stringify(this.configsData.originalFingerprints))
    },
    setOFP2UpdateRequest() {
      this.updateRequestData['original-fingerprints'] = JSON.parse(JSON.stringify(this.configsData.originalFingerprints))
    },
    setOFP2DeleteRequest() {
      this.deleteRequestData['original-fingerprints'] = JSON.parse(JSON.stringify(this.configsData.originalFingerprints))
    },
    setOFP2EnabledStateChangeRequest() {
      this.enabledStateChangeRequestData['original-fingerprints'] = JSON.parse(JSON.stringify(this.configsData.originalFingerprints))
    },
    async handleTargetCtxSearch(event, node) {
      event.stopPropagation()
      this.searchRequestData.searchTypeRadio = 1
      const startpos = {
        config: node.data.configFullPath,
        'context-pos-path': node.data.pos
      }
      return await this.handleSearch(startpos)
    },
    async handleSearch(startpos) {
      const reqData = {
        'web-server-options': {
          group_id: this.serverOptions.group_id,
          host_id: this.serverOptions.host_id,
          srv_name: this.serverOptions.srv_name
        },
        keywords: this.searchRequestData.searchKeywords,
        'is-regexp-rule': this.searchRequestData.isRegExp,
        'is-only-in-current': this.searchRequestData.onlyInCurrentConfig,
        'starting-position-list': []
      }
      this.setOFP2RequestData(reqData)
      if (this.searchRequestData.searchTypeRadio === 2) {
        reqData['starting-position-list'] = this.searchResponse.posList
      } else {
        if (!startpos) {
          reqData['starting-position-list'][0] = {
            config: this.currentConfig,
            'context-pos-path': []
          }
        } else {
          reqData['starting-position-list'][0] = startpos
        }
      }
      const resp = await searchCtxPoses(reqData)
      if (resp.code === 0) {
        this.searchResponse.total = resp.data.length
        this.searchResponse.index = this.searchResponse.total - 1
        this.searchResponse.posList = resp.data
      } else {
        this.searchResponse.total = -1
        this.searchResponse.index = 0
        this.searchResponse.posList = []
        return resp.code
      }
      for (let i = 0; i < this.searchResponse.total; i++) {
        this.setConfigTreeDefExpandedNodeByPos(this.searchResponse.posList[i])
      }
      // this.$nextTick(async() => {
      //   return await this.changeCurrentConfStruct()
      // })
      this.$nextTick(async() => {
        return await this.handleSearchNext()
      })
    },
    setConfigTreeDefExpandedNodeByPos(pos) {
      if (!pos) return
      if (pos['context-pos-path'] && Array.isArray(pos['context-pos-path'])) {
        for (let i = 0; i < pos['context-pos-path'].length - 1; i++) {
          this.configTreeExpandedKeysMapPush(this.configsData.configPathTreeNodeKeyMap[pos.config], pos['context-pos-path'].slice(0, i + 1).toString())
        }
        // this.configTreeExpandedKeysMapPush(this.configsData.configPathTreeNodeKeyMap[pos.config], pos['context-pos-path'].slice(0, pos['context-pos-path'].length).toString())
      }
    },
    async handleSearchIndexChange(pos) {
      let p
      if (this.configsData.configPathTreeNodeKeyMap[pos.config] !== this.currentConfig) {
        p = await this.changeConfStructTo(pos.config)
      } else {
        for (let i = 0; i < pos['context-pos-path'].length - 1; i++) {
          if (!this.$refs.configTree.getNode(pos['context-pos-path'].slice(0, i + 1).toString()).expanded) {
            p = await this.changeCurrentConfStruct()
            break
            // this.$refs.configTree.getNode(pos['context-pos-path'].slice(0, i + 1).toString()).expanded = true
          }
        }
      }
      this.$nextTick(() => {
        this.$refs.configTree.setCurrentKey(pos['context-pos-path'].toString())
        // const ctxnode = this.$refs.configTree.getNode(pos['context-pos-path'].toString())
        // if (ctxnode) {
        //   this.$refs.configTree.setCurrentKey(ctxnode.data.id)
        //   this.scrollToNode(ctxnode.data.id)
        // }
      })
      this.$nextTick(() => {
        this.scrollToNode(pos['context-pos-path'].toString())
      })
      return p
    },
    async handleSearchPrev() {
      if (this.searchResponse.index > 0) {
        this.searchResponse.index--
      } else {
        this.searchResponse.index = this.searchResponse.total - 1
      }
      return await this.handleSearchIndexChange(this.searchResponse.posList[this.searchResponse.index])
    },
    async handleSearchNext() {
      if (this.searchResponse.total <= 0) {
        this.searchResponse.index = -1
        return
      }
      if (this.searchResponse.total - 1 > this.searchResponse.index) {
        this.searchResponse.index++
      } else {
        this.searchResponse.index = 0
      }
      return await this.handleSearchIndexChange(this.searchResponse.posList[this.searchResponse.index])
    },
    scrollToNode(nodeId) {
      const nodeElement = document.getElementById(nodeId)
      if (nodeElement) {
        nodeElement.scrollIntoView({ behavior: 'smooth', block: 'center' })
      }
    },
    handleSearchHighlight(node, data) {
      for (let i = 0; i < this.searchResponse.total; i++) {
        const pos = this.searchResponse.posList[i]
        if (this.currentConfig === this.configsData.configPathTreeNodeKeyMap[pos.config] && data.id === pos['context-pos-path'].toString()) {
          if (this.searchResponse.posList[this.searchResponse.index]['context-pos-path'].toString() === data.id) return 'background-color: #FFFF77'
          return 'background-color: #EEFFBB'
        }
      }
      return ''
    },
    formatExtendLabel(contextNode) {
      const proxy = []
      // if (contextNode['proxy_pass'] === undefined) return ''
      switch (contextNode['context-type']) {
        case 'dir_http_proxy_pass':
        case 'dir_stream_proxy_pass':
          if (contextNode['proxy-pass'] === undefined || contextNode['proxy-pass'].addresses === undefined || contextNode['proxy-pass'].addresses === null) return ''
          for (const address of contextNode['proxy-pass'].addresses) {
            proxy.push(address['domain-name'] + ':' + address.port.toString())
          }
          break
        default: return ''
      }
      if (proxy.length > 0) {
        return proxy.join(', ')
      }
      return ''
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
::v-deep .hidden-el-radio .el-radio__inner {
  display: none;
}
::v-deep .el-tree-node__content {
  height: auto;
  white-space: pre-wrap;
  word-wrap: break-word; /* 允许长单词或URL地址换行到下一行 */
  word-break: break-all; /* 允许在单词内换行 */
}
.hljs {
  max-height: 600px;
  width: 100%;
  overflow-y: scroll;
  overflow-x: hidden!important;
  font-size: 16px;
  white-space: pre-wrap;
  word-wrap: break-word; /* 允许长单词或URL地址换行到下一行 */
  word-break: break-all; /* 允许在单词内换行 */
}
</style>
