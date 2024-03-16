import type { Config } from "tailwindcss";

const config: Config = {
	content: [
		"./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
		"./src/components/**/*.{js,ts,jsx,tsx,mdx}",
		"./src/app/**/*.{js,ts,jsx,tsx,mdx}",
	],
	theme: {
		extend: {
			backgroundImage: {
				"gradient-radial": "radial-gradient(var(--tw-gradient-stops))",
				"gradient-conic":
					"conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))",
			},
			colors: {
				"my-black": "var(--my-black)",
				"my-light-black": "var(--my-light-black)",
				"my-white": "var(--my-white)",
				"my-red": "var(--my-red)",
				"my-green": "var(--my-green)",
				"my-yellow": "var(--my-yellow)",
				"myc-black": "var(--myc-black)",
				"myc-light-black": "var(--myc-light-black)",
				"myc-white": "var(--myc-white)",
				"myc-red": "var(--myc-red)",
				"myc-green": "var(--myc-green)",
				"myc-yellow": "var(--myc-yellow)",
				"myc-fff": "var(--myc-fff)",
			},
			keyframes: {
				wiggle: {
					"0%, 100%": { transform: "rotate(-3deg)" },
					"50%": { transform: "rotate(3deg)" },
				},
				fadein: {
					"0%": {
						opacity: "0",
						transform: "translateY(5px)",
						height: "0",
					},
					"100%": {
						opacity: "1",
						transform: "translateY(0)",
						height: "inherit",
					},
				},
				fadeout: {
					"0%": {
						opacity: "1",
						transform: "translateY(-2px)",
						height: "inherit",
					},
					"100%": {
						opacity: "0",
						transform: "translateY(0)",
						height: "0",
					},
				},
			},
			animation: {
				wiggle: "wiggle 1s ease-in-out infinite",
				fadein: "fadein .8s forwards",
				fadeout: "fadeout .5s forwards",
			},
		},
	},
	plugins: [],
};
export default config;
