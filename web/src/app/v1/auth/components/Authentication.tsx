import { type LoadButtonProps, LoadButton } from "@/app/components/button";
import {
	Forms,
	PasswordFrame,
	type PasswordFrameProps,
	TextInputFrame,
	type TextInputFrameProps,
} from "@/app/components/forms";
import type { NonNullablePick } from "@/app/components/types";

export type BasicAuthenticationFormProps = {
	password: PasswordFrameProps;
	userId: TextInputFrameProps;
	sendButton: NonNullablePick<LoadButtonProps, "active" | "onClick">;
};

export const BasicAuthenticationForm = (
	props: BasicAuthenticationFormProps,
) => {
	const loadButtonProps: LoadButtonProps = {
		color: "green",
		...props.sendButton,
	};
	return (
		<Forms.Container>
			<h2 className="text-center text-3xl">Authentication</h2>
			<Forms.Content>
				<TextInputFrame {...props.userId} />
			</Forms.Content>
			<Forms.Content>
				<PasswordFrame {...props.password} />
			</Forms.Content>
			<div className="mx-auto w-1/3  py-4">
				<LoadButton {...loadButtonProps}>{"login"}</LoadButton>
			</div>
		</Forms.Container>
	);
};
