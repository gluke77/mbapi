# [MessageBird](https://www.messagebird.com) SMS API test

Usage
```
$MBAPI_KEY=<your API key> MBAPI_PORT=8088 MBAPI_QUEUE_SIZE=5 ./mbapi
```

`MBAPI_QUEUE_SIZE=0` -- use syncronous (unbuffered) sender

`MBAPI_QUEUE_SIZE>0` -- use asyncronous sender. All requests that do not fit into the buffer will be discarded

In another shell
```
$curl -X POST --data '{"originator": "yournamehere", "recipient":"31612345678", "message":"hello"}' http://localhost:8088
```
