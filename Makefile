bootstrap:
	-@git remote add dokku dokku@dokku-dash-vm.local:dokku-dash
	@vagrant ssh -c "dokku apps:create dokku-dash"
	@vagrant ssh -c "dokku config:set dokku-dash DOKKU_API_SOCKET='unix:///app/dokku-api/api.sock'"
	@vagrant ssh -c "dokku docker-options:add dokku-dash deploy '-v /tmp/dokku-api:/app/dokku-api'"

deploy:
	@git push dokku master

updateapi:
	@vagrant ssh -c 'cd $$GOPATH/src/github.com/nikelmwann/dokku-api/dokku-api && go get && go install && sudo supervisorctl restart dokku-api'
