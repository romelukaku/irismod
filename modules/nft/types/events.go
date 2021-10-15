package types

// NFT module event types
var (
	EventTypeIssueClass    = "issue_class"
	EventTypeTransfer      = "transfer_nft"
	EventTypeEditNFT       = "edit_nft"
	EventTypeMintNFT       = "mint_nft"
	EventTypeBurnNFT       = "burn_nft"
	EventTypeTransferClass = "transfer_class"

	AttributeValueCategory = ModuleName

	AttributeKeySender    = "sender"
	AttributeKeyCreator   = "creator"
	AttributeKeyRecipient = "recipient"
	AttributeKeyOwner     = "owner"
	AttributeKeyTokenID   = "token_id"
	AttributeKeyTokenURI  = "token_uri"
	AttributeKeyClassID   = "class_id"
	AttributeKeyClassName = "class_name"
)
