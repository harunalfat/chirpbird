import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import * as socketIo from 'socket.io-client';
import { Channel } from '../model/channel';
import { Event } from '../model/event';
import { Message, ServerEventWrapper } from '../model/message';
import { User } from '../model/user';
import { ISocketService } from './i-socket-service';


const SERVER_URL = 'http://localhost:8080';

@Injectable()
export class SocketService implements ISocketService {
    searchUsersByUsernameRpc(username: string): Promise<User[]> {
        throw new Error('Method not implemented.');
    }
    fetchChannelByNameRpc(channelName: string): Promise<Channel> {
        throw new Error('Method not implemented.');
    }

    private socket;

    fetchChannelMessagesRpc(channelId: string): Promise<Message[]> {
        throw new Error('Method not implemented.');
    }

    subscribe<Message>(channelId: string, cb: (message: ServerEventWrapper<Message>) => void): void {
        throw new Error('Method not implemented.');
    }

    public initSocket(): void {
        this.socket = socketIo(SERVER_URL);
    }

    isConnected(): boolean {
        throw new Error('Method not implemented.');
    }

    public send<T>(channelId: string, message: T): void {
        this.socket.emit(channelId, message);
    }

    public onMessage(): Observable<Message> {
        return new Observable<Message>(observer => {
            this.socket.on('message', (data: Message) => observer.next(data));
        });
    }

    public onEvent(event: Event): Observable<any> {
        return new Observable<Event>(observer => {
            this.socket.on(event, () => observer.next());
        });
    }
}
