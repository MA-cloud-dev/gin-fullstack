import service from '@/utils/request'

export const getCampusProductList = (params) => {
  return service({
    url: '/campusProduct/getCampusProductList',
    method: 'get',
    params
  })
}

export const findCampusProduct = (params) => {
  return service({
    url: '/campusProduct/findCampusProduct',
    method: 'get',
    params
  })
}

export const updateCampusProductStatus = (data) => {
  return service({
    url: '/campusProduct/updateCampusProductStatus',
    method: 'post',
    data
  })
}
