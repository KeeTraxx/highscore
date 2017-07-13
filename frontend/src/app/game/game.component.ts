import {Component, OnInit} from "@angular/core";
import {ActivatedRoute, Router} from "@angular/router";
import {Game} from "../game";
import {BackendService} from "../backend.service";

@Component({
  selector: 'app-game',
  templateUrl: './game.component.html',
  styleUrls: ['./game.component.css']
})
export class GameComponent implements OnInit {
  game: Game = {names: [{name: ''}]};
  constructor(private activatedRoute: ActivatedRoute, private backend: BackendService, private router:Router) {
  }

  ngOnInit() {
    this.activatedRoute.params.subscribe(params => {
      if (params['id'] !== 'new') {
        this.backend.getGame(+params['id']).subscribe(game => this.game = game)
      } else {
        this.game = {names: [{name: ''}]};
      }
    })
  }

  saveGame(game: Game) {
    this.backend.saveGame(game)
      .subscribe(game => {
        console.log(game);
        this.router.navigate(['/']);
      })
  }

  compareId(a: any, b: any) {
    return a && b && a.id == b.id;
  }


}
