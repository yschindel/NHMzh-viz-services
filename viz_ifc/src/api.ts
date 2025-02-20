import express from "express";
import multer from "multer";
import cors from "cors";
import { saveToMinIO, minioClient, createFileName } from "./minio";
import { sendKafkaMessage } from "./kafka";
import { producer } from "./index";

const router = express.Router();
const upload = multer({ storage: multer.memoryStorage() });

let project: string;
let timestamp: string;
let fileName: string;
let location: string;

// Enable CORS
router.use(cors());

// POST endpoint to upload IFC file
// @ts-ignore
router.post("/upload", upload.single("file"), async (req, res) => {
	try {
		if (!req.file) {
			return res.status(400).json({ error: "No file uploaded" });
		}

		project = req.body.project || "default";
		timestamp = req.body.timestamp || new Date().toISOString();
		fileName = req.file.originalname;

		// Create a unique filename
		location = createFileName(project, req.file.originalname, timestamp);

		const metadata = {
			"X-Amz-Meta-Project-Name": project,
			"X-Amz-Meta-Filename": req.file.originalname,
			"X-Amz-Meta-Created-At": timestamp,
		};
		// Save to MinIO
		await saveToMinIO(minioClient, process.env.MINIO_IFC_BUCKET || "ifc-files", location, req.file.buffer, metadata);

		res.status(200).json({
			message: "File uploaded successfully",
			location,
			project,
			fileName,
			timestamp,
		});
	} catch (error) {
		console.error("Error uploading file:", error);
		res.status(500).json({ error: "Failed to upload file" });
	} finally {
		// Send the message to Kafka topic
		if (producer) {
			await sendKafkaMessage(producer, location);
			console.log("Sent message to Kafka:", location);
		}
	}
});

export default router;
