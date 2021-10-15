package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/irisnet/irismod/modules/nft/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                types.ModuleName,
		Short:              "Querying commands for the NFT module",
		DisableFlagParsing: true,
	}

	queryCmd.AddCommand(
		GetCmdQueryClass(),
		GetCmdQueryClasses(),
		GetCmdQueryCollection(),
		GetCmdQuerySupply(),
		GetCmdQueryOwner(),
		GetCmdQueryNFT(),
	)

	return queryCmd
}

// GetCmdQuerySupply queries the supply of a nft collection
func GetCmdQuerySupply() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "supply [class-id]",
		Long:    "total supply of a collection or owner of NFTs.",
		Example: fmt.Sprintf("$ %s query nft supply <class-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var owner sdk.AccAddress
			ownerStr, err := cmd.Flags().GetString(FlagOwner)
			if err != nil {
				return err
			}

			if len(ownerStr) > 0 {
				owner, err = sdk.AccAddressFromBech32(ownerStr)
				if err != nil {
					return err
				}
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Supply(context.Background(), &types.QuerySupplyRequest{
				ClassId: args[0],
				Owner:   owner.String(),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	cmd.Flags().AddFlagSet(FsQuerySupply)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryOwner queries all the NFTs owned by an account
func GetCmdQueryOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "owner [address]",
		Long:    "Get the NFTs owned by an account address.",
		Example: fmt.Sprintf("$ %s query nft owner <address> --class-id=<class-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if _, err := sdk.AccAddressFromBech32(args[0]); err != nil {
				return err
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			classID, err := cmd.Flags().GetString(FlagClassID)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Owner(context.Background(), &types.QueryOwnerRequest{
				ClassId:    classID,
				Owner:      args[0],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	cmd.Flags().AddFlagSet(FsQueryOwner)
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "nfts")

	return cmd
}

// GetCmdQueryCollection queries all the NFTs from a collection
func GetCmdQueryCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "collection [class-id]",
		Long:    "Get all the NFTs from a given collection.",
		Example: fmt.Sprintf("$ %s query nft collection <class-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Collection(
				context.Background(),
				&types.QueryCollectionRequest{
					ClassId:    args[0],
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "nfts")

	return cmd
}

// GetCmdQueryClasses queries all classes
func GetCmdQueryClasses() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "classes",
		Long:    "Query all classinations of all collections of NFTs.",
		Example: fmt.Sprintf("$ %s query nft classes", version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Classes(context.Background(), &types.QueryClassesRequest{Pagination: pageReq})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all classes")
	return cmd
}

// GetCmdQueryClass queries the specified class
func GetCmdQueryClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "class [class-id]",
		Long:    "Query the class by the specified class id.",
		Example: fmt.Sprintf("$ %s query nft class <class-id>", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Class(
				context.Background(),
				&types.QueryClassRequest{ClassId: args[0]},
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp.Class)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryNFT queries a single NFTs from a collection
func GetCmdQueryNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token [class-id] [nft-id]",
		Long:    "Query a single NFT from a collection.",
		Example: fmt.Sprintf("$ %s query nft token <class-id> <nft-id>", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if err := types.ValidateTokenID(args[1]); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.NFT(context.Background(), &types.QueryNFTRequest{
				ClassId: args[0],
				TokenId: args[1],
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp.NFT)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
