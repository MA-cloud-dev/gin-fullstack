import service from '@/utils/request'

export const getCampusOperationLogList = (params) => {
  return service({
    url: '/campusOperationLog/getCampusOperationLogList',
    method: 'get',
    params
  })
}

export const findCampusOperationLog = (params) => {
  return service({
    url: '/campusOperationLog/findCampusOperationLog',
    method: 'get',
    params
  })
}
