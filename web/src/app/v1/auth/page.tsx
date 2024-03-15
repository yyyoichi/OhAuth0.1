"use client";
import { useAuthProps } from "./lib/useAuthProps";

export default function Home() {
	useAuthProps();
	return (
		<div>
			<h1>Hello world</h1>
		</div>
	);
}
