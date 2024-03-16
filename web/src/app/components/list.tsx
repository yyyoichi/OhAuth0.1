import { DotIcon } from "./icons";
import type { MyColorPallet } from "./types";

export type MyUlProps = React.ComponentProps<"ul">;

export const MyUl = ({ children, className = "", ...props }: MyUlProps) => {
	const ulProps = {
		className: `list-none ${className}`,
		...props,
	};
	return <ul {...ulProps}>{children}</ul>;
};

export type MyUlLiProps = React.ComponentProps<"li">;
export const MyUlLi = ({ children, className = "", ...props }: MyUlLiProps) => {
	const liProps = {
		className: `flex items-center gap-1 ${className}`,
		...props,
	};
	return (
		<li {...liProps}>
			<DotIcon />
			{children}
		</li>
	);
};
