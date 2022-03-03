MODULES = user_service article_service

db_reset:
	for module in $(MODULES); do $(MAKE) -C $$module db_reset; done

vendor:
	for module in $(MODULES); do cd $$module && go mod vendor && cd -; done

test:
	for module in $(MODULES); do cd $$module && go test -cover -tags integration ./... && cd -; done
