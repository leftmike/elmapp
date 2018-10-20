module View exposing (..)

import Api exposing (User)
import LoginTypes as Login

type alias Session =
    { name : String
    }

type View = MainView Session | CounterView Session | LoginView Login.State

login : User -> View
login user =
    MainView {name = user.username}
