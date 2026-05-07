# CDN Client

This module provides a simple API wrapper for the [Titanic! CDN](https://github.com/osuTitanic/cdn).

## Usage

Create a new client with the default production streaming host (`cdn.DefaultStreamingBaseUrl`):

```go
func main() {
	client, err := cdn.NewDefaultClient(cdn.WithAccessKey("access-key"))
	if err != nil {
		panic(err)
	}

	// ...
}
```

You can also use `NewClient` to target the redirect host (`cdn.DefaultRedirectBaseUrl`) or a local CDN server:

```go
client, err := cdn.NewClient(
	"http://localhost:6969",
	cdn.WithAccessKey("access-key"),
)
if err != nil {
	panic(err)
}
```

You can download objects without authentication:

```go
object, err := client.GetObject(context.Background(), "path/to/object", nil)
if err != nil {
	return err
}
defer object.Body.Close()

_, err = io.Copy(destination, object.Body)
```

When using the redirect host / s3 passthrough, `GetObject` returns the `307` response with the `Location` header:

```go
client, err := cdn.NewClient(cdn.DefaultRedirectBaseUrl)
if err != nil {
	return err
}

object, err := client.GetObject(context.Background(), "path/to/object", nil)
if err != nil {
	return err
}
fmt.Println(object.Location())
```

You can call admin endpoints using the authenticated client:

1. Use `AdminSession` to fetch the current access key permissions
	```go
	session, err := client.AdminSession(context.Background())
	if err != nil {
		return err
	}
	fmt.Println(session.Name, session.Permissions)
	```

2. Use `AdminListFiles` to list objects under a prefix
	```go
	files, err := client.AdminListFiles(context.Background(), cdn.ListFilesOptions{
		Prefix: "path/to/prefix",
		Limit:  100,
	})
	if err != nil {
		return err
	}
	
	for _, file := range files.Items {
		fmt.Println(file.Key, file.Size)
	}
	```

3. Use `AdminUploadFile` to upload or replace an object
	```go
	upload, err := client.AdminUploadFile(
		context.Background(),
		"path/to/object",
		fileReader,
		&cdn.UploadFileOptions{
			ContentType:  "application/octet-stream",
			CacheControl: "public, max-age=31536000, immutable",
		},
	)
	if err != nil {
		return err
	}
	fmt.Println(upload.Key, upload.ETag)
	```

4. Use `AdminDeleteFile` to delete an object
	```go
	if err := client.AdminDeleteFile(context.Background(), "path/to/object"); err != nil {
		return err
	}
	```
