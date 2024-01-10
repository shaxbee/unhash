GOCOVERPKG := github.com/shaxbee/unhash/...

include makefiles/go.mk

BENCHSTAT := go run golang.org/x/perf/cmd/benchstat@latest 

benchmark:
	go test -run '^$$' -bench=. -benchmem -count 10 github.com/shaxbee/unhash/e2e >$(TMPDIR)/benchmark.txt
	$(BENCHSTAT) -col '/algo' $(TMPDIR)/benchmark.txt