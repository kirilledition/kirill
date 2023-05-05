package cmd

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

func validatePDBId(ids []string) ([]string, error) {
	var res []string
	isNotDigit := func(c rune) bool { return c < '0' || c > '9' }

	for i, val := range ids {
		if len(val) != 4 || strings.IndexFunc(val, isNotDigit) == -1 {
			return nil, fmt.Errorf("error in pdb %d: %s", i+1, val)
		}
		res = append(res, strings.ToLower(val))

	}
	return res, nil
}

func readPDBIdList(input []string) ([]string, error) {
	var ids []string
	_, err := os.Stat(input[0])
	if os.IsNotExist(err) {
		logger.Println("Assuming input is a list of PDB IDs")
		ids, err := validatePDBId(input)
		if err != nil {
			return nil, err
		}
		return ids, nil
	}
	logger.Println("Assuming input as a file with a list of PDB IDs")

	file, err := os.Open(input[0])
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ids = append(ids, scanner.Text())
	}
	ids, err = validatePDBId(ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

type PDBClient struct {
	scheme string
	host   string
	path   string
	client *http.Client
}

func (c *PDBClient) fetch(id string, outputPath string) error {
	url := url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   path.Join(c.path, id+".pdb.gz"),
	}
	filename := strings.ToUpper(id) + ".pdb"
	filename = path.Join(outputPath, filename)

	resp, err := c.client.Get(url.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer body.Close()

	buf, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filename, buf, 0644); err != nil {
		return err
	}

	logger.Printf("Loaded %s to %s", id, filename)

	return nil
}

func fetchPDB(input []string, outputPath string, client *PDBClient) {
	ids, err := readPDBIdList(input)
	if err != nil {
		logger.Fatalln(err)
		os.Exit(1)
	}

	for _, id := range ids {
		if err := client.fetch(id, outputPath); err != nil {
			if err != nil {
				logger.Fatalln(err)
				os.Exit(1)
			}
		}
	}
}

var fetchpdbCmd = &cobra.Command{
	Use:   "fetchpdb [PDB IDs or input file]",
	Short: "Fetch protein structures from the Protein Data Bank",
	Long: `fetchpdb is a command-line tool to download protein structures from the Protein Data Bank (PDB).
It accepts a list of PDB IDs or an input file containing PDB IDs, one per line.

Example usage:

1. Download structures for a list of PDB IDs:
   kirill fetchpdb 1abc 2def 3ghi

2. Download structures from an input file (each PDB ID on a separate line):
   kirill fetchpdb pdb_ids.txt

3. Download structures from an input file and save them to a specific output directory:
   kirill fetchpdb pdb_ids.txt -o /path/to/output`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		outputPath, _ := cmd.Flags().GetString("output")

		var logFile *os.File
		var err error

		logPath := path.Join(outputPath, "fetchpdb")
		logger, logFile, err = getLogger(logPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer logFile.Close()

		logger.Println(getCommandLine())

		client := &PDBClient{
			scheme: "https",
			host:   "files.rcsb.org",
			path:   "download",
			client: &http.Client{},
		}

		fetchPDB(args, outputPath, client)
	},
}

func init() {
	rootCmd.AddCommand(fetchpdbCmd)

	fetchpdbCmd.Flags().StringP("output", "o", ".", "Output directory")
}
