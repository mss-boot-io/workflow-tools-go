/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/4/7 15:23
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/4/7 15:23
 */

package cmd

import (
	"errors"
	"github.com/mss-boot-io/workflow-tools/cmd/change"
	"github.com/mss-boot-io/workflow-tools/cmd/gitops"
	"github.com/mss-boot-io/workflow-tools/cmd/tag"
	"os"

	"github.com/spf13/cobra"

	"github.com/mss-boot-io/workflow-tools/cmd/aws"
	"github.com/mss-boot-io/workflow-tools/cmd/dep"
	"github.com/mss-boot-io/workflow-tools/cmd/mk"
)

var rootCmd = &cobra.Command{
	Use:          "go-workflow-tools",
	Short:        "go-workflow-tools is a tool for actions workflow",
	SilenceUsage: true,
	Long:         `go-workflow-tools`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(mk.StartCmd)
	rootCmd.AddCommand(dep.StartCmd)
	rootCmd.AddCommand(aws.StartCmd)
	rootCmd.AddCommand(change.StartCmd)
	rootCmd.AddCommand(gitops.StartCmd)
	rootCmd.AddCommand(tag.StartCmd)
}

//Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
