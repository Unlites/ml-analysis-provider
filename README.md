# ML Analysis Provider

Service that stores and provides analyzes of queries to LLM and LLM's answers with using ELK stack for advanced analysis capabilities.

## Install

Can be installed via docker-compose with ENV=local environment for testing or deployed to kubernetes with configured ELK SSL and CA auth using Helm and Makefile commands.  

## Usage

Manage analysis of LLM using controller API according to OpenAPI documentation. Visit Kibana page to check stored analyzes and work with data visualisation.

## Workflow

Get data from Provider
```sh
User -> controller -> NATS Queue -> worker -> request ID (or IDs) from ElasticSearch
                                           -> request data from PostgreSQL with given ID (or IDs) from ElasticSearch
```

Store data to Provider
```sh
User -> controller -> NATS Queue -> worker -> Postgres <- Logstash pulls data from PG -> ElasticSearch
```

View and analysing
```sh
User -> Kibana -> ElasticSearch
```

## Stack

 - Go
 - Docker
 - Kubernetes
 - OpenAPI (server code generation with oapi-codegen)
 - ElasticSearch
 - Logstash
 - Kibana
 - NATS
 - Postgres


