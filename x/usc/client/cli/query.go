package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/gaia/v7/x/usc/types"
	"github.com/spf13/cobra"
)

// NewQueryCmd returns the cli query commands for this module.
func NewQueryCmd() *cobra.Command {
	distQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the USC module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	distQueryCmd.AddCommand(
		newCmdQueryParams(),
		newCmdQueryPool(),
	)

	return distQueryCmd
}

// newCmdQueryParams implements the params query command.
func newCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query module params",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			req := types.QueryParamsRequest{}

			res, err := queryClient.Params(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// newCmdQueryPool implements the pool query command.
func newCmdQueryPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool",
		Short: "Query module pool balance",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			req := types.QueryPoolRequest{}

			res, err := queryClient.Pool(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
