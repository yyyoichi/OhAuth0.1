import {
	type ApiMockConfig,
	type GetServiceClient,
	ServiceClient,
} from "@/utils/api";
import { inMock } from "@/utils/config";

export type Config = {
	mode?: string;
	serviceClientConfig?: ApiMockConfig;
};

export class AuthExternal {
	readonly getServiceClient: GetServiceClient;
	constructor(cfg: Config = {}) {
		this.getServiceClient = inMock(cfg.mode)
			? ServiceClient.mget(cfg.serviceClientConfig)
			: ServiceClient.get;
	}
}
