export type FadeAnimProps = Pick<React.ComponentProps<"label">, "children"> & {
	in: boolean;
};
export const FadeAnim = ({ children, ...props }: FadeAnimProps) => {
	const c = props.in ? "animate-fadein" : "animate-fadeout";
	return <div className={`${c} overflow-hidden`}>{children}</div>;
};
