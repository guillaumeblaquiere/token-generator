# token-generator
Allow to generate id_token based on service account.

This tool can help for testing API when `gcloud` is not installed on a machine, with postman or curl
 for example. With curl, you can perform this
```
curl -H "Authorization: Bearer $(curl localhost:8080?aud=https://my-service)" https://my-service/...
```
 
It can be used in collocation with an application where getting a token is too complex 
(old application or out-of-date framework). For this:

* Run the token-generator server on a free port of the server
* Run your (old) application
* Simply perform a GET with an HTTP library on localhost and on the token-generator port for
getting the signed id_token. 

# Run the server
Run the token generator with these params:

 * **port** for changing the port. Optional, default is 8080
 * **file** for setting the service account json secret file url path. If missing use this one 
 configured in `GOOGLE_APPLICATION_CREDENTIALS` ev var
 
 Example
```
 token-generator -port 8081
```

## CAUTION
**Never expose this service on a public IP. Your credential can be stolen!**

# Use the server

Query to the opened port with the audience of the service to reach in the `aud` query parameter

Example
```
curl http://localhost:8080?aud=http://my-service
```

The return is the id_token signed by Google.

# License

This library is licensed under Apache 2.0. Full license text is available in
[LICENSE](https://github.com/guillaumeblaquiere/token-generator/tree/master/LICENSE).
