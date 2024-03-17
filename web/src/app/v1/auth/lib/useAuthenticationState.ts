import { useCallback, useState } from "react";
import { external } from "./external";

export const useAuthenticationState = () => {
	const [basic, setBasic] = useState({
		userId: "",
		password: "",
	});
	const [jwt, setJwt] = useState("");
	const [inputCoution, setInputCoution] = useState({
		userId: "",
		password: "",
	});

	const validate = useCallback(() => {
		const userIdCoution: string[] = [];
		const passwordCoution: string[] = [];
		if (basic.userId === "") {
			userIdCoution.push("'User Id' is required.");
		}
		if (basic.password === "") {
			passwordCoution.push("'Password' is required.");
		}
		if (userIdCoution.length === 0 && passwordCoution.length === 0) {
			return true;
		}
		// invalid !
		setInputCoution({
			userId: userIdCoution.join(""),
			password: passwordCoution.join(""),
		});
		return false;
	}, [basic]);

	const doBasicAuthenticatioin = useCallback(
		async (clientId: string) => {
			const resp = await external.postAuthentication({
				clientId,
				...basic,
			});
			if (resp instanceof Error) {
				return resp;
			}
			setJwt(resp.jwt);
			return null;
		},
		[basic],
	);

	return {
		validate,
		doBasicAuthenticatioin,
		clearInputCoution: () => setInputCoution({ userId: "", password: "" }),
		setBasic,
		jwt,
		...basic,
		inputCoution,
	};
};
