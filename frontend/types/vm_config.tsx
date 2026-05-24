export interface DiskConfig {
    id: string;
    size: number;
    serial: string;
    isBootDisk?: boolean;
}

export interface VMConfig {
    name: string;
    osImage: string;
    disks: DiskConfig[];
    sshKeyFile?: File;
    yamlFile?: File;
    yamlConfig: string;
}
