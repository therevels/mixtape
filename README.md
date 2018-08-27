mixtape
=======

Configuration
-------------

1.	You must register your application (including redirect URI) with Spotify in your [developer dashboard](https://developer.spotify.com/dashboard/applications).

2.	Configure the following environment variables:

	-	`SPOTIFY_ID` - the client ID
	-	`SPOTIFY_SECRET` - the client secret

Local development
------------------

1.	Clone the repo:

	```console
	git clone git@github.com:mixtape.git $GOPATH/src/github.com/therevels/mixtape
	```

2.	Copy the example `.env` file and edit the environment variables (see [configuration instructions](#configuration)\):

	```console
	cp .env.example .env
	vim .env
	```

3.	Build the image locally:

	```console
	docker-compose build
	```

4.	Start the app containers locally:

	```console
	docker-compose up
	```

5.	Visit https://localhost:8088 in a browser
