package helpers

import "fmt"

// Génère le corps HTML de l'email d'invitation à rejoindre une organisation.
func InvitationEmailHTML(orgName, inviterEmail, frontendURL, token string) string {
	link := fmt.Sprintf("%s/invitations/accept?token=%s", frontendURL, token)
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family:sans-serif;max-width:600px;margin:auto;padding:32px">
  <h2>Invitation à rejoindre %s</h2>
  <p>%s vous invite à rejoindre l'organisation <strong>%s</strong> sur Skema.</p>
  <p>
    <a href="%s" style="background:#111;color:#fff;padding:12px 24px;border-radius:6px;text-decoration:none;display:inline-block">
      Accepter l'invitation
    </a>
  </p>
  <p style="color:#888;font-size:13px">Ce lien expire dans 48 heures.</p>
</body>
</html>`, orgName, inviterEmail, orgName, link)
}
