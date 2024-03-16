import { useState, useEffect } from "react";
import { ToggleViewIcon, type ToggleViewIconProps } from "./icons";
import {
	type MyInputLabelProps,
	type MyInputProps,
	type MyInputDescriptionProps,
	type MyInputCautionProps,
	MyInputLabel,
	MyInputDescription,
	MyInputCaution,
	MyInput,
} from "./input";
import type { Essential, NonNullablePick } from "./types";

export const Forms = {
	// wrap contents
	Container: ({ children }: { children: React.ReactNode }) => (
		<div className="flex flex-col gap-3 rounded-md bg-my-heavy-white py-2 ">
			{children}
		</div>
	),
	// wrap section in contents
	Content: ({ children }: { children: React.ReactNode }) => (
		<div className="px-4 py-1">{children}</div>
	),
};

///////////////////////////
// test input components //
///////////////////////////

export type TextInputFrameProps = {
	label: Pick<MyInputLabelProps, "value">;
	input: Essential<MyInputProps, "value" | "onChange" | "readOnly">;
	description: Pick<MyInputDescriptionProps, "value">;
	coution: Pick<MyInputCautionProps, "value">;
};
export const TextInputFrame = (props: TextInputFrameProps) => {
	const labelProps: MyInputLabelProps = {
		...props.label,
	};
	const inputProps: MyInputProps = {
		type: "text",
		...props.input,
	};
	const descriptionProps: MyInputDescriptionProps = {
		...props.description,
	};
	const coutionProps: MyInputDescriptionProps = {
		...props.coution,
	};
	return (
		<>
			<MyInputLabel {...labelProps} />
			<MyInputDescription {...descriptionProps} />
			<MyInputCaution {...coutionProps} />
			<MyInput {...inputProps} />
		</>
	);
};

/////////////////////////
// password components //
/////////////////////////

export type PasswordFrameProps = {
	description: Pick<MyInputDescriptionProps, "value">;
	coution: Pick<MyInputCautionProps, "value">;
} & PasswordFieldProps;

export const PasswordFrame = (props: PasswordFrameProps) => {
	const labelProps: MyInputLabelProps = {
		value: "Password",
	};
	const descriptionProps: MyInputDescriptionProps = {
		...props.description,
	};
	const coutionProps: MyInputDescriptionProps = {
		...props.coution,
	};

	return (
		<>
			<MyInputLabel {...labelProps} />
			<MyInputDescription {...descriptionProps} />
			<MyInputCaution {...coutionProps} />
			<PasswordField {...props} />
		</>
	);
};

type PasswordFieldProps = {
	input: Essential<MyInputProps, "value" | "onChange" | "readOnly"> & {
		visible?: boolean;
	};
};

const PasswordField = ({
	input: { visible, ...props },
}: PasswordFieldProps) => {
	const [visiblePassword, setVisiblePassword] = useState<boolean>(Boolean);
	const passwordToggleViewProps: PasswordToggleViewProps = {
		icon: {
			view: visiblePassword,
			onClick: () => setVisiblePassword((v) => !v),
		},
	};
	useEffect(() => {
		if (typeof visible !== "undefined") {
			setVisiblePassword(visible);
		}
	}, [visible]);
	const inputProps: MyInputProps = {
		type: visiblePassword ? "text" : "password",
		className: visiblePassword ? "text-[1.3rem]" : "",
		...props,
	};
	return (
		<div className="flex items-center bg-my-white pr-2">
			<MyInput {...inputProps} />
			<PasswordToggleView {...passwordToggleViewProps} />
		</div>
	);
};

type PasswordToggleViewProps = {
	icon: Pick<ToggleViewIconProps, "view" | "onClick">;
};
const PasswordToggleView = (props: PasswordToggleViewProps) => {
	const toggleViewIconProps: ToggleViewIconProps = {
		width: 22,
		height: 22,
		className: "fill-my-black",
		...props.icon,
	};
	return <ToggleViewIcon {...toggleViewIconProps} />;
};
