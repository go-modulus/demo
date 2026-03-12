package service_test

import (
	"boilerplate/internal/auth"
	"boilerplate/internal/auth/service"
	"boilerplate/internal/infra/test/expect"
	"boilerplate/internal/infra/test/spec"
	"testing"
	"time"
)

type tokenParserCfg struct {
	auth.ModuleConfig
}

func TestTokenParser_GenerateNewToken(t *testing.T) {
	cfg := &auth.ModuleConfig{
		PublicKey:  "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAw+FxoK2nyVSwh3TW0Usg\ntbpQi6zAaQKWJZgJlmrHFfuV9h4JxfZeOvrr/nJOzhPFGKpKhBmrzNjPkIWKbnWh\nO0KLT/SAAx44odRwjESS2C7B7vt1Es8UJhZvauTFRtG7/K/VMB85KEAZ0rfHdYLE\nNGtPQ0YAT5wtGxtPiCYBQVuDcSR4reI6xn06ly8ieCOgYB6rzxmfwH71nbGf9ABE\nSFCEA0C3SJQseflvl0Pqbaz1EGOHz55pXDZsm9+Q4tmacm5FNM8kT7N3sBtoFCjB\nlGpvDzSjS0IzDlkT16OtFj6a9lJnANTItf1Xj1sWCG0xqIXIjs0LPaYtk94Wyskn\neQIDAQAB\n-----END PUBLIC KEY-----",
		PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQEAw+FxoK2nyVSwh3TW0UsgtbpQi6zAaQKWJZgJlmrHFfuV9h4J\nxfZeOvrr/nJOzhPFGKpKhBmrzNjPkIWKbnWhO0KLT/SAAx44odRwjESS2C7B7vt1\nEs8UJhZvauTFRtG7/K/VMB85KEAZ0rfHdYLENGtPQ0YAT5wtGxtPiCYBQVuDcSR4\nreI6xn06ly8ieCOgYB6rzxmfwH71nbGf9ABESFCEA0C3SJQseflvl0Pqbaz1EGOH\nz55pXDZsm9+Q4tmacm5FNM8kT7N3sBtoFCjBlGpvDzSjS0IzDlkT16OtFj6a9lJn\nANTItf1Xj1sWCG0xqIXIjs0LPaYtk94WyskneQIDAQABAoIBAAufqC+BvpAq72LK\nYyj+TU2bZcY0FSLmHWTTtdXwNiouyMJXG/tcWFElrLVnsRd3gt2o1U0rfM9mGkMY\nOZO+dTfgBgZTmvkUZQXtZlRciV48wBNfZp4cHbT45dmxA6flMEOZZ4w4fhrOWvnv\nH+3wOZZmu8hhMgmHWgHeyARrRO3M1rxXQVKwMDphOcMlIBRwxYRQsL7tWuPA7hX/\n5ZNLW3aXd5/8IynOY9sDYAQBKLv6swrKkh99OYGnm8CMPaKZsdDBhzgf9XA9ouSR\noDs3PnqYY4BBHiys7Yqf/2SOvyiSb1+U9iw7OCa+yCRV+TZl6E9ybtwgPuworaiU\nVTZrglECgYEA9MMpp7qfn2kV+GUX+SUV7XzcwyQJSTwi2/B5gjvOOsmvMkz9RvQ2\nyQMBpkaqHTn54VcyEOl9OJ3Vdt7IOjwCD7hkH5F9Z2zsEFRCpC9015yKWdBsa/re\n9KoCLj0mDEtf9NysIn+x2nuv6kxP6t8KEUe2u5iNqOgK5yrO94KNJb0CgYEAzN++\nze0J74URJd2kaSOaLP2TEbQTT+1dfkOzPW92HirnNn7rOKTVnD8ht1G1VZXzQbQ3\nQ0lwci7vUPbSVw6hApFFMfhr5lfyDPcZBnRfBIHKfEmaTn+k/YQqAH1ztPrzMLrn\nH39Eehyed4Xqe4L4LWN1cLfnG3puEGpgTNm6zm0CgYBv6TAJlcXYMEcwXKC6dN9y\nx610t+xbBNj3cRtNlaS0snSdbiA8KftGq048xYCQfmqnQqQMoYV0to3cnP41yiwz\nHd8BpBcPi/jfendB9MTatKN5b1ezg3Afs//tPl5ALtJ/9cnquDIMsJL9cMj1nedP\ngVemrJjQys/5ZFRfTNzWjQKBgEwKvCJo2eg6Jrw8QRr5KO+MCvtmMEjZXHtSG4Qx\nC9F0sS8L+riijdqZoCUPwdOLfaekgWKLLp5jB1aw1i+T8XUngFxkzX/IosHnMTWx\nGddtaT+qfgim3hFu7bwS1FCXWI58wO5y6XK9jp/kZ70CRqVqJhv5VmFfltym7yl3\nIxwdAoGAXkAtp1jFrFD6oKjtyuZTfSoVxqhOlsDnRIEO0s7tMNlzq+8B+3CyfTWw\n0Qg7ep42azS8Xwsk2Vz0ZDH6BelQhRVD4vmcEwMHpZZPYhpnJ8LnzSbK056NizwU\nFgtfWp86CTM/38DK0HvrgLGwB9PPZLmtxYyp6S3x2/yWeI8xRMY=\n-----END RSA PRIVATE KEY-----",
	}
	parser, err := service.NewTokenParser(cfg)

	if err != nil {
		t.Fatal("Cannot create token parser", err.Error())
		return
	}
	token, claims, err := parser.GenerateNewToken(
		"test",
		"123",
		[]string{"moderator"},
		[]string{"*"},
		time.Hour,
	)

	spec.When(t, "New token is generated")
	spec.NoError(t, "No error", err)
	spec.Then(t, "token is not empty", expect.NotEqual("", token))
	spec.Then(
		t, "claims are valid",
		expect.NotNil(claims),
		expect.Equal("test", claims.Subject),
		expect.Equal("moderator", claims.Roles[0]),
		expect.Equal("*", claims.Permissions[0]),
		expect.True(claims.ExpiresAt.After(time.Now().Add(59*time.Minute))),
	)
}

