import service from '@/utils/request'

export const getCampusAdminStaffList = (params) => {
  return service({
    url: '/campusAdminStaff/getCampusAdminStaffList',
    method: 'get',
    params
  })
}

export const findCampusAdminStaff = (params) => {
  return service({
    url: '/campusAdminStaff/findCampusAdminStaff',
    method: 'get',
    params
  })
}

export const updateCampusAdminStaffStatus = (data) => {
  return service({
    url: '/campusAdminStaff/updateCampusAdminStaffStatus',
    method: 'post',
    data
  })
}
