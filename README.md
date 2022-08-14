# Rate Limiter

this api is a rate limit for handling number of requests with differents ways.

i use PostgreSQL as database.

### Setup DataBase:
```
1. go to db package.
2. in GetConnection function put your database data.
```

### types of limiters :

#### 1. ByIp
```
this function limit number of requests base on IP Adddress.
```

#### 2. ByAppKey
```
this function limit number of requests base on content of header in X-App-Key.
```

#### 3. combinational
```
this limit number of requests by use of IP and content of header in X-App-Key.
```