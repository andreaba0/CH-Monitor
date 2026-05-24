import type { Metadata } from "next";
import { JetBrains_Mono } from "next/font/google";
import "./globals.css";
import { Toaster } from "@/components/sonner";

const jetbrainsMono = JetBrains_Mono({
    variable: "--font-jetbrains-mono",
    subsets: ["latin"],
});

export const metadata: Metadata = {
    title: "VMM Console",
    description: "VMM Console",
};

export default function RootLayout({
    menu,
    body,
}: Readonly<{
    menu: React.ReactNode;
    body: React.ReactNode;
}>) {
    return (
        <html
            lang="en"
            className={`${jetbrainsMono.variable} h-full antialiased`}
        >
            <body>
                <div className="min-h-screen bg-[#0d1117] flex">
                    <Toaster position="top-right" />
                    {menu}
                    {/* Main Content */}
                    {body}
                </div>
            </body>
        </html>
    );
}
