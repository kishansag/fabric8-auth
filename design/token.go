package design

import (
	d "github.com/goadesign/goa/design"
	a "github.com/goadesign/goa/design/apidsl"
)

// externalToken represents a token object
var externalToken = a.MediaType("application/vnd.externalToken+json", func() {
	a.TypeName("ExternalToken")
	a.Description("Tokens from external providers such as GitHub or OpenShift")
	a.Attributes(func() {
		a.Attribute("access_token", d.String, "The token associated with the identity for the specific external provider")
		a.Attribute("scope", d.String, "The scope associated with the token")
		a.Attribute("token_type", d.String, "The type of the toke, example : bearer")
		a.Attribute("username", d.String, "The username of the identity loaded from the specific external provider. Optional attribute.")
		a.Required("access_token", "scope", "token_type")
	})

	a.View("default", func() {
		a.Attribute("access_token")
		a.Attribute("scope")
		a.Attribute("token_type")
		a.Attribute("username")
		a.Required("access_token", "scope", "token_type")
	})

})

var _ = a.Resource("token", func() {

	a.BasePath("/token")

	a.Action("Retrieve", func() {
		a.Security("jwt")
		a.Routing(
			a.GET(""),
		)
		a.Params(func() {
			a.Param("for", d.String, "The resource for which the external token is being fetched, example https://github.com/fabric8-services/fabric8-auth or https://api.starter-us-east-2.openshift.com")
			a.Param("force_pull", d.Boolean, "Pull the user's details for the specific connected account, example, the user's updated github username would be fetched from github. If this is not set or false, then the user profile will be pulled only if the stored user's details did not have the username")
			a.Required("for")
		})
		a.Description("Get the external token for resources belonging to external providers like Github and OpenShift")
		a.Response(d.OK, externalToken)
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})

	a.Action("Delete", func() {
		a.Security("jwt")
		a.Routing(
			a.DELETE(""),
		)
		a.Params(func() {
			a.Param("for", d.String, "The resource for which the external token is being deleted, example https://github.com/fabric8-services/fabric8-auth or https://api.starter-us-east-2.openshift.com")
			a.Required("for")
		})
		a.Description("Delete the external token for resources belonging to external providers like Github and OpenShift")
		a.Response(d.OK)
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})

	a.Action("Exchange", func() {
		a.Routing(
			a.POST(""),
		)
		a.Payload(tokenExchange)
		a.Description("Obtain a security token")
		a.Response(d.OK, func() {
			a.Media(OauthToken)
		})
		a.Response(d.Unauthorized, JSONAPIErrors)
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
	})

	a.Action("keys", func() {
		a.Routing(
			a.GET("keys"),
		)
		a.Params(func() {
			a.Param("format", d.String, func() {
				a.Enum("pem", "jwk")
				a.Description("Key format. If set to \"jwk\" (used by default) then JSON Web Key format will be used. If \"pem\" then a PEM-like format (PEM without header and footer) will be used.")
			})
		})
		a.Description("Returns public keys which should be used to verify tokens")
		a.Response(d.OK, func() {
			a.Media(PublicKeys)
		})
	})

	a.Action("generate", func() {
		a.Routing(
			a.GET("generate"),
		)
		a.Description("Generate a set of Tokens for different Auth levels. NOT FOR PRODUCTION. Only available if server is running in dev mode")
		a.Response(d.OK, func() {
			a.Media(a.CollectionOf(AuthToken))
		})
		a.Response(d.Unauthorized, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
	})

	a.Action("refresh", func() {
		a.Routing(
			a.POST("refresh"),
		)
		a.Payload(refreshToken)
		a.Description("Refresh access token")
		a.Response(d.OK, func() {
			a.Media(AuthToken)
		})
		a.Response(d.Unauthorized, JSONAPIErrors)
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
	})

	a.Action("link", func() {
		a.Security("jwt")
		a.Routing(
			a.GET("link"),
		)
		a.Params(func() {
			a.Param("for", d.String, "Resource we need to link accounts for. Multiple resources should be separated by comma.", func() {
				a.Example("https://github.com,https://api.starter-us-east-2.openshift.com")
			})
			a.Param("redirect", d.String, "URL to be redirected to after successful account linking. If not set then will redirect to the referrer instead.")
			a.Required("for")
		})
		a.Description("Get a redirect location which should be used to initiate account linking between the user account and an external resource provider such as GitHub")
		a.Response(d.OK, func() {
			a.Media(redirectLocation)
		})
		a.Response(d.Unauthorized, JSONAPIErrors)
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
	})

	a.Action("callback", func() {
		a.Routing(
			a.GET("/link/callback"),
		)
		a.Params(func() {
			a.Param("code", d.String, "Code provided by an external oauth2 resource provider")
			a.Param("state", d.String, "State generated by the link request")
			a.Required("code", "state")
		})
		a.Description("Callback from an external oauth2 resource provider such as GitHub as part of user's account linking")
		a.Response(d.TemporaryRedirect)
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
	})
})

// PublicKeys represents an public keys payload
var PublicKeys = a.MediaType("application/vnd.publickeys+json", func() {
	a.TypeName("PublicKeys")
	a.Description("Public Keys")
	a.Attributes(func() {
		a.Attribute("keys", a.ArrayOf(d.Any))
		a.Required("keys")
	})
	a.View("default", func() {
		a.Attribute("keys")
	})
})

var refreshToken = a.Type("RefreshToken", func() {
	a.Attribute("refresh_token", d.String, "Refresh token")
})

var tokenExchange = a.Type("TokenExchange", func() {
	a.Attribute("grant_type", d.String, func() {
		a.Enum("client_credentials", "authorization_code")
		a.Description("Grant type. If set to \"client_credentials\" then this token exchange request is for a Protection API Token (PAT). PAT can be used to authenticate the corresponding Service Account. If the Grant Type is \"authorization_code\" we can use a authorization_code to get access_token")
	})
	a.Attribute("client_id", d.String, "Service Account ID. Used to obtain a PAT for this service account.")
	a.Attribute("client_secret", d.String, "Service Account secret. Used to obtain a PAT for this service account.")
	a.Attribute("redirect_uri", d.String, "Must be identical to the redirect URI provided while getting the authorization_code")
	a.Attribute("code", d.String, "this is the authorization_code you received from /api/authorize endpoint")
	a.Required("grant_type", "client_id")
})

// AuthToken represents an authentication JWT Token
var AuthToken = a.MediaType("application/vnd.authtoken+json", func() {
	a.TypeName("AuthToken")
	a.Description("JWT Token")
	a.Attributes(func() {
		a.Attribute("token", tokenData)
		a.Required("token")
	})
	a.View("default", func() {
		a.Attribute("token")
	})
})

var tokenData = a.Type("TokenData", func() {
	a.Attribute("access_token", d.String, "Access token")
	a.Attribute("expires_in", d.Any, "Access token expires in seconds")
	a.Attribute("refresh_expires_in", d.Any, "Refresh token expires in seconds")
	a.Attribute("refresh_token", d.String, "Refresh token")
	a.Attribute("token_type", d.String, "Token type")
	a.Attribute("not-before-policy", d.Any, "Token is not valid if issued before this date")
	a.Required("expires_in")
	a.Required("refresh_expires_in")
	a.Required("not-before-policy")
})

// OauthToken represents an Oauth 2.0 token payload
var OauthToken = a.MediaType("application/vnd.oauthtoken+json", func() {
	a.TypeName("OauthToken")
	a.Description("Oauth 2.0 token payload")
	a.Attributes(func() {
		a.Attribute("access_token", d.String, "Access token")
		a.Attribute("expiry", d.String, "Expiry")
		a.Attribute("refresh_token", d.String, "RefreshToken")
		a.Attribute("token_type", d.String, "Token type")
	})
	a.View("default", func() {
		a.Attribute("access_token")
		a.Attribute("expiry")
		a.Attribute("refresh_token")
		a.Attribute("token_type")
	})
})

// redirectLocation represents an redirect location
var redirectLocation = a.MediaType("application/vnd.redirectlocation+json", func() {
	a.TypeName("RedirectLocation")
	a.Description("Redirect Location")
	a.Attributes(func() {
		a.Attribute("redirect_location", d.String, "The location which should be used to redirect browser")
		a.Required("redirect_location")
	})
	a.View("default", func() {
		a.Attribute("redirect_location")
		a.Required("redirect_location")
	})
})
