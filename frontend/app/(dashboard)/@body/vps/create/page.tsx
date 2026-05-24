"use client";

import { useState } from "react";
import {
    ArrowLeft,
    Box,
    HardDrive,
    Cpu,
    Key,
    Upload,
    Plus,
    Trash2,
    Network,
} from "lucide-react";
import Link from "next/link";

interface OSImage {
    id: string;
    name: string;
    version: string;
    description: string;
    defaultDiskSize: number;
}

interface VMType {
    id: string;
    name: string;
    vcpus: number;
    memory: number;
    description: string;
}

interface DiskConfig {
    id: string;
    size: number;
    serial: string;
    isBootDisk?: boolean;
}

const availableOSImages: OSImage[] = [
    {
        id: "ubuntu-22.04",
        name: "Ubuntu Server",
        version: "22.04 LTS",
        description: "Ubuntu 22.04 LTS (Jammy Jellyfish)",
        defaultDiskSize: 20,
    },
    {
        id: "ubuntu-24.04",
        name: "Ubuntu Server",
        version: "24.04 LTS",
        description: "Ubuntu 24.04 LTS (Noble Numbat)",
        defaultDiskSize: 20,
    },
    {
        id: "debian-12",
        name: "Debian",
        version: "12 (Bookworm)",
        description: "Debian 12 Bookworm",
        defaultDiskSize: 20,
    },
    {
        id: "fedora-40",
        name: "Fedora Server",
        version: "40",
        description: "Fedora Server 40",
        defaultDiskSize: 25,
    },
    {
        id: "alpine-3.19",
        name: "Alpine Linux",
        version: "3.19",
        description: "Alpine Linux 3.19 (lightweight)",
        defaultDiskSize: 5,
    },
    {
        id: "arch-latest",
        name: "Arch Linux",
        version: "Latest",
        description: "Arch Linux rolling release",
        defaultDiskSize: 20,
    },
];

const vmTypes: VMType[] = [
    {
        id: "t2.small",
        name: "t2.small",
        vcpus: 1,
        memory: 2048,
        description: "1 vCPU, 2 GB RAM",
    },
    {
        id: "t2.medium",
        name: "t2.medium",
        vcpus: 2,
        memory: 4096,
        description: "2 vCPU, 4 GB RAM",
    },
    {
        id: "t2.large",
        name: "t2.large",
        vcpus: 2,
        memory: 8192,
        description: "2 vCPU, 8 GB RAM",
    },
    {
        id: "m5.large",
        name: "m5.large",
        vcpus: 4,
        memory: 8192,
        description: "4 vCPU, 8 GB RAM",
    },
    {
        id: "m5.xlarge",
        name: "m5.xlarge",
        vcpus: 4,
        memory: 16384,
        description: "4 vCPU, 16 GB RAM",
    },
    {
        id: "c5.large",
        name: "c5.large",
        vcpus: 2,
        memory: 4096,
        description: "2 vCPU, 4 GB RAM (Compute optimized)",
    },
    {
        id: "c5.xlarge",
        name: "c5.xlarge",
        vcpus: 4,
        memory: 8192,
        description: "4 vCPU, 8 GB RAM (Compute optimized)",
    },
    {
        id: "r5.large",
        name: "r5.large",
        vcpus: 2,
        memory: 16384,
        description: "2 vCPU, 16 GB RAM (Memory optimized)",
    },
];

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

