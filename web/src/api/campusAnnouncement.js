import service from '@/utils/request'

export const getCampusAnnouncementList = (params) => {
  return service({
    url: '/campusAnnouncement/getCampusAnnouncementList',
    method: 'get',
    params
  })
}

export const findCampusAnnouncement = (params) => {
  return service({
    url: '/campusAnnouncement/findCampusAnnouncement',
    method: 'get',
    params
  })
}

export const createCampusAnnouncement = (data) => {
  return service({
    url: '/campusAnnouncement/createCampusAnnouncement',
    method: 'post',
    data
  })
}

export const updateCampusAnnouncement = (data) => {
  return service({
    url: '/campusAnnouncement/updateCampusAnnouncement',
    method: 'put',
    data
  })
}

export const deleteCampusAnnouncement = (data) => {
  return service({
    url: '/campusAnnouncement/deleteCampusAnnouncement',
    method: 'delete',
    data
  })
}

export const updateCampusAnnouncementStatus = (data) => {
  return service({
    url: '/campusAnnouncement/updateCampusAnnouncementStatus',
    method: 'post',
    data
  })
}

