import { ChevronRight, Plus } from "lucide-react";
import Link from "next/link";

const vms = [
    {
        id: "vm-001",
        name: "ubuntu-web-server",
        status: "running",
        cpus: 4,
        memory: 4096,
        diskFiles: [
            {
                name: "ubuntu-root.img",
                size: 42949672960,
                path: "/var/lib/vms/ubuntu-web-server/ubuntu-root.img",
                serial: "CH-XY8K9L2M3N4P",
                isBootDisk: true,
            },
        ],
        uptime: "2d 14h 32m",
        cpuUsage: 23,
        memoryUsage: 67,
        cpuHistory: [
            { time: "14:00", value: 18 },
            { time: "14:05", value: 22 },
            { time: "14:10", value: 25 },
            { time: "14:15", value: 19 },
            { time: "14:20", value: 24 },
            { time: "14:25", value: 28 },
            { time: "14:30", value: 23 },
        ],
        memoryHistory: [
            { time: "14:00", value: 62 },
            { time: "14:05", value: 64 },
            { time: "14:10", value: 66 },
            { time: "14:15", value: 65 },
            { time: "14:20", value: 68 },
            { time: "14:25", value: 69 },
            { time: "14:30", value: 67 },
        ],
        networkInterfaces: [
            {
                id: "net-0",
                type: "public",
                macAddress: "52:54:00:12:34:56",
                tapDevice: "vmtap0",
                ipAddress: "192.168.1.10",
            },
        ],
        kernelPath: "/boot/vmlinuz-5.15.0",
        yamlPath: "/etc/cloud-hypervisor/vms/ubuntu-web-server.yaml",
        yamlConfig:
            "# VM Configuration\nname: ubuntu-web-server\ncpus:\n  boot_vcpus: 4\n  max_vcpus: 4",
        createdAt: "2026-05-10T10:30:00Z",
    },
    {
        id: "vm-002",
        name: "debian-database",
        status: "running",
        cpus: 8,
        memory: 8192,
        diskFiles: [
            {
                name: "debian-root.img",
                size: 107374182400,
                path: "/var/lib/vms/debian-database/debian-root.img",
                serial: "CH-AB1CD2EF3GH4",
                isBootDisk: true,
            },
            {
                name: "data-disk.img",
                size: 214748364800,
                path: "/var/lib/vms/debian-database/data-disk.img",
                serial: "CH-QW5ER6TY7UI8",
                isBootDisk: false,
            },
        ],
        uptime: "5d 3h 18m",
        cpuUsage: 45,
        memoryUsage: 82,
        cpuHistory: [
            { time: "14:00", value: 42 },
            { time: "14:05", value: 48 },
            { time: "14:10", value: 46 },
            { time: "14:15", value: 44 },
            { time: "14:20", value: 47 },
            { time: "14:25", value: 49 },
            { time: "14:30", value: 45 },
        ],
        memoryHistory: [
            { time: "14:00", value: 78 },
            { time: "14:05", value: 80 },
            { time: "14:10", value: 81 },
            { time: "14:15", value: 79 },
            { time: "14:20", value: 83 },
            { time: "14:25", value: 84 },
            { time: "14:30", value: 82 },
        ],
        networkInterfaces: [
            {
                id: "net-0",
                type: "private",
                vpcId: "vpc-backend",
                macAddress: "52:54:00:ab:cd:ef",
                tapDevice: "vpc-backend-0",
                ipAddress: "10.0.1.10",
            },
            {
                id: "net-1",
                type: "public",
                macAddress: "52:54:00:78:90:12",
                tapDevice: "vmtap1",
                ipAddress: "192.168.1.20",
            },
        ],
        kernelPath: "/boot/vmlinuz-5.15.0",
        yamlPath: "/etc/cloud-hypervisor/vms/debian-database.yaml",
        yamlConfig:
            "# VM Configuration\nname: debian-database\ncpus:\n  boot_vcpus: 8\n  max_vcpus: 8",
        createdAt: "2026-05-05T08:15:00Z",
    },
    {
        id: "vm-003",
        name: "alpine-test",
        status: "stopped",
        cpus: 2,
        memory: 1024,
        diskFiles: [
            {
                name: "alpine.img",
                size: 10737418240,
                path: "/var/lib/vms/alpine-test/alpine.img",
                serial: "CH-ZX9CV8BN7ML6",
                isBootDisk: true,
            },
        ],
        uptime: "0d 0h 0m",
        cpuUsage: 0,
        memoryUsage: 0,
        cpuHistory: [],
        memoryHistory: [],
        networkInterfaces: [
            {
                id: "net-0",
                type: "public",
                macAddress: "52:54:00:34:56:78",
                tapDevice: "vmtap2",
            },
        ],
        kernelPath: "/boot/vmlinuz-alpine",
        yamlPath: "/etc/cloud-hypervisor/vms/alpine-test.yaml",
        yamlConfig:
            "# VM Configuration\nname: alpine-test\ncpus:\n  boot_vcpus: 2\n  max_vcpus: 2",
        createdAt: "2026-05-14T16:45:00Z",
    },
];

