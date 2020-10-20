package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show release list",
	Long:  `Show release list`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Requires a git repo in format: SCHEME:HOST/PATH. ex: https://github.com/abcd/xyz")
		}

		// URL parse method requires scheme for it to parse properly.
		if !strings.HasPrefix(args[0], "http") {
			return errors.New("Scheme missing. Provide a git repo in format ex:  https://github.com/abcd/xyz")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		checkFlags()
		listReleases(args)
	},
}

var prerelease string
var short bool

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&prerelease, "prerelease", "", "false", "Wether to fetch pre-releases")
	listCmd.Flags().BoolVarP(&short, "short", "", false, "Print only tag versions without extra formatting")
}

func checkFlags() {
	if !(prerelease == "true" || prerelease == "false") {
		fmt.Println("The prerelease flag must be either of true or false. Default is false.")
		os.Exit(1)
	}
}

func listReleases(args []string) {
	var repoName = args[0]

	scheme, host, path, _, err := splitURL(repoName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Ref: https://docs.github.com/en/free-pro-team@latest/rest/reference/repos#list-releases
	URLObj := url.URL{
		Scheme:   scheme,
		Host:     "api." + host,
		Path:     "repos" + path + "/releases",
		RawQuery: "per_page=20",
	}

	releaseURL := URLObj.String()

	resp, err := http.Get(releaseURL)
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}

	var data []map[string]interface{}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(err)
	}

	var header = []string{"tag_name", "pre-release", "name"}
	var rowData [][]string
	for i, r := range data {
		var row []string
		var releaseData = make(map[string]string)
		for k, v := range r {
			if k == "tag_name" || k == "name" || k == "prerelease" {
				releaseData[k] = convString(v)
				continue
			}
			delete(data[i], k)
		}
		row = append(row, releaseData["tag_name"])
		row = append(row, releaseData["prerelease"])
		row = append(row, releaseData["name"])

		if !(prerelease == "false" && releaseData["prerelease"] == "true") {
			rowData = append(rowData, row)
		}

	}

	if short {
		for _, row := range rowData {
			fmt.Println(row[0])
		}
	} else {
		writeTable(header, rowData)
	}
}

func splitURL(URL string) (string, string, string, map[string][]string, error) {
	u, err := url.Parse(URL)
	var queryStrings = make(map[string][]string)

	if err != nil {
		return "", "", "", queryStrings, err
	}

	for k, v := range u.Query() {
		queryStrings[k] = v
	}

	return u.Scheme, u.Host, u.Path, queryStrings, nil
}

func convString(input interface{}) string {
	switch dataType := input.(type) {
	case string:
	case float64:
		return input.(string)
	case bool:
		return strconv.FormatBool(input.(bool))
	default:
		fmt.Println("unsupported type", dataType)
	}

	return input.(string)
}

func writeTable(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

// func getReleases(total int,  per_page int) {

// }
