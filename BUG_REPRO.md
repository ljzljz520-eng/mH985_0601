# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	trainingdesk/cmd/trainingd	[no test files]
ok  	trainingdesk/internal/analytics	0.001s
ok  	trainingdesk/internal/api	0.011s
ok  	trainingdesk/internal/archive	0.020s
ok  	trainingdesk/internal/catalog	0.027s
ok  	trainingdesk/internal/command	0.001s
ok  	trainingdesk/internal/exporter	0.001s
--- FAIL: Test985BusinessRegression (0.01s)
    regression_test.go:32: selected record = model.Record{ID:"record-a", StoreID:"flagship", Title:"Opening", Content:"Open", Category:"", Status:"draft", Version:1, SortKey:1, Owner:"", Reviewer:"", CreatedSeq:1, UpdatedSeq:1, PublishedSeq:0, ArchivedSeq:0}
FAIL
FAIL	trainingdesk/internal/flow008	0.041s
ok  	trainingdesk/internal/importer	0.009s
ok  	trainingdesk/internal/ledger	0.001s
ok  	trainingdesk/internal/lifecycle	0.001s
ok  	trainingdesk/internal/maintenance	0.002s
?   	trainingdesk/internal/model	[no test files]
ok  	trainingdesk/internal/notification	0.001s
ok  	trainingdesk/internal/policy	0.002s
ok  	trainingdesk/internal/query	0.001s
ok  	trainingdesk/internal/report	0.001s
ok  	trainingdesk/internal/review	0.006s
ok  	trainingdesk/internal/store	0.006s
ok  	trainingdesk/internal/validation	0.001s
ok  	trainingdesk/internal/workflow	0.002s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/trainingd): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/trainingd): exit `0`
