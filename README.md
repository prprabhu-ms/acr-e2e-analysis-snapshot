# 🚧 acr-e2e-analysis 🚧

An analysis of browser test health in the CI for [Azure Communication Services UI library repository](https://github.com/azure/communication-ui-library).

## Snapshot ❗❓

This repository is a a ~ single git-sha snapshot of an earlier repository at
https://github.com/prprabhu-ms/acr-e2e-analysis

Objects in `data/*` are so large that they must be stored in [GitHub's gitLFS](https://docs.github.com/en/repositories/working-with-files/managing-large-files/about-git-large-file-storage) and repeatedly updating those items costs 💲. I deleted that repository and left this snapshot for posterity.

If you want to continue using this analysis tooling, fork the repository, and either update the data locally without pushing to GitHub, or better still, offload the data to some blob store to stop paying GitHub LFS charges.

## Crash course in using this repository

### Analyze existing data

[`data/raw/`](./data/raw/) should contain prefetched historical data for CI. This snapshot instead contains the folder compressed as [`./data.zip`](./data.zip). First, decompress it into [`data/`](./data) so the rest of the tooling can find it.

[`analysis/`](./analysis/) contains a set of Jupyter notebooks to analyze this data.

Best way to regenerate these notebooks is to run [`analysis/refresh.ps1`](./analysis/refresh.ps1) on PowerShell. This script will populate the [`data/cleaned`](./data/cleaned/) with cleaned up data and recreate all the other Jupyter notebooks. See [`README.md`](./analysis/README.md) for instructions for setting up your Jupyter environment.

### Fetching CI data

[fetch/](./fetch/) contains a Go binary to fetch public data about the repository. See instructions in the [`README.md`](./fetch/README.md) for setup instructions and the best way to fetch the data.