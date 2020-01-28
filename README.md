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


## Authentication issue 

If Captcha has been triggered you cannot use Jira's REST API for the particular user
When Captcha is a trigger your will receive a `X-Seraph-LoginReason: AUTHENTICATION_DENIED`

For example in html format:

```html
<p>Basic Authentication Failure - Reason : AUTHENTICATION_DENIED</p>
```

To fix this you can turn off CATPCHA: 
1. Login in as the jira administrator > settings (COG icon) > system > edit settings
2. In "Maximum Authentication Attempts Allowed" leave blank

## Configuring Jira for HTTP over TLS

```bash
docker ps
docker exec -it containerid bash
cd conf
keytool -genkeypair -alias tomcat -keystore $PWD/.keystore -keyalg RSA
> password: changeit # to match tomcat's default keypass
> first name & last name: localhost
```

Add the following attributes to the connector located in `conf/web.xml`:

- SSLEnabled="true"
- scheme="https"
- secure="true"
- keystoreFile="location of .keystore"

```xml
<Connector port="8080" relaxedPathChars="[]|" relaxedQueryChars="[]|{}^&#x5c;&#x60;&quot;&lt;&gt;"
    maxThreads="150" minSpareThreads="25" connectionTimeout="20000" enableLookups="false"
    maxHttpHeaderSize="8192" protocol="HTTP/1.1" useBodyEncodingForURI="true" redirectPort="8443"
    acceptCount="100" disableUploadTimeout="true" bindOnInit="false"
    SSLEnabled="true" scheme="https" secure="true"
    keystoreFile="/opt/jira/conf/.keystore"/>
```

Force your web application to work with SSL (i.e. no http)

```xml
<security-constraint>
    <web-resource-collection>
        <web-resource-name>all-except-attachments</web-resource-name>
        <url-pattern>*.jsp</url-pattern>
        <url-pattern>*.jspa</url-pattern>
        <url-pattern>/browse/*</url-pattern>
        <url-pattern>/issues/*</url-pattern>
    </web-resource-collection>
    <user-data-constraint>
        <transport-guarantee>CONFIDENTIAL</transport-guarantee>
    </user-data-constraint>
</security-constraint>
```

Restart the server:

* Stop the containers
* `docker-compose up`