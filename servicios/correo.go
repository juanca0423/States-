package servicios

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

func buildVerificationLink(token string) string {
	return fmt.Sprintf("https://eeffs.com/verificar?token=%s", token)
}

func buildResetLink(token string) string {
	return fmt.Sprintf("https://eeffs.com/reset-password?token=%s", token)
}

func buildVerificationHtml(nombre string, token string) string {
	link := buildVerificationLink(token)
	return fmt.Sprintf(`
			<div style="font-family: sans-serif; max-width: 600px; margin: auto; color: #333;">
				<h2 style="color: #007bff;">¡Hola %s!</h2>
				<p>Para terminar de configurar tu cuenta en <strong>States</strong>, por favor confirma tu correo haciendo clic abajo:</p>
				<div style="text-align: center; margin: 30px 0;">
					<a href="%s" style="background-color: #007bff; color: white; padding: 12px 25px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold;">
						Verificar mi correo
					</a>
				</div>
				<hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
				<p style="font-size: 0.8em; color: #888; text-align: center;">
					<strong>Nota:</strong> Este es un correo automático, por favor no respondas a esta dirección.
				</p>
				<p style="font-size: 0.8em; color: #888; text-align: center;">
					Si no creaste esta cuenta en States, puedes ignorar este mensaje de forma segura.
				</p>
			</div>
		`, nombre, link)
}

func EnviarCorreoVerificacion(destinatario string, nombre string, token string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(apiKey)

	// El link que apunta a tu servidor de producción
	link := buildVerificationLink(token)

	params := &resend.SendEmailRequest{
		From:    "States <no-reply@eeffs.com>",
		To:      []string{destinatario},
		Subject: "Confirma tu cuenta en States",
		Html: fmt.Sprintf(`
			<div style="font-family: sans-serif; max-width: 600px; margin: auto; color: #333;">
				<h2 style="color: #007bff;">¡Hola %s!</h2>
				<p>Para terminar de configurar tu cuenta en <strong>States</strong>, por favor confirma tu correo haciendo clic abajo:</p>
				<div style="text-align: center; margin: 30px 0;">
					<a href="%s" style="background-color: #007bff; color: white; padding: 12px 25px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold;">
						Verificar mi correo
					</a>
				</div>
				<hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
				<p style="font-size: 0.8em; color: #888; text-align: center;">
					<strong>Nota:</strong> Este es un correo automático, por favor no respondas a esta dirección.
				</p>
				<p style="font-size: 0.8em; color: #888; text-align: center;">
					Si no creaste esta cuenta en States, puedes ignorar este mensaje de forma segura.
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

func EnviarCorreoRecuperacion(destinatario string, nombre string, token string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(apiKey)

	link := buildResetLink(token)

	params := &resend.SendEmailRequest{
		From:    "States <no-reply@eeffs.com>",
		To:      []string{destinatario},
		Subject: "Recuperación de contraseña en States",
		Html: fmt.Sprintf(`
			<div style="font-family: sans-serif; max-width: 600px; margin: auto; color: #333;">
				<h2 style="color: #007bff;">¡Hola %s!</h2>
				<p>Has solicitado restablecer tu contraseña en <strong>States</strong>. Haz clic en el siguiente botón para cambiarla:</p>
				<div style="text-align: center; margin: 30px 0;">
					<a href="%s" style="background-color: #007bff; color: white; padding: 12px 25px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold;">
						Restablecer contraseña
					</a>
				</div>
				<hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
				<p style="font-size: 0.8em; color: #888; text-align: center;">
					<strong>Nota:</strong> Si no solicitaste este cambio, puedes ignorar este correo de forma segura.
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
