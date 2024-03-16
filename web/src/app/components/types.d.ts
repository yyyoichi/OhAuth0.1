export type NonNullablePick<T, K extends keyof T> = {
	[P in K]-?: NonNullable<T[P]>;
};
export type Essential<T, K extends keyof T> = NonNullablePick<T, K> &
	Omit<T, K>;
export type MyColorPallet =
	| "red"
	| "yellow"
	| "light-black"
	| "black"
	| "white"
	| "green";
