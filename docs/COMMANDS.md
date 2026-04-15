# gpc Command Reference

Auto-generated from `gpc --help` on 2026-04-15.

gpc is a fast, lightweight, and scriptable CLI for the Google Play Developer API.

It provides complete coverage of the Android Publisher API v3, letting you manage
apps, releases, in-app products, subscriptions, reviews, and more from your terminal.

Usage:
  gpc [command]

GETTING STARTED:
  auth        Manage authentication
  config      Manage CLI configuration
  doctor      Check your gpc setup
  version     Print the version of gpc
  completion  Generate the autocompletion script for the specified shell

APP MANAGEMENT:
  apps   Manage apps
  edits  Manage edit sessions

RELEASE PIPELINE:
  releases             Manage releases
  tracks               Manage release tracks
  apks                 Manage APKs
  bundles              Manage app bundles
  deobfuscation        Manage deobfuscation files
  expansionfiles       Manage expansion files
  countryavailability  Manage country availability

MONETIZATION:
  iap              Manage in-app products
  subscriptions    Manage subscriptions
  baseplans        Manage subscription base plans
  offers           Manage subscription offers
  onetimeproducts  Manage one-time products
  purchaseoptions  Manage purchase options
  otpoffers        Manage one-time product offers
  pricing          Pricing utilities

STORE PRESENCE:
  listings    Manage store listings
  images      Manage store listing images
  details     Manage app details
  testers     Manage testers
  reviews     Manage reviews
  datasafety  Manage data safety declarations

APP VITALS:
  vitals  [beta] App vitals and quality metrics

ORDERS & PURCHASES:
  orders     Manage orders
  purchases  Manage purchases

ACCOUNT MANAGEMENT:
  users   Manage users
  grants  Manage grants

DEVICE & RECOVERY:
  devices               Manage device tier configs
  apprecovery           Manage app recovery actions
  externaltransactions  Manage external transactions

APK VARIANTS:
  generatedapks    Manage generated APKs
  systemapks       Manage system APK variants
  internalsharing  Manage internal app sharing

Flags:
  -h, --help             help for gpc
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "gpc [command] --help" for more information about a command.

## apks

```
List, upload, or add externally hosted APKs for an app.

Usage:
  apks [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  add-externally-hosted  Add an externally hosted APK
  list                   List all APKs
  upload                 Upload an APK

Flags:
  -h, --help             help for apks
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "apks [command] --help" for more information about a command.
```

## apprecovery

```
List, create, deploy, cancel, or add targeting to app recovery actions.

Usage:
  apprecovery [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  add-targeting  Add targeting to an app recovery action
  cancel         Cancel an app recovery action
  create         Create an app recovery action
  deploy         Deploy an app recovery action
  list           List app recovery actions

Flags:
  -h, --help             help for apprecovery
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "apprecovery [command] --help" for more information about a command.
```

## apps

```


Usage:
  apps [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  get  Get app details

Flags:
  -h, --help             help for apps
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "apps [command] --help" for more information about a command.
```

## auth

```
Authenticate with the Google Play Developer API using a service account key file.

Usage:
  auth [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  login   Authenticate with a service account key file
  logout  Remove stored credentials
  status  Show authentication status

Flags:
  -h, --help             help for auth
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "auth [command] --help" for more information about a command.
```

## baseplans

```


Usage:
  baseplans [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  activate              Activate a base plan
  batch-migrate-prices  Batch migrate prices for base plans
  batch-update-states   Batch update base plan states
  deactivate            Deactivate a base plan
  delete                Delete a base plan
  migrate-prices        Migrate prices for a base plan

Flags:
  -h, --help             help for baseplans
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "baseplans [command] --help" for more information about a command.
```

## bundles

```
List or upload Android App Bundles (AAB) for an app.

Usage:
  bundles [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  list    List all bundles
  upload  Upload an app bundle

Flags:
  -h, --help             help for bundles
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "bundles [command] --help" for more information about a command.
```

## completion

```
Generate the autocompletion script for gpc for the specified shell.
See each sub-command's help for details on how to use the generated script.


Usage:
  completion [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  bash        Generate the autocompletion script for bash
  fish        Generate the autocompletion script for fish
  powershell  Generate the autocompletion script for powershell
  zsh         Generate the autocompletion script for zsh

Flags:
  -h, --help             help for completion
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "completion [command] --help" for more information about a command.
```

## config

```
Get and set gpc configuration values stored in ~/.gpc/config.json.

Usage:
  config [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  get   Get a configuration value
  list  List all configuration values
  path  Print the config file path
  set   Set a configuration value

Flags:
  -h, --help             help for config
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "config [command] --help" for more information about a command.
```

## countryavailability

```
Get country availability for a track.

Usage:
  countryavailability [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  get  Get country availability for a track

Flags:
  -h, --help             help for countryavailability
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "countryavailability [command] --help" for more information about a command.
```

## datasafety

