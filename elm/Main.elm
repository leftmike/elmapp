module Main exposing (main)

{-
To Do:
- remove Counter
- add Register
- handle urls and routes
-}

import Browser exposing (Document)
import Counter
import Html exposing (..)
import Html.Attributes exposing (..)
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
        View.LoginView state ->
            { title = Login.title state ++ " - Elm App"
            , body =
                [ viewHeader model
                , Login.view state
                , viewFooter model
                ]
            }

viewHeader : Model -> Html Msg
viewHeader model =
    let
        (home, login, register) =
            case model.activeView of
                View.MainView _ -> (" active", "", "")
                View.CounterView _ -> ("", "", "")
                View.LoginView _ -> ("", " active", "")
    in
    nav [ class "navbar navbar-light" ]
        [ div [ class "container" ]
            [ a [ class "navbar-brand", href "/#/" ] [ text "conduit" ]
            , ul [ class "nav navbar-nav pull-xs-right" ]
                [ li [ class "nav-item" ]
                    [ a [ class ("nav-link" ++ home), href "/#/"] [ text "Home" ] ]
                , li [ class "nav-item" ]
                    [ a [ class ("nav-link" ++ login), href "/#/login"] [ text "Sign in" ] ]
                , li [ class "nav-item" ]
                    [ a [ class ("nav-link" ++ register), href "/#/register"] [ text "Sign up" ] ]
                ]
            ]
        ]

viewFooter : Model -> Html Msg
viewFooter model =
    footer []
        [ div [ class "container" ]
            [ a [ class "logo-font"
                , href "https://github.com/leftmike/elmapp" ]
                [ text "Fork on GitHub" ]
            ]
        ]

viewWithSession : Model -> View.Session -> (Model -> Html Msg) -> Document Msg
viewWithSession model session v =
    { title = "Elm App"
    , body =
        [ div []
             [ viewHeaderSession model session
             , v model
             ]
         ]
    }

-- XXX: change to viewSessionHeader
viewHeaderSession : Model -> View.Session -> Html Msg
viewHeaderSession model session =
    div []
        [ button [ Events.onClick (ChangeView (View.MainView session)) ] [ text "Main View" ]
        , button [ Events.onClick (ChangeView (View.CounterView session)) ] [ text "Counter View" ]
        , button [ Events.onClick (ChangeView (View.LoginView Login.init))]
            [ text ("Logout " ++ session.name) ]
        ]

viewMain : Model -> Html Msg
viewMain model =
    text "Main View"
