import { defineStore } from 'pinia';
import { computed, ref } from 'vue';
import { v4 as uuidv4 } from 'uuid';

export const userinfoStore = defineStore('userinfo', {
  state: () => ({
    userCode: "",
  }),
  persist: {
    storage: sessionStorage,
    beforeRestore: (ctx) => {
      console.log(`即将恢复 '${ctx.store.$id}'`)
    },
    afterRestore: (ctx) => {
      if (ctx.store.$state.userCode === '') {
        ctx.store.$state.userCode = uuidv4()
        console.log(ctx.store.$state.userCode)
        sessionStorage.setItem("userinfo",JSON.stringify({"userCode":ctx.store.$state.userCode}))
      }
    },
    debug: true,
  },
})