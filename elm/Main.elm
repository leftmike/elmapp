module Main exposing (main)

{-
To Do:
-}

import Browser exposing (Document)
import Counter
import Html exposing (..)
import Html.Events as Events
import Login
import Message exposing (Msg(..))
import Model exposing (..)
import View

main =
    Browser.document {init = init, update = update, view = view, subscriptions = subscriptions}

init : () -> (Model, Cmd Msg)
init _ = (
    { activeView = View.LoginView Login.init
    , counter = Counter.init
    }, Cmd.none)

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
        ChangeView newView -> ({model | activeView = newView}, Cmd.none)
        UpdateCounter _ -> Counter.update msg model
        UpdateLogin loginMsg -> Login.update loginMsg model

subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none

view : Model -> Document Msg
view model =
    case model.activeView of
        View.MainView session -> viewWithSession model session viewMain
        View.CounterView session -> viewWithSession model session Counter.view
        View.LoginView state -> Login.view state

viewWithSession : Model -> View.Session -> (Model -> Html Msg) -> Document Msg
viewWithSession model session v =
    { title = "Elm App"
    , body =
        [ div []
             [ viewHeader model session
             , v model
             ]
         ]
    }

viewHeader : Model -> View.Session -> Html Msg
viewHeader model session =
    div []
        [ button [ Events.onClick (ChangeView (View.MainView session)) ] [ text "Main View" ]
        , button [ Events.onClick (ChangeView (View.CounterView session)) ] [ text "Counter View" ]
        , button [ Events.onClick (ChangeView (View.LoginView Login.init))]
            [ text ("Logout " ++ session.name) ]
        ]

viewMain : Model -> Html Msg
viewMain model =
    text "Main View"
