# Events

The nft module emits the following events:

## Handlers

### MsgIssueClass

| Type        | Attribute Key | Attribute Value  |
| :---------- | :------------ | :--------------- |
| issue_class | class_id      | {nftClassID}     |
| issue_class | class_name    | {nftClassName}   |
| issue_class | creator       | {creatorAddress} |
| message     | module        | nft              |
| message     | sender        | {senderAddress}  |

### MsgTransferNFT

| Type         | Attribute Key | Attribute Value    |
| :----------- | :------------ | :----------------- |
| transfer_nft | token_id      | {tokenID}          |
| transfer_nft | class_id      | {nftClassID}       |
| transfer_nft | sender        | {senderAddress}    |
| transfer_nft | recipient     | {recipientAddress} |
| message      | module        | nft                |
| message      | sender        | {senderAddress}    |

### MsgEditNFT

| Type     | Attribute Key | Attribute Value |
| :------- | :------------ | :-------------- |
| edit_nft | token_id      | {tokenID}       |
| edit_nft | class_id      | {nftClassID}    |
| edit_nft | token_uri     | {tokenURI}      |
| edit_nft | owner         | {ownerAddress}  |
| message  | module        | nft             |
| message  | sender        | {senderAddress} |

### MsgMintNFT

| Type     | Attribute Key | Attribute Value    |
| :------- | :------------ | :----------------- |
| mint_nft | token_id      | {tokenID}          |
| mint_nft | class_id      | {nftClassID}       |
| mint_nft | token_uri     | {tokenURI}         |
| mint_nft | recipient     | {recipientAddress} |
| message  | module        | nft                |
| message  | sender        | {senderAddress}    |

### MsgBurnNFTs

| Type     | Attribute Key | Attribute Value |
| :------- | :------------ | :-------------- |
| burn_nft | class_id      | {nftClassID}    |
| burn_nft | token_id      | {tokenID}       |
| burn_nft | owner         | {ownerAddress}  |
| message  | module        | nft             |
| message  | sender        | {senderAddress} |

### MsgTransferClass

| Type         | Attribute Key | Attribute Value    |
| :----------- | :------------ | :----------------- |
| transfer_class | class_id      | {nftClassID}       |
| transfer_class | sender        | {senderAddress}    |
| transfer_class | recipient     | {recipientAddress} |
| message      | module        | nft                |
| message      | sender        | {senderAddress}    |