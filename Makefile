format:
	@for s in $(shell git status -s |grep -e '.go$$'|grep -v '^D ' |awk '{print $$NF}'); do \
		echo "go fmt $$s";\
		go fmt $$s;\
	done

v:
	GOPROXY=https://goproxy.io go mod vendor


clean:
	@git clean -f -d -X
