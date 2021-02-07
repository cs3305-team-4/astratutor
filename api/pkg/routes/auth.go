package routes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cs3305-team-4/api/pkg/services"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InjectAuthRoutes(subrouter *mux.Router) {

	subrouter.PathPrefix("/login").HandlerFunc(authLogin).Methods("POST")
}

type LoginDTO struct {
	ID       string `json:"id" validate:"len=0"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

type LoginResponseDTO struct {
	JWT string `json:"jwt"`
}

type AuthClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func (ac *AuthClaims) Valid() error {
	return nil
}

func authLogin(w http.ResponseWriter, r *http.Request) {
	var login LoginDTO

	if !ParseBody(w, r, &login) {
		return
	}

	acc, err := services.ReadAccountByEmail(login.Email, nil)
	if err != nil {
		restError(w, r, errors.New("invalid email or password"), http.StatusForbidden)
		return
	}

	hash, err := services.ReadPasswordHashByAccountID(acc.ID)
	if err != nil {
		restError(w, r, errors.New("unknown error occured during authentication"), http.StatusForbidden)
		return
	}

	if hash.ValidMatch(login.Password) {
		claims := &AuthClaims{
			StandardClaims: jwt.StandardClaims{
				Issuer:    "grindsapp",
				Subject:   acc.ID.String(),
				Audience:  "grindsapp",
				IssuedAt:  int64(time.Now().Unix()),
				ExpiresAt: int64(time.Now().Add(time.Second * time.Duration(viper.GetUint64("auth.jwt.ttl"))).Unix()),
			},
			Email: acc.Email,
		}

		jwtStr, err := claims.serializeSignJWT()
		if err != nil {
			restError(w, r, errors.New("unknown error occured during authentication"), http.StatusForbidden)
			return
		}

		if err = json.NewEncoder(w).Encode(&LoginResponseDTO{JWT: jwtStr}); err != nil {
			restError(w, r, err, http.StatusInternalServerError)
			return
		}
	} else {
		restError(w, r, errors.New("invalid email or password"), http.StatusForbidden)
		return
	}

}

func parseVerifyJWT(jwtStr string) (*jwt.Token, error) {
	token, err := new(jwt.Parser).ParseWithClaims(
		jwtStr,
		&AuthClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("expected rsa signed jwt")
			}

			pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(viper.GetString("auth.jwt.public_key")))
			if err != nil {
				log.Errorf("error parsing jwt public key: %s", err)
				return nil, err
			}

			return pubKey, nil
		},
	)
	if err != nil {
		log.Errorf("error parsing jwt: %s", err)
		return nil, err
	}

	return token, nil
}

func (a *AuthClaims) serializeSignJWT() (string, error) {
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(viper.GetString("auth.jwt.private_key")))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, a)

	jwtStr, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return jwtStr, nil
}

type AuthContext struct {
	Claims  *AuthClaims
	Account *services.Account
}

func (ac *AuthContext) Authenticated() bool {
	return ac.Claims != nil && ac.Account != nil
}

type AuthContextKeyType string

const (
	authContextKey AuthContextKeyType = "routes.auth.context"
)

func ParseRequestAuth(r *http.Request) (*AuthContext, error) {
	val := r.Context().Value(authContextKey)
	if context, ok := val.(*AuthContext); ok {
		return context, nil
	}
	return nil, errors.New("auth context not valid type")
}

func ReadRequestAuthContext(r *http.Request) (*AuthContext, error) {
	val, err := ParseRequestAuth(r)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, errors.New("auth context not present on request")
	}
	if !val.Authenticated() {
		return nil, errors.New("auth context not authenticated")
	}
	return val, nil
}

func authRequired(next http.Handler) http.Handler {
	return authMiddleware(func(w http.ResponseWriter, r *http.Request, ac *AuthContext) error {
		// The auth middleware pulls their account and checks if it's suspended, so we don't need to do any checking here.
		return nil
	}, true)(next)
}

func authMiddleware(userSuppliedAuthCtxValidator func(w http.ResponseWriter, r *http.Request, ac *AuthContext) error, required bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// First we need to check if they used this middleware multiple times (i.e an AuthContext is already present on the request)
			authContext, err := ReadRequestAuthContext(r)

			// Auth context is already present, run validation for this route
			if err == nil {
				err = userSuppliedAuthCtxValidator(w, r, authContext)
				if err != nil {
					restError(w, r, err, http.StatusForbidden)
					return
				}

				next.ServeHTTP(w, r)
				return
			}

			// There was no auth context present on this request, we need to parse the JWT and generate a context
			if _, ok := r.Header["Authorization"]; ok {
				auth := r.Header["Authorization"][0]
				bearerCheck := strings.Split(auth, " ")

				if len(bearerCheck) != 2 || strings.ToLower(bearerCheck[0]) != "bearer" {
					restError(w, r, errors.New("only bearer tokens can be used for authorization"), http.StatusForbidden)
					return
				}

				jwtStr := bearerCheck[1]

				token, err := parseVerifyJWT(jwtStr)
				if err != nil {
					restError(w, r, errors.New("error verifying jwt"), http.StatusForbidden)
					return
				}

				if claims, ok := token.Claims.(*AuthClaims); ok {
					uuid, err := uuid.Parse(claims.Subject)
					if err != nil {
						restError(w, r, errors.New("could not parse sub claim on jwt, expected account uuid"), http.StatusForbidden)
						return
					}

					account, err := services.ReadAccountByID(uuid, nil)
					if err != nil {
						restError(w, r, fmt.Errorf("error occured when looking for account: %s", err), http.StatusForbidden)
						return
					}

					if account.Suspended == true {
						restError(w, r, errors.New("this account has been suspended"), http.StatusForbidden)
						return
					}

					authContext := &AuthContext{
						Claims:  claims,
						Account: account,
					}

					err = userSuppliedAuthCtxValidator(w, r, authContext)
					if err != nil {
						restError(w, r, err, http.StatusForbidden)
						return
					}

					ctx := context.WithValue(r.Context(), authContextKey, authContext)

					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}

				restError(w, r, errors.New("jwt claims invalid"), http.StatusForbidden)
				return
			}

			if required {
				restError(w, r, errors.New("endpoint requires authorization header"), http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), authContextKey, &AuthContext{})))
		})
	}
}
