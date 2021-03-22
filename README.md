# Tiny Streaming Data (TSD) Format

TSD is a simple and extensible file/data format optimized for applications that involve binary data streaming where every byte matters. It is inspired by interchange formats such as RIFF and IFF.

TSD solves the problem of packing different types of binary data in a single file or stream without delimeters. This is commonly encountered in packing [protobuf data](https://developers.google.com/protocol-buffers/docs/techniques#streaming).

## Format Specification

In TSD, much like RIFF, binary data blobs are packed in 'chunks', each with an identifier and a length prefix. In TSD, headers and padding are removed and protobuf varints are used instead of fixed-length fourCC's for efficiency.

A TSD file is an indefinitely repeating structure:

`[uvarint chunk ID][uvarint chunk length][bytes of binary data]...[ID][length][data]...`

Readers of TSD files use the chunk ID and length to iterate data and handle it specially. Consider a file type used for streaming a user's social media activity. We can link different chunk IDs to different types of activity:

1. Text Post (ASCII text)
2. Like (ID of liked post)
3. Photo (JPG data)

The file will look like: 

`[1][12]Hello world![3][42178]d^T<DF>=s<C5>DН<DD>E>;sԚ...[2][6]776541`

Because we use [varints](https://developers.google.com/protocol-buffers/docs/encoding#varints), packing data like `Hello world!` (with a chunk ID of 1 and a length of 12) only takes two additional bytes for a header. Chunk IDs can be any 64-bit unsigned integer, but chunks that occur more frequently should have lower ID values to take up less space.

A client reading this data will be able to use a switch statement on the chunk ID (1, 2, or 3) and read back each data blob knowing what to do with it. If, in the future, we want to add stories to our social media platform, we can introduce chunkID #4, and clients that don't yet support stories will be able to skip over that chunk.

Chunk ID 0 is reserved as a continuation signal - this facilitates writing data without knowing its size ahead of time by being able to split it into pieces. TSD like `[1][6]Hello [0][6]world!` can simply be read together as `Hello world!` with ID 1.

## What's in the Box

This repo contains a Go library to facilitate reading, writing, and inspecting Tiny Streaming Data. In the future I may add support for more languages, but it shouldn't be difficult to build your own reader/writer in any language that is also supported by protobuf.
