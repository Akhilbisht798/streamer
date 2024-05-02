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

const ffmpegProcesses = {};

app.get("/", (req, res) => {
    res.json({message: "Hello world"})
})

app.post("/rtmpLink", (req, res) => {
	const body = req.body;
	LINK = body.link;
	res.json({ "code": 200, message: "successfully set Link" })
})

class Ffmpeg {
    constructor (link) {
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
	    link,
	];
	this.ffmpegProcess = spawn('ffmpeg', options);
	this.ffmpegProcess.stdout.on('data', (data) => {
	    console.log(`ffmpeg stdout: ${data}`);
	});

	this.ffmpegProcess.stderr.on('data', (data) => {
	    console.error(`ffmpeg stderr: ${data}`);
	});

	this.ffmpegProcess.on('close', (code) => {
	    console.log(`ffmpeg process exited with code ${code}`);
	});
    }
}

io.on('connection', socket => {
    console.log('Socket Connected', socket.id);
    socket.on('link', (data) => {
	console.log("Link set for transfering data")
	if (ffmpegProcesses[socket.id]) {
	    ffmpegProcesses[socket.id].kill('SIGINT')
	}
	ffmpegProcesses[socket.id] = new Ffmpeg(data)
    });

    socket.on('stream', stream => {
        console.log('Binary Stream Incommming...')
	const ffmpegProcess = ffmpegProcesses[socket.id];
	if (ffmpegProcess) {
	    ffmpegProcess.ffmpegProcess.stdin.write(stream, (err) => {
		console.log('Err', err)
	    })
	} else {
	    console.log("No ffmpeg process found for socket: ", socket.id);
	}
    })
    socket.on('disconnect', () => {
	console.log(`Closing socket connection: `, socket.id)
	if (ffmpegProcesses[socket.id]) {
	    ffmpegProcesses[socket.id].kill('SIGINT')
	    delete ffmpegProcesses[socket.id]
	}
    })
})

server.listen(3000, () => console.log(`HTTP Server is runnning on PORT 3000`))
