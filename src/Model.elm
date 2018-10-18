module Model exposing (..)

import CounterTypes as Counter
import View

type alias Model =
    { activeView : View.View
    , counter : Counter.Model
    }
