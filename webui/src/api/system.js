import request from '@/utils/request'

const systemApi = {
  State: '/v1/state'
}

export const getSystemState = () => {
  return request({
    url: systemApi.State,
    method: 'get',
    headers: {
      'Content-Type': 'application/json;charset=UTF-8'
    }
  })
}
