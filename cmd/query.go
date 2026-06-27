package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/AmadlaOrg/hery/entity/query"
	"github.com/AmadlaOrg/hery/entity/resolve"
	"github.com/spf13/cobra"
)

var QueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query entities",
	Long: "Query entities from the SQLite cache, from a file/stdin with --from, or\n" +
		"from a directory of .hery files with --dir (resolves _extends/_requires like\n" +
		"`hery compose --dir`, then queries in memory — no shell pipe needed).\n\n" +
		"Selection (--type/--meta/--tag) filters the entity set; --jq transforms each\n" +
		"result. Output defaults to a table; pass -o json (or yaml) for pipelines, or\n" +
		"wrap results in a HERY envelope with --hery.\n\n" +
		"Exit codes: 0 = results returned, 2 = no match, 1 = error.",
	// query is a data pipe: keep cobra from printing usage/errors so failures
	// surface cleanly on stderr and exit codes stay meaningful.
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		out, code, err := runQuery(cmd)
		if err != nil {
			fmt.Fprintln(os.Stderr, "hery query:", err)
			os.Exit(1)
		}
		if out != "" {
			fmt.Println(out)
		}
		if code != 0 {
			os.Exit(code)
		}
		return nil
	},
}

// runQuery executes the query and returns the formatted output, the intended
// process exit code (0 = results, 2 = empty), and any runtime error.
func runQuery(cmd *cobra.Command) (string, int, error) {
	from, _ := cmd.Flags().GetString("from")
	dir, _ := cmd.Flags().GetString("dir")
	outputFlag, _ := cmd.Flags().GetString("output")
	heryFlag, _ := cmd.Flags().GetBool("hery")

	if dir != "" && from != "" {
		return "", 0, fmt.Errorf("--dir and --from are mutually exclusive")
	}

	// A HERY envelope is structured and can't be rendered as a table. With the
	// default table format, promote to JSON; if the user explicitly asked for a
	// table, that's a contradiction.
	if heryFlag && outputFlag == "table" {
		if cmd.Flags().Changed("output") {
			return "", 0, fmt.Errorf("--hery cannot be combined with -o table (use -o json or -o yaml)")
		}
		outputFlag = "json"
	}

	opts := query.SelectionOpts{}
	opts.Type, _ = cmd.Flags().GetString("type")
	opts.Meta, _ = cmd.Flags().GetString("meta")
	opts.Tag, _ = cmd.Flags().GetString("tag")
	opts.JQ, _ = cmd.Flags().GetString("jq")

	var results []map[string]any
	var err error
	switch {
	case dir != "":
		results, err = queryFromDir(dir, opts)
	case from != "":
		results, err = queryFromInput(from, opts)
	default:
		results, err = queryFromCache(opts)
	}
	if err != nil {
		return "", 0, err
	}

	output, err := formatResults(results, outputFlag, heryFlag)
	if err != nil {
		return "", 0, err
	}

	// Exit 2 on an empty selection so shell pipelines can branch on a match,
	// mirroring grep. Errors are reported separately with exit 1.
	if len(results) == 0 {
		return output, 2, nil
	}
	return output, 0, nil
}

// queryFromInput reads entities from a file (or stdin when from == "-") and
// runs the query in memory.
func queryFromInput(from string, opts query.SelectionOpts) ([]map[string]any, error) {
	r := os.Stdin
	if from != "-" {
		f, err := os.Open(from)
		if err != nil {
			return nil, fmt.Errorf("failed to open input %q: %w", from, err)
		}
		defer f.Close()
		r = f
	}

	docs, err := query.LoadDocs(r)
	if err != nil {
		return nil, err
	}
	return query.QueryDocs(docs, opts)
}

// queryFromDir resolves a directory of .hery files the same way `hery compose
// --dir` does (_extends merge, _requires ordering, layering), then queries the
// resolved graph in memory. It round-trips through the multi-doc YAML stream so
// the semantics are identical to `hery compose --dir <dir> | hery query --from -`.
func queryFromDir(dir string, opts query.SelectionOpts) ([]map[string]any, error) {
	res, err := resolve.New().Resolve(dir)
	if err != nil {
		return nil, err
	}
	out, err := resolve.MarshalAll(res.Layers)
	if err != nil {
		return nil, err
	}
	docs, err := query.LoadDocs(bytes.NewReader(out))
	if err != nil {
		return nil, err
	}
	return query.QueryDocs(docs, opts)
}

// queryFromCache reads entities from the SQLite cache (default behavior).
func queryFromCache(opts query.SelectionOpts) ([]map[string]any, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	dbPath := filepath.Join(homeDir, ".cache", "hery", "hery.db")

	db := database.New(dbPath)
	if err := db.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to open cache database: %w", err)
	}
	defer db.Close()

	return query.New(db).Query(opts)
}

// formatResults renders results in the requested format, optionally wrapped in
// a HERY envelope. The envelope applies to json/yaml only; table output stays
// human-readable and ignores --hery.
func formatResults(results []map[string]any, outputFlag string, hery bool) (string, error) {
	switch outputFlag {
	case "table":
		return query.FormatTable(results)
	case "yaml":
		if hery {
			return query.FormatYAMLValue(query.WrapHERY(results))
		}
		return query.FormatYAML(results)
	case "json":
		if hery {
			return query.FormatJSON(query.WrapHERY(results))
		}
		return query.FormatJSON(results)
	default:
		return "", fmt.Errorf("invalid output format %q (want table, json or yaml)", outputFlag)
	}
}

func init() {
	QueryCmd.Flags().String("type", "", "Filter by entity type (glob pattern)")
	QueryCmd.Flags().String("meta", "", "Filter by metadata content (substring)")
	QueryCmd.Flags().String("tag", "", "Filter by tag (substring in meta)")
	QueryCmd.Flags().String("jq", "", "jq expression for transformation")
	QueryCmd.Flags().StringP("from", "f", "", "Read entities from a YAML/JSON file ('-' for stdin) instead of the cache")
	QueryCmd.Flags().String("dir", "", "Resolve and query a directory of .hery files (like 'compose --dir'); mutually exclusive with --from")
	QueryCmd.Flags().StringP("output", "o", "table", "Output format: table, json or yaml")
	QueryCmd.Flags().Bool("hery", false, "Wrap output in a HERY envelope (json/yaml only)")
}
