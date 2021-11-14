[![cicd](https://github.com/jmpa-oss/up/actions/workflows/cicd.yml/badge.svg)](https://github.com/jmpa-oss/up/actions/workflows/cicd.yml)
[![dispatch](https://github.com/jmpa-oss/up/actions/workflows/dispatch.yml/badge.svg)](https://github.com/jmpa-oss/up/actions/workflows/dispatch.yml)
[![README.md](https://github.com/jmpa-oss/up/actions/workflows/README.md.yml/badge.svg)](https://github.com/jmpa-oss/up/actions/workflows/README.md.yml)
[![template-cleanup](https://github.com/jmpa-oss/up/actions/workflows/template-cleanup.yml/badge.svg)](https://github.com/jmpa-oss/up/actions/workflows/template-cleanup.yml)
[![update](https://github.com/jmpa-oss/up/actions/workflows/update.yml/badge.svg)](https://github.com/jmpa-oss/up/actions/workflows/update.yml)

<p align="center">
	<img src="docs/logo.png">
</p>

# up

```diff
+ A Go abstraction over the Up Bank API. 
```

## Workflows

workflow|description
---|---
[cicd](.github/workflows/cicd.yml)|Runs the CI/CD for the repository.
[dispatch](.github/workflows/dispatch.yml)|Pushes repository_dispatch events out to repositories built from this template.
[README.md](.github/workflows/README.md.yml)|Updates the README.md with new changes.
[template-cleanup](.github/workflows/template-cleanup.yml)|Cleans up the repository when a child is first created; triggers from the first commit to the repository.
[update](.github/workflows/update.yml)|Updates repository with changes from parent template.

