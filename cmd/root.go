package cmd

import (
	"github.com/andresdefi/gpc/cmd/apks"
	"github.com/andresdefi/gpc/cmd/apprecovery"
	"github.com/andresdefi/gpc/cmd/apps"
	"github.com/andresdefi/gpc/cmd/auth"
	"github.com/andresdefi/gpc/cmd/baseplans"
	"github.com/andresdefi/gpc/cmd/bundles"
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
	"github.com/andresdefi/gpc/internal/exitcode"
	"github.com/spf13/cobra"
)

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
	rootCmd.PersistentFlags().StringP("output", "o", "", "Output format: json or table (default: auto-detect)")

	// Auth & config
	rootCmd.AddCommand(auth.NewCmd())
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
}

// Execute runs the root command and returns the exit code.
func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		// Check if this is an API error with a status code.
		if apiErr, ok := err.(*exitcode.ExitError); ok {
			return apiErr.Code
		}
		return exitcode.Error
	}
	return exitcode.Success
}
