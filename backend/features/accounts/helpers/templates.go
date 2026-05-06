package helpers

import "fmt"

/*
 * VerificationEmailHTML génère le corps HTML de l'email de vérification.
 *
 * Attend  : l'URL du frontend et le token brut de vérification.
 * Retourne: le contenu HTML prêt à l'envoi.
 */

func VerificationEmailHTML(frontendURL, token string) string {
	link := fmt.Sprintf("%s/verify-email?token=%s", frontendURL, token)
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family:sans-serif;max-width:600px;margin:auto;padding:32px">
  <h2>Vérifiez votre adresse email</h2>
  <p>Cliquez sur le lien ci-dessous pour activer votre compte Skema :</p>
  <a href="%s" style="display:inline-block;padding:12px 24px;background:#000;color:#fff;text-decoration:none;border-radius:6px">
    Vérifier mon email
  </a>
  <p style="margin-top:24px;color:#666;font-size:14px">
    Ce lien expire dans 24 heures. Si vous n'avez pas créé de compte, ignorez cet email.
  </p>
</body>
</html>`, link)
}

/*
 * ResetEmailHTML génère le corps HTML de l'email de réinitialisation de mot de passe.
 *
 * Attend  : l'URL du frontend et le token brut de réinitialisation.
 * Retourne: le contenu HTML prêt à l'envoi.
 */

func ResetEmailHTML(frontendURL, token string) string {
	link := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="font-family:sans-serif;max-width:600px;margin:auto;padding:32px">
  <h2>Réinitialisation de votre mot de passe</h2>
  <p>Cliquez sur le lien ci-dessous pour définir un nouveau mot de passe :</p>
  <a href="%s" style="display:inline-block;padding:12px 24px;background:#000;color:#fff;text-decoration:none;border-radius:6px">
    Réinitialiser mon mot de passe
  </a>
  <p style="margin-top:24px;color:#666;font-size:14px">
    Ce lien expire dans 1 heure. Si vous n'avez pas fait cette demande, ignorez cet email.
  </p>
</body>
</html>`, link)
}
