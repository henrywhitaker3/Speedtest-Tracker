validate: lint test

test:
	@./node_modules/.bin/mocha \
		--require should \
		--reporter spec

test-cov: lib-cov
	@CONCAT_COV=1 ./node_modules/.bin/mocha \
		--require should \
		--reporter html-cov > coverage.html

lint:
	@./node_modules/.bin/jshint \
		--verbose \
		--config .jshintrc \
		lib/*.js \
		test/*.js

lib-cov:
	jscoverage lib lib-cov

.PHONY: test lint
