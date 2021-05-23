# link-tracker

## Build locally
After cloning the repository, execute the following command under the root directory

```shell
make run
```

## Create a link

`curl -POST http://localhost:8080/link -d '{"link":"http://www.google.com", "password":"123"}'`

## Open a link
You can either use cURL and follow the redirection with -L or opening a browser and navigate to the link. For a link with id 1 please visit: http://localhost:8080/link/1?password=123
