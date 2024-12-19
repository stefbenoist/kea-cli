package main

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {

	rootCmd.AddCommand(&cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

$ source <(kea-cli completion bash)

# To load completions for each session, execute once:
Linux:
  $ kea-cli completion bash > /etc/bash_completion.d/kea-cli
MacOS:
  $ kea-cli completion bash > /usr/local/etc/bash_completion.d/kea-cli

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ kea-cli completion zsh > "${fpath[1]}/_kea-cli"

# You will need to start a new shell for this setup to take effect.

Fish:

$ kea-cli completion fish | source

# To load completions for each session, execute once:
$ kea-cli completion fish > ~/.config/fish/completions/kea-cli.fish

Powershell:

PS> kea-cli completion powershell | Out-String | Invoke-Expression

# To load completions for every new session, run:
PS> kea-cli completion powershell > kea-cli.ps1
# and source this file from your powershell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletion(os.Stdout)
			}
		},
	})
}
