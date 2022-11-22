## Auth

This package uses [go-gh](https://github.com/cli/go-gh) with default auth.

You must have `gh` installed and authenticated using a config file. On Codespaces, you need to override the default authentication via `GITHUB_TOKEN`.

```sh
    # In a throwaway terminal
    $ unset GITHUB_TOKEN
    $ gh auth login
    # ... follow prompts
```

Then run

```sh
    go run main.go -help
```

### Incremental fetch

GitHub rate limits request to its REST API. This tool is not super smart about how the data is fetched. In fact, I think that the GitHub REST API is fairly limiting in how much the API calls can be optimized.

The tool implements some flags to batch requests to GitHub and sleep between consecutive API requests to avoid getting throttled. In addition, the tool implements retry with backoff for retryable error codes.

If you are fetching a lot of data (at the current rate of CI data generation, anything more than a couple days worth of CI data), you can get good results with

```sh
    go run .\main.go -outdir ..\data\raw -batch 50 -after 2022-11-02
```

to fetch all the data for CI runs since 2022-11-02, or

```sh
    go run .\main.go -outdir ..\data\raw -batch 50 -after 2022-11-02 -before 2022-12-02
```

to fetch all data for CI runs between 2022-11-02 and 2022-12-02.

### You said how long?

Note that this will take many hours as the tool waits for 5 minutes regularly to avoid GitHub API throttling. As a data point, fetching a month's worth of data recently took about 7 hours to complete.

[Happy slacking off friend!](https://xkcd.com/303/)

For the full list of supported command-line options, run

```sh
    go run main.go -help
```