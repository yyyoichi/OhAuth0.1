"use client";
import dynamic from "next/dynamic";
import { AuthExternal } from "./lib/external";

const Client = dynamic(() => import("./components"), { ssr: false });
export default function Home() {
	return (
		<div>
			<h1>Hello world</h1>
			<Client external={new AuthExternal()} />
		</div>
	);
}
