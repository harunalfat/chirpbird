import { Base } from "./base";
import { Channel } from "./channel";

export interface User extends Base {
    username?: string;
    avatar?: string;
    channels?: Channel[];
}
