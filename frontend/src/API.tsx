import axios, { AxiosError } from 'axios';
import type { State } from './models/State';

export const GetState = async (): Promise<State> => {
    try {
        const response = await axios.get(`${import.meta.env.VITE_BACKEND_URL}/status`);
        return response.data as State;
    } catch (error) {
        return (error as AxiosError).request.response
    }
    
    // let s: string = `{"servers":{"astroneer":{"startedAt":0,"extensions":[],"endsAt":0},"satisfactory":{"startedAt":0,"extensions":[],"endsAt":1751464089}}}`
    // return JSON.parse(s) as State;
};

export const Start = async (name: string): Promise<string> => {
    try {
        await axios.get(`${import.meta.env.VITE_BACKEND_URL}/server/${name}/start`);
        return ""
    } catch (error) {
        return (error as AxiosError).request.response
    }
}

export const Extend = async (name: string): Promise<string> => {
    try {
        await axios.get(`${import.meta.env.VITE_BACKEND_URL}/server/${name}/extend`);
        return ""
    } catch (error) {
        return (error as AxiosError).request.response
    }
};