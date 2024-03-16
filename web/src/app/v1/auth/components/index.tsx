"use client";
import type { ServiceClient } from "@/utils/api";
import { ServiceClientProps } from "../lib/serviceClientProps";
import { MyButton } from "@/app/components/button";
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

export type V1AuthPageProps = {
	serviceClient: ServiceClient;
};
export function V1AuthPage({ serviceClient }: V1AuthPageProps) {
	const sc = new ServiceClientProps(serviceClient);
	const authenticationProsp: BasicAuthenticationFormProps = {
		password: {
			coution: {
				value: "caution",
			},
			description: {
				value: "description",
			},
			input: {
				onChange: (e) => {
					console.log(e);
				},
				readOnly: false,
				value: "password",
				visible: true,
			},
		},
		userId: {
			label: {
				value: "User Id",
			},
			input: {
				onChange: (e) => console.log(e),
				readOnly: false,
				value: "user-1",
			},
			description: {
				value: "description",
			},
			coution: {
				value: "coution",
			},
		},
		sendButton: {
			active: true,
			onClick: (e) => console.log("send"),
		},
	};
	const authorizationProps: AuthorizationFormProps = {
		okButton: {
			active: true,
			onClick: (e) => console.log(e),
		},
		cancelButton: {
			active: true,
			onClick: (e) => console.log(e),
		},
	};
	const [isAuthenticatioin, toggleView] = useState(true);
	return (
		<>
			{/* <div className="bg-my-white">{sc.name()}</div>
			<div className="bg-my-light-black">{sc.name()}</div>
			<div className="bg-my-black">{sc.name()}</div>
			<div className="bg-my-red">{sc.name()}</div>
			<div className="bg-my-green">{sc.name()}</div>
			<div className="bg-my-yellow">{sc.name()}</div>
			<MyButton color="white">{"submit"}</MyButton>
			<MyButton color="black">{"submit"}</MyButton>
			<MyButton color="light-black">{"submit"}</MyButton>
			<MyButton color="red">{"submit"}</MyButton>
			<MyButton color="green">{"submit"}</MyButton>
			<MyButton color="yellow">{"submit"}</MyButton> */}
			<h2 className="my-3 text-center text-3xl">
				{"Sign in to "}
				<span className="text-my-red">{sc.name()}</span>
			</h2>
			<FadeAnim in={isAuthenticatioin}>
				<BasicAuthenticationForm {...authenticationProsp} />
			</FadeAnim>
			<FadeAnim in={!isAuthenticatioin}>
				<AuthorizationForm {...authorizationProps} />
			</FadeAnim>
			<MyButton
				color="yellow"
				className="cursor-pointer my-1"
				onClick={() => toggleView((v) => !v)}
			>
				{"toggle"}
			</MyButton>
		</>
	);
}
