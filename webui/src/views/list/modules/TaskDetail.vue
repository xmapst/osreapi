<template>
  <div>
    <a-list
      rowKey="name"
      :grid="{gutter: 6, xxl: 6, xl: 4, lg: 4, md: 3, sm: 2, xs: 1}"
      :dataSource="dataSource"
      class="card-list"
    >
      <a-list-item slot="renderItem" slot-scope="item">
        <template>
          <a-card size="small" :hoverable="true">
            <div @click="detailTaskStep(task, item.name)">
              <a-card-meta>
                <a slot="title">{{ item.name }}</a>
                <a-avatar class="card-avatar" slot="avatar" :src="naUrl" size="large"/>
                <div class="meta-content" slot="description">
                  {{ item.msg }}
                </div>
              </a-card-meta>
              <a-row>
                <a-col :span="7">{{ i18n.t('task.state') }}:</a-col>
                <a-col :span="17" v-text="i18n.t('task.'+item.state)" />
              </a-row>
              <a-row>
                <a-col :span="7">{{ i18n.t('task.code') }}:</a-col>
                <a-col :span="17" v-text="item.code" />
              </a-row>
              <a-row>
                <a-col :span="7">{{ i18n.t('task.st') }}:</a-col>
                <a-col :span="17" v-text="item.times.st" />
              </a-row>
              <a-row>
                <a-col :span="7">{{ i18n.t('task.et') }}:</a-col>
                <a-col :span="17" v-text="item.times.et" />
              </a-row>
            </div>
            <template class="ant-card-actions" slot="actions">
              <a v-if="item.state === 2 && [1001, 1004, 1005].includes(code)" @click="managerTaskStep(task, item.name, 'pause')">{{ i18n.t('task.pause') }}</a>
              <a v-if="item.state === 3 && [1001, 1004, 1005].includes(code)" @click="managerTaskStep(task, item.name, 'resume')">{{ i18n.t('task.resume') }}</a>
              <a v-if="[1, 2, 3].includes(item.state) && [1001, 1004, 1005].includes(code)" @click="managerTaskStep(task, item.name, 'kill')">{{ i18n.t('task.kill') }}</a>
              <a @click="detailTaskStep(task, item.name)">{{ i18n.t('task.more') }}</a>
            </template>
          </a-card>
        </template>
      </a-list-item>
    </a-list>
  </div>
</template>

<script>
import imgUrl from '@/assets/bx-analyse.png'
import i18n from '@/locales/index'
import { getTaskDetail, managerTaskStep } from '@/api/task'
import TaskStepDetail from './TaskStepDetail'
import { onUnmounted, ref } from 'vue'
import { Ellipsis } from '@/components'

const timer = ref(null)

export default {
  components: {
    Ellipsis
  },
  name: 'TaskDetail',
  props: {
    task: {
      type: null,
      default: null
    }
  },
  data () {
    return {
      naUrl: '',
      i18n: i18n,
      code: 0,
      dataSource: []
    }
  },
  created () {
    this.naUrl = imgUrl
    this.reload()
  },
  mounted () {
    timer.value = setInterval(() => {
      this.reload()
    }, 1000)

    onUnmounted(() => {
      clearInterval(timer.value)
      timer.value = null
    })
  },
  methods: {
    onOk () {
      console.log('监听了 modal ok 事件')
      return new Promise(resolve => {
        resolve(true)
      })
    },
    onCancel () {
      console.log('监听了 modal cancel 事件')
      return new Promise(resolve => {
        resolve(true)
      })
    },
    reload () {
      getTaskDetail(this.task).then(res => {
        this.code = res.code
        if (res.data !== null) {
          this.dataSource = []
          this.dataSource = res.data
        } else {
          this.$message.error(res.message)
        }
      }).catch(err => {
        this.$message.error(err.message)
      })
    },
    managerTaskStep (task, step, action) {
      this.$confirm({
        title: i18n.t('task.operate'),
        content: `${action} ${task} ${step}`,
        onOk: () => {
          managerTaskStep(this.task, step, action).then(res => {
            this.$message.info(`${action} ${this.task} ${step} ${i18n.t('success')}`)
          }).catch(err => {
            this.$message.error(err.message)
          })
          this.reload()
        }
      })
    },
    detailTaskStep (task, step) {
      this.$dialog(TaskStepDetail,
          // component props
          {
            task,
            step,
            on: {
              ok () {
                console.log('ok 回调')
              },
              cancel () {
                console.log('cancel 回调')
              },
              close () {
                console.log('modal close 回调')
              }
            }
          },
          // modal props
          {
            title: i18n.t('task.log'),
            width: '70%',
            centered: false,
            maskClosable: true,
            footer: ''
          }
      )
    }
  }
}
</script>

<style lang="less" scoped>
@import "~@/components/index.less";

.card-list {
  :deep(.ant-card-body:hover) {
    .ant-card-meta-title>a {
      color: @primary-color;
    }
  }

  :deep(.ant-card-meta-title) {
    margin-bottom: 12px;

    &>a {
      display: inline-block;
      max-width: 100%;
      color: rgba(0,0,0,.85);
    }
  }

  :deep(.meta-content) {
    position: relative;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    height: 64px;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;

    margin-bottom: 1em;
  }
}

.card-avatar {
  width: 48px;
  height: 48px;
  border-radius: 48px;
}

.ant-card-actions {
  background: #f7f9fa;

  li {
    float: left;
    text-align: center;
    margin: 12px 0;
    color: rgba(0, 0, 0, 0.45);
    width: 50%;

    &:not(:last-child) {
      border-right: 1px solid #e8e8e8;
    }

    a {
      color: rgba(0, 0, 0, .45);
      line-height: 22px;
      display: inline-block;
      width: 100%;
      &:hover {
        color: @primary-color;
      }
    }
  }
}

.new-btn {
  background-color: #fff;
  border-radius: 2px;
  width: 100%;
  height: 188px;
}

</style>
