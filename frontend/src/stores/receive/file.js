import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

export const fileReceiveStore = defineStore('fileReceive', {
  state: () => ({
    fileinfo: [],
    receiveBuffer: [],
    progress: [],
    count: 0,
  }),
  actions: {
    countIncrement() {
      this.count+=1;
    },
    setProgress(data,index) {
      this.progress[index] = data;
    },
  }
})