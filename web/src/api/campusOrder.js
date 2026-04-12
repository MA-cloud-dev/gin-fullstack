import service from '@/utils/request'

export const getCampusOrderList = (params) => {
  return service({
    url: '/campusOrder/getCampusOrderList',
    method: 'get',
    params
  })
}

export const findCampusOrder = (params) => {
  return service({
    url: '/campusOrder/findCampusOrder',
    method: 'get',
    params
  })
}