func TestTokenParser_Parse(t *testing.T) {
	cfg := &auth.ModuleConfig{
		PublicKey:  "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAw+FxoK2nyVSwh3TW0Usg\ntbpQi6zAaQKWJZgJlmrHFfuV9h4JxfZeOvrr/nJOzhPFGKpKhBmrzNjPkIWKbnWh\nO0KLT/SAAx44odRwjESS2C7B7vt1Es8UJhZvauTFRtG7/K/VMB85KEAZ0rfHdYLE\nNGtPQ0YAT5wtGxtPiCYBQVuDcSR4reI6xn06ly8ieCOgYB6rzxmfwH71nbGf9ABE\nSFCEA0C3SJQseflvl0Pqbaz1EGOHz55pXDZsm9+Q4tmacm5FNM8kT7N3sBtoFCjB\nlGpvDzSjS0IzDlkT16OtFj6a9lJnANTItf1Xj1sWCG0xqIXIjs0LPaYtk94Wyskn\neQIDAQAB\n-----END PUBLIC KEY-----",
		PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQEAw+FxoK2nyVSwh3TW0UsgtbpQi6zAaQKWJZgJlmrHFfuV9h4J\nxfZeOvrr/nJOzhPFGKpKhBmrzNjPkIWKbnWhO0KLT/SAAx44odRwjESS2C7B7vt1\nEs8UJhZvauTFRtG7/K/VMB85KEAZ0rfHdYLENGtPQ0YAT5wtGxtPiCYBQVuDcSR4\nreI6xn06ly8ieCOgYB6rzxmfwH71nbGf9ABESFCEA0C3SJQseflvl0Pqbaz1EGOH\nz55pXDZsm9+Q4tmacm5FNM8kT7N3sBtoFCjBlGpvDzSjS0IzDlkT16OtFj6a9lJn\nANTItf1Xj1sWCG0xqIXIjs0LPaYtk94WyskneQIDAQABAoIBAAufqC+BvpAq72LK\nYyj+TU2bZcY0FSLmHWTTtdXwNiouyMJXG/tcWFElrLVnsRd3gt2o1U0rfM9mGkMY\nOZO+dTfgBgZTmvkUZQXtZlRciV48wBNfZp4cHbT45dmxA6flMEOZZ4w4fhrOWvnv\nH+3wOZZmu8hhMgmHWgHeyARrRO3M1rxXQVKwMDphOcMlIBRwxYRQsL7tWuPA7hX/\n5ZNLW3aXd5/8IynOY9sDYAQBKLv6swrKkh99OYGnm8CMPaKZsdDBhzgf9XA9ouSR\noDs3PnqYY4BBHiys7Yqf/2SOvyiSb1+U9iw7OCa+yCRV+TZl6E9ybtwgPuworaiU\nVTZrglECgYEA9MMpp7qfn2kV+GUX+SUV7XzcwyQJSTwi2/B5gjvOOsmvMkz9RvQ2\nyQMBpkaqHTn54VcyEOl9OJ3Vdt7IOjwCD7hkH5F9Z2zsEFRCpC9015yKWdBsa/re\n9KoCLj0mDEtf9NysIn+x2nuv6kxP6t8KEUe2u5iNqOgK5yrO94KNJb0CgYEAzN++\nze0J74URJd2kaSOaLP2TEbQTT+1dfkOzPW92HirnNn7rOKTVnD8ht1G1VZXzQbQ3\nQ0lwci7vUPbSVw6hApFFMfhr5lfyDPcZBnRfBIHKfEmaTn+k/YQqAH1ztPrzMLrn\nH39Eehyed4Xqe4L4LWN1cLfnG3puEGpgTNm6zm0CgYBv6TAJlcXYMEcwXKC6dN9y\nx610t+xbBNj3cRtNlaS0snSdbiA8KftGq048xYCQfmqnQqQMoYV0to3cnP41yiwz\nHd8BpBcPi/jfendB9MTatKN5b1ezg3Afs//tPl5ALtJ/9cnquDIMsJL9cMj1nedP\ngVemrJjQys/5ZFRfTNzWjQKBgEwKvCJo2eg6Jrw8QRr5KO+MCvtmMEjZXHtSG4Qx\nC9F0sS8L+riijdqZoCUPwdOLfaekgWKLLp5jB1aw1i+T8XUngFxkzX/IosHnMTWx\nGddtaT+qfgim3hFu7bwS1FCXWI58wO5y6XK9jp/kZ70CRqVqJhv5VmFfltym7yl3\nIxwdAoGAXkAtp1jFrFD6oKjtyuZTfSoVxqhOlsDnRIEO0s7tMNlzq+8B+3CyfTWw\n0Qg7ep42azS8Xwsk2Vz0ZDH6BelQhRVD4vmcEwMHpZZPYhpnJ8LnzSbK056NizwU\nFgtfWp86CTM/38DK0HvrgLGwB9PPZLmtxYyp6S3x2/yWeI8xRMY=\n-----END RSA PRIVATE KEY-----",
	}
	parser, err := service.NewTokenParser(cfg)

	if err != nil {
		t.Fatal("Cannot create token parser", err.Error())
		return
	}
	t.Run(
		"ExtractClaims", func(t *testing.T) {

			token, _, err := parser.GenerateNewToken(
				"test",
				"123",
				[]string{"moderator"},
				[]string{"*"},
				time.Hour,
			)
			if err != nil {
				t.Fatal("Cannot generate token", err.Error())
				return
			}
			claims, err := parser.Parse(token)

			spec.Given(t, "A valid token")
			spec.When(t, "The token is parsed")
			spec.NoError(t, "No error", err)
			spec.Then(
				t, "Claims are valid",
				expect.NotNil(claims),
				expect.Equal("test", claims.Subject),
				expect.Equal("moderator", claims.Roles[0]),
				expect.Equal("*", claims.Permissions[0]),
				expect.True(claims.ExpiresAt.After(time.Now().Add(59*time.Minute))),
			)
		},
	)

	t.Run(
		"ExpiredToken", func(t *testing.T) {
			token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlcyI6WyJtb2RlcmF0b3IiXSwicGVybWlzc2lvbnMiOlsiKiJdLCJpc3MiOiJWRVJWIiwic3ViIjoidGVzdCIsImF1ZCI6WyIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0xMDAwMDAwMDAwMDEiXSwiZXhwIjoxNjg0MjE1MDMyLCJuYmYiOjE2ODQyMTUwMzEsImlhdCI6MTY4NDIxNTAzMSwianRpIjoiYTk0NjdmNTMtZTlhZC00NDBkLWJkZDgtYTRhZWNkMzVmYjM1In0.YtL_oielCkv9fFr-fGW2wNnGRCPh5-6pHY8OEQQCD5lvjndOJhM8XQI9hEDq2M04q2vnLPIOjKvlqPR4ECBC9usTVYHtJ5xzag6Mtt0X10hHtzx_SMXAg951xP5ebhTeQtmiGR4ouUGjCniQOSvwUB8hnZbZUthSpl_VMghyXlH-oMOAYZSCxZ8_id_Q8GgAqtAmE515Ar6r49PwRxqMGKH-8JyDtOkUlYh9HUR60sxVwY5Jt_s7uB9t6LPAqwd0Szh0jiETg64keS2XaqHhd9_kVDPQTKIcCqeG83qUT-b3KnjcwVf7VAauMGcWTjlces5Le7hiOQlc5PDV_lLZYA"
			claims, err := parser.Parse(token)

			spec.Given(t, "An expired token")
			spec.When(t, "The token is parsed")
			spec.HasCommonError(t, "Expiration error is present", err, service.TokenExpired)
			spec.Then(t, "Claims are nil", expect.Nil(claims))
		},
	)

	t.Run(
		"WrongSignature", func(t *testing.T) {
			token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlcyI6WyJtb2RlcmF0b3IiXSwicGVybWlzc2lvbnMiOlsiKiJdLCJpc3MiOiJWRVJWIiwic3ViIjoidGVzdCIsImF1ZCI6WyIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0xMDAwMDAwMDAwMDEiXSwiZXhwIjoxNjg0MjE1MDMyLCJuYmYiOjE2ODQyMTUwMzEsImlhdCI6MTY4NDIxNTAzMSwianRpIjoiYTk0NjdmNTMtZTlhZC00NDBkLWJkZDgtYTRhZWNkMzVmYjM1In0.YtL_oielCkv9fFr-fGW2wNnGRCPh5-6pHY8OEQQCD5lvjndOJhM8XQI9hEDq2M04q2vnLPIOjKvlqPR4ECBC9usTVYHtJ5xzag6Mtt0X10hHtzx_SMXAg951xP5ebhTeQtmiGR4ouUGjCniQOSvwUB8hnZbZUthSpl_VMghyXlH-oMOAYZSCxZ8_id_Q8GgAqtAmE515Ar6r49PwRxqMGKH-8JyDtOkUlYh9HUR60sxVwY5Jt_s7uB9t6LPAqwd0Szh0jiETg64keS2XaqHhd9_kVDPQTKIcCqeG83qUT-b3KnjcwVf7VAuMGcWTjlces5Le7hiOQlc5PDV_lLZYA"
			claims, err := parser.Parse(token)

			spec.Given(t, "The token with wrong signature")
			spec.When(t, "The token is parsed")
			spec.HasCommonError(t, "Invalid token error is present", err, service.IncorrectToken)
			spec.Then(t, "claims are nil", expect.Nil(claims))
		},
	)
}
