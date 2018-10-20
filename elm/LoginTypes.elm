module LoginTypes exposing (..)

import Api
import Http

type Msg
    = EmailInput String
    | PasswordInput String
    | LoginClick
    | LoginResponse (Result Http.Error Api.User)

type alias State = {email : String, password : String, errors : List String, requested : Bool}
