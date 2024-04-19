import http from 'http'
import { spawn } from 'child_process'
import express from 'express'
import { Server as SocketIO } from "socket.io";
import cors from "cors";
import dotenv from "dotenv";

dotenv.config()

const app = express();
app.use(cors())
const server = http.createServer(app);
const io = new SocketIO(server, {
    cors: {
	origin: "*",
	methods: ['GET', 'POST']
    }
})

app.get("/", (req, res) => {
    res.json({message: "Hello world"})
})

const options = [
	'-i',
	'-',
	'-c:v', 'libx264',
	'-preset', 'ultrafast',
	'-tune', 'zerolatency',
	'-r', `${25}`,
	'-g', `${25 * 2}`,
	'-keyint_min', 25,
	'-crf', '25',
	'-pix_fmt', 'yuv420p',
	'-sc_threshold', '0',
	'-profile:v', 'main',
	'-level', '3.1',
	'-c:a', 'aac',
	'-b:a', '128k',
	'-ar', 128000 / 4,
	'-f', 'flv',
	process.env.LINK,
];

const ffmpegProcess = spawn('ffmpeg', options);

ffmpegProcess.stdout.on('data', (data) => {
    console.log(`ffmpeg stdout: ${data}`);
});

ffmpegProcess.stderr.on('data', (data) => {
    console.error(`ffmpeg stderr: ${data}`);
});

ffmpegProcess.on('close', (code) => {
    console.log(`ffmpeg process exited with code ${code}`);
});

io.on('connection', socket => {
    console.log('Socket Connected', socket.id);
    socket.on('stream', stream => {
        console.log('Binary Stream Incommming...')
        ffmpegProcess.stdin.write(stream, (err) => {
            console.log('Err', err)
        })
    })
    socket.on('shutdown', () => {
	console.log(`Recived shutdown request`)
	server.close(() => console.log(`Http server closed`))
	ffmpegProcess.kill('SIGINT')
	process.exit(0)
    })
})

server.listen(3000, () => console.log(`HTTP Server is runnning on PORT 3000`))
