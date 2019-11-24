package cmd

import (
	"context"
	"fmt"
	"msc/elastic/datafaker"
	"strconv"

	"github.com/olivere/elastic/v7"
	"github.com/spf13/cobra"
)

func fakeData(c *cobra.Command, args []string) error {
	i := 10
	fakeDatas, err := datafaker.FakeDatas(i)
	if err != nil {
		return err
	}

	fmt.Println(fakeDatas)
	return nil
}
func index(c *cobra.Command, args []string) error {
	url := config.ElasticSearch.Scheme + config.ElasticSearch.HostName + ":" + strconv.Itoa(config.ElasticSearch.Port)
	client, err := elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		return err
	}

	number := config.Indexing.ItemNumbs

	fakedDatas, err := datafaker.FakeDatas(number)
	if err != nil {
		return err
	}
	ctx := context.Background()
	_, _, err = client.Ping(url).Do(ctx)
	if err != nil {
		return err
	}
	indexName := config.Indexing.IndexName
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}
	if !exists {
		_, err := client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
		if err != nil {
			return err
		}

	}
	indicer := client.Index().Index(indexName)
	for _, item := range fakedDatas {
		_, err := indicer.BodyJson(item).Do(ctx)
		if err != nil {
			return err
		}
	}
	return nil

}
