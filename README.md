<!-- markdownlint-disable MD041 MD010 -->

## `up-go`

<img align="right" width="250px" src="docs/logo.png">

```diff
+ 📚 A Go abstraction over the Up Bank API.
+    https://developer.up.com.au/docs.
```

```diff
- This is a WORK-IN-PROGRESS. Use at your own risk!

! This repository has no ties to Up.
```

<a href="LICENSE" target="_blank"><img src="https://img.shields.io/github/license/jmpa-io/up-go.svg" alt="GitHub License"></a>
[![CI/CD](https://github.com/jmpa-io/up-go/actions/workflows/cicd.yml/badge.svg)](https://github.com/jmpa-io/up-go/actions/workflows/cicd.yml)
[![Automerge](https://github.com/jmpa-io/up-go/actions/workflows/.github/workflows/dependabot-automerge.yml/badge.svg)](https://github.com/jmpa-io/up-go/actions/workflows/.github/workflows/dependabot-automerge.yml)

## `API Coverage`

The following API endpoints are currently covered by this package:

- [x] [List accounts](https://developer.up.com.au/#get_accounts).
- [ ] [Get an account by id](https://developer.up.com.au/#get_accounts_id).
- [ ] [List attachments](https://developer.up.com.au/#get_attachments).
- [ ] [Get an attachment by id](https://developer.up.com.au/#get_attachments_id).
- [ ] [List categories](https://developer.up.com.au/#get_categories).
- [ ] [Get a category by id](https://developer.up.com.au/#get_categories).
- [ ] [Add a category to a transaction](https://developer.up.com.au/#patch_transactions_transactionId_relationships_category).
- [x] [List tags](https://developer.up.com.au/#get_tags).
- [ ] [Add tags to a transaction](https://developer.up.com.au/#post_transactions_transactionId_relationships_tags).
- [ ] [Remote tags from a transaction](https://developer.up.com.au/#delete_transactions_transactionId_relationships_tags).
- [x] [List transactions](https://developer.up.com.au/#get_transactions).
- [ ] [Get a transaction by id](https://developer.up.com.au/#get_transactions_id).
- [ ] [List transactions by account](https://developer.up.com.au/#get_accounts_accountId_transactions).
- [x] [Utility - Ping](https://developer.up.com.au/#get_util_ping).
- [ ] [List webhooks](https://developer.up.com.au/#get_webhooks).
- [ ] [Create a webhook](https://developer.up.com.au/#post_webhooks).
- [ ] [Get a webhook by id](https://developer.up.com.au/#get_webhooks_id).
- [ ] [Delete webhook](https://developer.up.com.au/#delete_webhooks_id).
- [ ] [Ping a webhook](https://developer.up.com.au/#post_webhooks_webhookId_ping).
- [ ] [List webhook logs](https://developer.up.com.au/#get_webhooks_webhookId_logs).

## `Usage`

Below is a basic example of how to get started with this package:

```go
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jmpa-io/up-go"
)

func main() {

	// setup tracing.
	ctx := context.TODO()

	// retrieve token.
	token := os.Getenv("UP_TOKEN")

	// setup client.
	c, err := up.New(ctx, token, up.WithLogLevel(slog.LevelWarn))
	if err != nil {
		fmt.Printf("failed to setup client: %v\n", err)
		os.Exit(1)
	}

	// ping!
	p, err := c.Ping(ctx)
	if err != nil {
		fmt.Printf("failed to ping: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", p.Meta.StatusEmoji)
}
```

For more explicit examples, see the `cmd/*/main.go` files for details.

## `License`

This work is published under the MIT license.

Please see the [`LICENSE`](./LICENSE) file for details.
