import request from '@/utils/request'

const taskApi = {
  Task: '/v1/task'
}

export const getTaskList = () => {
  return request({
    url: taskApi.Task,
    method: 'get',
    headers: {
      'Content-Type': 'application/json;charset=UTF-8'
    }
  })
}

export const getTaskDetail = (task) => {
  return request({
    url: taskApi.Task + '/' + task,
    method: 'get',
    headers: {
      'Content-Type': 'application/json;charset=UTF-8'
    }
  })
}

export const getTaskStepDetail = (task, step) => {
  return request({
    url: taskApi.Task + '/' + task + '/step/' + step,
    method: 'get',
    headers: {
      'Content-Type': 'application/json;charset=UTF-8'
    }
  })
}

export const managerTask = (task, action) => {
  return request({
    url: taskApi.Task + '/' + task,
    method: 'put',
    params: {
      action: action
    },
    headers: {
      'Content-Type': 'application/json;charset=UTF-8'
    }
  })
}

export const managerTaskStep = (task, step, action) => {
  return request({
    url: taskApi.Task + '/' + task + '/step/' + step,
    method: 'put',
    params: {
      action: action
    },
    headers: {
      'Content-Type': 'application/json;charset=UTF-8'
    }
  })
}
