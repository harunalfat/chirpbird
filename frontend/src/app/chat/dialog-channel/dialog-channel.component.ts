import { Component, Inject, OnInit } from '@angular/core';
import { FormControl, Validators } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

@Component({
  selector: 'app-dialog-channel',
  templateUrl: './dialog-channel.component.html',
  styleUrls: ['./dialog-channel.component.css']
})
export class DialogChannelComponent implements OnInit {
  channelFormControl = new FormControl('', [Validators.required]);

  constructor(public dialogRef: MatDialogRef<DialogChannelComponent>,
    @Inject(MAT_DIALOG_DATA) public params: any) {
    this.channelFormControl.setValue(null);
  }

  ngOnInit() {
  }

  public onSave(): void {
    this.dialogRef.close({
      channelName: this.channelFormControl.value,
    });
  }

}
