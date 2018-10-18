module Login exposing (..)

import Dict
import LoginTypes as Types
import Html exposing (..)
import Html.Attributes as Attributes
import Html.Events as Events
import Message exposing (Msg(..))
import Model exposing (..)
import View

init = {name = "", password = "", bad = False}

getState : Model -> Types.State
getState model =
    case model.activeView of
        View.LoginView state -> state
        _ -> init

update: Types.Msg -> Model -> (Model, Cmd Msg)
update msg model =
    let
        state = getState model
    in
        case msg of
            Types.NameInput name ->
                ({model | activeView =
                    View.LoginView {state | name = name, bad = False}}, Cmd.none)
            Types.PasswordInput password ->
                ({model | activeView =
                    View.LoginView {state | password = password, bad = False}}, Cmd.none)
            Types.LoginClick -> updateLogin model state

accounts =
    Dict.fromList [("setup", "default"), ("mike", "password"), ("test", "test"), ("none", "")]

checkPassword : String -> String -> Bool
checkPassword name password =
    case Dict.get name accounts of
        Just pw -> pw == password
        Nothing -> False

updateLogin: Model -> Types.State -> (Model, Cmd Msg)
updateLogin model state =
    if checkPassword state.name state.password then
        ({model | activeView = View.login state.name}, Cmd.none)
    else
        ({model | activeView = View.LoginView {name = "", password = "", bad = True}}, Cmd.none)

nameInput = Types.NameInput >> UpdateLogin
passwordInput = Types.PasswordInput >> UpdateLogin
loginClick = UpdateLogin Types.LoginClick

view : Types.State -> Html Msg
view state =
    div []
        [ viewInput "Name" "text" state.name nameInput
        , viewInput "Password" "password" state.password passwordInput
        , button [ Events.onClick loginClick ] [ text "Login" ]
        , div [ Attributes.style "color" "red" ]
            [ text (if state.bad then "Bad name or password." else "") ]
        ]

viewInput : String -> String -> String -> (String -> Msg) -> Html Msg
viewInput label type_ value event =
    div []
        [ text label
        , input
            [ Attributes.type_ type_
            , Attributes.value value
            , Events.onInput event
            ] []
        ]
