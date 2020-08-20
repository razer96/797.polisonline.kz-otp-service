package elastic

import (
	"context"
	"encoding/json"
	"errors"
	"insurance-otp-service/logger"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type elasticClient struct {
	client *elasticsearch.Client
}

var (
	es = &elasticClient{}
	elasticHost = os.Getenv("ELASTIC_HOST")
	log       = logger.GetLogger()
	err error
)

func init() {
	if elasticHost == "" {
		log.Error("Error getting elastic host variable")
		os.Exit(1)
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			elasticHost,
		},
	}
	
	es.client, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Error("Error connecting to ElasticSearch", err)
		os.Exit(1)
	}
}

func GetESClient() *elasticClient {
	return es
}

func (e elasticClient) Insert(index string, docID string, body interface{}) error {
	
	log.Info("Insertion into ES started")
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Error("Error marshaling", err)
		return errors.New("Error marshaling body" + err.Error())
	}
	req := esapi.IndexRequest{
		Index: index,
		DocumentID: docID,
		Body: strings.NewReader(string(jsonBody)),
	}

	res, err := req.Do(context.Background(), e.client)

	if err != nil {
		log.Error("Insertion error", err)
		return errors.New("Error getting response: " + err.Error())
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Error("Insertion into ES error", res.Status(), res.String())
		return errors.New("[" + res.Status() + "] Error indexing document ID="+docID)
	}
	log.Info("Insertion into ES completed")
	return nil
}