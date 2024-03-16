import { LoadButton, type LoadButtonProps } from "@/app/components/button";
import { Forms } from "@/app/components/forms";
import {
	MyInputDescription,
	type MyInputDescriptionProps,
	MyInputLabel,
	type MyInputLabelProps,
} from "@/app/components/input";
import { MyUl, MyUlLi } from "@/app/components/list";
import type { NonNullablePick } from "@/app/components/types";

export type AuthorizationFormProps = {
	okButton: NonNullablePick<LoadButtonProps, "active" | "onClick">;
	cancelButton: NonNullablePick<LoadButtonProps, "active" | "onClick">;
};
export const AuthorizationForm = (props: AuthorizationFormProps) => {
	const labelProps: MyInputLabelProps = {
		value: "Do you want to allow access?",
	};
	const descriptionProps: MyInputDescriptionProps = {
		value:
			"An application or service has requested access to your OhAuth0.1 account and resources.",
	};
	const okButtonProps: LoadButtonProps = {
		color: "green",
		...props.okButton,
	};
	const cancelButtonProps: LoadButtonProps = {
		color: "light-black",
		...props.cancelButton,
	};
	return (
		<Forms.Container>
			<h2 className="text-center text-3xl">Authentication</h2>
			<Forms.Content>
				<MyInputLabel {...labelProps} />
				<MyInputDescription className="text-wrap" {...descriptionProps} />
				<MyUl className="py-4">
					<MyUlLi>{"View your profile"}</MyUlLi>
				</MyUl>
			</Forms.Content>
			<Forms.Content>
				<div className="mx-auto w-1/3  py-4 flex flex-col gap-3">
					<LoadButton {...cancelButtonProps}>{"Cancel"}</LoadButton>
					<LoadButton {...okButtonProps}>{"OK"}</LoadButton>
				</div>
			</Forms.Content>
		</Forms.Container>
	);
};
