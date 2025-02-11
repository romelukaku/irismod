# State

## NFT

Nft defines the tokenData of non-fungible tokens, mainly including ID, owner, and tokenURI.Nft can be transferred through `MsgTransferNFT`, or you can edit `tokenURI` information through `MsgEditNFT` transaction. The name of the collection and the id of nft identify the unique assets in the system. The `NFT` Interface inherits the NFT struct and includes getter functions for the asset data. It also includes a Stringer function in order to print the struct. The interface may change if tokenData is moved to it’s own module as it might no longer be necessary for the flexibility of an interface.

```go
// NFT non fungible token interface
type NFT interface {
    GetID() string              // unique identifier of the NFT
    GetName() string            // return the name of NFT
    GetOwner() sdk.AccAddress   // gets owner account of the NFT
    GetURI() string             // tokenData field: URI to retrieve the of chain tokenData of the NFT
    GetData() string            // return the Data of NFT
}
```

## Collections

As all NFTs belong to a specific `Collection`, however, considering the performance issue, we did not store the structure, but used `{classID}/{tokenID}` as the key to identify each nft ’s own collection, use `{class}` as the key to store the number of nft in the current collection, which is convenient for statistics and query.collection is defined as follows

```go
// Collection of non fungible tokens
type Collection struct {
    Class Class     `json:"class"`  // Class of the collection; not exported to clients
    NFTs  []NFT `json:"nfts"`   // NFTs that belongs to a collection
}
```

## Owners

Owner is a data structure specifically designed for nft owned by statistical model owners. The ownership of an NFT is set initially when an NFT is minted and needs to be updated every time there's a transfer or when an NFT is burned,defined as follows:

```go
// Owner of non fungible tokens
type Owner struct {
    Address       string            `json:"address"`
    IDCollections []IDCollection    `json:"id_collections"`
}
```

An `IDCollection` is similar to a `Collection` except instead of containing NFTs it only contains an array of `NFT` IDs. This saves storage by avoiding redundancy.

```go
// IDCollection of non fungible tokens
type IDCollection struct {
    ClassId string   `json:"class_id"`
    TokenIds []string `json:"token_ids"`
}

```
