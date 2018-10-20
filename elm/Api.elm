module Api exposing (User, userDecoder, loginPost, processError)

import Http
import Json.Decode as Decode exposing (Decoder)
import Json.Encode as Encode

type alias User = {email : String, username : String, bio : String, image : String, token : String}

userDecoder : Decoder User
userDecoder =
    Decode.map5 User
        (Decode.field "email" Decode.string)
        (Decode.field "username" Decode.string)
        (Decode.field "bio" Decode.string)
        (Decode.field "image" Decode.string)
        (Decode.field "token" Decode.string)

loginPost : String -> String -> Http.Request User
loginPost email password =
    let
        body =
            Encode.object
                [ ("user",
                    Encode.object
                        [ ("email", Encode.string email)
                        , ("password", Encode.string password)
                        ]
                )]
    in
        Http.post "/api/users/login" (Http.jsonBody body) (Decode.field "user" userDecoder)

processError : Http.Error -> List String
processError error =
    case error of
        Http.BadUrl url -> ["bad url: " ++ url]
        Http.Timeout -> ["request timed out"]
        Http.NetworkError -> ["network error"]
        Http.BadStatus response ->
            if response.status.code == 422 then
                processBadStatus response
            else
                [response.status.message, response.body]
        Http.BadPayload err response -> [err, response.status.message, response.body]

decodeErrors : Decoder (List String)
decodeErrors =
    Decode.field "errors" (Decode.field "body" (Decode.list Decode.string))

processBadStatus : Http.Response String -> List String
processBadStatus response =
    case Decode.decodeString decodeErrors response.body of
        Ok list -> list
        Err error -> [response.status.message, response.body, Decode.errorToString error]
