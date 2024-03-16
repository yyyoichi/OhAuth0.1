"use client";
import type { ServiceClient } from "@/utils/api";
import { ServiceClientProps } from "../lib/serviceClientProps";

export type V1AuthPageProps = {
	serviceClient: ServiceClient;
};
export function V1AuthPage({ serviceClient }: V1AuthPageProps) {
	const props = new ServiceClientProps(serviceClient);
	return (
		<div className="flex">
			<div>{props.name()}</div>
			{/* biome-ignore lint/a11y/useKeyWithClickEvents: <explanation> */}
			<div
				onClick={(e) => {
					console.log(e);
					props.redirect("code");
				}}
			>
				send
			</div>
		</div>
	);
}
