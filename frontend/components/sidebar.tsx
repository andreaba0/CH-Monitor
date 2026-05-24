"use client";

import { useState } from "react";
import { Menu, Server, Network, Plus, X } from "lucide-react";
import { usePathname } from "next/navigation";
import Link from "next/link";

export function Sidebar() {
    const pathname = usePathname();
    const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

    const isVPSActive = pathname.startsWith("/vps/");
    const isVPCActive = pathname.startsWith("/vpc/");

    function renderMenuIcon() {
        return (
            <>
                {/* 1. Mobile Menu Button - Visible on mobile, hidden on desktop (md and up) */}
                <button
                    onClick={() => setMobileMenuOpen(true)}
                    className="block md:hidden fixed top-4 left-4 z-40 p-2 bg-[#161b22] border border-[#30363d] rounded"
                >
                    <Menu className="w-5 h-5 text-[#e6edf3]" />
                </button>

                {/* 2. Backdrop Shadow - Only clickable/visible on mobile layouts */}
                {mobileMenuOpen && (
                    <div
                        onClick={() => setMobileMenuOpen(false)}
                        className="block md:hidden fixed inset-0 bg-black/50 z-40"
                    />
                )}
            </>
        );
    }

    return (
        <div className="min-h-screen bg-[#0d1117] flex">
            {renderMenuIcon()}
            {/* Side Menu */}
            <aside
                className={`w-64 bg-[#161b22] border-r border-[#30363d] flex flex-col fixed md:sticky md:top-0 h-full md:h-screen z-50 transition-transform ${
                    mobileMenuOpen
                        ? "translate-x-0"
                        : "-translate-x-full md:translate-x-0"
                }`}
            >
                <div className="p-4 border-b border-[#30363d]">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                            <div className="w-8 h-8 bg-blue-500/10 rounded flex items-center justify-center">
                                <svg
                                    className="w-5 h-5 text-blue-400"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                >
                                    <path
                                        strokeLinecap="round"
                                        strokeLinejoin="round"
                                        strokeWidth={2}
                                        d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z"
                                    />
                                </svg>
                            </div>
                            <h1 className="text-sm text-[#e6edf3]">
                                Cloud-Hypervisor
                            </h1>
                        </div>
                        <button
                            onClick={() => setMobileMenuOpen(false)}
                            className="md:hidden p-1 hover:bg-[#30363d] rounded transition-colors"
                        >
                            <X className="w-5 h-5 text-[#8b949e]" />
                        </button>
                    </div>
                </div>

                <nav className="flex-1 p-3">
                    <div className="space-y-1">
                        <button
                            className={`w-full rounded transition-colors ${
                                isVPSActive
                                    ? "bg-blue-500/10 text-blue-400"
                                    : "text-[#8b949e] hover:text-[#e6edf3] hover:bg-[#30363d]"
                            }`}
                        >
                            <Link
                                className="w-full flex gap-3 px-3 py-2 text-sm"
                                href="/vps/list"
                            >
                                <Server className="w-4 h-4" />
                                VPS
                            </Link>
                        </button>
                        <button
                            className={`w-full rounded transition-colors ${
                                isVPCActive
                                    ? "bg-purple-500/10 text-purple-400"
                                    : "text-[#8b949e] hover:text-[#e6edf3] hover:bg-[#30363d]"
                            }`}
                        >
                            <Link
                                className="flex items-center gap-3 px-3 py-2 text-sm"
                                href="/vpc/list"
                            >
                                <Network className="w-4 h-4" />
                                VPC
                            </Link>
                        </button>
                    </div>
                </nav>

                <div className="p-3 border-t border-[#30363d]">
                    {isVPSActive && (
                        <button className="w-full rounded transition-colors">
                            <Link
                                className="px-3 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm flex items-center justify-center gap-2"
                                href="/vps/create"
                            >
                                <Plus className="w-4 h-4" />
                                Deploy VM
                            </Link>
                        </button>
                    )}
                    {isVPCActive && (
                        <button className="w-full rounded transition-colors">
                            <Link
                                className="px-3 py-2 bg-purple-600 hover:bg-purple-700 text-white text-sm flex items-center justify-center gap-2"
                                href="/vpc/create"
                            >
                                <Plus className="w-4 h-4" />
                                Create VPC
                            </Link>
                        </button>
                    )}
                </div>
            </aside>
        </div>
    );
}
