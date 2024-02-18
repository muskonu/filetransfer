<script setup>
import { UploadFilled } from '@element-plus/icons-vue'
import { ref, watch } from 'vue';
import { storeToRefs } from 'pinia';
import { Search } from '@element-plus/icons-vue';
import { userinfoStore } from '@/stores/setting/index.js';
import { fileReceiveStore } from '@/stores/receive/file.js';
import { ElMessage, cascaderEmits } from 'element-plus';

const CmdRegister = 1;
const CmdAnswer = 2;
const CmdOffer = 3;
const CmdCandidate = 4;

const BadPayloadResponse = 101;
const CloseResponse = 102;
const ReceiveOffer = 103;
const ReceiveAnswer = 104;
const ReceiveCandidate = 105;

// 创建 RTCPeerConnection 对象
const pc = new RTCPeerConnection({
  iceServers: [
    {
      urls: import.meta.env.VITE_STUN_URL,
    },
  ]
})

const dataChannel = pc.createDataChannel('fileTransfer', {
  ordered: true, // 保证到达顺序
  negotiated: true, // 双向通信
  maxRetransmits: 30,
  id: 0,
})
dataChannel.binaryType = 'arraybuffer';
dataChannel.bufferedAmountLowThreshold = 128 * 1024 - 1;

dataChannel.onopen = (ev) => {
  ElMessage.success('连接已建立');
}

dataChannel.onclose = (ev) => {
  ElMessage.error('连接已关闭')
}

const fstore = fileReceiveStore();
const download = ref(false);
const { fileinfo, receiveBuffer, progress, count } = storeToRefs(fstore);
const { countIncrement, setProgress } = fstore;
var offset = 0;

dataChannel.onmessage = (event) => {
  let data = event.data;
  if (download.value === false) {
    fileinfo.value.push(JSON.parse(data));
    download.value = true;
    countIncrement();
    progress.value.push(0);
    const mergedBuffer = new ArrayBuffer(fileinfo.value[count.value - 1].size);
    receiveBuffer.value.push(mergedBuffer);
    return;
  }
  //组合arraybuffer
  let buffer = new Uint8Array(data);
  var view = new Uint8Array(receiveBuffer.value[count.value - 1], offset);
  view.set(buffer);
  offset += buffer.length;
  setProgress(Math.floor(offset / fileinfo.value[count.value - 1].size * 100), count.value - 1);
  if (offset === fileinfo.value[count.value - 1].size) {
    download.value = false;
    offset = 0;
  }
}

// 监听 RTCPeerConnection 的 onicecandidate 事件，当 ICE 服务器返回一个新的候选地址时，就会触发该事件
pc.onicecandidate = (ev) => {
  if (ev.candidate) {
    socket.send(JSON.stringify({
      command: CmdCandidate,
      source: userCode.value,
      target: remoteCode.value,
      body: ev.candidate,
    }));
    return;
  }
  inputRemoteCode.value = remoteCode.value;
};

const createOffer = async () => {
  // 创建 offer
  const offer = await pc.createOffer()
  // 设置本地描述
  await pc.setLocalDescription(offer)
}

const inputRemoteCode = ref('');
const remoteCode = ref('');
const socket = new WebSocket(import.meta.env.VITE_WEBSOCKET_URL);

socket.onopen = function () {
  socket.send(JSON.stringify({
    command: CmdRegister,
    source: userCode.value,
  }))
};

socket.onerror = function (event) {
  ElMessage({
    message: 'init error',
    type: 'error',
  })
}

socket.onmessage = function (event) {
  let info = JSON.parse(event.data);
  switch (info.code) {
    case BadPayloadResponse:
      console.log("BadPayloadResponse")
      BadPayloadResponseHandler();
      break;
    case CloseResponse:
      console.log("CloseResponse");
      CloseResponseHandler();
      break;
    case ReceiveOffer:
      console.log("ReceiveOffer");
      ReceiveOfferHandler(info.message);
      break;
    case ReceiveAnswer:
      console.log("ReceiveAnswer");
      ReceiveAnswerHandler(info.message);
      break;
    case ReceiveCandidate:
      console.log("ReceiveCandidate");
      ReceiveCandidateHandler(info.message);
      break;
  }
}

function BadPayloadResponseHandler() {
  ElMessage({
    message: '请输入正确的user code',
    type: 'error',
  })
}

function CloseResponseHandler() {
  ElMessage({
    message: '连接已关闭',
    type: 'error',
  })
}

async function ReceiveOfferHandler(val) {
  var info = JSON.parse(val);
  remoteCode.value = info.source;
  const offer = JSON.parse(info.message);
  await pc.setRemoteDescription(offer);
  const answer = await pc.createAnswer();
  await pc.setLocalDescription(answer);
  socket.send(JSON.stringify({
    command: CmdAnswer,
    source: userCode.value,
    target: info.source,
    body: pc.localDescription,
  }));
  inputValid.value = true;
}

