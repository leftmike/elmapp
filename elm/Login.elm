module Login exposing (init, update, view)

import Api
import Dict
import LoginTypes as Types
import Html exposing (..)
import Html.Attributes as Attributes
import Html.Events as Events
import Http
import Message exposing (Msg(..))
import Model exposing (..)
import View

init = {email = "", password = "", errors = [], requested = False}

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
            Types.EmailInput email ->
                if state.requested then
                    (model, Cmd.none)
                else
                    ({model | activeView =
                        View.LoginView {state | email = email}}, Cmd.none)

            Types.PasswordInput password ->
                if state.requested then
                    (model, Cmd.none)
                else
                    ({model | activeView =
                        View.LoginView {state | password = password}}, Cmd.none)

            Types.LoginClick ->
                if state.requested then
                    (model, Cmd.none)
                else
                    ({model | activeView =
                        View.LoginView {state | requested = True}}, loginRequested model state)

            Types.LoginResponse (Ok user) ->
                ({model | activeView = View.login user}, Cmd.none)

            Types.LoginResponse (Err error) ->
                ({model | activeView =
                    View.LoginView {state | errors = Api.processError error, requested = False}},
                    Cmd.none)

loginRequested : Model -> Types.State -> Cmd Msg
loginRequested model state =
    Http.send loginResponse (Api.loginPost state.email state.password)

emailInput = Types.EmailInput >> UpdateLogin
passwordInput = Types.PasswordInput >> UpdateLogin
loginClick = UpdateLogin Types.LoginClick
loginResponse = Types.LoginResponse >> UpdateLogin

view : Types.State -> Html Msg
view state =
    div []
        [ div [] [ text "Sign in" ]
        , div [] [ text "Need an account?" ]
        , viewErrors state.errors
        , viewInput "Email" "text" state.email emailInput
        , viewInput "Password" "password" state.password passwordInput
        , button [ Events.onClick loginClick ] [ text "Sign in" ]
        ]

viewInput : String -> String -> String -> (String -> Msg) -> Html Msg
viewInput placeholder type_ value event =
    div []
        [ input
            [ Attributes.type_ type_
            , Attributes.value value
            , Attributes.placeholder placeholder
            , Events.onInput event
            ] []
        ]

viewErrors : List String -> Html Msg
viewErrors errors =
    if List.isEmpty errors then
        text ""
    else
        ul [ Attributes.style "color" "red" ] (List.map viewError errors)

viewError : String -> Html Msg
viewError error =
    li [] [ text error ]
