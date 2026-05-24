export default function RootLayout({
    menu,
    body,
}: Readonly<{
    menu: React.ReactNode;
    body: React.ReactNode;
}>) {
    return (
        <div className="min-h-screen bg-[#0d1117] flex">
            {menu}
            {/* Main Content */}
            {body}
        </div>
    );
}
