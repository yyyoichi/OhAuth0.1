import type { ServiceClient as sc } from "@/utils/api";

export class ServiceClientProps {
	constructor(readonly sc: sc) {}
	redirect(code: string) {
		const url = `${this.sc.redirectUri}?code=${code}`;
		window.location.assign(url);
	}
	name = () => this.sc.name;
}
