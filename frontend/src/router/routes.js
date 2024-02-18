//定义路由组件
import sendVue from '@/views/send.vue'
import receiveVue from '@/views/receive.vue'

const routes = [
    {
        path: '/',
        name: 'index',
        title: '首页',
        component: sendVue, 
    },
    {
        path: '/send',
        name: 'send',
        title: '发送文件',
        component: sendVue, 
    },
    {
        path: '/receive',
        name: 'receive',
        title: '接收文件',
        component: receiveVue,
    },
]
export default routes
