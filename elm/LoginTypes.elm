module LoginTypes exposing (..)

type Msg = NameInput String | PasswordInput String | LoginClick

type alias State = {name : String, password : String, bad : Bool}
