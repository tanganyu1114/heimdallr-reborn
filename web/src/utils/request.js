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
    Message({
      showClose: true,
      message: error,
      type: 'error'
    })
    return error
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

    if (response.data.code === 0 || response.headers.success === 'true') {
      return response.data
    } else {
      Message({
        showClose: true,
        message: response.data.msg || decodeURI(response.headers.msg),
        type: response.headers.msgtype || 'error'
      })
      if (response.data.data && response.data.data.reload) {
        store.commit('user/LoginOut')
      }
      return response.data.msg ? response.data : response
    }
  },
  error => {
    closeLoading()
    Message({
      showClose: true,
      message: error,
      type: 'error'
    })
    return error
  }
)

export default service
