import { Base } from "./base";
import { User } from "./user";

export interface Channel extends Base {
    name: string
    creatorId: string
    isPrivate: boolean
    hashIdentifier: string
    isNewlyAdded?: boolean
    participants: {
        username: string
        id: string
    }[]
}