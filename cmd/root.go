package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/andresdefi/gpc/cmd/apks"
	"github.com/andresdefi/gpc/cmd/apprecovery"
	"github.com/andresdefi/gpc/cmd/apps"
	"github.com/andresdefi/gpc/cmd/auth"
	"github.com/andresdefi/gpc/cmd/baseplans"
	"github.com/andresdefi/gpc/cmd/bundles"
	cfgcmd "github.com/andresdefi/gpc/cmd/config"
	"github.com/andresdefi/gpc/cmd/countryavailability"
	"github.com/andresdefi/gpc/cmd/datasafety"
	"github.com/andresdefi/gpc/cmd/deobfuscation"
	"github.com/andresdefi/gpc/cmd/details"
	"github.com/andresdefi/gpc/cmd/devices"
	"github.com/andresdefi/gpc/cmd/edits"
	"github.com/andresdefi/gpc/cmd/expansionfiles"
	"github.com/andresdefi/gpc/cmd/externaltransactions"
	"github.com/andresdefi/gpc/cmd/generatedapks"
	"github.com/andresdefi/gpc/cmd/grants"
	"github.com/andresdefi/gpc/cmd/iap"
	"github.com/andresdefi/gpc/cmd/images"
	"github.com/andresdefi/gpc/cmd/internalsharing"
	"github.com/andresdefi/gpc/cmd/listings"
	"github.com/andresdefi/gpc/cmd/offers"
	"github.com/andresdefi/gpc/cmd/onetimeproducts"
	"github.com/andresdefi/gpc/cmd/orders"
	"github.com/andresdefi/gpc/cmd/otpoffers"
	"github.com/andresdefi/gpc/cmd/pricing"
	"github.com/andresdefi/gpc/cmd/purchaseoptions"
	"github.com/andresdefi/gpc/cmd/purchases"
	"github.com/andresdefi/gpc/cmd/releases"
	"github.com/andresdefi/gpc/cmd/reviews"
	"github.com/andresdefi/gpc/cmd/subscriptions"
	"github.com/andresdefi/gpc/cmd/systemapks"
	"github.com/andresdefi/gpc/cmd/testers"
	"github.com/andresdefi/gpc/cmd/tracks"
	"github.com/andresdefi/gpc/cmd/users"
	"github.com/andresdefi/gpc/cmd/vitals"
	"github.com/andresdefi/gpc/internal/exitcode"
	"github.com/andresdefi/gpc/internal/output"
	"github.com/spf13/cobra"
)

// commandGroup defines a logical grouping of commands for help display.
type commandGroup struct {
	Title    string
	Commands []string
}

var groups = []commandGroup{
	{
		Title:    "GETTING STARTED",
		Commands: []string{"auth", "config", "version", "completion"},
	},
	{
		Title:    "APP MANAGEMENT",
		Commands: []string{"apps", "edits"},
	},
	{
		Title:    "RELEASE PIPELINE",
		Commands: []string{"releases", "tracks", "apks", "bundles", "deobfuscation", "expansionfiles", "countryavailability"},
	},
	{
		Title:    "MONETIZATION",
		Commands: []string{"iap", "subscriptions", "baseplans", "offers", "onetimeproducts", "purchaseoptions", "otpoffers", "pricing"},
	},
	{
		Title:    "STORE PRESENCE",
		Commands: []string{"listings", "images", "details", "testers", "reviews", "datasafety"},
	},
	{
		Title:    "APP VITALS",
		Commands: []string{"vitals"},
	},
	{
		Title:    "ORDERS & PURCHASES",
		Commands: []string{"orders", "purchases"},
	},
	{
		Title:    "ACCOUNT MANAGEMENT",
		Commands: []string{"users", "grants"},
	},
	{
		Title:    "DEVICE & RECOVERY",
		Commands: []string{"devices", "apprecovery", "externaltransactions"},
	},
	{
		Title:    "APK VARIANTS",
		Commands: []string{"generatedapks", "systemapks", "internalsharing"},
	},
}