export default function DeployVM() {
    const [vmName, setVmName] = useState("");
    const [selectedOS, setSelectedOS] = useState("");
    const [selectedVMType, setSelectedVMType] = useState("");
    const [disks, setDisks] = useState<DiskConfig[]>([]);
    const [sshKeyFile, setSshKeyFile] = useState<File | null>(null);
    const [selectedVPCs, setSelectedVPCs] = useState<string[]>([]);

    const generateSerial = () => {
        const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
        let serial = "CH-";
        for (let i = 0; i < 12; i++) {
            serial += chars.charAt(Math.floor(Math.random() * chars.length));
        }
        return serial;
    };

    const handleOSChange = (osId: string) => {
        const os = availableOSImages.find((img) => img.id === osId);
        if (!os) return;

        const bootDisk: DiskConfig = {
            id: `disk-${Date.now()}`,
            size: os.defaultDiskSize,
            serial: generateSerial(),
            isBootDisk: true,
        };

        setSelectedOS(osId);
        setDisks([bootDisk]);
    };

    const addDisk = () => {
        const selectedOSData = availableOSImages.find(
            (img) => img.id === selectedOS,
        );
        const newDisk: DiskConfig = {
            id: `disk-${Date.now()}`,
            size: selectedOSData?.defaultDiskSize || 20,
            serial: generateSerial(),
            isBootDisk: false,
        };
        setDisks([...disks, newDisk]);
    };

    const removeDisk = (id: string) => {
        const disk = disks.find((d) => d.id === id);
        if (disk?.isBootDisk) return;
        setDisks(disks.filter((disk) => disk.id !== id));
    };

    const updateDiskSize = (id: string, size: number) => {
        setDisks(
            disks.map((disk) => (disk.id === id ? { ...disk, size } : disk)),
        );
    };

    const handleSshKeyUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;
        setSshKeyFile(file);
    };

    const toggleVPC = (vpcId: string) => {
        if (selectedVPCs.includes(vpcId)) {
            setSelectedVPCs(selectedVPCs.filter((id) => id !== vpcId));
        } else {
            setSelectedVPCs([...selectedVPCs, vpcId]);
        }
    };

    const handleDeploy = () => {
        const vmType = vmTypes.find((t) => t.id === selectedVMType);
        if (!vmType) return;
    };

    const selectedOSData = availableOSImages.find(
        (img) => img.id === selectedOS,
    );

    const isValid =
        vmName &&
        selectedOS &&
        selectedVMType &&
        disks.length > 0 &&
        sshKeyFile;

    return (
        <div className="p-6 pt-16 md:pt-6">
            <button className="transition-colors mb-4">
                <Link
                    href="/vps/list"
                    className="text-xs flex items-center gap-2 text-[#8b949e] hover:text-[#e6edf3]"
                >
                    <ArrowLeft className="w-3 h-3" />
                    Back to Dashboard
                </Link>
            </button>

            <div className="mb-6">
                <h1 className="text-xl text-[#e6edf3] mb-1">
                    Deploy New Virtual Machine
                </h1>
                <p className="text-[#8b949e] text-xs">
                    Configure and deploy a new VM instance
                </p>
            </div>

            <div className="space-y-6">
                {/* VM Name */}
                <div className="bg-[#161b22] border border-[#30363d] rounded p-5">
                    <label className="block text-sm text-[#e6edf3] mb-3">
                        VM Name <span className="text-red-400">*</span>
                    </label>
                    <input
                        type="text"
                        value={vmName}
                        onChange={(e) => setVmName(e.target.value)}
                        className="w-full px-3 py-2 bg-[#0d1117] border border-[#30363d] rounded text-[#e6edf3] text-sm focus:outline-none focus:border-blue-500 transition-colors"
                        placeholder="my-vm"
                    />
                </div>

                {/* Operating System */}
                <div className="bg-[#161b22] border border-[#30363d] rounded p-5">
                    <label className="block text-sm text-[#e6edf3] mb-3">
                        Operating System <span className="text-red-400">*</span>
                    </label>
                    <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                        {availableOSImages.map((os) => (
                            <label
                                key={os.id}
                                className={`cursor-pointer border rounded p-4 transition-colors ${
                                    selectedOS === os.id
                                        ? "border-blue-500 bg-blue-500/10"
                                        : "border-[#30363d] hover:border-[#484f58]"
                                }`}
                            >
                                <input
                                    type="radio"
                                    name="osImage"
                                    value={os.id}
                                    checked={selectedOS === os.id}
                                    onChange={(e) =>
                                        handleOSChange(e.target.value)
                                    }
                                    className="sr-only"
                                />
                                <div className="flex flex-col items-center text-center">
                                    <Box className="w-8 h-8 text-blue-400 mb-2" />
                                    <div className="text-sm text-[#e6edf3] mb-0.5">
                                        {os.name}
                                    </div>
                                    <div className="text-xs text-[#8b949e] mb-1">
                                        {os.version}
                                    </div>
                                    <div className="text-[10px] text-[#8b949e]">
                                        {os.defaultDiskSize} GB default
                                    </div>
                                </div>
                            </label>
                        ))}
                    </div>
                </div>

                {/* VM Type */}
                <div className="bg-[#161b22] border border-[#30363d] rounded p-5">
                    <label className="block text-sm text-[#e6edf3] mb-3">
                        VM Type <span className="text-red-400">*</span>
                    </label>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                        {vmTypes.map((type) => (
                            <label
                                key={type.id}
                                className={`cursor-pointer border rounded p-4 transition-colors ${
                                    selectedVMType === type.id
                                        ? "border-blue-500 bg-blue-500/10"
                                        : "border-[#30363d] hover:border-[#484f58]"
                                }`}
                            >
                                <input
                                    type="radio"
                                    name="vmType"
                                    value={type.id}
                                    checked={selectedVMType === type.id}
                                    onChange={(e) =>
                                        setSelectedVMType(e.target.value)
                                    }
                                    className="sr-only"
                                />
                                <div className="flex items-center gap-3">
                                    <Cpu className="w-5 h-5 text-blue-400 shrink-0" />
                                    <div className="flex-1">
                                        <div className="text-sm text-[#e6edf3] mb-0.5">
                                            {type.name}
                                        </div>
                                        <div className="text-xs text-[#8b949e]">
                                            {type.description}
                                        </div>
                                    </div>
                                </div>
                            </label>
                        ))}
                    </div>
                </div>

                {/* Disks */}
                <div className="bg-[#161b22] border border-[#30363d] rounded p-5">
                    <div className="flex items-center justify-between mb-3">
                        <label className="text-sm text-[#e6edf3] flex items-center gap-2">
                            <HardDrive className="w-4 h-4" />
                            Storage Configuration{" "}
                            <span className="text-red-400">*</span>
                        </label>
                        <button
                            type="button"
                            onClick={addDisk}
                            disabled={!selectedOS}
                            className="px-3 py-1.5 bg-[#30363d] hover:bg-[#484f58] text-[#e6edf3] rounded text-xs flex items-center gap-1.5 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            <Plus className="w-3.5 h-3.5" />
                            Add Additional Disk
                        </button>
                    </div>

                    {selectedOS && disks.length > 0 && (
                        <p className="text-xs text-[#8b949e] mb-4">
                            Boot disk size can be adjusted. Add more disks for
                            additional storage.
                        </p>
                    )}

                    {disks.length > 0 ? (
                        <div className="space-y-3">
                            {disks.map((disk, index) => (
                                <div
                                    key={disk.id}
                                    className={`border rounded p-4 ${
                                        disk.isBootDisk
                                            ? "border-blue-500/30 bg-blue-500/5"
                                            : "border-[#30363d] bg-[#0d1117]"
                                    }`}
                                >
                                    <div className="flex items-center gap-2 mb-3">
                                        <HardDrive className="w-4 h-4 text-[#8b949e] shrink-0" />
                                        <div className="flex-1">
                                            <div className="flex items-center gap-2">
                                                <div className="text-sm text-[#e6edf3]">
                                                    Disk {index + 1}
                                                </div>
                                                {disk.isBootDisk && (
                                                    <span className="text-[9px] px-1.5 py-0.5 bg-blue-500/20 text-blue-400 border border-blue-500/30 rounded">
                                                        BOOT
                                                    </span>
                                                )}
                                            </div>
                                            <div className="text-xs text-[#8b949e]">
                                                Serial: {disk.serial}
                                            </div>
                                        </div>
                                        {!disk.isBootDisk && (
                                            <button
                                                type="button"
                                                onClick={() =>
                                                    removeDisk(disk.id)
                                                }
                                                className="p-1.5 hover:bg-[#30363d] rounded transition-colors text-red-400"
                                            >
                                                <Trash2 className="w-4 h-4" />
                                            </button>
                                        )}
                                    </div>
                                    <div className="flex items-center gap-2">
                                        <input
                                            type="number"
                                            value={disk.size}
                                            onChange={(e) =>
                                                updateDiskSize(
                                                    disk.id,
                                                    parseInt(e.target.value) ||
                                                        1,
                                                )
                                            }
                                            min={
                                                selectedOSData?.defaultDiskSize ||
                                                5
                                            }
                                            className="w-28 px-3 py-2 bg-[#161b22] border border-[#30363d] rounded text-[#e6edf3] text-sm focus:outline-none focus:border-blue-500 transition-colors"
                                        />
                                        <span className="text-sm text-[#8b949e]">
                                            GB (min:{" "}
                                            {selectedOSData?.defaultDiskSize ||
                                                5}{" "}
                                            GB)
                                        </span>
                                    </div>
                                </div>
                            ))}
                        </div>
                    ) : (
                        <div className="border border-dashed border-[#30363d] rounded p-8 text-center">
                            <HardDrive className="w-8 h-8 text-[#8b949e] mx-auto mb-2" />
                            <p className="text-sm text-[#8b949e]">
                                Select an OS to configure storage
                            </p>
                            <p className="text-xs text-[#8b949e] mt-1">
                                A boot disk will be created automatically
                            </p>
                        </div>
                    )}
                </div>

                {/* VPC Selection */}
                <div className="bg-[#161b22] border border-[#30363d] rounded p-5">
                    <div className="flex items-center justify-between mb-3">
                        <label className="text-sm text-[#e6edf3] flex items-center gap-2">
                            <Network className="w-4 h-4" />
                            Virtual Private Clouds
                        </label>
                        {vpcs.length === 0 && (
                            <button className="text-xs text-blue-400 hover:text-blue-300 transition-colors">
                                <Link href="/vpc/create">Create VPC</Link>
                            </button>
                        )}
                    </div>

                    {vpcs.length > 0 ? (
                        <>
                            <p className="text-xs text-[#8b949e] mb-3">
                                Select VPCs to attach to this VM (optional)
                            </p>
                            <div className="space-y-2">
                                {vpcs.map((vpc) => (
                                    <label
                                        key={vpc.id}
                                        className={`cursor-pointer border rounded p-3 transition-colors flex items-center gap-3 ${
                                            selectedVPCs.includes(vpc.id)
                                                ? "border-purple-500 bg-purple-500/10"
                                                : "border-[#30363d] hover:border-[#484f58]"
                                        }`}
                                    >
                                        <input
                                            type="checkbox"
                                            checked={selectedVPCs.includes(
                                                vpc.id,
                                            )}
                                            onChange={() => toggleVPC(vpc.id)}
                                            className="w-4 h-4 rounded border-[#30363d] bg-[#0d1117] text-purple-600 focus:ring-purple-500 focus:ring-offset-0"
                                        />
                                        <div className="flex-1">
                                            <div className="text-sm text-[#e6edf3] mb-0.5">
                                                {vpc.name}
                                            </div>
                                            <div className="text-xs text-[#8b949e] font-mono">
                                                {vpc.networkClass}
                                            </div>
                                        </div>
                                        <span className="text-xs text-[#8b949e]">
                                            {vpc.id}
                                        </span>
                                    </label>
                                ))}
                            </div>
                            {selectedVPCs.length > 0 && (
                                <div className="mt-3 text-xs text-[#8b949e]">
                                    {selectedVPCs.length} VPC
                                    {selectedVPCs.length > 1 ? "s" : ""}{" "}
                                    selected
                                </div>
                            )}
                        </>
                    ) : (
                        <div className="border border-dashed border-[#30363d] rounded p-8 text-center">
                            <Network className="w-8 h-8 text-[#8b949e] mx-auto mb-2" />
                            <p className="text-sm text-[#8b949e] mb-1">
                                No VPCs available
                            </p>
                            <p className="text-xs text-[#8b949e] mb-3">
                                Create VPCs to set up private networks for your
                                VMs
                            </p>
                            <button className="px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded text-sm inline-flex items-center gap-2 transition-colors">
                                <Link href="/vpc/create">
                                    <Plus className="w-4 h-4" />
                                    Create VPC
                                </Link>
                            </button>
                        </div>
                    )}
                </div>

                {/* SSH Key */}
                <div className="bg-[#161b22] border border-[#30363d] rounded p-5">
                    <label className="text-sm text-[#e6edf3] mb-3 flex items-center gap-2">
                        <Key className="w-4 h-4" />
                        SSH Public Key <span className="text-red-400">*</span>
                    </label>
                    <div className="border border-dashed border-[#30363d] rounded p-4">
                        <div className="flex items-center gap-3">
                            <Upload className="w-5 h-5 text-[#8b949e] shrink-0" />
                            <div className="flex-1">
                                <p className="text-sm text-[#8b949e] mb-2">
                                    Upload SSH public key (.pub)
                                </p>
                                <label className="cursor-pointer">
                                    <input
                                        type="file"
                                        onChange={handleSshKeyUpload}
                                        className="hidden"
                                        accept=".pub"
                                    />
                                    <span className="px-4 py-2 bg-[#30363d] hover:bg-[#484f58] text-[#e6edf3] rounded text-sm transition-colors inline-block">
                                        Choose File
                                    </span>
                                </label>
                            </div>
                        </div>
                        {sshKeyFile && (
                            <div className="mt-3 flex items-center gap-2 bg-green-500/10 border border-green-500/20 rounded px-3 py-2">
                                <Key className="w-4 h-4 text-green-400" />
                                <span className="text-sm text-green-400">
                                    {sshKeyFile.name}
                                </span>
                            </div>
                        )}
                    </div>
                </div>

                {/* Actions */}
                <div className="flex items-center justify-end gap-3 pt-4">
                    <button className="px-4 py-2 bg-[#30363d] hover:bg-[#484f58] text-[#e6edf3] rounded text-sm transition-colors">
                        <Link href="/vps/list">Cancel</Link>
                    </button>
                    <button
                        onClick={handleDeploy}
                        disabled={!isValid}
                        className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded text-sm transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        Deploy Virtual Machine
                    </button>
                </div>
            </div>
        </div>
    );
}