export default async function DashboardPage() {
    const totalVMs = vms.length;
    const runningVMs = vms.filter((vm) => vm.status === "running").length;

    const getStatusColor = (status: string) => {
        switch (status) {
            case "running":
                return "bg-green-500/10 text-green-400 border-green-500/20";
            case "stopped":
                return "bg-[#30363d] text-[#8b949e] border-[#30363d]";
            case "paused":
                return "bg-yellow-500/10 text-yellow-400 border-yellow-500/20";
            case "starting":
                return "bg-blue-500/10 text-blue-400 border-blue-500/20";
            default:
                return "bg-[#30363d] text-[#8b949e] border-[#30363d]";
        }
    };

    return (
        <div className="p-6 pt-16 md:pt-6">
            <div className="flex items-center justify-between mb-6">
                <div>
                    <h1 className="text-xl text-[#e6edf3] mb-1">
                        Virtual Machines
                    </h1>
                    <p className="text-[#8b949e] text-xs">
                        Manage and monitor your VM instances
                    </p>
                </div>
                <button className="px-4 py-2 rounded transition-colors bg-blue-600 hover:bg-blue-700">
                    <Link
                        href="/vps/create"
                        className="text-sm flex items-center gap-2 text-white"
                    >
                        <Plus className="w-4 h-4" />
                        Create VPS
                    </Link>
                </button>
            </div>

            <div>
                <div className="flex items-center justify-between mb-4">
                    <h2 className="text-base text-[#e6edf3]">VM Instances</h2>
                    <div className="text-sm text-[#8b949e]">
                        {totalVMs} total • {runningVMs} running
                    </div>
                </div>

                {vms.length > 0 ? (
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                        {vms.map((vm) => (
                            <div
                                key={vm.id}
                                className="bg-[#161b22] border border-[#30363d] rounded p-5 hover:border-[#484f58] transition-colors"
                            >
                                <div className="flex items-start justify-between mb-4">
                                    <div className="flex-1">
                                        <div className="flex items-center gap-2 cursor-pointer group">
                                            <Link href={`/vm/${vm.id}`}>
                                                <h3 className="text-base text-[#e6edf3] group-hover:text-blue-400 transition-colors">
                                                    {vm.name}
                                                </h3>
                                            </Link>
                                            <ChevronRight className="w-4 h-4 text-[#8b949e] group-hover:text-blue-400 transition-colors" />
                                        </div>
                                        <p className="text-xs text-[#8b949e] mt-1">
                                            {vm.id}
                                        </p>
                                    </div>
                                    <span
                                        className={`px-2.5 py-1 rounded text-xs border ${getStatusColor(vm.status)}`}
                                    >
                                        {vm.status.toUpperCase()}
                                    </span>
                                </div>

                                <div className="grid grid-cols-2 gap-3 mb-4 pb-4 border-b border-[#30363d]">
                                    <div>
                                        <div className="text-xs text-[#8b949e] mb-1">
                                            vCPUs
                                        </div>
                                        <div className="text-sm text-[#e6edf3]">
                                            {vm.cpus} cores
                                        </div>
                                    </div>
                                    <div>
                                        <div className="text-xs text-[#8b949e] mb-1">
                                            Memory
                                        </div>
                                        <div className="text-sm text-[#e6edf3]">
                                            {vm.memory} MB
                                        </div>
                                    </div>
                                    <div>
                                        <div className="text-xs text-[#8b949e] mb-1">
                                            Disks
                                        </div>
                                        <div className="text-sm text-[#e6edf3]">
                                            {vm.diskFiles.length}
                                        </div>
                                    </div>
                                    <div>
                                        <div className="text-xs text-[#8b949e] mb-1">
                                            Network
                                        </div>
                                        <div className="flex flex-wrap gap-1">
                                            {vm.networkInterfaces.map(
                                                (iface, idx) => (
                                                    <span
                                                        key={idx}
                                                        className={`text-[9px] px-1.5 py-0.5 rounded ${
                                                            iface.type ===
                                                            "public"
                                                                ? "bg-blue-500/10 text-blue-400 border border-blue-500/20"
                                                                : "bg-purple-500/10 text-purple-400 border border-purple-500/20"
                                                        }`}
                                                    >
                                                        {iface.type}
                                                    </span>
                                                ),
                                            )}
                                        </div>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                ) : (
                    <div className="bg-[#161b22] border border-[#30363d] rounded p-12 text-center">
                        <div className="text-[#8b949e] text-base mb-2">
                            No virtual machines deployed
                        </div>
                        <div className="text-xs text-[#8b949e]">
                            Click &quot;Deploy VM&quot; to get started
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}
