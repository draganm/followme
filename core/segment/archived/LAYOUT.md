# Memory-Mapped File (mmaped) Layout

The memory-mapped file is read-only and consists of two main sections:

## Index Section

The index section contains a lexicographically ordered list of UUIDv6 entries. Each entry is a unique 128-bit identifier and is accompanied by a 64-bit offset that indicates the position within the file where the associated message type and data start. The index section ends with an empty entry that contains only zeros.

## Data Section

The data section contains the actual content of the messages, divided into the message type and message data. The offset from the index section points to the corresponding entry in this section.

The run-length encoding (RLE) method is used for both the message type and data, which begins with a length byte followed by the actual data.

For the message type, the maximum length of the run is 255 bytes, while for the message data, the maximum run length is 2^32-1 bytes.

The separation between the index and data sections allows for efficient data storage and retrieval.

## File Reading Process

To read a message from the file, locate the UUIDv6 in the index section and use the corresponding offset to find the message type and data in the data section.

## Conclusion

The mmaped file format provides an efficient method for data management. The use of UUIDv6 identifiers, offsets, and run-length encoding enables quick data location, access, and space efficiency.

```
|-------------------------------------------|
|              INDEX SECTION                |
|-------------------------------------------|
| UUIDv6 (128 bits)   | Offset (64 bits)    |
|---------------------|---------------------|
| 1b4e28ba-2fa1-11... | 000000000000000A56  |
| 6ba7b810-9dad-11... | 00000000000008090   |
| 00000000-0000-00... | 000000000000000000  |  <- Termination entry
|-------------------------------------------|

|---------------------------------------------------------------------------|
|                          DATA SECTION                                     |
|---------------------------------------------------------------------------|
| Offset (from Index Section) | Msg Type RLE  | Msg Data RLE                |
|-----------------------------|---------------|-----------------------------|
| 000000000000000A56          | 05  | ABCDE   | 000000000A   | Hello World  |
| 00000000000008090           | 04  | WXYZ    | 000000000B   | Goodbye      |
|------------------------------------------------------------|--------------|

```