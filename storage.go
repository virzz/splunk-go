package splunk

type storage struct {
	KV *kv
	// CSV *csv
}

var (
	Storage = &storage{
		KV: &kv{owner: "nobody", app: "search"},
		// CSV: &csv{owner: "nobody", app: "search"},
	}
)