```
Update data safety declarations for an app.

Usage:
  datasafety [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  update  Update data safety declarations

Flags:
  -h, --help             help for datasafety
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "datasafety [command] --help" for more information about a command.
```

## deobfuscation

```
Upload ProGuard mapping or native debug symbol files for an APK.

Usage:
  deobfuscation [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  upload  Upload a deobfuscation file

Flags:
  -h, --help             help for deobfuscation
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "deobfuscation [command] --help" for more information about a command.
```

## details

```
Get or update app details (contact info, default language).

Usage:
  details [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  get     Get app details
  update  Update app details

Flags:
  -h, --help             help for details
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "details [command] --help" for more information about a command.
```

## devices

```
List, get, or create device tier configurations for an app.

Usage:
  devices [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  create-tier-config  Create a device tier config
  get-tier-config     Get a device tier config
  list-tier-configs   List device tier configs

Flags:
  -h, --help             help for devices
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "devices [command] --help" for more information about a command.
```

## doctor

```
Validate that gpc is configured correctly.

Checks:
  - CLI version and Go runtime
  - Config file exists and is readable
  - Service account credentials are stored
  - OAuth2 token can be obtained
  - Google Play API is reachable

Usage:
  doctor [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

Flags:
  -h, --help             help for doctor
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "doctor [command] --help" for more information about a command.
```

## edits

```
Create, validate, commit, or delete edit sessions. Edits are transactional containers for app changes.

Usage:
  edits [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  commit    Commit an edit session
  create    Create a new edit session
  delete    Delete an edit session
  get       Get an edit session
  validate  Validate an edit session

Flags:
  -h, --help             help for edits
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "edits [command] --help" for more information about a command.
```

## expansionfiles

```
Get, update, or upload expansion files (OBBs) for an app.

Usage:
  expansionfiles [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  get     Get an expansion file
  update  Update an expansion file
  upload  Upload an expansion file

Flags:
  -h, --help             help for expansionfiles
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "expansionfiles [command] --help" for more information about a command.
```

## externaltransactions

```
Create, get, or refund external transactions for an app.

Usage:
  externaltransactions [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  create  Create an external transaction
  get     Get an external transaction
  refund  Refund an external transaction

Flags:
  -h, --help             help for externaltransactions
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "externaltransactions [command] --help" for more information about a command.
```

## generatedapks

```
List or download generated APKs for an app.

Usage:
  generatedapks [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  download  Download a generated APK
  list      List generated APKs

Flags:
  -h, --help             help for generatedapks
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "generatedapks [command] --help" for more information about a command.
```

## grants

```
Create, update, or delete user grants.

Usage:
  grants [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  create  Create a grant
  delete  Delete a grant
  update  Update a grant

Flags:
  -h, --help             help for grants
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "grants [command] --help" for more information about a command.
```

## iap

```


Usage:
  iap [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  batch-delete  Batch delete in-app products
  batch-get     Batch get in-app products
  batch-update  Batch update in-app products
  create        Create an in-app product
  delete        Delete an in-app product
  get           Get an in-app product
  list          List in-app products
  update        Update an in-app product

Flags:
  -h, --help             help for iap
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "iap [command] --help" for more information about a command.
```

## images

```
List, upload, or delete images for store listings.

Usage:
  images [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  delete      Delete a specific image
  delete-all  Delete all images for a type
  list        List images for a listing
  upload      Upload an image

Flags:
  -h, --help             help for images
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "images [command] --help" for more information about a command.
```

## internalsharing

```
Upload APKs or bundles for internal app sharing.

Usage:
  internalsharing [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  upload-apk     Upload an APK for internal sharing
  upload-bundle  Upload a bundle for internal sharing

Flags:
  -h, --help             help for internalsharing
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "internalsharing [command] --help" for more information about a command.
```

## listings

```
List, get, update, or delete store listings for an app.

Usage:
  listings [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  delete      Delete a store listing by language
  delete-all  Delete all store listings
  get         Get a store listing by language
  list        List all store listings
  update      Update a store listing

Flags:
  -h, --help             help for listings
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "listings [command] --help" for more information about a command.
```

## offers

```


Usage:
  offers [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  activate             Activate an offer
  batch-get            Batch get offers
  batch-update         Batch update offers
  batch-update-states  Batch update offer states
  create               Create an offer
  deactivate           Deactivate an offer
  delete               Delete an offer
  get                  Get an offer
  list                 List offers
  update               Update an offer

Flags:
  -h, --help             help for offers
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "offers [command] --help" for more information about a command.
```

## onetimeproducts

```


Usage:
  onetimeproducts [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  batch-delete  Batch delete one-time products
  batch-get     Batch get one-time products
  batch-update  Batch update one-time products
  create        Create a one-time product
  delete        Delete a one-time product
  get           Get a one-time product
  list          List one-time products
  update        Update a one-time product

Flags:
  -h, --help             help for onetimeproducts
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "onetimeproducts [command] --help" for more information about a command.
```

## orders

