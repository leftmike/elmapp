module Counter exposing (..)

import CounterTypes as Types
import Html exposing (..)
import Html.Events as Events
import Message exposing (Msg(..))
import Model exposing (..)
import String

init = {clickCount = 0}

update : Msg -> Model -> (Model, Cmd Msg)
update msg ({counter} as model) =
    case msg of
        UpdateCounter counterMsg ->
            case counterMsg of
                Types.Click ->
                    ({model | counter = {clickCount = counter.clickCount + 1}}, Cmd.none)
        _ -> (model, Cmd.none)

click = UpdateCounter Types.Click

view : Model -> Html Msg
view model =
    div []
        [ button [ Events.onClick click ] [ text "Click" ]
        , text (String.fromInt model.counter.clickCount)
        ]
