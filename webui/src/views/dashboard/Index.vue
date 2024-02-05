<template>
  <div>
    <a-card :bordered="false">
      <a-row :gutter="24">
        <a-col :sm="6" :xs="24">
          <info :title="$t('pool.worker')" :value="pool.size" :bordered="true" />
        </a-col>
        <a-col :sm="6" :xs="24">
          <info :title="$t('pool.running')" :value="pool.running" :bordered="true" />
        </a-col>
        <a-col :sm="6" :xs="24">
          <info :title="$t('pool.waiting')" :value="pool.waiting" :bordered="true" />
        </a-col>
        <a-col :sm="6" :xs="24">
          <info :title="$t('pool.total')" :value="pool.total"/>
        </a-col>
      </a-row>
    </a-card>

    <div>
      <a-row :gutter="24" class="system_state">
        <a-col :span="12">
          <a-card
            v-if="system.os"
            class="card_item"
            :bordered="false"
            :title="$t('system.runtime')"
          >
            <div>
              <a-row :gutter="10">
                <a-col :span="12">os:</a-col>
                <a-col :span="12" v-text="system.os.go_os" />
              </a-row>
              <a-row :gutter="10">
                <a-col :span="12">cpu nums:</a-col>
                <a-col :span="12" v-text="system.os.num_cpu" />
              </a-row>
              <a-row :gutter="10">
                <a-col :span="12">compiler:</a-col>
                <a-col :span="12" v-text="system.os.compiler" />
              </a-row>
              <a-row :gutter="10">
                <a-col :span="12">go version:</a-col>
                <a-col :span="12" v-text="system.os.go_version" />
              </a-row>
              <a-row :gutter="10">
                <a-col :span="12">goroutine nums:</a-col>
                <a-col :span="12" v-text="system.os.num_goroutine" />
              </a-row>
            </div>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card
            v-if="system.disk"
            class="card_item"
            :bordered="false"
            :title="$t('system.disk')"
          >
            <div>
              <a-row :gutter="10">
                <a-col :span="12">
                  <a-row :gutter="10">
                    <a-col :span="12">total (MB)</a-col>
                    <a-col :span="12" v-text="system.disk.total_mb" />
                  </a-row>
                  <a-row :gutter="10">
                    <a-col :span="12">used (MB)</a-col>
                    <a-col :span="12" v-text="system.disk.used_mb" />
                  </a-row>
                  <a-row :gutter="10">
                    <a-col :span="12">total (GB)</a-col>
                    <a-col :span="12" v-text="system.disk.total_gb" />
                  </a-row>
                  <a-row :gutter="10">
                    <a-col :span="12">used (GB)</a-col>
                    <a-col :span="12" v-text="system.disk.used_gb" />
                  </a-row>
                </a-col>
                <a-col :span="12">
                  <a-progress
                    type="dashboard"
                    :percent="system.disk.used_percent"
                    strokeColor="butt"
                  />
                </a-col>
              </a-row>
            </div>
          </a-card>
        </a-col>
      </a-row>
      <a-row :gutter="15" class="system_state">
        <a-col :span="12">
          <a-card
            v-if="system.cpu"
            class="card_item"
            :body-style="{ height: '180px', 'overflow-y': 'scroll' }"
            :bordered="false"
            :title="$t('system.cpu')"
          >
            <div>
              <a-row :gutter="15">
                <a-col :span="12">physical number of cores:</a-col>
                <a-col :span="12" v-text="system.cpu.cores" />
              </a-row>
              <a-row v-for="(item, index) in system.cpu.cpus" :key="index" :gutter="10">
                <a-col :span="12">core {{ index }}:</a-col>
                <a-col
                  :span="12"
                ><a-progress
                  type="line"
                  :percent="+item.toFixed(0)"
                  strokeColor="butt"
                /></a-col>
              </a-row>
            </div>
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card
            v-if="system.ram"
            class="card_item"
            :bordered="false"
            :title="$t('system.ram')"
          >
            <div>
              <a-row :gutter="10">
                <a-col :span="12">
                  <a-row :gutter="10">
                    <a-col :span="12">total (MB)</a-col>
                    <a-col :span="12" v-text="system.ram.total_mb" />
                  </a-row>
                  <a-row :gutter="10">
                    <a-col :span="12">used (MB)</a-col>
                    <a-col :span="12" v-text="system.ram.used_mb" />
                  </a-row>
                  <a-row :gutter="10">
                    <a-col :span="12">total (GB)</a-col>
                    <a-col :span="12" v-text="system.ram.total_mb / 1024" />
                  </a-row>
                  <a-row :gutter="10">
                    <a-col :span="12">used (GB)</a-col>
                    <a-col
                      :span="12"
                      v-text="(system.ram.used_mb / 1024).toFixed(2)"
                    />
                  </a-row>
                </a-col>
                <a-col :span="12">
                  <a-progress
                    type="dashboard"
                    :percent="system.ram.used_percent"
                    strokeColor="butt"
                  />
                </a-col>
              </a-row>
            </div>
          </a-card>
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script>
import { onUnmounted, ref } from 'vue'
import { getSystemState } from '@/api/system'
import { getPoolState } from '@/api/pool'
import { Radar } from '@/components'
import Info from './components/Info'

