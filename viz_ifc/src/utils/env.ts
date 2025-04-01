import { log } from "./logger";

export function getEnv(key: string): string {
	const value = process.env[key];
	if (!value) {
		log.error(`Environment variable ${key} is not set`);
		process.exit(1);
	}
	return value;
}
