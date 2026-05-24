import { Network, Plus, Trash2 } from "lucide-react";
import Link from "next/link";

const vpcs = [
    {
        id: "vpc-001",
        name: "production-vpc",
        networkClass: "10.0.0.0/16",
        createdAt: new Date().toISOString(),
    },
    {
        id: "vpc-002",
        name: "staging-vpc",
        networkClass: "10.1.0.0/16",
        createdAt: new Date().toISOString(),
    },
];

export default async function VPCListPage() {
    const formatDate = (dateString: string) => {
        const date = new Date(dateString);
        return date.toLocaleString("en-US", {
            year: "numeric",
            month: "short",
            day: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    };

    return (
        <div className="p-6 pt-16 md:pt-6">
            <div className="flex items-center justify-between mb-6">
                <div>
                    <h1 className="text-xl text-[#e6edf3] mb-1">
                        Virtual Private Clouds
                    </h1>
                    <p className="text-[#8b949e] text-xs">
                        Manage isolated private networks for your VMs
                    </p>
                </div>
                <button className="px-4 py-2 rounded transition-colors bg-blue-600 hover:bg-blue-700">
                    <Link
                        href="/vpc/create"
                        className="text-sm flex items-center gap-2 text-white"
                    >
                        <Plus className="w-4 h-4" />
                        Create VPC
                    </Link>
                </button>
            </div>

            {vpcs.length > 0 ? (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {vpcs.map((vpc) => (
                        <div
                            key={vpc.id}
                            className="bg-[#161b22] border border-[#30363d] rounded p-5 hover:border-[#484f58] transition-colors"
                        >
                            <div className="flex items-start justify-between mb-4">
                                <div className="flex-1">
                                    <div className="flex items-center gap-2 mb-1">
                                        <Network className="w-4 h-4 text-purple-400" />
                                        <h3 className="text-base text-[#e6edf3]">
                                            {vpc.name}
                                        </h3>
                                    </div>
                                    <p className="text-xs text-[#8b949e]">
                                        {vpc.id}
                                    </p>
                                </div>
                                <button
                                    className="p-1.5 hover:bg-[#30363d] rounded transition-colors text-red-400"
                                    title="Delete VPC"
                                >
                                    <Trash2 className="w-4 h-4" />
                                </button>
                            </div>

                            <div className="space-y-3 pt-3 border-t border-[#30363d]">
                                <div>
                                    <div className="text-xs text-[#8b949e] mb-1">
                                        Network Class
                                    </div>
                                    <div className="text-sm text-[#e6edf3] font-mono">
                                        {vpc.networkClass}
                                    </div>
                                </div>
                                <div>
                                    <div className="text-xs text-[#8b949e] mb-1">
                                        Created
                                    </div>
                                    <div className="text-xs text-[#e6edf3]">
                                        {formatDate(vpc.createdAt)}
                                    </div>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            ) : (
                <div className="bg-[#161b22] border border-[#30363d] rounded p-12 text-center">
                    <Network className="w-12 h-12 text-[#8b949e] mx-auto mb-3" />
                    <div className="text-[#8b949e] text-base mb-2">
                        No VPCs created
                    </div>
                    <div className="text-xs text-[#8b949e] mb-4">
                        Create a VPC to set up isolated private networks
                    </div>
                    <button className="px-4 py-2 rounded transition-colors bg-blue-600 hover:bg-blue-700">
                        <Link
                            href="/vpc/create"
                            className="text-sm inline-flex items-center gap-2 text-white"
                        >
                            <Plus className="w-4 h-4" />
                            Create Your First VPC
                        </Link>
                    </button>
                </div>
            )}
        </div>
    );
}
