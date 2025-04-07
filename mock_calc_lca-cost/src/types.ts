export interface IfcFileData {
	project: string;
	filename: string;
	timestamp: string;
	data?: LcaData[] | CostData[];
}

export interface LcaData {
	id: string;
	sequence: number;
	mat_kbob: string;
	gwp_relative: number;
	gwp_absolute: number;
	penr_relative: number;
	penr_absolute: number;
	upb_relative: number;
	upb_absolute: number;
}

export interface CostData {
	id: string;
	cost: number;
	cost_unit: number;
}
