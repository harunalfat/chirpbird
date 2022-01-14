import { Component, Inject, OnInit } from '@angular/core';
import { FormControl, Validators } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { IStoreUserService } from '../shared/services/i-store-user.service';

@Component({
  selector: 'app-dialog-user',
  templateUrl: './dialog-user.component.html',
})
export class DialogUserComponent implements OnInit {
  usernameFormControl = new FormControl('', [Validators.required]);
  previousUsername: string;

  constructor(public dialogRef: MatDialogRef<DialogUserComponent>,
    @Inject(MAT_DIALOG_DATA) public params: any,
    private storedUser: IStoreUserService) {
      this.previousUsername = this.storedUser.getStoredUser() ? this.storedUser.getStoredUser() : (params.username ? params.username : undefined);
      this.usernameFormControl.setValue(storedUser.getStoredUser() ? storedUser.getStoredUser().username : (params.username ? params.username : ""));
  }

  ngOnInit() {
  }

  public onSave(): void {
    if (this.usernameFormControl.valid) {
      this.dialogRef.close({
        username: this.usernameFormControl.value,
        dialogType: this.params.dialogType,
        previousUsername: this.previousUsername
      });
    }
  }

}
