package KafkaMq

type Msg struct {
	Topic string
	Msg   string
}

type Config struct {
	Topics                 []string          `json:"topics"`
	Topic                  map[string]string `json:"topic"`
	Host                   []string          `json:"host"`
	Group                  string            `json:"group"`
	IsDebug                bool              `json:"is_debug"`
	AutoCommit             bool              `json:"auto_commit"`
	AutoCommitIntervalMS   int64             `json:"auto_commit_interval_ms"`
	SessionTimeoutMs       int64             `json:"session_timeout_ms"`
	FromBeginning          bool              `json:"from_beginning"`
	Version                string            `json:"version"`
	MessagesQueueLength    string            `json:"messages_queue_length"`
	MaxPartitionFetchBytes int64             `json:"max_partition_fetch_bytes"`
}

type ConsumerI interface {
	Consume(msg Msg)
}

type SaslConfig struct {
	SaslMechanism string `json:"sasl.mechanism"`
	SaslUser      string `json:"sasl.user"`
	SaslPassword  string `json:"sasl.password"`
}

type NetConfig struct {
	ConnectTimeoutMS    int   `json:"connect.timeout.ms,string"`
	TimeoutMS           int   `json:"timeout.ms,string"`
	TimeoutMSForEachAPI []int `json:"timeout.ms.for.eachapi"`
	KeepAliveMS         int   `json:"keepalive.ms,string"`
}

type TLSConfig struct {
	Cert               string `json:"cert"`
	Key                string `json:"key"`
	CA                 string `json:"ca"`
	InsecureSkipVerify bool   `json:"insecure.skip.verify,string"`
	ServerName         string `json:"servername"`
}

type ConsumerConfig struct {
	NetConfig
	*SaslConfig
	BootstrapServers     string `json:"bootstrap.servers"`
	ClientID             string `json:"client.id"`
	GroupID              string `json:"group.id"`
	RetryBackOffMS       int    `json:"retry.backoff.ms,string"`
	MetadataMaxAgeMS     int    `json:"metadata.max.age.ms,string"`
	SessionTimeoutMS     int32  `json:"session.timeout.ms,string"`
	FetchMaxWaitMS       int32  `json:"fetch.max.wait.ms,string"`
	FetchMaxBytes        int32  `json:"fetch.max.bytes,string"`
	FetchMinBytes        int32  `json:"fetch.min.bytes,string"`
	FromBeginning        bool   `json:"from.beginning,string"`
	AutoCommit           bool   `json:"auto.commit,string"`
	AutoCommitIntervalMS int    `json:"auto.commit.interval.ms,string"`
	OffsetsStorage       int    `json:"offsets.storage,string"`
	Version              string `json:"version"`

	MetadataRefreshIntervalMS int `json:"metadata.refresh.interval.ms,string"`

	TLSEnabled bool       `json:"tls.enabled,string"`
	TLS        *TLSConfig `json:"tls"`
}
