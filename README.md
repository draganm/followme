# Follow me
Extremely simple event bus

## How it works
* All events are kept in one queue stored on the disk.
    * The queue consists of maximally `n` segments (memory mapped files) with the size of `m`.
    * When a new event is received that would require writing past the end of `n`-th segment, the oldest segment gets deleted and a new empty segment will be created where the event will be written to.
* Senders send by sending a `POST` message to the `/api/events` endpoint
    * Body of the request will be stored as event, no processing of the body will be performed
    * Server will respond with `200` message and body in the following format:
    ```json
    {
        "id": "<uuidv6>"
    }
    ```
    where `uuidv6` is the time-sortable unique id of the event.
* Every receiver can GET to the `/api/events` endpoint to receive an `SSE` stream of the events.
    * When no request body is provided with the `GET` request, `SSE` events will contain all the events after the `Last-Event-ID` header will be returned.
    * Returned events will be of the type `event`
    * Optionally receiver request can contain a short JS in the following format:
        ```js
        function handleEvent(evt) {
            ...
        }
        ```
        This script will process every event since teh `Last-Event-ID` and can either return a value that will be returned to the receiver, or `null` which will result in no event returned.

        If the execution of the script throws an error, an `error` event will be emitted and the request will be closed.

## Limitations

### Maxial size of the messages
Since the qeue is stored as segments of the maximal size of `m`, the largest message that can be written into the queue must have `m` minus the message meta data (`uuidv6` of the message and `varint` encoded length of the message).

### Writing throughput
Writing to the queue imposes a total ordering of the events. This mean that writing events in parallel won't increase the throughtput of the service.
Also, after every messages a `msync` syscall will be peformed to make the change durable on the disk.

### Fault tolerance
At the moment, `followme` is meant to be non-redundant since process service.
This means that as soon as the process terminates, the service will be unavailable.
Ideally this service would be run in Kubernetes as a single repica `StatefulSet`, which would restart pods on crashes.



