# jira
default port: 8080

The default port is set in docker-compose.yml 8080:8080, which refers to the host:container mapping.
To change the port, change the port value of the host.

## requirement
- docker
- docker-compose

## How to run with docker-compose

- start jira & mysql

```
    docker-compose up
```

- start jira & mysql in daemon mode

```
    docker-compose up -d
```

- default db(mysql) configured to be:

```bash
    host=mysql-jira
    port=3306
    db=jira
    user=jira
    passwd=password
```

## Connecting jira to database during setup

After running with docker compose, access the jira setup through localhost:8080

Choose to set up jira yourself:

* DB type, Database name, user, and password can be found in docker-compose.yml
* Host can be found by looking at the IP of the database service. Find the network connecting jira-server and the database by typing `docker network ls` and look for a network called `*_network-bridge`. Then use the network ID or the full network name in `docker network inspect <network>` and locate the database container name (which is defined in docker-compose.yml). After locating the container you will be able to see the IPV4 addr.

### Example

```txt
Database connection: My Own database
Database type: MySQL 5.7+
HostName: Look up docker jira network
port: 3306
Database: jira
Username: jira
Password: password
```

## Setting up OAuth for REST APIs

1. run `generate_keys.sh`
2. Go to Jira settings > Products > Application links > Create a new link
3. Application Name: app, Application Type: Generic Application, Tick create incoming link, and press Continue
4. Consumer key: This will be used in the Go client, Consumer Name: jira app, Public key: copy from `jira_publickey.pem`, and Continue to finish Application link creation