import { Injectable } from "@angular/core";
import * as Centrifuge from "centrifuge";
import { Observable } from "rxjs";
import { environment } from "src/environments/environment";
import { Channel } from "../model/channel";
import { Event } from "../model/event";
import { Message, ServerEventWrapper } from "../model/message";
import { User } from "../model/user";
import { ISocketService } from "./i-socket-service";

const FETCH_MESSAGE = "fetch_message"
const FETCH_CHANNEL = "fetch_channel"
const SEARCH_USERS = "search_users"

@Injectable({
    providedIn: 'root'
})
export class CentrifugeService implements ISocketService {

    private BASE_URL = environment.backendWsUrl;
    private centrifuge: Centrifuge

    constructor(){}

    async searchUsersByUsernameRpc(username: string): Promise<User[]> {
        const data: ServerEventWrapper<string> = {
            data: username,
        }
        
        const response = await this.centrifuge.namedRPC(SEARCH_USERS, data)
        return response.data
    }

    subscribe<T>(channelId: string, cb: (message: ServerEventWrapper<T>) => void): void {
        this.centrifuge.subscribe(channelId, cb)
    }

    initSocket(token: string): void {
        const url = this.BASE_URL + `?userId=${token}`;
        this.centrifuge = new Centrifuge(url, {
            debug: true,
            
        });
        this.centrifuge.connect();
    }

    isConnected(): boolean {
        return this.centrifuge.isConnected()
    }

    send<T>(channelId: string, message: T): void {
        this.centrifuge.publish(channelId, message)
    }

    onMessage(): Observable<Message> {
        return new Observable<Message>(observer => {
            observer.next()
        })
    }

    onEvent(event: Event): Observable<any> {
        return new Observable<Event>(observer => {
            this.centrifuge.on(event, (data) => observer.next(data))
        })
    }

    async fetchChannelByNameRpc(channelName: string): Promise<Channel> {
        const data : ServerEventWrapper<string> = {
            data: channelName,
        }
        const response = await this.centrifuge.namedRPC(FETCH_CHANNEL, data)
        return response;
    }

    async fetchChannelMessagesRpc(channelId: string): Promise<Message[]> {
        var maxTry = 5;
        while (maxTry > 0 && !this.isConnected()) {
            maxTry--;
            await new Promise((resolve) => setTimeout(resolve, 1000))
        }
        
        if (!this.isConnected()) return [];
        
        const data : ServerEventWrapper<string> = {
            data: channelId,
        }

        const response = await this.centrifuge.namedRPC(FETCH_MESSAGE, data);
        return response.data
    }

}