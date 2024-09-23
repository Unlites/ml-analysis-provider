build_prod_images:
	docker build -t ml-analysis-provider_controller:prod controller/ --build-arg=ENVIRONMENT=prod
	docker build -t ml-analysis-provider_worker:prod worker/ --build-arg=ENVIRONMENT=prod

install_charts:
	helm install postgres-release oci://registry-1.docker.io/bitnamicharts/postgresql -f postgres/values.yml && \
	helm repo add elastic https://helm.elastic.co && \
	helm install elasticsearch-release elastic/elasticsearch -f elk/elastic_values.yml && \
	helm install logstash-release elastic/logstash -f elk/logstash_values.yml && \
	helm install kibana-release elk/kibana -f elk/kibana_values.yml && \
	helm install controller-release controller/controller-chart && \
	helm install worker-release wprker/worker-chart
