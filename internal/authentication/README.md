# Authentication

This module contains the shared authentication helpers used across Titanic services:

- [Password verification](#password-verification)
- [Website sessions](#website-sessions)
- [API access & refresh tokens](#api-tokens)
- [CSRF tokens](#csrf-tokens)
- [Authorization parsing](#authorization-parsing)

## Usage

### Website sessions

The website uses session IDs stored in HTTP-only cookies for authentication. The session store is backed by Redis & handles expiration and cleanup of sessions.  
Below are some examples of how to implement the website login, user resolution & logout flows using the session store:

```go
// WebsiteLogin uses the website session store for frontend authentication
func WebsiteLogin(app *state.State, user *schemas.User, w http.ResponseWriter, r *http.Request) error {
	store := authentication.NewWebsiteSessionStore(app.Redis)
	expiry := 30 * 24 * time.Hour

	session, err := store.Create(
		r.Context(),
		user.Id,
		time.Now(),
		expiry,
	)
	if err != nil {
		return err
	}

	cookie := authentication.NewWebsiteSessionCookie(
		app.Config,
		r,
		session.Id,
		expiry,
	)
	http.SetCookie(w, cookie)
	return nil
}
```

```go
// ResolveWebsiteUser validates the website session cookie on every request & returns the associated user
func ResolveWebsiteUser(app *state.State, r *http.Request) (*schemas.User, error) {
	cookie, err := r.Cookie(authentication.WebsiteSessionCookieName)
	if err != nil {
		return nil, err
	}

	store := authentication.NewWebsiteSessionStore(app.Redis)
	session, err := store.Validate(r.Context(), cookie.Value, time.Now())
	if err != nil || session == nil {
		return nil, err
	}

	return app.Users.ById(session.UserId)
}
```

```go
// WebsiteLogout deletes the website session & expires the cookie
func WebsiteLogout(app *state.State, w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(authentication.WebsiteSessionCookieName)
	if err == nil {
		store := authentication.NewWebsiteSessionStore(app.Redis)
		store.Delete(r.Context(), cookie.Value)
	}

	expiredCookie := authentication.NewExpiredCookie(authentication.WebsiteSessionCookieName, app.Config)
	http.SetCookie(w, expiredCookie)
	return nil
}
```

### API tokens

The API uses jwt access & refresh tokens for authentication. The access token is short-lived and used for authenticating API requests, while the refresh token is long-lived and used to obtain new access tokens. The refresh tokens are also stored in redis to allow for revocation.  
Below are some examples of how to implement the API login, token validation & refresh flows:

```go
func ApiLogin(app *state.State, user *schemas.User) (*authentication.TokenPair, error) {
	expiryAccessToken := 15 * time.Minute
	expiryRefreshToken := 30 * 24 * time.Hour

	pair, err := authentication.GenerateTokenPair(
		app.Config.FrontendSecretKey,
		user,
		time.Now(),
		expiryAccessToken,
		expiryRefreshToken,
		authentication.TokenSourceAPI,
	)
	if err != nil {
		return nil, err
	}

	store := authentication.NewSessionStore(app.Redis)
	claims, err := authentication.ValidateTokenType(
		pair.RefreshToken,
		app.Config.FrontendSecretKey,
		authentication.TokenTypeRefresh,
	)
	if err != nil {
		return nil, err
	}

	// Register refresh token in redis store to allow for revocation
	_, err = store.Upsert(context.Background(), claims, time.Now())
	if err != nil {
		return nil, err
	}

	return pair, nil
}
```

```go
// Example of how to validate the access token from the Authorization header of an API request
claims, err := authentication.ValidateTokenType(
	token,
	app.Config.FrontendSecretKey,
	authentication.TokenTypeAccess,
)
```

```go
// Example of how to validate the refresh token from the Authorization header of a token refresh request
claims, err := authentication.ValidateTokenType(
	refreshToken,
	app.Config.FrontendSecretKey,
	authentication.TokenTypeRefresh,
)
if err != nil {
	return err
}

store := authentication.NewSessionStore(app.Redis)
ok, err := store.Validate(context.Background(), claims, time.Now())
if err != nil || !ok {
	return err
}
```

### Password verification

This should be self-explanatory.  
osu! clients use md5 hashes to log in, so there are 2 methods of password verification:

```go
if !authentication.VerifyPasswordHash(password, user.Bcrypt) {
	return errors.New("invalid credentials")
}
```

```go
if !authentication.VerifyPasswordHashFromMd5(md5Password, user.Bcrypt) {
	return errors.New("invalid credentials")
}
```

### CSRF tokens

The frontend uses rotating CSRF tokens stored in redis to protect against CSRF attacks. The tokens are associated with a user ID and have a short expiration time.  
Here is an example of how to generate & validate CSRF tokens:

```go
store := authentication.NewCSRFStore(app.Redis)
token, err := store.Upsert(context.Background(), user.Id)
if err != nil {
	return err
}

ok, err := store.Validate(context.Background(), user.Id, token)
if err != nil || !ok {
	return err
}
```

### Authorization parsing

The API uses the `Authorization: <scheme> <data>` header format to pass the access token for authentication, where the scheme is typically "Bearer" for API requests or "Basic" for login requests. The `ParseAuthorization` function can be used to parse the header and extract the scheme and data.

```go
authorization := authentication.ParseAuthorization(r.Header.Get("Authorization"))

if authorization.Scheme == "basic" {
	username, password, err := authentication.ParseBasicCredentials(authorization.Data)
	if err != nil {
		return err
	}

	_ = username
	_ = password
}
```
