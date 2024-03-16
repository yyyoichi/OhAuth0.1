import type { ComponentProps } from "react";
import type { MyColorPallet } from "./types";

type MyLoaderProps = {
	color: MyColorPallet;
} & ComponentProps<"div">;
export const MyLoader = ({
	className = "",
	color,
	...props
}: MyLoaderProps) => {
	return (
		<div className={"flex items-center justify-center p-1"}>
			<div
				className={`h-4 w-4 animate-spin rounded-full border-b-2 border-t-2  ${getColorClassName(
					color,
				)} ${className}`}
				{...props}
			/>
		</div>
	);
};

const getColorClassName = (color: MyLoaderProps["color"]) => {
	switch (color) {
		case "red":
			return "border-my-red";
		case "yellow":
			return "border-my-yellow";
		case "light-black":
			return "border-my-light-black";
		case "black":
			return "border-my-black";
		case "white":
			return "border-my-white";
		case "green":
			return "border-my-green";
	}
};
