<template>
    <div id="codep">
        <el-row :gutter="10">
        <el-col :span="24">Your Code:</el-col>
    </el-row>
    <el-row :gutter="10">
        <el-col :span="20">
                <el-input v-model="userCode" disabled />
        </el-col>
        <el-col :span="4" style="text-align: right;">
            <el-button type="primary" :icon="CopyDocument" round id="row" v-clipboard:copy="userCode"
                v-clipboard:success="onCopySuccess" v-clipboard:error="onCopyError">Copy
            </el-button>
        </el-col>
    </el-row>
    </div>
    
</template>
  
<script setup>
import { ElMessage } from 'element-plus';
import { ref } from 'vue';
import { storeToRefs } from 'pinia'
import { CopyDocument } from '@element-plus/icons-vue';
import { userinfoStore } from '@/stores/setting/index.js';

const store = userinfoStore()
const { userCode } = storeToRefs(store)

const onCopySuccess = () => {
    ElMessage.success('复制成功')
}

const onCopyError = () => {
    ElMessage.error('复制失败')
}
</script>

<style scoped>
#row {
    padding: 5px 8px;
}

#codep {
    color: black;
    text-align: center;
    box-shadow: 0px -20px 20px -20px gray;
}

</style>@/stores/setting/setting.js