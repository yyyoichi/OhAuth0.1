const HOST = "http://localhost:8088";

export type GetServiceClient = (param: {
	clientId: string;
}) => Promise<ServiceClient | Error>;

export type PostAuthentication = (param: {
	clientId: string;
	userId: string;
	password: string;
}) => Promise<Authentication | Error>;

export type PostAuthorization = (param: {
	jwt: string;
	clientId: string;
	scope: string;
}) => Promise<Authorization | Error>;

export class ServiceClient {
	static get: GetServiceClient = async (param) => {
		const url = `${HOST}/api/v1/clients/${param.clientId}`;
		const resp = await fetch(url);
		const body = await json<{
			client_id: string;
			name: string;
			scope: string;
		}>(resp);
		if (body instanceof Error) {
			return body;
		}
		return new ServiceClient(body.client_id, body.name, body.scope);
	};
	static mget = (config?: ApiMockConfig): GetServiceClient => {
		const c = defaultConfig(config);
		return (param) => {
			return new Promise((resolve, _reject) => {
				setTimeout(() => {
					const Err = error(c.status);
					if (Err !== null) {
						resolve(new Err(param.clientId));
					}
					resolve(new ServiceClient(param.clientId, "name", "profile:view"));
				}, c.ms);
			});
		};
	};
	private constructor(
		readonly clientId: string = "",
		readonly name: string = "",
		readonly scope: string = "",
	) {}
}

export class Authentication {
	static post: PostAuthentication = async (param) => {
		const url = `${HOST}/api/v1/authentication`;
		const resp = await fetch(url, {
			method: "POST",
			body: JSON.stringify({
				cilent_id: param.clientId,
				user_id: param.userId,
				password: param.password,
			}),
		});
		const body = await json<{ jwt: string }>(resp);
		if (body instanceof Error) {
			return body;
		}
		return new Authentication(body.jwt);
	};
	static mpost = (config?: ApiMockConfig): PostAuthentication => {
		const c = defaultConfig(config);
		return (param) => {
			return new Promise((resolve, _reject) => {
				const jwt = `${param.userId}.${param.password}.${param.clientId}`;
				setTimeout(() => {
					const Err = error(c.status);
					if (Err !== null) {
						resolve(new Err(jwt));
					}
					resolve(new Authentication(jwt));
				}, c.ms);
			});
		};
	};
	private constructor(readonly jwt: string = "") {}
}

export class Authorization {
	static post: PostAuthorization = async (param) => {
		const url = `${HOST}/api/v1/authorization`;
		const resp = await fetch(url, {
			method: "POST",
			body: JSON.stringify({
				jwt: param.jwt,
				client_id: param.clientId,
				response_type: "code",
				scope: param.scope,
			}),
		});
		const body = await json<{ code: string }>(resp);
		if (body instanceof Error) {
			return body;
		}
		body as {
			code: string;
		};
		return new Authorization(body.code);
	};
	static mpost = (config?: ApiMockConfig): PostAuthorization => {
		const c = defaultConfig(config);
		return (param) => {
			return new Promise((resolve, _reject) => {
				const code = `${param.clientId}-${param.scope}-jwt-${param.jwt}`;
				setTimeout(() => {
					const Err = error(c.status);
					if (Err !== null) {
						resolve(new Err(code));
					}
					resolve(new Authorization(code));
				}, c.ms);
			});
		};
	};
	private constructor(readonly code: string = "") {}
}
const defaultConfig = (c?: ApiMockConfig) => {
	const df = {
		status: HttpStatus.Ok,
		ms: 500,
	};
	if (!c) {
		return df;
	}
	if (c.status) {
		df.status = c.status;
	}
	if (c.ms) {
		df.ms = c.ms;
	}
	return df;
};
export type ApiMockConfig = {
	status?: HttpStatus;
	ms?: number;
};

export enum HttpStatus {
	Ok = 200,
	BadRequest = 400,
	NotFound = 404,
	InternalServer = 500,
}

export class BadRequestError extends Error {}
export class NotFoundError extends Error {}
export class InternalServerError extends Error {}
const json = async <T>(resp: Response) => {
	const body = await resp.json();
	const Err = error(resp.status);
	if (Err !== null) {
		return new Err(body.status);
	}
	return body as T;
};
const error = (s: HttpStatus) => {
	switch (s) {
		case HttpStatus.BadRequest:
			return BadRequestError;
		case HttpStatus.NotFound:
			return NotFoundError;
		case HttpStatus.InternalServer:
			return InternalServerError;
	}
	return null;
};
