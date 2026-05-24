export default function DashboardLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return (
        <main className="flex-1 overflow-auto md:ml-0">
            <div className="max-w-7xl mx-auto">{children}</div>
        </main>
    );
}
