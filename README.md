Portfello
=========

Portfello is an open source solution for aggregating AIS' (Account Information Services) into one clearly readable view.

What's it for?
--------------

If your family has a number of accounts in different banks, credit cards, adolescent kids, and planning and tracking
your budgets becomes a hassle then you can use this app to easily plan and track your expenses. 

Included WebUI and mobile App allows users to check:

- balance on their different accounts
- balance on their spouses credit cards
- funds on their children prepaid debit cards
- how many prepaid credits there are left on their phones
- share their expenses history with spouses
- plan and observe home budgets!

How it works?
-------------

Portfello Server can be either self-hosted - to maximize your security, or used as a service @ www.portfel.lo

After setting up self-hosted server/signing up at www.portel.lo, you are able to create and share as many Wallets 
as you wish.  Each Wallet can represent any of available sync sources, either bank accounts, prepaid cards (supported 
ones only) or prepaid phones. To set up synchronisation of your Wallet with Banking Provider you have to provide an 
API key for accessing that provider, and an encryption key to securely store your API key. 

After a Wallet has started synchronising your WebUI and App will show current balance and any changes made to it,
either expenses or deposits, and update the balance accordingly. Depending on Banking Provider there may be a slight 
delay between making of a payment and registering the change in Portfello.

A Wallet which is successfully tracking expenses can then be shared with other users of the same Portfello server,
but only the Wallet Owner can make changes to it, such as:

- share with other users
- remove a user from a shared Wallet
- delete a Wallet.

How secure is it?
-----------------

An Account Information Services (ie. your Banking Provider) API Access Key (a secret used to limit access to your 
Bank's data) is stored encrypted using an encryption key which only you have access to. Don't worry if you have somehow
lost your API key - your Wallet(s) will become inactive and you can delete them at any time or, if you won't, they will
be removed automatically after 7 days without synchronisation. Deleting a Porftello Wallet does not alter your Bank in
any way, it only serves as a mirror of your balance and has no impact on yours or anybody else's Real Money.

Your balance and expenses history information that's kept locally on Server is encrypted using the same key and as
such can only be accessed by you or those with whom you share your Portfello Wallet - nobody else. Access to your
account is secured by systems of award-winning Auth0 provided free of charge by Okta.com. Portfello Server and App(s)
don't care, don't need and don't bother with your password but it Better Be Good anyway as obtaining access to your
Porftello account will grant interlopers a clear view into what your spending is.

How to set up self-hosted server?
---------------------------------

This requires a few systems:

- a computer or VM or Kubernetes cluster to run the server on
- a valid Domain Name
- Auth0 account for authentication
- PostgresQL database for storage
- Redis server (optional) for caching
- Prometheus server (optional) for metrics
- Graphana (optional) for presentation of Prometheus metrics

If you have access to a Kubernetes instance you can deploy these services (except for Auth0) using provided helm charts,
otherwise it's up to you how and where they are running and how secure they are. 

After securing access to the services listed above you will need to configure each one. Clone this repository, locate 
`.portfello.yaml` file in the root of the repo folder, and fill out each line while following next steps.

Configure Authentication
------------------------

The `auth` section in configuration file looks more or less like so:

```yaml
auth:
  domain: "YOUR_DOMAIN.auth0.com"
  client_id: "YOUR_CLIENT_ID"
  client_secret: "YOUR_CLIENT_SECRET"
  audience: "YOUR_AUDIENCE"
  connection_id: "con_CONNECTION_ID"
```

Set up an account on [Auth0](manage.auth0.com), then create:

- an API, call it "Portfello API"
- an API Audience, this is up to you
- an App, call it "Portfello Server", make it Machin-to-Machine type
- another App, called "Portfello UI", make it Web-App type
- assign these apps permission to access this API
- create user roles, see Authorization section for a full list
- create additional connections (eg. Google, Facebook, etc) in addition to the default Database connection

Then find the API Domain, Client ID and Client Secret of "Portfello Server" App and paste them into respective fields.
Copy the "Portfello UI" Client ID and provide it the first time you use Portfello Mobile App.

Verify configuration by running `go run main.go config` - this should check if minimum required configuration is good.

Configure Database
------------------

Find the connection details of your database server (hostname, username, database name, port number and password)
and put them into `database_dsn` field in `.portfello.yaml`. For testing you can use the default value and use
`docker-compose up db` to start the database server. 

After configuring `database_dsn` run `go run main.go migrate up` to create tables in a fresh database. Running this 
command repeatedly is safe. To clear the database and restart from scratch run `go run main.go migrate drop` and 
`go run main.go migrate up` again. For testing only. The `migrate` command will complain if `database_dsn` is invalid.

Configure API
-------------

To configure the API server go to the `graphql` section, it should look like so:

```yaml
graphql:
  host: "http://localhost"
  port: "8080"
  enable_playground: true
```

Set `port` to whatever port number your server will listen on and, unless testing locally, set `enable_playground` to
`false`. The `host` field should match the domain at which your Server will be available, eg.: `https://your.domain`.
If you set scheme to `https` then on first run Portfello will request a LetsEncrypt certificate for this domain.
If you set scheme to plain `http` then it's up to you to terminate the connection with a secure link.

After starting the server (eg. using a container in cloud or by running locally with `go run main.go serve` ) you 
should be able to visit the host+port combination in your browser and see the homepage at https://your.domain:8080 , 
the graph playground at https://your.domain/playground and the login page at https://your.domain/login - please visit
each to make sure configuration so far is correct.

Configure optional services
---------------------------

TBD.

First Run
=========

Migrate your database and start the API server if you haven't already:

```bash
$ go run main.go migrate up
$ go run main.go serve
```

After starting your server verifying that basic screens work fine it is time to set up the first account. That cannot 
be done through web interface and you have to use CLI again:

```bash
$ go run main.go provision --superuser your@email.com 
```

This will create an account for managing this Portfello Server installation. It will also be the only account able to
create new users initially - You are free to enable registration in Auth0 but required configuration is not documented
yet.

If the "your@email.com" user was successfully created on Auth0 then in your terminal you should see a message showing
your system-wide cryptographic public key - copy it and save it in configuration file. It is part of the data encryption 
mechanism and if you happen to lose it then every user of your installation will lose their personal data. 

After updating the config file stop the API server and start it again:

```bash
$ go run main.go serve
```

Now you should be able to log in using either the WebUI or the Mobile App into the system and add new users.

Next steps
==========

After logging in test out a Wallet by using a Banking Provider's sandbox such as this one: 
[MBank](https://developer.api.mbank.pl/documentation/sandbox-v2#section/How-to-test-requests-for-AIS). 
Generate a key and configure your first Wallet using it, then observe how the data changes in UI, all without
using your real account.

After playing around with different accounts and Wallets, stop the server, reset the database and run again through
initial setup described earlier to start with a clean state:

```bash
$ go run main.go migrate drop
$ go run main.go migrate up
$ go run main.go provision --superuser your@real.email
$ go run main.go serve
```

Metrics
-------

TBD.

Encryption
----------

TBD.
