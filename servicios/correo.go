package servicios

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

func EnviarCorreoVerificacion(destinatario string, nombre string, token string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(apiKey)

	// El link que apunta a tu servidor de producción
	link := fmt.Sprintf("https://eeffs.com/verificar?token=%s", token)

	params := &resend.SendEmailRequest{
		From:    "States <onboarding@resend.dev>", // Cambia a info@eeffs.com tras verificar dominio
		To:      []string{destinatario},
		Subject: "Confirma tu cuenta en States",
		Html: fmt.Sprintf(`
			<div style="font-family: sans-serif; max-width: 600px; margin: auto;">
				<h2>¡Hola %s!</h2>
				<p>Para terminar de configurar tu cuenta en <strong>States</strong>, por favor confirma tu correo haciendo clic abajo:</p>
				<a href="%s" style="background-color: #007bff; color: white; padding: 12px 25px; text-decoration: none; border-radius: 5px; display: inline-block;">
					Verificar mi correo
				</a>
				<p style="margin-top: 20px; font-size: 0.8em; color: #666;">
					Si no creaste esta cuenta, puedes ignorar este correo.
				</p>
			</div>
		`, nombre, link),
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return err
	}

	return nil
}
