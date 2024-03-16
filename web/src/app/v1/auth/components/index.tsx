import type { AuthExternal } from "../lib/external";
import { useClientState } from "../lib/useClientState";

export type PageProps = {
	external: AuthExternal;
};
export default function Page(props: PageProps) {
	const client = useClientState(props.external);
	return (
		<div className="flex">
			<div>{client.client ? client.client?.name : "loading"}</div>
			{/* biome-ignore lint/a11y/useKeyWithClickEvents: <explanation> */}
			<div
				onClick={(e) => {
					console.log(e);
					client.redirect("hoge");
				}}
			>
				send
			</div>
		</div>
	);
}