async function ReceiveAnswerHandler(info) {
  const answer = JSON.parse(info)
  await pc.setRemoteDescription(answer)
}

async function ReceiveCandidateHandler(val) {
  const candidate = JSON.parse(val);
  await pc.addIceCandidate(new RTCIceCandidate(candidate));
}

const ustore = userinfoStore();
const { userCode } = storeToRefs(ustore);;
const inputValid = ref(false);

async function connection() {
  remoteCode.value = inputRemoteCode.value;
  if (pc.localDescription === null) {
    await createOffer()
  }
  socket.send(JSON.stringify({
    command: CmdOffer,
    source: userCode.value,
    target: remoteCode.value,
    body: pc.localDescription,
  }))
  inputValid.value = true;
}


/***********************
upload                *
                     *
*********************/

const input = ref(null);

function FileSize(number) {
  if (number < 1024) {
    return `${number} bytes`;
  } else if (number >= 1024 && number < 1048576) {
    return `${(number / 1024).toFixed(1)} KB`;
  } else if (number >= 1048576) {
    return `${(number / 1048576).toFixed(1)} MB`;
  }
}

function filesDisplay() {
  fileData.value = [];
  const curFiles = input.value.files;
  for (const file of curFiles) {
    let tmp = {};
    tmp.name = file.name;
    tmp.size = FileSize(file.size);
    fileData.value.push(tmp);
  }
}

const fileData = ref([]);

const handleDelete = (index, row) => {
  let dt = new DataTransfer();
  for (let i = 0; i < input.value.files.length; i++) {
    let file = input.value.files[i];
    if (index !== i)
      dt.items.add(file) // exclude the file. thus removing index file.
  }
  input.value.files = dt.files;
  fileData.value.splice(index, 1);
}

//发送文件
const dissend = ref(false);
const BytePerChunk = 128 * 1024;
var currentChunk = 0;
const fileReader = new FileReader();
var file;
const snedIndex = ref(0);

function readChunk() {
  let start = BytePerChunk * currentChunk;
  let end = Math.min(file.size, start + BytePerChunk);
  fileReader.readAsArrayBuffer(file.slice(start, end));
}

fileReader.onload = function () {
  dataChannel.send(fileReader.result);
  currentChunk++;
  if (BytePerChunk * currentChunk >= file.size) {
    snedIndex.value++;
    currentChunk = 0;
    fileData.value.shift();
  }
}

dataChannel.onbufferedamountlow = (event) => {
  readChunk();
}

function sendFiles() {
  snedIndex.value = 0;
  if (0 < input.value.files.length) {
    dissend.value = true;
    currentChunk = 0;
    file = input.value.files[0];
    dataChannel.send(JSON.stringify({
      name: file.name,
      size: file.size,
      user: userCode.value,
    }))
    readChunk();
  }
}

watch(snedIndex, (newIndex) => {
  if (newIndex < input.value.files.length) {
    file = input.value.files[newIndex];
    dataChannel.send(JSON.stringify({
      name: file.name,
      size: file.size,
      user: userCode.value,
    }))
    readChunk();
  } else {
    let dt = new DataTransfer();
    input.value.files = dt.files;
    dissend.value = false;
  }
})

</script>

<template>
  <div id="header">
    <el-row :gutter="10">
      <el-col :span="18"><el-input v-model="inputRemoteCode" class="code-input" placeholder="Type User Code"
          :disabled="inputValid" /></el-col>
      <el-col :span="6" style="text-align: right;"><el-button type="primary" @click="connection" :icon="Search" round
          id="row">Search</el-button></el-col>
    </el-row>
  </div>
  <input name="fileInput" ref="input" id="fileInput" type="file" placeholder="选择你的文件" style="opacity:0" multiple
    @change="filesDisplay">
  <label for="fileInput">
    <el-card shadow="hover" id="upload-container">
      <el-icon size=80 id="upload-icon"><upload-filled /></el-icon>
      <div class="el-upload__text" style="font-weight: bold;">
        click here to <em>upload</em>
      </div>
    </el-card>
  </label>
  <el-table :data="fileData" style="width: 100%" empty-text="-">
    <el-table-column prop="name" />
    <el-table-column prop="size" />
    <el-table-column>
      <template #default="scope">
        <el-button size="small" type="danger" @click="handleDelete(scope.$index, scope.row)">Delete
        </el-button>
      </template>
    </el-table-column>
  </el-table>
  <div id="send">
    <el-button color="#626aef" plain @click="sendFiles" :disabled="dissend">Send</el-button>
  </div>
</template>

<style scoped>
.demo-progress .el-progress--line {
  margin-bottom: 15px;
}

#upload-container {
  height: 200px;
}

#upload-icon {
  color: #66adf9;
}

#send {
  margin-top: 50px;
}

#row {
  padding: 5px 8px;
}
</style>