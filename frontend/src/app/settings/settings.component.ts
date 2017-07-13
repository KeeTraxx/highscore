import { Component, OnInit } from '@angular/core';
import {User} from "../user";
import {BackendService} from "../backend.service";

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.css']
})
export class SettingsComponent implements OnInit {

  private user:User;

  constructor(private backend:BackendService) { }

  ngOnInit() {
    this.backend.getUser().subscribe(user => this.user = user);
  }

  saveSettings(user:User) {
    this.backend.saveSettings(user);
  }

}
