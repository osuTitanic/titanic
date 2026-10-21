# Authentication

This module contains the shared authentication helpers used across Titanic services:

- [Password verification](#password-verification)
- [Website sessions](#website-sessions)
- [API access & refresh tokens](#api-tokens)
- [CSRF tokens](#csrf-tokens)
- [Authorization parsing](#authorization-parsing)
- [Permission scopes](#permission-scopes)

## Usage

### Website sessions

The website uses session IDs stored in HTTP-only cookies for authentication. The session store is backed by redis and handles session expiration & cleanup.

The following examples show the login, user resolution, and logout flow:

```go
// WebsiteLogin uses the website session store for frontend authentication
func WebsiteLogin(app *state.State, user *schemas.User, w http.ResponseWriter, r *http.Request) error {
	expiry := 30 * 24 * time.Hour

	session, err := app.SessionStore.Create(
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
// ResolveWebsiteUser validates the website session cookie & returns its user
func ResolveWebsiteUser(app *state.State, r *http.Request) (*schemas.User, error) {
	cookie, err := r.Cookie(authentication.WebsiteSessionCookieName)
	if err != nil {
		return nil, err
	}

	session, err := app.SessionStore.Validate(r.Context(), cookie.Value, time.Now())
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
		if err := app.SessionStore.Delete(r.Context(), cookie.Value); err != nil {
			return err
		}
	}

	expiredCookie := authentication.NewExpiredWebsiteSessionCookie(app.Config, r)
	http.SetCookie(w, expiredCookie)
	return nil
}
```

### API tokens

The API uses jwt access & refresh tokens. Access tokens are short-lived and authenticate API requests. Refresh tokens are long-lived, obtain new access tokens, and are stored in redis so they can be revoked.

The following examples show token creation & validation:

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

	claims, err := authentication.ValidateTokenType(
		pair.RefreshToken,
		app.Config.FrontendSecretKey,
		authentication.TokenTypeRefresh,
	)
	if err != nil {
		return nil, err
	}

	// Register the refresh token in redis so it can be revoked
	_, err = app.SessionStoreApi.Upsert(context.Background(), claims, time.Now())
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

ok, err := app.SessionStoreApi.Validate(context.Background(), claims, time.Now())
if err != nil || !ok {
	return err
}
```

### Password verification

Website logins use the plain-text password supplied by the user:

```go
if !authentication.VerifyPasswordHash(password, user.Bcrypt) {
	return errors.New("invalid credentials")
}
```

osu! clients submit an md5 password hash, so client authentication uses this helper:

```go
if !authentication.VerifyPasswordHashFromMd5(md5Password, user.Bcrypt) {
	return errors.New("invalid credentials")
}
```

### CSRF tokens

The frontend uses rotating CSRF tokens stored in redis. Each token belongs to a user and expires after the store's configured TTL.

Here is an example of how to generate & validate CSRF tokens:

```go
token, err := app.CSRFStore.Upsert(context.Background(), user.Id)
if err != nil {
	return err
}

ok, err := app.CSRFStore.Validate(context.Background(), user.Id, token)
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
