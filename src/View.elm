module View exposing (..)

import LoginTypes as Login

type alias Session =
    { name : String
    }

type View = MainView Session | CounterView Session | LoginView Login.State

login : String -> View
login name =
    MainView {name = name}
