import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MaterialModule } from '../shared/material/material.module';
import { ChatComponent } from './chat.component';
import { DialogPrivateComponent } from './dialog-private/dialog-private.component';
import { DialogUserComponent } from './dialog-user/dialog-user.component';
import { CentrifugeService } from './shared/services/centrifuge-service';
import { ISocketService } from './shared/services/i-socket-service';
import { IStoreUserService } from './shared/services/i-store-user.service';
import { ServerStoreUserService } from './shared/services/server-store-user.service';
import { SocketService } from './shared/services/socket.service';



@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    MatAutocompleteModule,
    ReactiveFormsModule,
    MaterialModule
  ],
  declarations: [ChatComponent, DialogUserComponent, DialogPrivateComponent],
  providers: [
    {provide: ISocketService, useClass: CentrifugeService},
    {provide: IStoreUserService, useClass: ServerStoreUserService},
  ],
  entryComponents: [DialogUserComponent, DialogPrivateComponent]
})
export class ChatModule { }
