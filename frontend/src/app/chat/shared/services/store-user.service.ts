import { Injectable } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { Message } from '../model/message';
import { User } from '../model/user';
import { IStoreUserService } from './i-store-user.service';

const CHANNEL_NAMES = "channelNames";

@Injectable({
  providedIn: 'root'
})
export class StoreUserService implements IStoreUserService {

  constructor() { }

  private initialChannelSource = new Subject<string>();
  private changeChannelSource = new Subject<string>();

  getInitChannelObservable(): Observable<any> {
    return this.changeChannelSource.asObservable();
  }

  getChangeChannelObservable(): Observable<any> {
    return this.initialChannelSource.asObservable()
  }

  /**
   * getStoredUser
   */
  public getStoredUser() {
    let storedUser = sessionStorage.getItem("userName");
    return storedUser ? storedUser : "";
  }

  public getAllUsers(): string[] {
    return ['apa', 'apakah', 'apee']
  }

  /**
   * storeUser
   */
  public storeUser(userName): Promise<User> {
    sessionStorage.setItem("userName", userName);
    return Promise.resolve({
      username: userName,
    })
  }

  public getAllChannelNames(): string[] {
    const res = JSON.parse(sessionStorage.getItem(CHANNEL_NAMES));
    return res || []
  }

 public addChannel(channelName: string) {
    const allChannels = this.getAllChannelNames();
    if (!allChannels.includes(channelName)) {
      allChannels.push(channelName);
    }

    allChannels.sort((a,b) => a < b ? -1 : 1)

    sessionStorage.setItem(CHANNEL_NAMES, JSON.stringify(allChannels))
 }

 public storeAllMessages(messages: Message[], channelName: string) {
    sessionStorage.setItem(`channels:${channelName}`, JSON.stringify(messages))
 }

 public storeMessage(message: Message, channelName: string) {
    const channelMessages = this.getMessages(channelName)
    channelMessages.push(message)
    this.storeAllMessages(channelMessages, channelName)
 }

 public getMessages(channelName: string): Message[] {
    const res = JSON.parse(sessionStorage.getItem(`channels:${channelName}`))
    return res ? res : []
 }

 public announceInitialChannel(channelName: string) {
   this.initialChannelSource.next(channelName);
 }

 public announceChangeChannel(channelName: string) {
   this.changeChannelSource.next(channelName);
 }

}
