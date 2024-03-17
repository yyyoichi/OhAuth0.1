"use client";
import type { ServiceClient } from "@/utils/api";
import { ServiceClientProps } from "../lib/serviceClientProps";
import {
	BasicAuthenticationForm,
	type BasicAuthenticationFormProps,
} from "./Authentication";
import {
	AuthorizationForm,
	type AuthorizationFormProps,
} from "./Authorization";
import { useState } from "react";
import { FadeAnim } from "@/app/components/animation";
import { useAuthenticationState } from "../lib/useAuthenticationState";
import { external } from "../lib/external";

export type V1AuthPageProps = {
	serviceClient: ServiceClient;
};
export function V1AuthPage({ serviceClient }: V1AuthPageProps) {
	const sc = new ServiceClientProps(serviceClient);
	const props = getV1AuthProps(sc);
	return (
		<>
			<h2 className="my-3 text-center text-3xl">
				{"Sign in to "}
				<span className="text-my-red">{sc.name()}</span>
			</h2>
			<FadeAnim in={props.inAuthenticationPage}>
				<BasicAuthenticationForm {...props.authenticationProsp} />
			</FadeAnim>
			<FadeAnim in={!props.inAuthenticationPage}>
				<AuthorizationForm {...props.authorizationProps} />
			</FadeAnim>
		</>
	);
}

export const getV1AuthProps = (sc: ServiceClientProps) => {
	const [loading, setLoading] = useState(false);
	const [inAuthenticationPage, setInAuthenticationPage] = useState(true);
	const Auth = useAuthenticationState();

	const inputIsReadOnly = loading;
	const buttonIsActive = !loading;

	const entryDoAuthentication = () => {
		setLoading(true);
		Auth.clearInputCoution();

		if (!Auth.validate()) {
			setLoading(false);
			return;
		}
		Auth.doBasicAuthenticatioin(sc.clientId())
			.then((result) => {
				if (result instanceof Error) {
					window.alert(result.message);
					return;
				}
				setInAuthenticationPage(false);
			})
			.finally(() => setLoading(false));
	};

	const authenticationProsp: BasicAuthenticationFormProps = {
		userId: {
			label: {
				value: "User Id",
			},
			description: {
				value: "",
			},
			input: {
				onChange: (e) => {
					e.preventDefault();
					Auth.setBasic((prev) => ({ ...prev, userId: e.target.value }));
				},
				readOnly: inputIsReadOnly,
				value: Auth.userId,
			},
			coution: {
				value: Auth.inputCoution.userId,
			},
		},
		password: {
			coution: {
				value: Auth.inputCoution.password,
			},
			description: {
				value: "",
			},
			input: {
				onKeyDown: (e) => {
					if (e.key === "Enter") {
						e.preventDefault();
						entryDoAuthentication();
					}
				},
				onChange: (e) => {
					e.preventDefault();
					Auth.setBasic((prev) => ({ ...prev, password: e.target.value }));
				},
				readOnly: inputIsReadOnly,
				value: Auth.password,
				visible: false, // default password is invisible.
			},
		},
		sendButton: {
			active: buttonIsActive,
			onClick: (e) => {
				e.preventDefault();
				entryDoAuthentication();
			},
		},
	};
	const authorizationProps: AuthorizationFormProps = {
		okButton: {
			active: buttonIsActive,
			onClick: (e) => {
				e.preventDefault();
				setLoading(true);
				external
					.postAuthorization({
						clientId: sc.clientId(),
						jwt: Auth.jwt,
						scope: sc.scope(),
					})
					.then((resp) => {
						if (resp instanceof Error) {
							window.alert(resp.message);
							return;
						}
						sc.redirect(resp.code);
					})
					.finally(() => setLoading(false));
			},
		},
		cancelButton: {
			active: buttonIsActive,
			onClick: (e) => {
				e.preventDefault();
				sc.redirect("");
			},
		},
	};

	return {
		authenticationProsp,
		authorizationProps,
		inAuthenticationPage,
	};
};
