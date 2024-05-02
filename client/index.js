const video = document.getElementById("stream")
const webcamBtn = document.getElementById("webcam")
const screenShareBtn = document.getElementById("screen")
const liveBtn = document.getElementById("go-live")

const socket = io("")

webcamBtn.addEventListener("click", () => {
    getWebcamStream()
})

screenShareBtn.addEventListener("click", () => {
    getScreenStream()
})

liveBtn.addEventListener("click", () => {
    sendStream()
})

const state = { webcam: null, screen: null, stream: null }


async function getWebcamStream() {
    if (state.webcam !== null) {
        state.webcam.getTracks().forEach(track => track.stop());
        state.webcam = null
        if (state.screen !== null) {
            video.srcObject = state.screen
            state.stream = state.screen
        } else {
            video.srcObject = null
            state.stream = null
        }
        return
    }
    const webcamStream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
    state.webcam = webcamStream
    if (state.screen === null) {
        video.srcObject = webcamStream
        state.stream = webcamStream
        return
    }
    mergeStream()
}

async function getScreenStream() {
    if (state.screen !== null) {
        state.screen.getTracks().forEach(track => track.stop());
        state.screen = null
        if (state.webcam !== null) {
            video.srcObject = state.webcam
            state.stream = state.webcam
        } else {
            video.srcObject = null
            state.stream = null
        }
        return
    }
    const screenStream = await navigator.mediaDevices.getDisplayMedia({ video: true })
    state.screen = screenStream
    changeTracksHeight()

    if (state.webcam === null) {
        video.srcObject = screenStream
        state.stream = screenStream
        return
    }
    mergeStream()
}


async function mergeStream() {
    const merger = new VideoStreamMerger()

    merger.addStream(state.screen, {
        x: 0, 
        y: 0,
        width: merger.width,
        height: merger.height,
        mute: true,
    })
    merger.addStream(state.webcam, {
        x: 0,
        y: merger.height - 100,
        width: 150,
        height: 150,
        mute: false,
    });
    merger.start()

    video.srcObject = merger.result;
    state.stream = merger.result;
}

async function sendStream() {

    const mediaRecorder = new MediaRecorder(state.screen, {
        audioBitsPerSecond: 128000,
        videoBitsPerSecond: 2500000,
        framerate: 25
    });
    mediaRecorder.ondataavailable = e => {
        socket.emit("stream", e.data)
    }

    mediaRecorder.start(25)
}

