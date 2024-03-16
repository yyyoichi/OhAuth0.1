import type { ComponentProps } from "react";

export type MyInputProps = ComponentProps<"input">;
export const MyInput = ({ className = "", ...props }: MyInputProps) => (
	<input
		className={`block w-full rounded-sm bg-my-white px-2 text-lg text-my-black outline-none focus:outline-none ${className}`}
		{...props}
	/>
);

export type MyInputLabelProps = {
	value: string;
} & ComponentProps<"div">;
export const MyInputLabel = ({
	className = "",
	...props
}: MyInputLabelProps) => {
	return (
		<div className={`text-lg ${className}`} {...props}>
			{props.value}
		</div>
	);
};

export type MyInputDescriptionProps = {
	value: string;
} & ComponentProps<"div">;
export const MyInputDescription = ({
	className = "",
	...props
}: MyInputDescriptionProps) => (
	<div className={`text-sm opacity-60 ${className}`} {...props}>
		{props.value}
	</div>
);

export type MyInputCautionProps = {
	value: string;
} & ComponentProps<"div">;
export const MyInputCaution = ({
	className = "",
	...props
}: MyInputCautionProps) => (
	<div className={`text-sm text-my-red opacity-80 ${className}`} {...props}>
		{props.value}
	</div>
);
