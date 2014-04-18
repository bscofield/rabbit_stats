# rabbit_stats

A small Go executable to use RabbitMQ's HTTP API to collect a few stats on a specified queue and send them to Librato.

## Environment variables

Most of the configuration in the script is handled by environment variables. You'll need to set the following to use it:

* LIBRATO_EMAIL - email address of the Librato account to use
* LIBRATO_KEY - API token for Librato access
* LIBRATO_SOURCE - source string to use when sending stats to Librato
* RABBIT_DOMAIN - Domain (including protocol) that hosts RabbitMQ
* RABBIT_VHOST - RabbitMQ vhost to check
* RABBIT_QUEUE - RabbitMQ queue to check
* RABBIT_USER - username for RabbitMQ access
* RABBIT_PASSWORD - password for RabbitMQ access
