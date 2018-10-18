module Message exposing (Msg(..))

import CounterTypes as Counter
import LoginTypes as Login
import View exposing (View)

type Msg =
      ChangeView View
    | UpdateCounter Counter.Msg
    | UpdateLogin Login.Msg
