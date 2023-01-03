package KafkaMq

func DefaultConsumerConfig() *ConsumerConfig {
	c := &ConsumerConfig{
		NetConfig: NetConfig{
			ConnectTimeoutMS:    30000,
			TimeoutMS:           30000,
			TimeoutMSForEachAPI: make([]int, 0),
			KeepAliveMS:         7200000,
		},
		ClientID:             "ngo2",
		GroupID:              "",
		SessionTimeoutMS:     30000,
		RetryBackOffMS:       100,
		MetadataMaxAgeMS:     300000,
		FetchMaxWaitMS:       500,
		FetchMaxBytes:        10 * 1024 * 1024,
		FetchMinBytes:        1,
		FromBeginning:        false,
		AutoCommit:           true,
		AutoCommitIntervalMS: 5000,
		OffsetsStorage:       1,
		Version:              "0.9.0.0",
	}

	return c
}
