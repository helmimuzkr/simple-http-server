test-server: 
		go clean --testcache
		go test ./server/. -cover -v
