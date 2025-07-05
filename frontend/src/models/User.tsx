export interface UserInfo {
    name: string;
    canStart: boolean;
    canExtend: boolean;
    canStop: boolean;
    allowedServers: string[];
}