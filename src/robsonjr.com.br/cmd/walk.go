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
	"time"
)

type walkCmdArgDef struct {
	Url string
	Dest string
}

var walkCmdArg = walkCmdArgDef{}

var walkCmd = &cobra.Command{
	Use:   "walk",
	Short: "Walk through the root url",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if walkCmdArg.Url == "" {
			return errors.New("url can't be null")
		}

 		if walkCmdArg.Dest == "" {
			return errors.New("url can't be null")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		url := walkCmdArg.Url
		outputPath := walkCmdArg.Dest

		// total concurrent jobs
		totalJobs := 2
		currJobs := 0

		// storage for urls
		pendingUrls := make([]string, 0)
		workingUrls := make(map[string]bool)

		// notify channels
		chNewUrl := make(chan string)
		chNotifyEnd := make(chan string)
		chTick := make(chan bool)

		// avoid starvation from pendingUrls processing
		go func() {
			ticker := time.NewTicker(500 * time.Millisecond)

			for {
				time.Sleep(1600 * time.Millisecond)
				chTick <- true

				if !signals.ShouldContinue() {
					break
				}
			}
			ticker.Stop()
		}()

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
					// finding the next unprocessed url
					nextUrl := getNextUrlToWalk(&pendingUrls, &workingUrls)

					if nextUrl != "" {
						// updating the number of jobs
						currJobs++

						// properly starting the job
						go walk(chNewUrl, chNotifyEnd, nextUrl, outputPath)
					}
				}
			}

			// listening for events
			select {
			case url := <-chNewUrl:
				// enqueuing new url for work if it is not
				if validation.IsChildrenUrl(walkCmdArg.Url, url) {
					if _, ok := workingUrls[url]; !ok {
						pendingUrls = append(pendingUrls, url)
					}
				}
			case url := <-chNotifyEnd:
				// updating our structure with workload end
				if _, ok := workingUrls[url]; !ok {
					workingUrls[url] = true
				}

				currJobs--
			case <-chTick:
					continue
			}
		}

		return nil
	},
}

func walk(chNewUrl chan string, chNotifyEnd chan string, url string, outputPath string) {
	fmt.Printf("downloading: %v\n", url)

	defer func() {
		chNotifyEnd <- url
	}()

	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		fmt.Println("ERROR: Failed to crawl:", url)
		return
	}

	b := resp.Body
	doc, _ := ioutil.ReadAll(resp.Body)
	defer b.Close() // close Body when the function completes

	anchors := anchors.GetWalkValidPages(string(doc))
	for _, anchor := range anchors {
		chNewUrl <- anchor
	}
}

func getNextUrlToWalk(pendingUrls *[]string, workingUrls *map[string]bool) string {
	nextUrl := ""

	for len(*pendingUrls) > 0 && nextUrl == "" {
		// we are getting the first element from the pending pool, and adding it to working
		tryUrl := (*pendingUrls)[0]

		// removing the url we just got from the pending work
		*pendingUrls = (*pendingUrls)[1:]

		// we should skip this job if it is already being processed
		if _, ok := (*workingUrls)[tryUrl]; !ok {
			(*workingUrls)[tryUrl] = false
			nextUrl = tryUrl
		}
	}

	return nextUrl
}

func init() {
	walkCmd.Flags().StringVar(&walkCmdArg.Url,"url", "", "the url to be crawled")
	walkCmd.MarkFlagRequired("url")
	walkCmd.Flags().StringVar(&walkCmdArg.Dest, "dest", "", "the destination path to store the data")
	walkCmd.MarkFlagRequired("dest")
}
