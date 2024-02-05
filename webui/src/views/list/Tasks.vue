<template>
  <div>
    <a-list
      rowKey="id"
      :grid="{gutter: 6, xxl: 6, xl: 4, lg: 4, md: 3, sm: 2, xs: 1}"
      :dataSource="dataSource"
      class="card-list"
    >
      <a-list-item slot="renderItem" slot-scope="item">
        <template>
          <a-card size="small" :hoverable="true">
            <div @click="detailTask(item.id)">
              <a-card-meta>
                <a slot="title">{{ item.id }}</a>
                <a-avatar class="card-avatar" slot="avatar" :src="naUrl" size="large"/>
                <div class="meta-content" slot="description">{{ item.msg }}</div>
              </a-card-meta>
              <a-row>
                <a-col :span="7">{{ $t('task.state') }}:</a-col>
                <a-col :span="17" v-text="$t('task.'+item.state)" />
              </a-row>
              <a-row>
                <a-col :span="7">{{ $t('task.code') }}:</a-col>
                <a-col :span="17" v-text="item.code" />
              </a-row>
              <a-row>
                <a-col :span="7">{{ $t('task.st') }}:</a-col>
                <a-col :span="17" v-text="item.times.st" />
              </a-row>
              <a-row>
                <a-col :span="7">{{ $t('task.et') }}:</a-col>
                <a-col :span="17" v-text="item.times.et" />
              </a-row>
            </div>
            <template class="ant-card-actions" slot="actions">
              <a v-if="[1, 2].includes(item.state)" @click="managerTask(item.id, 'pause')">{{ $t('task.pause') }}</a>
              <a v-if="[3].includes(item.state)" @click="managerTask(item.id, 'resume')">{{ $t('task.resume') }}</a>
              <a v-if="[1, 2, 3].includes(item.state)" @click="managerTask(item.id, 'kill')">{{ $t('task.kill') }}</a>
              <a @click="detailTask(item.id)">{{ $t('task.more') }}</a>
            </template>
          </a-card>
        </template>
      </a-list-item>
    </a-list>
  </div>
</template>

<script>
import imgUrl from '@/assets/bx-analyse.png'
import { getTaskList, managerTask } from '@/api/task'
import { onUnmounted, ref } from 'vue'
import TaskDetail from '@/views/list/modules/TaskDetail'

const timer = ref(null)

export default {
  name: 'CardList',
  components: {
    TaskDetail
  },
  data () {
    return {
      naUrl: '',
      dataSource: []
    }
  },
  created () {
    this.naUrl = imgUrl
    this.reload()
  },
  mounted() {
    timer.value = setInterval(() => {
      this.reload()
    }, 1000)

    onUnmounted(() => {
      clearInterval(timer.value)
      timer.value = null
    })
  },
  methods: {
    reload () {
      getTaskList().then(res => {
        if (res.data !== null) {
          this.dataSource = []
          for (let i = 0; i < res.data.tasks.length; i++) {
            const task = res.data.tasks[i]
            res.data.tasks[i].msg = task.msg.length > 100 ? task.msg.substring(0, 100) + '...' : task.msg
          }
          this.dataSource = res.data.tasks
        }
      }).catch(err => {
        this.$message.error(err.message)
      })
    },
    managerTask (task, action) {
      this.$confirm({
        title: this.$t('task.operate'),
        content: `${action} ${task}`,
        onOk: () => {
          managerTask(task, action).then(res => {
            this.$message.info(`${action} ${task} ${this.$t('success')}`)
          }).catch(err => {
            this.$message.error(err.message)
          })
          this.reload()
        }
      })
    },
    detailTask (task) {
      this.$dialog(TaskDetail,
        // component props
        {
          task,
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
          title: this.$t('task.detail'),
          width: '90%',
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
