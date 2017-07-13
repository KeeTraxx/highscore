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
  MdSidenavModule, MdIconModule, MdToolbarModule, MdAutocompleteModule, MdTableModule, MdSnackBarModule, MdDialogModule
} from "@angular/material";
import {BackendService} from "./backend.service";
import { SettingsComponent } from './settings/settings.component';
import { SignupComponent } from './signup/signup.component';

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
    SettingsComponent,
    SignupComponent
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
    MdSnackBarModule,
    MdDialogModule
  ],
  providers: [
    BackendService
  ],
  bootstrap: [AppComponent],
  entryComponents: [SignupComponent]
})
export class AppModule { }
