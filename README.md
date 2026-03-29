# chainer

A UTXO-based blockchain implementation written in Go.

## What's implemented

**Crypto** — Ed25519 keys and signatures, addresses derived as the last 20 bytes of the public key.

**Blocks** — header (version, height, previous block hash, merkle root, timestamp) plus a list of transactions. Blocks are hashed via SHA256 of the serialized header.

**Transactions** — standard UTXO model: inputs reference outputs of previous transactions (by hash and index), outputs specify recipient and amount. Each input is signed by the owner's private key. Verification checks all input signatures.

**Serialization** — protobuf.

## Structure

```
proto/     — .proto schemas and generated code
core/      — blocks and transactions
crypto/    — keys, signatures, addresses
utils/     — test helpers
```

## Dependencies

- Go 1.25+
- [Task](https://taskfile.dev) (optional, for Taskfile)
- protoc + protoc-gen-go
