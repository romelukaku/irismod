package testutil

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	nftcli "github.com/irisnet/irismod/modules/nft/client/cli"
)

// MsgRedelegateExec creates a redelegate message.
func IssueClassExec(clientCtx client.Context, from string, class string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		class,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, nftcli.GetCmdIssueClass(), args)
}

func BurnNFTExec(clientCtx client.Context, from string, classID string, tokenID string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		classID,
		tokenID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, nftcli.GetCmdBurnNFT(), args)
}

func MintNFTExec(clientCtx client.Context, from string, classID string, tokenID string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		classID,
		tokenID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, nftcli.GetCmdMintNFT(), args)
}

func EditNFTExec(clientCtx client.Context, from string, classID string, tokenID string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		classID,
		tokenID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, nftcli.GetCmdEditNFT(), args)
}

func TransferNFTExec(clientCtx client.Context, from string, recipient string, classID string, tokenID string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		recipient,
		classID,
		tokenID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, nftcli.GetCmdTransferNFT(), args)
}

func QueryClassExec(clientCtx client.Context, classID string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		classID,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, nftcli.GetCmdQueryClass(), args)
}

func QueryCollectionExec(clientCtx client.Context, classID string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		classID,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, nftcli.GetCmdQueryCollection(), args)
}

func QueryClassesExec(clientCtx client.Context, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, nftcli.GetCmdQueryClasses(), args)
}

func QuerySupplyExec(clientCtx client.Context, class string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		class,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, nftcli.GetCmdQuerySupply(), args)
}

func QueryOwnerExec(clientCtx client.Context, address string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		address,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, nftcli.GetCmdQueryOwner(), args)
}

func QueryNFTExec(clientCtx client.Context, classID string, tokenID string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		classID,
		tokenID,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, nftcli.GetCmdQueryNFT(), args)
}

func TransferClassExec(clientCtx client.Context, from string, recipient string, classID string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		recipient,
		classID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}

	args = append(args, extraArgs...)
	return clitestutil.ExecTestCLICmd(clientCtx, nftcli.GetCmdTransferClass(), args)
}
