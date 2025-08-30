export interface State {
    servers: Map<string, ServerState>
}

export interface ServerState {
    startedAt: number;
    extensions: number[];
    endsAt: number;
    status: string;
}