package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagTokenName = "name"
	FlagTokenURI  = "uri"
	FlagTokenData = "data"
	FlagRecipient = "recipient"
	FlagOwner     = "owner"

	FlagClassName        = "name"
	FlagClassID          = "class-id"
	FlagSchema           = "schema"
	FlagSymbol           = "symbol"
	FlagMintRestricted   = "mint-restricted"
	FlagUpdateRestricted = "update-restricted"
)

var (
	FsIssueClass    = flag.NewFlagSet("", flag.ContinueOnError)
	FsMintNFT       = flag.NewFlagSet("", flag.ContinueOnError)
	FsEditNFT       = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferNFT   = flag.NewFlagSet("", flag.ContinueOnError)
	FsQuerySupply   = flag.NewFlagSet("", flag.ContinueOnError)
	FsQueryOwner    = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferClass = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsIssueClass.String(FlagSchema, "", "Class data structure definition")
	FsIssueClass.String(FlagClassName, "", "The name of the class")
	FsIssueClass.String(FlagSymbol, "", "The symbol of the class")
	FsIssueClass.Bool(FlagMintRestricted, false, "mint restricted of nft under class")
	FsIssueClass.Bool(FlagUpdateRestricted, false, "update restricted of nft under class")

	FsMintNFT.String(FlagTokenURI, "", "URI for supplemental off-chain tokenData (should return a JSON object)")
	FsMintNFT.String(FlagRecipient, "", "Receiver of the nft, if not filled, the default is the sender of the transaction")
	FsMintNFT.String(FlagTokenData, "", "The origin data of the nft")
	FsMintNFT.String(FlagTokenName, "", "The name of the nft")

	FsEditNFT.String(FlagTokenURI, "[do-not-modify]", "URI for the supplemental off-chain token data (should return a JSON object)")
	FsEditNFT.String(FlagTokenData, "[do-not-modify]", "The token data of the nft")
	FsEditNFT.String(FlagTokenName, "[do-not-modify]", "The name of the nft")

	FsTransferNFT.String(FlagTokenURI, "[do-not-modify]", "URI for the supplemental off-chain token data (should return a JSON object)")
	FsTransferNFT.String(FlagTokenData, "[do-not-modify]", "The token data of the nft")
	FsTransferNFT.String(FlagTokenName, "[do-not-modify]", "The name of the nft")

	FsQuerySupply.String(FlagOwner, "", "The owner of the nft")

	FsQueryOwner.String(FlagClassID, "", "The name of the collection")
}
