import {Injectable} from "@angular/core";
import {Game} from "./game";
import {Observable, BehaviorSubject} from "rxjs";
import {Http, Headers, RequestOptions, Response} from "@angular/http";
import {Platform} from "./platform";
import {MdSnackBar} from "@angular/material";
import {User} from "./user";
import {Router} from "@angular/router";

@Injectable()
export class BackendService {

  public gameQuery: BehaviorSubject<string> = new BehaviorSubject('');
  private games: BehaviorSubject<Game[]> = new BehaviorSubject([]);
  private platforms: BehaviorSubject<Platform[]> = new BehaviorSubject<Platform[]>([]);
  private user: BehaviorSubject<User> = new BehaviorSubject(null);

  private options = new RequestOptions({
    headers: new Headers({'Content-Type': 'application/json'})
  });

  constructor(private http: Http, private snackbar: MdSnackBar, private router: Router) {
    this.gameQuery.subscribe(query => this.filterBy(query));

    this.http.get('/api/platforms')
      .map(res => res.json())
      .subscribe(platforms => this.platforms.next(platforms));

    this.http.get('/api/profile')
      .map(res => res.json())
      .subscribe(user => this.user.next(user));

  }

  public filterBy(query: string): void {
    this.http.get('/api/games?name=' + query)
      .map(res => res.json())
      .subscribe(data => this.games.next(data));
  }

  public getGames(): Observable<Game[]> {
    return this.games.asObservable().share();
  }

  public getGame(id: number): Observable<Game> {
    return this.http.get('/api/games/' + id).map(res => res.json())
  }

  public saveGame(game: Game): Observable<Game> {
    let stripped: Game = {
      id: game.id,
      names: game.names,
      platform: game.platform
    };

    if (game.id) {
      return this.http.patch('/api/games/' + game.id, JSON.stringify(stripped), this.options)
        .map(res => res.json())
        .map(json => {
          this.gameQuery.next(this.gameQuery.value);
          return json;
        });
    } else {
      return this.http.post('/api/games', JSON.stringify(stripped), this.options)
        .map(res => res.json())
        .map(json => {
          this.gameQuery.next(this.gameQuery.value);
          return json;
        });
    }
  }

  public saveScore(game: Game, score: number) {
    return this.http.post('/api/scores', {
      game_id: game.id,
      score
    }).map(res => res.json())
      .map(json => {
        this.gameQuery.next(this.gameQuery.value);
        this.snackbar.open('Score saved! Your best rank is now ' + json.rank + ' for this game!');
        return json;
      });
  }

  public getPlatforms(): Observable<Platform[]> {
    return this.platforms.asObservable().share();
  }

  public getUser(): Observable<User> {
    return this.user.asObservable().share();
  }

  public isLoggedIn(): Observable<boolean> {
    return this.getUser().share().map(user => !!user)
  }

  public logout(): void {
    this.http.post('/logout', '').subscribe(res => window.location.href = '/');
  }

  public saveSettings(user: User): void {
    this.http.patch('/api/settings', user, this.options).subscribe(res => {
      this.gameQuery.next(this.gameQuery.value);
      this.snackbar.open('Settings saved!');
      this.router.navigate(['/']);
    });
  }

  public localLogin(name: string, password: string): Observable<Response> {
    return this.http.post('/local/login', {name, password}, this.options);
  }

  public signup(registration: any): void {
    this.http.post('/api/register', registration, this.options)
      .subscribe(res => {
        this.snackbar.open('Registration successful. Please login...');
      }, (err:Response) => {
        console.log(err);
        this.snackbar.open('Registration FAILED! ' + err.text());
      })
  }

}