```
Get, batch-get, or refund orders for an app.

Usage:
  orders [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  batch-get  Batch get orders
  get        Get an order
  refund     Refund an order

Flags:
  -h, --help             help for orders
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "orders [command] --help" for more information about a command.
```

## otpoffers

```


Usage:
  otpoffers [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  activate             Activate a one-time product offer
  batch-delete         Batch delete one-time product offers
  batch-get            Batch get one-time product offers
  batch-update         Batch update one-time product offers
  batch-update-states  Batch update one-time product offer states
  cancel               Cancel a one-time product offer
  deactivate           Deactivate a one-time product offer
  list                 List one-time product offers

Flags:
  -h, --help             help for otpoffers
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "otpoffers [command] --help" for more information about a command.
```

## pricing

```


Usage:
  pricing [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  convert  Convert region prices

Flags:
  -h, --help             help for pricing
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "pricing [command] --help" for more information about a command.
```

## purchaseoptions

```


Usage:
  purchaseoptions [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  batch-delete         Batch delete purchase options
  batch-update-states  Batch update purchase option states

Flags:
  -h, --help             help for purchaseoptions
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "purchaseoptions [command] --help" for more information about a command.
```

## purchases

```
Manage in-app product purchases, subscriptions, and voided purchases.

Usage:
  purchases [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:
  subscriptions  Manage subscription purchases

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  products          Manage product purchases
  products-v2       Manage product purchases (v2)
  subscriptions-v2  Manage subscription purchases (v2)
  voided            Manage voided purchases

Flags:
  -h, --help             help for purchases
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "purchases [command] --help" for more information about a command.
```

## releases

```
List releases, deploy artifacts, promote between tracks, and manage rollouts.

Usage:
  releases [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  deploy   Deploy an APK or AAB to a track
  halt     Halt a staged rollout
  list     List releases for a track
  promote  Promote a release between tracks
  rollout  Update staged rollout fraction

Flags:
  -h, --help             help for releases
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "releases [command] --help" for more information about a command.
```

## reviews

```
List, get, or reply to user reviews.

Usage:
  reviews [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  get    Get a review
  list   List reviews
  reply  Reply to a review

Flags:
  -h, --help             help for reviews
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "reviews [command] --help" for more information about a command.
```

## subscriptions

```


Usage:
  subscriptions [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  archive       Archive a subscription
  batch-get     Batch get subscriptions
  batch-update  Batch update subscriptions
  create        Create a subscription
  delete        Delete a subscription
  get           Get a subscription
  list          List subscriptions
  update        Update a subscription

Flags:
  -h, --help             help for subscriptions
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "subscriptions [command] --help" for more information about a command.
```

## systemapks

```
List, get, create, or download system APK variants for an app.

Usage:
  systemapks [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  create    Create a system APK variant
  download  Download a system APK variant
  get       Get a system APK variant
  list      List system APK variants

Flags:
  -h, --help             help for systemapks
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "systemapks [command] --help" for more information about a command.
```

## testers

```
Get or update testers for a release track.

Usage:
  testers [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  get     Get testers for a track
  update  Update testers for a track

Flags:
  -h, --help             help for testers
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "testers [command] --help" for more information about a command.
```

## tracks

```
List, get, update, or create release tracks (internal, alpha, beta, production, custom).

Usage:
  tracks [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  create  Create a custom track
  get     Get a track
  list    List all tracks
  update  Update a track

Flags:
  -h, --help             help for tracks
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "tracks [command] --help" for more information about a command.
```

## users

```
List, create, update, or delete developer account users.

Usage:
  users [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  create  Create a user
  delete  Delete a user
  list    List users
  update  Update a user

Flags:
  -h, --help             help for users
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "users [command] --help" for more information about a command.
```

## version

```


Usage:
  version [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

Flags:
  -h, --help             help for version
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "version [command] --help" for more information about a command.
```

## vitals

```
View crash rates, ANR rates, startup performance, and other quality
metrics from the Play Developer Reporting API.

These commands provide the same data visible in Play Console's "Android Vitals"
section, accessible from your terminal or CI/CD pipelines.

Usage:
  vitals [command]

GETTING STARTED:

APP MANAGEMENT:

RELEASE PIPELINE:

MONETIZATION:

STORE PRESENCE:

APP VITALS:

ORDERS & PURCHASES:

ACCOUNT MANAGEMENT:

DEVICE & RECOVERY:

APK VARIANTS:

ADDITIONAL COMMANDS:
  anrs       View ANR (Application Not Responding) rate metrics
  battery    View battery usage metrics
  crashes    View crash rate metrics
  errors     View error counts and clusters
  overview   Show vitals overview for an app
  rendering  View slow rendering metrics
  startup    View app startup time metrics

Flags:
  -h, --help             help for vitals
  -o, --output string    Output format: json, table, csv, or yaml (default: auto-detect)
  -p, --package string   Android package name (e.g. com.example.app)

Use "vitals [command] --help" for more information about a command.
```

