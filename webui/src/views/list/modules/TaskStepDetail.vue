<template>
  <div>
    <a-textarea
      id="textarea_id"
      :auto-size="{ minRows: 30, maxRows: 30 }"
      showCount
      readOnly
      style="background: #000000; color: #ffffff; resize: none"
      v-model="dataSource"
    ></a-textarea>
  </div>
</template>

<script>
import { getTaskStepDetail } from '@/api/task'
import { onUnmounted, ref } from 'vue'

const timer = ref(null)

export default {
  name: 'TaskStepDetail',
  props: {
    task: {
      type: null,
      default: null
    },
    step: {
      type: null,
      default: null
    }
  },
  data () {
    return {
      dataLength: 0,
      dataSource: ''
    }
  },
  created () {
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
      getTaskStepDetail(this.task, this.step).then(res => {
        let dataSource = res.message
        let dataLength = 1
        if (res.data !== null) {
          const columns = []
          for (let i = 0; i < res.data.length; i++) {
            const log = res.data[i]
            columns.push(log.content)
          }
          dataSource = columns.join(`\n`)
          dataLength = res.data.length
        }
        if (dataLength !== this.dataLength) {
          this.dataLength = dataLength
          this.dataSource = dataSource
          this.$nextTick(() => {
            setTimeout(() => {
              const textarea = document.getElementById('textarea_id')
              textarea.scrollTop = textarea.scrollHeight
            }, 13)
          })
        }
      }).catch(err => {
        this.$message.error(err.message)
      })
    }
  }
}
</script>

<style lang="less" scoped>
@import "~@/components/index.less";
textarea,textarea.ant-input:hover,textarea:focus{
  border: 1px solid #DAE2F3;
  -webkit-box-shadow: none;
  box-shadow: none;
}
</style>
