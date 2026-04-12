import service from '@/utils/request'

export const getCampusUserList = (params) => {
  return service({
    url: '/campusUser/getCampusUserList',
    method: 'get',
    params
  })
}

export const findCampusUser = (params) => {
  return service({
    url: '/campusUser/findCampusUser',
    method: 'get',
    params
  })
}

export const updateCampusUserStatus = (data) => {
  return service({
    url: '/campusUser/updateCampusUserStatus',
    method: 'post',
    data
  })
}
