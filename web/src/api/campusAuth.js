import service from '@/utils/request'

export const getCampusAuthList = (params) => {
  return service({
    url: '/campusAuth/getCampusAuthList',
    method: 'get',
    params
  })
}

export const findCampusAuth = (params) => {
  return service({
    url: '/campusAuth/findCampusAuth',
    method: 'get',
    params
  })
}

export const reviewCampusAuth = (data) => {
  return service({
    url: '/campusAuth/reviewCampusAuth',
    method: 'post',
    data
  })
}

export const revokeCampusAuth = (data) => {
  return service({
    url: '/campusAuth/revokeCampusAuth',
    method: 'post',
    data
  })
}