var rootCmd = &cobra.Command{
	Use:   "gpc",
	Short: "Google Play Console CLI",
	Long: `gpc is a fast, lightweight, and scriptable CLI for the Google Play Developer API.

It provides complete coverage of the Android Publisher API v3, letting you manage
apps, releases, in-app products, subscriptions, reviews, and more from your terminal.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringP("package", "p", "", "Android package name (e.g. com.example.app)")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Output format: json, table, csv, or yaml (default: auto-detect)")

	// Enable fuzzy command suggestions.
	rootCmd.SuggestionsMinimumDistance = 2

	// Getting started
	rootCmd.AddCommand(auth.NewCmd())
	rootCmd.AddCommand(cfgcmd.NewCmd())
	rootCmd.AddCommand(newVersionCmd())

	// App management
	rootCmd.AddCommand(apps.NewCmd())
	rootCmd.AddCommand(edits.NewCmd())

	// Release pipeline
	rootCmd.AddCommand(tracks.NewCmd())
	rootCmd.AddCommand(releases.NewCmd())
	rootCmd.AddCommand(apks.NewCmd())
	rootCmd.AddCommand(bundles.NewCmd())
	rootCmd.AddCommand(deobfuscation.NewCmd())
	rootCmd.AddCommand(expansionfiles.NewCmd())
	rootCmd.AddCommand(countryavailability.NewCmd())

	// Monetization
	rootCmd.AddCommand(iap.NewCmd())
	rootCmd.AddCommand(subscriptions.NewCmd())
	rootCmd.AddCommand(baseplans.NewCmd())
	rootCmd.AddCommand(offers.NewCmd())
	rootCmd.AddCommand(onetimeproducts.NewCmd())
	rootCmd.AddCommand(purchaseoptions.NewCmd())
	rootCmd.AddCommand(otpoffers.NewCmd())
	rootCmd.AddCommand(pricing.NewCmd())

	// Store presence
	rootCmd.AddCommand(listings.NewCmd())
	rootCmd.AddCommand(images.NewCmd())
	rootCmd.AddCommand(details.NewCmd())
	rootCmd.AddCommand(testers.NewCmd())
	rootCmd.AddCommand(reviews.NewCmd())
	rootCmd.AddCommand(datasafety.NewCmd())

	// App Vitals
	rootCmd.AddCommand(vitals.NewCmd())

	// Orders & purchases
	rootCmd.AddCommand(orders.NewCmd())
	rootCmd.AddCommand(purchases.NewCmd())

	// Account management
	rootCmd.AddCommand(users.NewCmd())
	rootCmd.AddCommand(grants.NewCmd())

	// Device management
	rootCmd.AddCommand(devices.NewCmd())

	// App recovery
	rootCmd.AddCommand(apprecovery.NewCmd())

	// External transactions
	rootCmd.AddCommand(externaltransactions.NewCmd())

	// APK variants
	rootCmd.AddCommand(generatedapks.NewCmd())
	rootCmd.AddCommand(systemapks.NewCmd())
	rootCmd.AddCommand(internalsharing.NewCmd())

	// Override help to show grouped commands.
	rootCmd.SetHelpFunc(groupedHelp)
}

// Execute runs the root command and returns the exit code.
func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		if apiErr, ok := err.(*exitcode.ExitError); ok {
			return apiErr.Code
		}
		output.Errorf("%v", err)
		return exitcode.Error
	}
	return exitcode.Success
}

// groupedHelp renders commands organized by logical category instead of alphabetically.
func groupedHelp(cmd *cobra.Command, args []string) {
	fmt.Println(cmd.Long)
	fmt.Println()

	fmt.Println("Usage:")
	fmt.Printf("  %s [command]\n\n", cmd.Use)

	// Build a map of command name -> command for quick lookup.
	cmdMap := make(map[string]*cobra.Command)
	for _, c := range cmd.Commands() {
		cmdMap[c.Name()] = c
		for _, alias := range c.Aliases {
			cmdMap[alias] = c
		}
	}

	// Render each group.
	for _, g := range groups {
		fmt.Printf("%s:\n", g.Title)
		tw := tabwriter.NewWriter(cmd.OutOrStdout(), 2, 4, 2, ' ', 0)
		for _, name := range g.Commands {
			if c, ok := cmdMap[name]; ok {
				_, _ = fmt.Fprintf(tw, "  %s\t%s\n", c.Name(), c.Short)
			}
		}
		_ = tw.Flush()
		fmt.Println()
	}

	// Show any commands not in a group (e.g. help, completion).
	grouped := make(map[string]bool)
	for _, g := range groups {
		for _, name := range g.Commands {
			grouped[name] = true
		}
	}

	var ungrouped []*cobra.Command
	for _, c := range cmd.Commands() {
		if !grouped[c.Name()] && c.Name() != "help" {
			ungrouped = append(ungrouped, c)
		}
	}

	if len(ungrouped) > 0 {
		fmt.Println("ADDITIONAL COMMANDS:")
		tw := tabwriter.NewWriter(cmd.OutOrStdout(), 2, 4, 2, ' ', 0)
		for _, c := range ungrouped {
			_, _ = fmt.Fprintf(tw, "  %s\t%s\n", c.Name(), c.Short)
		}
		_ = tw.Flush()
		fmt.Println()
	}

	fmt.Println("Flags:")
	fmt.Println(cmd.Flags().FlagUsages())

	fmt.Printf("Use \"%s [command] --help\" for more information about a command.\n", cmd.Use)
}

