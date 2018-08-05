mixtape
=======

Configuration
-------------

1.	You must register your application (including redirect URI) with Spotify in your [developer dashboard](https://developer.spotify.com/dashboard/applications).

2.	Configure the following environment variables:

	-	`SPOTIFY_ID` - the client ID
	-	`SPOTIFY_SECRET` - the client secret
	-	`REDIRECT_URI` - the URI for Spotify to redirect back to after authentication and authorization

Local development
-----------------

1.	Clone the repo:

	```console
	git clone git@github.com:mixtape.git $GOPATH/src/github.com/therevels/mixtape
	```

2.	Set up a `.env` file with the environment variables described [configuration instructions](#configuration):

	```console
	cp .env.example .env
	```

3.	Use docker compose to start the app in a container:

	```console
	docker-compose up
	```

4.	Visit http://localhost:8088 in a browser
