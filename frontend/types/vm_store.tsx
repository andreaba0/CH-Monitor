import { create } from "zustand";
import { VMConfig } from "./vm_config";

export interface NetworkInterface {
    id: string;
    type: "public" | "private";
    vpcId?: string;
    macAddress: string;
    tapDevice: string;
    ipAddress?: string;
}

export interface DiskFile {
    name: string;
    size: number;
    path: string;
    serial: string;
    isBootDisk: boolean;
}

export interface VM {
    id: string;
    name: string;
    status: "running" | "stopped" | "paused" | "starting";
    cpus: number;
    memory: number;
    diskFiles: DiskFile[];
    uptime: string;
    cpuUsage: number;
    memoryUsage: number;
    cpuHistory: { time: string; value: number }[];
    memoryHistory: { time: string; value: number }[];
    networkInterfaces: NetworkInterface[];
    kernelPath: string;
    yamlPath: string;
    yamlConfig: string;
    createdAt: string;
}

export interface VMStore {
    vms: VM[];
    addVM: (config: VMConfig) => void;
    updateVM: (id: string, updates: Partial<VM>) => void;
    deleteVM: (id: string) => void;
    getVM: (id: string) => VM | undefined;
    startVM: (id: string) => void;
    stopVM: (id: string) => void;
    restartVM: (id: string) => void;
    addDisk: (vmId: string, size: number) => void;
    removeDisk: (vmId: string, serial: string) => void;
}

export const useVMStore = create<VMStore>((set, get) => ({
    vms: [
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
    ],

    addVM: (config: VMConfig) => {
        const vms = get().vms;

        // Parse YAML to extract basic info
        const cpusMatch = config.yamlConfig.match(/boot_vcpus:\s*(\d+)/);
        const memoryMatch = config.yamlConfig.match(/size:\s*(\d+)M/);

        const newVM: VM = {
            id: `vm-${String(vms.length + 1).padStart(3, "0")}`,
            name: config.name,
            status: "stopped",
            cpus: cpusMatch ? parseInt(cpusMatch[1]) : 2,
            memory: memoryMatch ? parseInt(memoryMatch[1]) : 2048,
            diskFiles: config.disks.map((disk, index) => ({
                name: `disk${index + 1}.img`,
                size: disk.size * 1073741824,
                path: `/var/lib/vms/${config.name}/disk${index + 1}.img`,
                serial: disk.serial,
                isBootDisk: index === 0,
            })),
            uptime: "0d 0h 0m",
            cpuUsage: 0,
            memoryUsage: 0,
            cpuHistory: [],
            memoryHistory: [],
            networkInterfaces: [],
            kernelPath: `/var/lib/cloud-hypervisor/images/${config.osImage}/vmlinuz`,
            yamlPath: `/etc/cloud-hypervisor/vms/${config.name}.yaml`,
            yamlConfig: config.yamlConfig,
            createdAt: new Date().toISOString(),
        };
        set({ vms: [...vms, newVM] });
    },

    updateVM: (id: string, updates: Partial<VM>) => {
        set((state) => ({
            vms: state.vms.map((vm) =>
                vm.id === id ? { ...vm, ...updates } : vm,
            ),
        }));
    },

    deleteVM: (id: string) => {
        set((state) => ({ vms: state.vms.filter((vm) => vm.id !== id) }));
    },

    getVM: (id: string) => {
        return get().vms.find((vm) => vm.id === id);
    },

    startVM: (id: string) => {
        get().updateVM(id, { status: "starting" });
        setTimeout(() => {
            get().updateVM(id, { status: "running", uptime: "0d 0h 1m" });
        }, 1500);
    },

    stopVM: (id: string) => {
        get().updateVM(id, {
            status: "stopped",
            uptime: "0d 0h 0m",
            cpuUsage: 0,
            memoryUsage: 0,
        });
    },

    restartVM: (id: string) => {
        get().stopVM(id);
        setTimeout(() => get().startVM(id), 1000);
    },

    addDisk: (vmId: string, size: number) => {
        const vm = get().vms.find((v) => v.id === vmId);
        if (!vm || vm.status !== "stopped") return;

        const generateSerial = () => {
            const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
            let serial = "CH-";
            for (let i = 0; i < 12; i++) {
                serial += chars.charAt(
                    Math.floor(Math.random() * chars.length),
                );
            }
            return serial;
        };

        const diskIndex = vm.diskFiles.length + 1;
        const newDisk: DiskFile = {
            name: `disk${diskIndex}.img`,
            size: size * 1073741824,
            path: `/var/lib/vms/${vm.name}/disk${diskIndex}.img`,
            serial: generateSerial(),
            isBootDisk: false,
        };

        get().updateVM(vmId, {
            diskFiles: [...vm.diskFiles, newDisk],
        });
    },

    removeDisk: (vmId: string, serial: string) => {
        const vm = get().vms.find((v) => v.id === vmId);
        if (!vm || vm.status !== "stopped") return;

        const disk = vm.diskFiles.find((d) => d.serial === serial);
        if (!disk || disk.isBootDisk) return;

        get().updateVM(vmId, {
            diskFiles: vm.diskFiles.filter((disk) => disk.serial !== serial),
        });
    },
}));
