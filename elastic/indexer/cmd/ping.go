package cmd

import (
	"context"
	"strconv"

	"github.com/olivere/elastic/v7"
	"github.com/spf13/cobra"
)

func ping(cmd *cobra.Command, args []string) error {
	url := config.ElasticSearch.Scheme + config.ElasticSearch.HostName + ":" + strconv.Itoa(config.ElasticSearch.Port)
	client, err := elastic.NewClient(elastic.SetURL(url))

	if err != nil {
		return err
	}
	ctx := context.Background()
	_, code, err := client.Ping(url).Do(ctx)
	if err != nil {
		return err
	}
	cmd.Printf("response with code %d", code)

	return nil
}
