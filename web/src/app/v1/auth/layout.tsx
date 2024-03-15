export default function AuthLayout({
	children,
}: Readonly<{
	children: React.ReactNode;
}>) {
	return (
		<main className="flex min-h-screen flex-col items-center justify-bet ween p-24">
			{children}
		</main>
	);
}
