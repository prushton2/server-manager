import axios from 'axios';
import type { State } from './models/State';

export const getState = async (): Promise<State> => {
    // try {
        // const response = await axios.get(`${import.meta.env.VITE_BACKEND_URL}/containerInfo`);
        // return response.data as State[];
    // } catch (error) {
    //     console.error('Error fetching container info:', error);
    //     throw new Error('Failed to fetch container info');
    // }
    let s: string = `{"servers":{"astroneer":{"startedAt":0,"extensions":[],"endsAt":2751344238},"satisfactory":{"startedAt":0,"extensions":[],"endsAt":0}}}`

    return JSON.parse(s) as State;
};