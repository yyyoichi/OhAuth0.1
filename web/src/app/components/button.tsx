import type { ComponentProps } from "react";
import type { MyColorPallet } from "./types";
import { MyLoader } from "./loader";

export type MyButtonProps = {
	color: MyColorPallet;
} & ComponentProps<"div">;
export const MyButton = ({
	className = "",
	color,
	...props
}: MyButtonProps) => {
	const buttonProps: ComponentProps<"div"> = {
		...props,
	};
	const getColorClassName = () => {
		switch (color) {
			case "red":
				return "bg-my-red text-myc-white";
			case "yellow":
				return "bg-my-yellow text-myc-white";
			case "light-black":
				return "bg-my-light-black text-my-white";
			case "black":
				return "bg-my-black text-my-white";
			case "white":
				return "bg-my-white text-my-black";
			case "green":
				return "bg-my-green text-myc-white";
		}
	};
	return (
		<div
			className={`rounded-md text-center font-bold hover:opacity-80 ${getColorClassName()} ${className}`}
			{...buttonProps}
		/>
	);
};

export type LoadButtonProps = {
	active: boolean;
} & MyButtonProps;
export const LoadButton = ({
	active,
	className = "",
	color,
	children,
	onClick,
	...props
}: LoadButtonProps) => {
	const getPairColor = (): MyColorPallet => {
		switch (color) {
			case "red":
				return "white";
			case "yellow":
				return "black";
			case "light-black":
				return "white";
			case "black":
				return "white";
			case "white":
				return "black";
			case "green":
				return "black";
		}
	};
	const Child = active ? children : <MyLoader color={getPairColor()} />;
	const onclick = active ? onClick : undefined;

	return (
		<MyButton
			className={`w-full text-center ${
				active ? "cursor-pointer" : ""
			} ${className}`}
			color={color}
			onClick={onclick}
			{...props}
		>
			{Child}
		</MyButton>
	);
};
