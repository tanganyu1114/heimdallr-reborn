import service from '@/utils/request'

// @Tags ProxyClaim
// @Summary 获取代理申请列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param page query integer false "页码（从 1 开始，默认: 1）"
// @Param page_size query integer false "每页数量（默认: 20，最大: 100）"
// @Param status query string false "按申请状态过滤"
// @Param applicant query string false "按申请人名称过滤"
// @Param oa_id query string false "按 OA ID 过滤（支持模糊匹配）"
// @Success 200 {string} string "{"code":0,"message":"success","data":{"total_count":100,"page":1,"page_size":20,"items":[]}}"
// @Router /gjallarhorn/v1/proxy-claims [get]
export const listProxyClaims = (params) => {
  return service({
    url: '/gjallarhorn/v1/proxy-claims',
    method: 'get',
    params
  })
}

// @Tags ProxyClaim
// @Summary 获取代理申请
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param oaID path string true "代理声明的唯一标识符"
// @Success 200 {string} string "{"code":0,"message":"success","data":{}}"
// @Router /gjallarhorn/v1/proxy-claims/:oaID [get]
export const getProxyClaim = (oaID) => {
  return service({
    url: `/gjallarhorn/v1/proxy-claims/${oaID}`,
    method: 'get'
  })
}

// @Tags ProxyClaim
// @Summary 获取代理申请详情（包含环境详情）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param oaID path string true "代理声明的唯一标识符"
// @Success 200 {string} string "{"code":0,"message":"success","data":{}}"
// @Router /gjallarhorn/v1/proxy-claims/:oaID/detail [get]
export const getProxyClaimDetail = (oaID) => {
  return service({
    url: `/gjallarhorn/v1/proxy-claims/${oaID}/detail`,
    method: 'get'
  })
}

// @Tags ProxyClaim
// @Summary 提交代理申请
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ApplyProxyClaimRequest true "创建或更新代理声明"
// @Success 200 {string} string "{"code":0,"message":"success","data":null}"
// @Router /gjallarhorn/v1/proxy-claims/apply [post]
export const applyProxyClaim = (data) => {
  return service({
    url: '/gjallarhorn/v1/proxy-claims/apply',
    method: 'post',
    data
  })
}

// @Tags ProxyClaim
// @Summary 审批代理申请
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param oaID path string true "待审批的代理声明 OA ID"
// @Success 200 {string} string "{"code":0,"message":"success","data":null}"
// @Router /gjallarhorn/v1/proxy-claims/:oaID/approve [post]
export const approveProxyClaim = (oaID) => {
  return service({
    url: `/gjallarhorn/v1/proxy-claims/${oaID}/approve`,
    method: 'post'
  })
}

// @Tags ProxyClaim
// @Summary 迁移代理申请
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.MigrateProxyClaimRequest true "迁移代理声明请求"
// @Success 200 {string} string "{"code":0,"message":"success","data":null}"
// @Router /gjallarhorn/v1/proxy-claims/migrate [post]
export const migrateProxyClaim = (data) => {
  return service({
    url: '/gjallarhorn/v1/proxy-claims/migrate',
    method: 'post',
    data
  })
}
