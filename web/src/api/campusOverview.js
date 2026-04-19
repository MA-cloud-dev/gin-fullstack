import service from '@/utils/request'

export const getCampusOverview = () => {
  return service({
    url: '/campusOverview/getCampusOverview',
    method: 'get'
  })
}
