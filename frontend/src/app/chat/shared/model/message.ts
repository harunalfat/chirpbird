import {User} from './user';
import {Action} from './action';

export interface Message {
    id?: string
    created?: Date
    sender?: User;
    channelId?: string;
    data?: any;
    action?: Action;
}

export interface ServerEventWrapper<T> {
    data: T
    gen?: any
    info?: {
        user: string
        client: string
    }
    offset?: number
    seq?: number
}