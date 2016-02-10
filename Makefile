bootstrap:
	-@git remote add dokku dokku@dokku-dash.local:dokku-dash
	@vagrant ssh -c "dokku apps:create dokku-dash"

deploy:
	@git push dokku master