const timer = ref(null)

export default {
  name: 'Workplace',
  components: {
    Radar,
    Info
  },
  data () {
    return {
      pool: {
        size: 0,
        total: 0,
        running: 0,
        waiting: 0
      },
      system: {
        os: {
          go_os: 'unknown',
          num_cpu: 0,
          compiler: 'unknown',
          go_version: 'unknown',
          num_goroutine: 0
        },
        ram: {
          used_mb: 0,
          total_mb: 0,
          used_percent: 0
        },
        cpu: {
          cores: 0,
          cpus: [
              0
          ]
        },
        disk: {
          used_mb: 0,
          used_gb: 0,
          total_mb: 0,
          total_gb: 0,
          used_percent: 0
        }
      }
    }
  },
  computed: {
  },
  created () {
    this.reload()
  },
  mounted () {
    timer.value = setInterval(() => {
      this.reload()
    }, 1000 * 3)

    onUnmounted(() => {
      clearInterval(timer.value)
      timer.value = null
    })
  },
  methods: {
    async reload() {
      await getSystemState().then(res => {
        this.system = res.data
      }).catch(err => {
        this.$message.error(err.message)
      })

      await getPoolState().then(res => {
        this.pool = res.data
      }).catch(err => {
        this.$message.error(err.message)
      })
    }
  }
}
</script>

<style lang="less" scoped>
@import './Index.less';

.system_state {
  padding: 10px;
}

.card_item {
  height: 280px;
}

.project-list {
  .card-title {
    font-size: 0;

    a {
      color: rgba(0, 0, 0, 0.85);
      margin-left: 12px;
      line-height: 24px;
      height: 24px;
      display: inline-block;
      vertical-align: top;
      font-size: 14px;

      &:hover {
        color: #1890ff;
      }
    }
  }

  .card-description {
    color: rgba(0, 0, 0, 0.45);
    height: 44px;
    line-height: 22px;
    overflow: hidden;
  }

  .project-item {
    display: flex;
    margin-top: 8px;
    overflow: hidden;
    font-size: 12px;
    height: 20px;
    line-height: 20px;

    a {
      color: rgba(0, 0, 0, 0.45);
      display: inline-block;
      flex: 1 1 0;

      &:hover {
        color: #1890ff;
      }
    }

    .datetime {
      color: rgba(0, 0, 0, 0.25);
      flex: 0 0 auto;
      float: right;
    }
  }

  .ant-card-meta-description {
    color: rgba(0, 0, 0, 0.45);
    height: 44px;
    line-height: 22px;
    overflow: hidden;
  }
}

.item-group {
  padding: 20px 0 8px 24px;
  font-size: 0;

  a {
    color: rgba(0, 0, 0, 0.65);
    display: inline-block;
    font-size: 14px;
    margin-bottom: 13px;
    width: 25%;
  }
}

.members {
  a {
    display: block;
    margin: 12px 0;
    line-height: 24px;
    height: 24px;

    .member {
      font-size: 14px;
      color: rgba(0, 0, 0, 0.65);
      line-height: 24px;
      max-width: 100px;
      vertical-align: top;
      margin-left: 12px;
      transition: all 0.3s;
      display: inline-block;
    }

    &:hover {
      span {
        color: #1890ff;
      }
    }
  }
}

.mobile {
  .project-list {
    .project-card-grid {
      width: 100%;
    }
  }

  .more-info {
    border: 0;
    padding-top: 16px;
    margin: 16px 0 16px;
  }

  .headerContent .title .welcome-text {
    display: none;
  }
}
</style>
