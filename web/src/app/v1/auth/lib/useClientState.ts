import { useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import type { AuthExternal } from "./external";
import type { ServiceClient } from "@/utils/api";

export const useClientState = (external: AuthExternal) => {
	const [client, setClient] = useState<ServiceClient | null>(null);
	const searchParams = useSearchParams();
	const clientId = searchParams.get("client_id");
	// biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
	useEffect(() => {
		let process = true;
		if (!clientId) return;
		setTimeout(() => {
			if (!process) return;
			external.getServiceClient({ clientId }).then((resp) => {
				if (!process) return;
				if (resp instanceof Error) {
					console.error(resp);
					return;
				}
				setClient(resp);
			});
		}, 100);
		return () => {
			process = false;
		};
	}, [clientId]);

	const redirect = (code: string) => {
		if (!client) return;
		const url = `${client.redirectUri}?code=${code}`;
		window.location.assign(url);
	};

	return {
		client,
		redirect,
	};
};
