import { type ApiMockConfig, ServiceClient } from "@/utils/api";
import { inMock } from "@/utils/config";
import { useSearchParams } from "next/navigation";
import { useEffect } from "react";

export type config = {
	mode?: string;
	serviceClientConfig?: ApiMockConfig;
};
export const useAuthProps = (cfg: config = {}) => {
	const getServiceClient = inMock(cfg.mode)
		? ServiceClient.mget(cfg.serviceClientConfig)
		: ServiceClient.get;

	const searchParams = useSearchParams();
	const clientId = searchParams.get("client_id");
	useEffect(() => {
		let process = true;
		if (!clientId) return;
		// ServiceClient.mget(new BadRequestError())
		getServiceClient({ clientId }).then((resp) => {
			if (!process) return;
			if (resp instanceof Error) {
				console.error(resp);
			}
			console.log(resp);
		});
		return () => {
			process = false;
		};
	}, [clientId, getServiceClient]);
};
