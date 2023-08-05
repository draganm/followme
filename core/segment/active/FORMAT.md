# Message Format

Each message stored an active segments is composed of the following fields:

1. **UUID** The first 16 bytes represent a UUID. This UUID  must be monotonically increasing for every consequent message (uuidv6). UUID of all zeroes signifies beginning of unused space in the segment.
2. **Type Length** The following 1 byte represents an 8-bit unsigned integer `typLen` that specifies the length of the Type field.
3. **Type** The next `typLen` bytes represent the message type. This is typically used to determine how the payload of the message should be interpreted.
4. **Message Length** The next 4 bytes represent a 32-bit big-endian unsigned integer `msgLen` that specifies the length of the message body.
5. **Message** The next `msgLen` bytes represent the message data. This is the payload of the message.
6. **Checksum** The next 4 bytes represent a 64-bit Highway hash of all previous field. Checksum is used to check the integrity

## Message Diagram

A single message in the segment would have the following structure:

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                             UUID                              |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                             UUID                              |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  Type Length  |                                               |
+-+-+-+-+-+-+-+-+                                               |
|                             Type                              |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                           Message Length                      |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                             Data                              |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                            Checksum                           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                            Checksum                           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```
