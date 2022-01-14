import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { Channel } from "../model/channel";
import { Message } from "../model/message";
import { User } from "../model/user";

@Injectable()
export abstract class IStoreUserService {
    abstract getStoredUser(): User
    abstract getAllUsers(): string[]
    abstract storeUser(user: User): Promise<User>
    abstract getAllChannels(): Channel[]
    abstract fetchChannelByName(channelName: string): Promise<Channel>
    abstract addChannel(channel: Channel, creatorId: string, isPrivate: boolean): Promise<Channel>
    abstract storeAllMessages(messages: Message[], channelId: string)
    abstract storeMessage(message: Message, channelId: string)
    abstract getMessages(channelId: string): Promise<Message[]>
    abstract announceInitialChannel(channel: Channel)
    abstract announceChangeChannel(channel: Channel)
    abstract getChangeChannelObservable(): Observable<Channel>
    abstract getInitChannelObservable(): Observable<Channel>
    abstract searchUsersByUsername(username: string): Promise<User[]>
}