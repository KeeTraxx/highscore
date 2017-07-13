import {BrowserModule} from "@angular/platform-browser";
import {RouterModule} from "@angular/router";
import {NgModule} from "@angular/core";
import {AppComponent} from "./app.component";
import {GameComponent} from "./game/game.component";
import {HomeComponent} from "./home/home.component";
import {HttpModule} from "@angular/http";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import {
  MdButtonModule, MdInputModule, MdSelectModule, MdListModule, MdCardModule,
  MdSidenavModule, MdIconModule, MdToolbarModule, MdAutocompleteModule, MdTableModule, MdSnackBarModule
} from "@angular/material";
import {BackendService} from "./backend.service";
import { SettingsComponent } from './settings/settings.component';

const ROUTES = [
  {
    path: 'game/:id',
    component: GameComponent
  },
  {
    path: 'settings',
    component: SettingsComponent
  },
  {
    path: '',
    component: HomeComponent
  }
]

@NgModule({
  declarations: [
    AppComponent,
    GameComponent,
    HomeComponent,
    SettingsComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    FormsModule,
    ReactiveFormsModule,
    HttpModule,
    RouterModule.forRoot(ROUTES, {useHash: true}),
    MdButtonModule,
    MdInputModule,
    MdSelectModule,
    MdListModule,
    MdCardModule,
    MdSidenavModule,
    MdIconModule,
    MdToolbarModule,
    MdSnackBarModule
  ],
  providers: [
    BackendService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
