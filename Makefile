bootstrap:
	-@git remote add dokku dokku@dokku-dash.local:dokku-dash
	@vagrant ssh -c "dokku apps:create dokku-dash"
	@vagrant ssh -c "dokku config:set dokku-dash DOKKU_API_SOCKET='unix:///tmp/dokku-api/api.sock/'"
	@vagrant ssh -c "dokku docker-options:add dokku-dash deploy '-v /tmp/dokku-api:/app/dokku-api'"

deploy:
	@git push dokku master
