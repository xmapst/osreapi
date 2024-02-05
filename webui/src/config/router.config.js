// eslint-disable-next-line
import { BasicLayout } from '@/layouts'
import { bxAnaalyse } from '@/core/icons'

export const asyncRouterMap = [
  {
    path: '/',
    name: 'index',
    component: BasicLayout,
    meta: { title: 'menu.home' },
    redirect: '/dashboard',
    children: [
      // dashboard
      {
        path: '/dashboard',
        name: 'dashboard',
        component: () => import('@/views/dashboard/Index'),
        meta: { title: 'menu.dashboard', icon: bxAnaalyse }
      },
      // list
      {
        path: '/list',
        name: 'list',
        component: () => import('@/views/list/Tasks'),
        meta: { title: 'menu.task', icon: 'table' }
      }
    ]
  },
  {
    path: '*',
    redirect: '/404',
    hidden: true
  }
]

/**
 * 基础路由
 * @type { *[] }
 */
export const constantRouterMap = [
  {
    path: '/404',
    component: () => import('@/views/404')
  }
]
