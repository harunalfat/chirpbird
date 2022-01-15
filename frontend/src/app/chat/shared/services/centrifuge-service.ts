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
const CREATE_CHANNEL = "create_channel"
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
        if (!channelId) {
            this.centrifuge.on("message", cb)
            return
        }
        this.centrifuge.subscribe(channelId, cb)
    }

    subscribeServer<T>(cb: (ctx: ServerEventWrapper<T>) => void): void {
        this.centrifuge.on("publish", cb)
    }

    async initSocket(token: string): Promise<void> {
        const url = this.BASE_URL + `?userId=${token}`;
        this.centrifuge = new Centrifuge(url, {
            debug: true,
        });
        this.centrifuge.connect();
        await this.waitForNodeConnected()
        
        await this.centrifuge.send("hello world")
    }

    isConnected(): boolean {
        return this.centrifuge.isConnected()
    }

    send<T>(channelId: string, message: T): void {
        if (!channelId) {
            this.centrifuge.send(message)
            return
        }
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

    private async waitForNodeConnected(): Promise<boolean> {
        var maxTry = 5;
        while (maxTry > 0 && !this.isConnected()) {
            maxTry--;
            await new Promise((resolve) => setTimeout(resolve, 1000))
        }
        
        return this.isConnected();
    }

    async fetchChannelMessagesRpc(channelId: string): Promise<Message[]> {
        const connected = await this.waitForNodeConnected()
        if (!connected) return [];
        
        const data : ServerEventWrapper<string> = {
            data: channelId,
        }

        const response = await this.centrifuge.namedRPC(FETCH_MESSAGE, data);
        return response.data
    }

    async addChannelRpc(channel: Channel): Promise<Channel> {
        const connected = await this.waitForNodeConnected()
        if (!connected) return;

        const data : ServerEventWrapper<Channel> = {
            data: channel,
        }

        const response = await this.centrifuge.namedRPC(CREATE_CHANNEL, data)
        return response.data
    }

}