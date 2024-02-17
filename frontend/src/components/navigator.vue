<template>
  <el-tabs class="demo-tabs">
    <el-tab-pane>
      <template #label>
        <router-link to="/send">
          <span class="custom-tabs-label">
            <el-icon>
              <Share />
            </el-icon>
            <span>Send</span>
          </span>
        </router-link>
      </template>
    </el-tab-pane>

    <el-tab-pane>
      <template #label>
        <router-link to="/receive">
          <span class="custom-tabs-label">
            <el-icon>
              <MessageBox />
            </el-icon>
            <span>Receive</span>
          </span>
        </router-link>
      </template>
    </el-tab-pane>

    <!-- <el-tab-pane>
      <template #label>
        <router-link to="/setting">
          <span class="custom-tabs-label">
            <el-icon>
              <Setting />
            </el-icon>
            <span>Setting</span>
          </span>
        </router-link>
      </template>
    </el-tab-pane> -->
  </el-tabs>
</template>

<script lang="ts" setup>
import { Share,Setting,MessageBox } from '@element-plus/icons-vue'
</script>
<style>
.demo-tabs .custom-tabs-label .el-icon {
  vertical-align: middle;
  color: #3f7fed;
  font-size: 32px;
  font-weight: 600;
}

.demo-tabs .custom-tabs-label span {
  vertical-align: middle;
  margin-left: 4px;
  padding: 16px;
  color: #4f5053;
  font-size: 20px;
  font-weight: 600;
}

a{
  text-decoration: none;
}

</style>