import {Component, ViewChild} from '@angular/core';
import {BackendService} from "./backend.service";
import {MdSnackBar, MdSidenav, MdDialog} from "@angular/material";
import {Router} from "@angular/router";
import {SignupComponent} from "./signup/signup.component";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  @ViewChild('sidenav') public sidenav: MdSidenav;
  name:string;
  password:string;
  title = 'app';

  constructor(public backend:BackendService, private snackbar:MdSnackBar, private router:Router, private dialog:MdDialog){
    router.events.subscribe(next => {
      console.log(next);
      this.sidenav.close();
    });
  }

  localLogin(name:string, password:string) {
    this.backend.localLogin(name, password).subscribe(res => {
      window.location.href = '/';
    }, err => {
       this.snackbar.open("Login failed.");
       name = "";
       password = "";
    });
  }

  openSignup() {
    this.dialog.open(SignupComponent);
  }

}
