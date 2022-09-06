package wevpn

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/openvpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Provider) OpenVPNConfig(connection models.Connection,
	settings settings.OpenVPN, ipv6Supported bool) (lines []string) {
	providerSettings := utils.OpenVPNProviderSettings{
		RemoteCertTLS: true,
		AuthUserPass:  true,
		Ciphers: []string{
			openvpn.AES256gcm,
		},
		Auth:          openvpn.SHA512,
		Ping:          30, //nolint:gomnd
		RenegDisabled: true,
		CA:            "MIIDQjCCAiqgAwIBAgIUPppqnRZfvGGrT4GjXFE4Q29QzgowDQYJKoZIhvcNAQELBQAwEzERMA8GA1UEAwwIQ2hhbmdlTWUwHhcNMTkxMTA1MjMzMzIzWhcNMjkxMTAyMjMzMzIzWjATMREwDwYDVQQDDAhDaGFuZ2VNZTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAL5DFBJlTqhXukJFWlI8TNW9+HEQCZXhyVFvQhJFF2xIGVNx51XzqxiRANjVJZJrA68kV8az0v2Dxj0SFnRWDR6pOjjdp2CyHFcgHyfv+4MrsreAtkue86bB/1ECPWaoIwtaLnwI6SEmFZl98RlI9v4M/8IE4chOnMrM/F22+2OXI//TduvTcbyOMUiiouIP8UG1FB3J5FyuaW6qPZz2G0efDoaOI+E9LSxE87OoFrII7UqdHlWxRb3nUuPU1Ee4rN/d4tFyP4AvPKfsGhVOwyGG21IdRnbXIuDi0xytkCGOZ4j2bq5zqudnp4Izt6yJgdzZpQQWK3kSHB3qTT/Yzl8CAwEAAaOBjTCBijAdBgNVHQ4EFgQUXYkoo4WbkkvbgLVdGob9RScRf3AwTgYDVR0jBEcwRYAUXYkoo4WbkkvbgLVdGob9RScRf3ChF6QVMBMxETAPBgNVBAMMCENoYW5nZU1lghQ+mmqdFl+8YatPgaNcUThDb1DOCjAMBgNVHRMEBTADAQH/MAsGA1UdDwQEAwIBBjANBgkqhkiG9w0BAQsFAAOCAQEAOr1XmyWBRYfTQPNvZZ+DjCfiRYzLOi2AGefZt/jETqPDF8deVbyL1fLhXZzuX+5Etlsil3PflJjpzc/FSeZRuYRaShtwF3j6I08Eww9rBkaCnsukMUcLtMOvhdAU8dUakcRA2wkQ7Z+TWdMBv5+/6MnX10An1fIz7bAy3btMEOPTEFLo8Bst1SxJtUMaqhUteSOJ1VorpK3CWfOFaXxbJAb4E0+3zt6Vsc8mY5tt6wAi8IqiN4WD79ZdvKxENK4FMkR1kNpBY97mvdf82rzpwiBuJgN5ywmH78Ghj+9T8nI6/UIqJ1y22IRYGv6dMif8fHo5WWhCv3qmCqqY8vwuxw==",             //nolint:lll
		Cert:          "MIIDTDCCAjSgAwIBAgIRAKxt8SMIXezjmHm2KDCAQdIwDQYJKoZIhvcNAQELBQAwEzERMA8GA1UEAwwIQ2hhbmdlTWUwHhcNMTkxMTA1MjMzMzI0WhcNMjkxMTAyMjMzMzI0WjAOMQwwCgYDVQQDDAN0Y3AwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCvEwY2erLhMm3Mpsnybm3G6zvGyeblUAaehQVEUs+KM2/5np0Ovx0y8Iz9pIC9ITaWM0B3dM6uBsNEtylZIe4Dd9aFujunSeCFsLRf8i9AbrUombpQ6P4jzYFBxwcEw//UShwa4HZI6JuSYikdpx/dyXdBH2skahwDVc8VUFdBLLSglfKGbuzP9GsdSwQCeBRWgA3dvIzIkQkBwfnt9WQKUfRAe8e5NybaAn8Yuu9sjLkQe6eyV7toxkZTcEXdABG2vtdTEzlAsQilZzIxg3jcdeEgMgRKngng+YNP0rR5nofZ1iDlp+vBj0nuqTTJLHMrRWPIc7bdYFD/f2J49WORAgMBAAGjgZ8wgZwwCQYDVR0TBAIwADAdBgNVHQ4EFgQUmSAFmCo1FAKVq8RQF7jMxMxcMtUwTgYDVR0jBEcwRYAUXYkoo4WbkkvbgLVdGob9RScRf3ChF6QVMBMxETAPBgNVBAMMCENoYW5nZU1lghQ+mmqdFl+8YatPgaNcUThDb1DOCjATBgNVHSUEDDAKBggrBgEFBQcDAjALBgNVHQ8EBAMCB4AwDQYJKoZIhvcNAQELBQADggEBADPqdEgL+0kou8P974QEaNg1XOAXpwP0NNqbkZ/Oj9+Lp96YAhAHOAJig+RWbBktK8zu8oUUGR1qLXAWCmirlXErVuBRnadTEh3A7SOuY02BcsYAtpQ2EU9j5K/LV7nTfagkVdWy7x/av361UD4t9fv1j4YYTh4XLRp7KVXs6AGZ7T1hqPYFMUIoPpFhPzFxH4euJjfazr4SkTR6k6Vhw3pyFd6HP65vcqpzHGxFytSa8HtltBk2DpzIf8yV9TEy+gOXFaaGss0YKQ5OU1ieqZRuLVEGiu17lByYiQGyemIETJbdkyiSg93dDJRxjaTk7c8CEdpipt07ndSIPldMtXA=", //nolint:lll
		TLSCrypt:      "7be66c0df0b8855e076d9e37b19f9ff3c1735ed537dee6dc786e51bdb8502f878077eeba0420a25e2b04814d22bbdcc0191a4fc396fdba1af6eb090a9d8664f18e70012ee98a2e32c28620a771d13cf3a619c417480c2c312562fffaebfd7ba73f57a28edde6c287365e6ce28291a29728da211cb53e01aa46b92f5f276c61fb46bd810b41219022c8f3d9e699fe9ade6bfcbb937fbbf6f49d741740e71c7c008a9a13c2432608038c6310b4f33588d8d234b3dffcf0823395267d73140d0e9a40e323ca92866c37073bfb072ab9de518bb9f2c65df7e219c2f114afbcf7c6e3c401cb08c3ed2901725b0601d2b5de89245719dd32506d52f149d14156215c1e",                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             //nolint:lll
		ExtraLines: []string{
			"redirect-gateway def1 bypass-dhcp",
		},
	}
	return utils.OpenVPNConfig(providerSettings, connection, settings, ipv6Supported)
}
