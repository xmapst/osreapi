import request from '@/utils/request'

const poolApi = {
  State: '/v1/pool'
}

export const getPoolState = () => {
  return request({
    url: poolApi.State,
    method: 'get',
    headers: {
      'Content-Type': 'application/json;charset=UTF-8'
    }
  })
}
