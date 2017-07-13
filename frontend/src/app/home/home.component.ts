import {Component, OnInit} from "@angular/core";
import {BackendService} from "../backend.service";
import {FormControl} from "@angular/forms";

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {

  formControl:FormControl = new FormControl();

  constructor(public backend:BackendService) { }

  ngOnInit() {
    this.formControl.valueChanges.debounceTime(1000).subscribe(query => {
      this.backend.gameQuery.next(query);
    })
  }

}
