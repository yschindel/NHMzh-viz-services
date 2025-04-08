import winston from "winston";

// Define log levels (matching OpenTelemetry severity levels)
const levels = {
	error: 0,
	warn: 1,
	info: 2,
	http: 3,
	debug: 4,
};

// Create the logger instance
const logger = winston.createLogger({
	level: process.env.NODE_ENV === "development" ? "debug" : "info",
	levels,
	format: winston.format.combine(
		winston.format.timestamp(),
		winston.format.json() // Use JSON format for better OpenTelemetry compatibility
	),
	transports: [
		new winston.transports.Console({
			format: winston.format.combine(winston.format.colorize(), winston.format.simple()),
		}),
		new winston.transports.File({
			filename: "logs/error.log",
			level: "error",
		}),
		new winston.transports.File({
			filename: "logs/combined.log",
		}),
	],
});

// Create a wrapper class that will make it easier to add OpenTelemetry later
class Logger {
	private static instance: Logger;
	private winstonLogger: winston.Logger;

	private constructor() {
		this.winstonLogger = logger;
	}

	public static getInstance(): Logger {
		if (!Logger.instance) {
			Logger.instance = new Logger();
		}
		return Logger.instance;
	}

	public error(message: string, ...meta: any[]): void {
		this.winstonLogger.error(message, ...meta);
	}

	public warn(message: string, ...meta: any[]): void {
		this.winstonLogger.warn(message, ...meta);
	}

	public info(message: string, ...meta: any[]): void {
		this.winstonLogger.info(message, ...meta);
	}

	public http(message: string, ...meta: any[]): void {
		this.winstonLogger.http(message, ...meta);
	}

	public debug(message: string, ...meta: any[]): void {
		this.winstonLogger.debug(message, ...meta);
	}
}

// Export a singleton instance
export const log = Logger.getInstance();
