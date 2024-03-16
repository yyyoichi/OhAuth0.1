import { V1AuthPage, type V1AuthPageProps } from "./components";
import { AuthExternal } from "./lib/external";

export default async function Page({
	searchParams,
}: {
	searchParams: { [key: string]: string | string[] | undefined };
}) {
	const external = new AuthExternal();
	let clientId: string;
	switch (typeof searchParams.client_id) {
		case "string":
			clientId = searchParams.client_id;
			break;
		case "undefined":
			clientId = "501";
			break;
		case "object":
			clientId = searchParams.client_id[0];
	}
	const serviceClient = await external.getServiceClient({ clientId });
	if (serviceClient instanceof Error) {
		return <div>Error: {serviceClient.message}</div>;
	}
	const pageProps: V1AuthPageProps = {
		serviceClient: {
			clientId: serviceClient.clientId,
			name: serviceClient.name,
			scope: serviceClient.scope,
			redirectUri: serviceClient.redirectUri,
		},
	};
	return <V1AuthPage {...pageProps} />;
}
