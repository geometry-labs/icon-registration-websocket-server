# Icon Registration Websocket Server

This service is used in the [icon-etl]() project
This service will listen for transaction and contract event registrations and respond with a stream of filtered data, all through websockets

## Example filter requests
Send this JSON after opening a websocket connection
```
{
        "connection_type": "ws",
        "endpoint": "wss://test",
        "transaction_events": [
            {
                "to_address": "cx274424e7f501ff60b13fa624ce70e03e7a2b0c12",
                "from_address": "cx274424e7f501ff60b13fa624ce70e03e7a2b0c69"
            }
        ],
        "log_events": [
            {
                "address": "cx1b97c1abfd001d5cd0b5a3f93f22cccfea77e34e",
                "keyword": "BetPlaced",
                "position":  0
            },
            {
                "address": "cx1b97c1abfd001d5cd0b5a3f93f22cccfea77e34e",
                "keyword": "BetPlaced",
                "position":  1
            },
            {
                "address": "cx1b97c1abfd001d5cd0b5a3f93f22cccfea77e34e",
                "keyword": "BetResult",
                "position":  0
            }
        ]
}
```

## Local build
```
go build -o main .
```

## Docker build
```
docker build . -t icon-registration-websocket-server:latest
docker run \
  -p "3000:3000" \
  -e ICON_REGISTRATION_WEBSOCKET_OUTPUT_TOPIC="outputs"
  -e ICON_REGISTRATION_WEBSOCKET_BROKER_URL="kafka:9092"
  -e ICON_REGISTRATION_WEBSOCKET_PORT="3000"
  kafka-websocket-server:latest
```

## Enviroment Variables

| Name | Description | Default | Required |
|------|-------------|---------|----------|
| ICON_REGISTRATION_WEBSOCKET_BROKER_URL | location of broker | NULL | True |
| ICON_REGISTRATION_WEBSOCKET_OUTPUT_TOPIC | name of topic where all filtered data will be | "output" | False |
| ICON_REGISTRATION_WEBSOCKET_PORT | port to expose for websocket connections | "3000" | False |
