import service from '@/utils/request'

export const getCampusReportList = (params) => {
  return service({
    url: '/campusReport/getCampusReportList',
    method: 'get',
    params
  })
}

export const findCampusReport = (params) => {
  return service({
    url: '/campusReport/findCampusReport',
    method: 'get',
    params
  })
}

export const handleCampusReport = (data) => {
  return service({
    url: '/campusReport/handleCampusReport',
    method: 'post',
    data
  })
}

