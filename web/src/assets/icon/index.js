import Vue from 'vue'
import SvgIcon from '@/components/svgicon'// svg component

// register globally
// eslint-disable-next-line
Vue.component('svg-icon', SvgIcon)

const req = require.context('./svg', false, /\.svg$/)
const requireAll = requireContext => requireContext.keys().map(requireContext)
requireAll(req)
// TODO:解决icon无法使用的问题
