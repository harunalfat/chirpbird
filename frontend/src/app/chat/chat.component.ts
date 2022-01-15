import { AfterViewInit, Component, ElementRef, OnInit, QueryList, ViewChild, ViewChildren } from '@angular/core';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { MatList, MatListItem } from '@angular/material/list';
import { TranslateService } from '@ngx-translate/core';
import { Subject } from 'rxjs';
import { DialogPrivateComponent } from './dialog-private/dialog-private.component';
import { DialogUserType } from './dialog-user/dialog-user-type';
import { DialogUserComponent } from './dialog-user/dialog-user.component';
import { Action } from './shared/model/action';
import { Channel } from './shared/model/channel';
import { Event } from './shared/model/event';
import { Message, ServerEventWrapper } from './shared/model/message';
import { User } from './shared/model/user';
import { ISocketService } from './shared/services/i-socket-service';
import { IStoreUserService } from './shared/services/i-store-user.service';



const AVATAR_URL = 'https://api.adorable.io/avatars/285';

@Component({
  selector: 'app-chat',
  templateUrl: './chat.component.html',
  styleUrls: ['./chat.component.css']
})
export class ChatComponent implements OnInit, AfterViewInit {
  action = Action;
  user: User;
  currentChannel: Channel;
  messages: Message[] = [];
  messageContent: string;
  storedUserName: string;
  dialogUserRef: MatDialogRef<DialogUserComponent> | null;
  dialogPrivateChannelRef: MatDialogRef<DialogPrivateComponent> | null;
  defaultDialogUserParams: any = {
    disableClose: true,
    data: {
      title: 'Welcome',
      title_pt: 'Bem-vindo',
      dialogType: DialogUserType.NEW
    }
  };
  
  public static returned: Subject<any> = new Subject();

  // getting a reference to the overall list, which is the parent container of the list items
  @ViewChild(MatList, { read: ElementRef, static: true }) matList: ElementRef;

  // getting a reference to the items/messages within the list
  @ViewChildren(MatListItem, { read: ElementRef }) matListItems: QueryList<MatListItem>;

  constructor(private socketService: ISocketService,
    private storedUser: IStoreUserService,
    public dialog: MatDialog, private translate: TranslateService) {
    translate.setDefaultLang('en');
    this.storedUser.getChangeChannelObservable().subscribe(async passedChannel => {
      if (passedChannel.isNewlyAdded) {
        this.socketService.subscribe(passedChannel.id, (message) => {
          this.onSubscribeMessage(message, passedChannel)
        })
      }

      this.messages = (await this.storedUser.getMessages(passedChannel.id)) || []
      this.currentChannel = passedChannel
    })
  }

  ngOnInit(): void {
    this.user = this.storedUser.getStoredUser();
    // Using timeout due to https://github.com/angular/angular/issues/14748
    setTimeout(() => {
      this.openNewUserPopup(this.defaultDialogUserParams);
    }, 0);
  }

  ngAfterViewInit(): void {
    // subscribing to any changes in the list of items / messages
    this.matListItems.changes.subscribe(elements => {
      this.scrollToBottom();
    });
  }

  // auto-scroll fix: inspired by this stack overflow post
  // https://stackoverflow.com/questions/35232731/angular2-scroll-to-bottom-chat-style
  private scrollToBottom(): void {
    try {
      this.matList.nativeElement.scrollTop = this.matList.nativeElement.scrollHeight;
    } catch (err) {
    }
  }

  onSubscribeMessage(message: ServerEventWrapper<Message>, channel: Channel) {
    this.storedUser.storeMessage(message, message.data.channelId)
    if (this.currentChannel.id === channel.id) {
      this.messages.push(message.data)
      this.storedUser.getMessages(message.data.channelId);
    }
  }

  private async initIoConnection(): Promise<void> {
    await this.socketService.initSocket(this.user.id);
    if (this.user.channels?.length < 1) {
      const lobby: Channel = {
        name: "Lobby",
        isPrivate: false,
        createdAt: new Date(),
        creatorId: this.user.id,
        hashIdentifier: "",
        participants: [],
      }

      const newChannel = await this.storedUser.addChannel(lobby, this.user.id, false)
      
      this.user.channels.push(newChannel);
    }

    for (const channel of this.user.channels) {

      if (channel.name === "Lobby") {
        this.currentChannel = channel
        this.messages = await this.storedUser.getMessages(this.currentChannel.id) || [];
        this.storedUser.announceInitialChannel(channel)
      }

      this.socketService.subscribe(channel.id, (message) => {
        this.onSubscribeMessage(message, channel)
      })
    }

    this.socketService.onEvent(Event.SUBSCRIBE)
      .subscribe(data => {
        console.log(`Server subscribe user [${this.user.username}] to channel [${data.channel.name}]`)
      })

    this.socketService.onEvent(Event.CONNECT)
      .subscribe((data) => {
        console.log('connected');
        console.log(data)
      });

    this.socketService.onEvent(Event.DISCONNECT)
      .subscribe((data) => {
        console.log('disconnected');
        console.log(data)
      });
  }

  public onClickPrivateMessage() {
    this.dialogPrivateChannelRef = this.dialog.open(DialogPrivateComponent, {
      data: this.user,
    });
    this.dialogPrivateChannelRef.afterClosed().subscribe(async dialogParams => {
      if (!dialogParams) return

      const friend = {
        username: dialogParams.username,
        id: dialogParams.id,
      }

      const joinName: string[] = [this.user.username, friend.username]
      joinName.sort((l,r) => l < r ? -1 : 1)

      const privateChannelName = joinName.join(" & ")
      const channelInput: Channel = {
        name: privateChannelName,
        participants: [{username: this.user.username, id: this.user.id}, friend],
        creatorId: this.user.id,
        hashIdentifier: "",
        isPrivate: true,
      }

      const channel = await this.storedUser.addChannel(channelInput, this.user.id, true)
      if (!channel) return

      this.messages = await this.storedUser.getMessages(channel.name)
      this.currentChannel = channel
      this.storedUser.announceInitialChannel(this.currentChannel)
      this.socketService.send(null, channel)
    })
  }

  private openNewUserPopup(params): void {
    this.dialogUserRef = this.dialog.open(DialogUserComponent, params);
    this.dialogUserRef.afterClosed().subscribe(async paramsDialog => {
      if (!paramsDialog) {
        return;
      }

      if (paramsDialog.dialogType === DialogUserType.NEW) {
        sessionStorage.clear();
        const result = await this.storedUser.storeUser({
          username: paramsDialog.username,
        })
        result.channels?.sort((a, b) => a.name < b.name ? -1 : 1)
        this.user = result;
        if (!result.channels) {
          this.user.channels = []
        }
        sessionStorage.setItem("channelList", JSON.stringify(this.user.channels))
        await this.initIoConnection();
        // subscribe to server pushing new channel
        this.socketService.subscribeServer<Channel>(async (ctx) => {
          if (ctx.channel !== "SERVER_NOTIFICATION") return

          this.socketService.subscribe<Message>(ctx.data.id, (message) => {
            this.onSubscribeMessage(message, ctx.data)
          })
        })
      }
      console.log("FINISH INITIALIZATION")
    });
  }

  public sendMessage(message: string, channelId: string): void {
    if (!message) {
      return;
    }

    const msg: Message = {
      sender: {
        id: this.user.id,
        username: this.user.username,
      },
      data: message,
      channelId: channelId,
      created: new Date(),
    }

    this.socketService.send(msg.channelId, msg);
    this.messageContent = null;
  }

  switchLanguage(language: string) {
    this.translate.use(language);
  }

}
