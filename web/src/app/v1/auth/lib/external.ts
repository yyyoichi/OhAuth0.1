import {
	type ApiMockConfig,
	type GetServiceClient,
	ServiceClient,
	type PostAuthentication,
	type PostAuthorization,
	Authentication,
	Authorization,
} from "@/utils/api";
import { inMock } from "@/utils/config";

export type Config = {
	mode?: string;
	serviceClientConfig?: ApiMockConfig;
	authenticationConfig?: ApiMockConfig;
	authorizationConfig?: ApiMockConfig;
};

export class AuthExternal {
	getServiceClient: GetServiceClient;
	postAuthentication: PostAuthentication;
	postAuthorization: PostAuthorization;
	constructor(private cfg: Config = {}) {
		this.getServiceClient = inMock(cfg.mode)
			? ServiceClient.mget(cfg.serviceClientConfig)
			: ServiceClient.get;
		this.postAuthentication = inMock(cfg.mode)
			? Authentication.mpost(cfg.authenticationConfig)
			: Authentication.post;
		this.postAuthorization = inMock(cfg.mode)
			? Authorization.mpost(cfg.authorizationConfig)
			: Authorization.post;
	}
	setInterval(ms: number) {
		const cfg = this.cfg;
		cfg.serviceClientConfig = {
			...cfg.serviceClientConfig,
			ms,
		};
		cfg.authenticationConfig = {
			...cfg.authenticationConfig,
			ms,
		};
		cfg.authorizationConfig = {
			...cfg.authorizationConfig,
			ms,
		};
		const n = new AuthExternal(cfg);
		this.getServiceClient = n.getServiceClient;
		this.postAuthentication = n.postAuthentication;
		this.postAuthorization = n.postAuthorization;
		this.cfg = cfg;
	}
	changeMode(mode: Config["mode"]) {
		const cfg = {
			...this.cfg,
			mode,
		};
		const n = new AuthExternal(cfg);
		this.getServiceClient = n.getServiceClient;
		this.postAuthentication = n.postAuthentication;
		this.postAuthorization = n.postAuthorization;
		this.cfg = cfg;
	}
}

export const external = new AuthExternal();
