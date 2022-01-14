
import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { Channel } from "../model/channel";
import { Event } from "../model/event";
import { Message, ServerEventWrapper } from "../model/message";
import { User } from "../model/user";

@Injectable()
export abstract class ISocketService {
    abstract initSocket(any: any): void
    abstract isConnected(): boolean
    abstract send<T>(channelId: string, message: T): void
    abstract subscribe<T>(channelId: string, cb: (message: ServerEventWrapper<T>) => void): void
    abstract onMessage(): void
    abstract onEvent(event: Event): Observable<any>
    abstract fetchChannelByNameRpc(channelName: string): Promise<Channel>
    abstract fetchChannelMessagesRpc(channelId: string): Promise<Message[]>
    abstract searchUsersByUsernameRpc(username: string): Promise<User[]>
}