import {GameName} from "./game-name";
import {Platform} from "./platform";
import {User} from "./user";
import {Score} from "./score";
export interface Game {
  id?: number
  names: GameName[]
  platform?: Platform
  updated_by?: User
  scores?: Score[]
}
