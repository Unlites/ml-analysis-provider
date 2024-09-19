install_charts:
	helm install postgres-release oci://registry-1.docker.io/bitnamicharts/postgresql -f postgres/values.yml && \
	helm repo add elastic https://helm.elastic.co && \
	helm install elasticsearch-release elastic/elasticsearch -f elk/elastic_values.yml
