# Follow me
Extremely simple event bus

## 

## How it works
* In the first draft of the service the queue will be stored in [bolted](https://github.com/draganm/bolted) map with `<uuidv6>-<type>` key and event payload as value. 
* In the future versions, all events are kept in one queue stored on the disk.
    * The queue consists of maximally `n` segments (memory mapped files) with the size of `m`.
    * When a new event is received that would require writing past the end of `n`-th segment, the oldest segment gets deleted and a new empty segment will be created where the event will be written to.
    * Every segment will be accompanied with one main and per-type index in a separate file. All indices are going to be in the format:
        * `<uuidv6 (16 bytes)>`
        * `<big endian encoded start position of the message in the segment (4 bytes)>`
        * `<big endian encoded end position of the message in the segment (4 bytes)>`
* Senders send by sending a `POST` message to the `/api/events?type=<type>` endpoint
    * Body of the request together with the type will be stored as event. No processing of the body will be performed, which means that the body can be any text type.
    * Server will respond with `200` message and body in the following format:
    ```json
    {
        "id": "<uuidv6>"
    }
    ```
    where `uuidv6` is the time-sortable unique id of the event.
* Every receiver can GET to the `/api/events` endpoint to receive an `SSE` stream of the events.
    * When no request parameter are provided with the `GET` request, `SSE` events will contain all the events after the `Last-Event-ID` header will be returned.
    * Optionally, the client can pass `typeMatcher=<regexp matching the wanted types>` matcher to filter our only the wanted types. 
    * Returned events will have the type and body of the stored event.

## Limitations

### Event's can't be arbitrary binary data
SSE does not support retrning binary data, intead events have to be valid utf-8 sequences.
If binary data is to be sent as an event (e.g. a protobuf message) it has to be sent in a text-encoded format such `base64`.

### Writing throughput
Writing to the queue imposes a total ordering of the events. This mean that writing events in parallel won't increase the throughtput of the service.
Also, after every messages a `msync` syscall will be peformed to make the change durable on the disk.

### Fault tolerance
At the moment, `followme` is meant to be non-redundant since process service.
This means that as soon as the process terminates, the service will be unavailable.
Ideally this service would be run in Kubernetes as a single repica `StatefulSet`, which would restart pods on crashes.



