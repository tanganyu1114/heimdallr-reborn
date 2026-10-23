import axios from 'axios' // 引入axios
import { Message } from 'element-ui'
import { store } from '@/store/index'
import context from '@/main.js'
const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_API,
  timeout: 99999
})
let acitveAxios = 0
let timer
const showLoading = () => {
  acitveAxios++
  if (timer) {
    clearTimeout(timer)
  }
  timer = setTimeout(() => {
    if (acitveAxios > 0) {
      context.$bus.emit('showLoading')
    }
  }, 400)
}

const closeLoading = () => {
  acitveAxios--
  if (acitveAxios <= 0) {
    clearTimeout(timer)
    context.$bus.emit('closeLoading')
  }
}
// http request 拦截器
service.interceptors.request.use(
  config => {
    if (!config.donNotShowLoading) {
      showLoading()
    }
    const token = store.getters['user/token']
    const user = store.getters['user/userInfo']
    config.data = JSON.stringify(config.data)
    config.headers = {
      'Content-Type': 'application/json',
      'x-token': token,
      'x-user-id': user.ID
    }
    return config
  },
  error => {
    closeLoading()
    const errorMsg = error.message || '请求配置错误'
    Message({
      showClose: true,
      message: errorMsg,
      type: 'error'
    })
    return Promise.reject(error)
  }
)

// http response 拦截器
service.interceptors.response.use(
  response => {
    closeLoading()
    if (response.headers['new-token']) {
      store.commit('user/setToken', response.headers['new-token'])
    }

    // 检查是否为二进制响应（Excel 文件下载等）
    const responseType = response.config.responseType
    if (responseType === 'arraybuffer' || responseType === 'blob') {
      // 直接返回响应对象，让调用方自行处理
      return response.data
    }

    // gjallarhorn 等接口成功时可能返回 null，直接返回空对象
    if (response.data === null || response.data === undefined) {
      return {}
    }

    // 海姆达尔后端格式：{code: 0, data: {...}, msg: "..."}
    if (response.data.code !== undefined) {
      if (response.data.code === 0 || response.headers.success === 'true') {
        return response.data
      } else {
        Message({
          showClose: true,
          message: response.data.msg || response.data.message || decodeURI(response.headers.msg || '') || '请求失败',
          type: response.headers.msgtype || 'error'
        })
        if (response.data.data && response.data.data.reload) {
          store.commit('user/LoginOut')
        }
        return response.data.msg ? response.data : response
      }
    }

    // gjallarhorn 等其他服务格式：直接返回数据（如 {total_count, page, items}）
    // 当 response.data 为 null 时（如 apply 接口成功返回），返回空对象避免报错
    return response.data || {}
  },
  error => {
    closeLoading()
    const errorMsg = error.response?.data?.msg || error.response?.data?.message || error.message || '请求失败'
    Message({
      showClose: true,
      message: errorMsg,
      type: 'error'
    })
    return Promise.reject(error)
  }
)

export default service
