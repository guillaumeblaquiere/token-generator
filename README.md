# token-generator
Allow to generate id_token based on service account

# Run the server
Run the token generator with these params:

 * **port** for changing the port. Optional, default is 8080
 * **file** for setting the service account json secret file url path. If missing use this one 
 configured in `GOOGLE_APPLICATION_CREDENTIALS` ev var
 
 Example
```
 token-generator -port 8081
```

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
