module search-api

go 1.24.3

require (
	github.com/bradfitz/gomemcache v0.0.0-20250403215159-8d39553ac7cf
	github.com/karlseguin/ccache/v3 v3.0.1
	github.com/rabbitmq/amqp091-go v1.10.0
)

replace github.com/karlseguin/ccache/v3 => ./internal/ccache
