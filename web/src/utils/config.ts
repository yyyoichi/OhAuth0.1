export const inMock = (mode = process.env.NEXT_PUBLIC_MODE) => {
	return mode === "mock";
};
