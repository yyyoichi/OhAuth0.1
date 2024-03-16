export default function AuthLayout({
	children,
}: Readonly<{
	children: React.ReactNode;
}>) {
	return (
		<main className="flex min-h-screen flex-col items-center justify-bet ween justify-end pb-16 md:mx-auto md:w-1/3 md:justify-center md:pb-auto">
			{children}
		</main>
	);
}
