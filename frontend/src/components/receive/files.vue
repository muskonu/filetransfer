<script setup>
import {
    Delete,
    Select,
} from '@element-plus/icons-vue';
import { onMounted, ref, watch } from 'vue';
import { storeToRefs } from 'pinia';
import { fileReceiveStore } from '@/stores/receive/file.js';

function FileSize(number) {
    if (number < 1024) {
        return `${number} bytes`;
    } else if (number >= 1024 && number < 1048576) {
        return `${(number / 1024).toFixed(1)} KB`;
    } else if (number >= 1048576) {
        return `${(number / 1048576).toFixed(1)} MB`;
    }
}

const files = ref([]);
const fstore = fileReceiveStore();
const { fileinfo, receiveBuffer, progress, count } = storeToRefs(fstore);

onMounted(() => {
    for (let i = 0; i < count.value; i++) {
        let p = progress.value[i];
        let status = "";
        if (p === 100) {
            status = "success";
        }
        files.value.push({
            name: fileinfo.value[i].name,
            size: FileSize(fileinfo.value[i].size),
            user: fileinfo.value[i].user,
            progress: p,
            status: status,
        })
    }
})

watch(count, () => {
    console.log(1);
    files.value.push({
        name: fileinfo.value[count.value - 1].name,
        size: FileSize(fileinfo.value[count.value - 1].size),
        user: fileinfo.value[count.value - 1].user,
        progress: 0,
        status: "",
    })
    console.log(files.value);
})

watch(progress, () => {
    files.value[count.value - 1].progress = progress.value[count.value - 1];
    if (progress.value[count.value - 1] === 100) {
        files.value[count.value - 1].status = "success";
        console.log("progress success");
    }
}, { deep: true })

function downloadFile(index) {
    console.log(index);
    console.log(fileinfo.value);
    const blob = new Blob([receiveBuffer.value[index]], { type: 'application/octet-stream' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = fileinfo.value[index].name;
    a.click();
    URL.revokeObjectURL(url);
}

</script>

<template>
    <el-table :data="files" style="width: 100%" empty-text="-">
        <el-table-column prop="name" label="File" />
        <el-table-column prop="size" label="Size" />
        <el-table-column prop="user" label="User" />
        <el-table-column label="operations">
            <template #default="scope">
                <el-button-group class="button-container">
                    <el-button type="primary" :icon="Select" @click="downloadFile(scope.$index)"/>
                </el-button-group>
            </template>
        </el-table-column>
        <el-table-column prop="progress" label="status">
            <template #default="scope">
                <el-progress type="circle" width=64 :percentage="scope.row.progress" :status="scope.row.status" />
            </template>
        </el-table-column>
    </el-table>
</template>

<style scoped></style>