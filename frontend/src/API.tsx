import axios, { AxiosError } from 'axios';
import type { State } from './models/State';
import type { UserInfo } from './models/User';

export const GetState = async (password: string): Promise<State> => {
    try {
        const response = await axios.post(`${import.meta.env.VITE_BACKEND_URL}/status`, {password: password});
        return response.data as State;
    } catch (error) {
        return (error as AxiosError).request.response
    }
    
    // let s: string = `{"servers":{"astroneer":{"startedAt":0,"extensions":[],"endsAt":0},"satisfactory":{"startedAt":0,"extensions":[],"endsAt":1751464089}}}`
    // return JSON.parse(s) as State;
};

export const Action = async (name: string, password: string, action: string): Promise<string> => {
    try {
        await axios.post(`${import.meta.env.VITE_BACKEND_URL}/server/${name}/${action}`, {password: password});
        return ""
    } catch (error) {
        return (error as AxiosError).request.response
    }
}

export const Authenticate = async (password: string): Promise<UserInfo | string> => {
    try {
        let response = await axios.post(`${import.meta.env.VITE_BACKEND_URL}/authenticate`, {password: password});
        return response.data as UserInfo
    } catch (error) {
        return (error as AxiosError).request.response
    }
};