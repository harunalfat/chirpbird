import { Component, Inject, OnInit } from '@angular/core';
import { FormControl, Validators } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ISocketService } from '../shared/services/i-socket-service';
import { debounceTime } from 'rxjs/operators'
import { User } from '../shared/model/user';

@Component({
  selector: 'app-dialog-private',
  templateUrl: './dialog-private.component.html',
})
export class DialogPrivateComponent implements OnInit {
  userFormControl = new FormControl('', [Validators.required]);
  filteredOptions: User[]

  constructor(public dialogRef: MatDialogRef<DialogPrivateComponent>,
    private socketService: ISocketService,
    @Inject(MAT_DIALOG_DATA) public params: any) {
    this.userFormControl.setValue(null);
  }

  ngOnInit() {
    this.userFormControl.valueChanges.pipe(debounceTime(500)).subscribe(async val => {
      this.filteredOptions = await this._compare(val)
    })
  }

  private async _compare(filterValue: string): Promise<User[]> {
    if (!filterValue) return;
    const loweredFilterValue = filterValue.toLowerCase();
    const users = await this.socketService.searchUsersByUsernameRpc(loweredFilterValue) || []
    return users
  }

  public onSave(): void {
    if (this.params.username === this.userFormControl.value) return

    const user = this.filteredOptions.find(u => u.username === this.userFormControl.value)
    if (!user) return

    this.dialogRef.close({
      id: user.id,
      username: user.username,
    });
  }

}
