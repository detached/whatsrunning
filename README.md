# Whatsrunning?

Show the deployed versions of your projects

Push new version information with

```
curl -X PUT 127.0.0.1:8080/api/project/MyGreatProject/stage/dev -d '{ "version": "1.0.5-SNAPSHOT"}'
```
