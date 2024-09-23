# ML Analysis Provider

Service that stores and provides analyzes of queries to LLM and LLM's answers with using ELK stack for advanced analysis capabilities.

## Install

Can be installed via docker-compose with ENV=local environment for testing or deployed to kubernetes with configured ELK SSL and CA auth using Helm and Makefile commands.  

## Usage

Manage analysis of LLM using controller API according to OpenAPI documentation. Visit Kibana page to check stored analyzes and work with data visualisation.

## Stack

 - Go
 - Docker
 - Kubernetes
 - ElasticSearch
 - Logstash
 - Kibana
 - NATS
 - Postgres


