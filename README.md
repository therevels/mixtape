mixtape
=======

Local development
-----------------

1.	Clone the repo

	```console
	git clone git@github.com:mixtape.git $GOPATH/src/github.com/therevels/mixtape
	```

2.	Build the docker image locally

	```console
	cd $GOPATH/src/github.com/therevels/mixtape
	docker build -t mixtape .
	```

3.	Run the docker image locally

	```console
	docker run -it -p 8088:8088 mixtape
	```

4.	Visit http://localhost:8088 in a browser
