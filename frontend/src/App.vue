<script setup>
import Navigator from '@/components/navigator.vue';
import Setting from '@/components/setting.vue';
</script>

<template>
    <div id="app">
        <el-container>
            <el-header>
                <Navigator></Navigator>
            </el-header>
            <el-main>
                <router-view v-slot="{ Component }">
                    <keep-alive>
                        <component :is="Component" />
                    </keep-alive>
                </router-view>
            </el-main>
            <el-footer>
                <Setting></Setting>
            </el-footer>
        </el-container>
    </div>
</template>

<style scoped>
#app {
    display: flex;
    overflow-x: auto;
}
</style>