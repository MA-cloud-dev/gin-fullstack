import service from '@/utils/request'

export const getCampusCategoryList = (params) => {
  return service({
    url: '/campusCategory/getCampusCategoryList',
    method: 'get',
    params
  })
}

export const findCampusCategory = (params) => {
  return service({
    url: '/campusCategory/findCampusCategory',
    method: 'get',
    params
  })
}

export const createCampusCategory = (data) => {
  return service({
    url: '/campusCategory/createCampusCategory',
    method: 'post',
    data
  })
}

export const updateCampusCategory = (data) => {
  return service({
    url: '/campusCategory/updateCampusCategory',
    method: 'put',
    data
  })
}

export const updateCampusCategoryStatus = (data) => {
  return service({
    url: '/campusCategory/updateCampusCategoryStatus',
    method: 'post',
    data
  })
}
