package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"robsonjr.com.br/utils/anchors"
	"robsonjr.com.br/utils/signals"
	"robsonjr.com.br/utils/validation"
)

var walkCmd = &cobra.Command{
	Use:   "walk",
	Short: "Walk through the root url",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			return err
		}
		if url == "" {
			return errors.New("url can't be null")
		}

		dest, err := cmd.Flags().GetString("dest")
		if err != nil {
			return err
		}
		if dest == "" {
			return errors.New("url can't be null")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := cmd.Flags().GetString("url")
		outputPath, _ := cmd.Flags().GetString("dest")

		// total concurrent jobs
		totalJobs := 2
		currJobs := 0

		// storage for urls
		pendingUrls := make([]string, 0)
		workingUrls := make(map[string]bool)

		// notify channels
		chNewUrl := make(chan string)
		chNotifyEnd := make(chan string)

		// first job, we should run with the root url
		pendingUrls = append(pendingUrls, url)

		// Subscribe to both channels
		for {
			if !signals.ShouldContinue() {
				if currJobs == 0 {
					break
				}
			} else if len(pendingUrls) == 0 && currJobs == 0 {
				// we have ended depth crawling from the initial root url
				break
			} else {
				// starting new job in case it is possible to parallelize it
				if len(pendingUrls) > 0 && currJobs < totalJobs {
					// we are getting the first element from the pending pool, and adding it to working
					popUrl := pendingUrls[0]

					// removing the url we just got from the pending work
					pendingUrls = pendingUrls[1:]

					// we should skip this job if it is already being processed
					if _, ok := workingUrls[popUrl]; !ok {
						workingUrls[popUrl] = false

						// updating the number of jobs
						currJobs++

						// properly starting the job
						go walk(chNewUrl, chNotifyEnd, popUrl, outputPath)
					}
				}
			}

			// listening for events
			select {
			case url := <-chNewUrl:
				// enqueuing new url for work if it is not
				if _, ok := workingUrls[url]; !ok {
					pendingUrls = append(pendingUrls, url)
				}
			case url := <-chNotifyEnd:
				// updating our structure with workload end
				if _, ok := workingUrls[url]; !ok {
					workingUrls[url] = true
				}

				currJobs--
			}
		}

		return nil
	},
}

func walk(chNewUrl chan string, chNotifyEnd chan string, url string, outputPath string) {
	fmt.Printf("downloading: %v\n", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("ERROR: Failed to crawl:", url)
		return
	}

	b := resp.Body
	doc, _ := ioutil.ReadAll(resp.Body)
	defer b.Close() // close Body when the function completes

	anchors := anchors.GetWalkValidPages(string(doc))
	for _, anchor := range anchors {
		if validation.IsChildrenUrl(url, anchor) {
			chNewUrl <- anchor
		}
	}

	chNotifyEnd <- url
}

func init() {
	walkCmd.Flags().String("url", "", "the url to be crawled")
	walkCmd.MarkFlagRequired("url")
	walkCmd.Flags().String("dest", "", "the destination path to store the data")
	walkCmd.MarkFlagRequired("dest")
}
