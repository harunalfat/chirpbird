import { Injectable } from "@angular/core";
import { Observable, Subject } from "rxjs";
import { Message } from "../model/message";
import { User } from "../model/user";
import { IStoreUserService } from "./i-store-user.service";
import { Axios, AxiosResponse } from "axios";
import { Channel } from "../model/channel";
import { ISocketService } from "./i-socket-service";
import { $ } from "protractor";
import { environment } from "src/environments/environment";

interface Response<T> {
    data: T
    errors: string[]
}

@Injectable({
    providedIn: 'root',
})
export class ServerStoreUserService implements IStoreUserService {

    constructor(
        private socketService: ISocketService,
    ) {}

    private initialChannelSource = new Subject<Channel>();
    private changeChannelSource = new Subject<any>();
    private axios = new Axios({
        baseURL: environment.backendUrl,
    })

    getChangeChannelObservable(): Observable<any> {
        return this.changeChannelSource.asObservable();
    }

    getInitChannelObservable(): Observable<any> {
        return this.initialChannelSource.asObservable();
    }

    getStoredUser() {
        const user = JSON.parse(sessionStorage.getItem("user"))
        return user
    }

    getAllUsers(): string[] {
        throw new Error("Method not implemented.");
    }

    async storeUser(user: User): Promise<User> {        
        const res = await this.axios.post<string>('/users', JSON.stringify(user))
        if (res.status != 201) {
            console.error("Failed to store user")
            return
        }

        const resp: Response<User> = JSON.parse(res.data)
        sessionStorage.setItem("user", JSON.stringify(resp.data))
        return resp.data
    }

    fetchChannelByName(channelName: string): Promise<Channel> {
        return this.socketService.fetchChannelByNameRpc(channelName)
    }

    getAllChannels(): Channel[] {
        const res = JSON.parse(sessionStorage.getItem("channelList"))
        return res
    }

    async addChannel(channel: Channel, creatorId: string, isPrivate = false): Promise<Channel> {
        channel.creatorId = creatorId
        channel.isPrivate = isPrivate
        channel.hashIdentifier = ""

        const channels = this.getAllChannels();
        if (channels.find(c => c.name === channel.name)) return null

        const resp = await this.socketService.addChannelRpc(channel)

        if (channels.length < 1) {
            channels.push(resp)
            sessionStorage.setItem("channelList", JSON.stringify(channels))
            return resp
        }
        if (channels.find(c => c.id !== resp.id)) {
            channels.push(resp)
            channels.sort((a,b) => a.name < b.name ? -1 : 1)
            sessionStorage.setItem("channelList", JSON.stringify(channels))
        }

        return resp
    }

    storeAllMessages(messages: Message[], channelId: string) {
        sessionStorage.setItem(channelId, JSON.stringify(messages))
    }

    async storeMessage(message: Message, channelId: string) {
        const messages = await this.getMessages(channelId, false) || []
        messages.push(message.data)
        this.storeAllMessages(messages, channelId)
    }

    async getMessages(channelId: string, fetchFromServer = true): Promise<Message[]> {
        var messages: Message[] = JSON.parse(sessionStorage.getItem(channelId)) || [];
        if (messages.length === 0 && fetchFromServer) {
            messages = (await this.socketService.fetchChannelMessagesRpc(channelId)) || []
            sessionStorage.setItem(channelId, JSON.stringify(messages))
        }

        return messages
    }

    announceInitialChannel(channel: Channel) {
        this.initialChannelSource.next(channel)
    }

    announceChangeChannel(data: any) {
        this.changeChannelSource.next(data)
    }

    searchUsersByUsername(username: string): Promise<User[]> {
        throw new Error("Method not implemented.");
    }
}