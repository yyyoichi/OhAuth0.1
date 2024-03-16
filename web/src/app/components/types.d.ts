export type NonNullablePick<T, K extends keyof T> = {
	[P in K]-?: NonNullable<T[P]>;
};

export type MyColorPallet =
	| "red"
	| "yellow"
	| "light-black"
	| "black"
	| "white"
	| "green";
