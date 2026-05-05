package mailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const resendApiUrl = "https://api.resend.com/emails"

type Email struct {
	To      string
	Subject string
	HTML    string
}

type Mailer struct {
	from   string
	apiKey string
}

type resendPayload struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

/*
 * New instancie un Mailer configuré avec l'expéditeur et la clé API Resend.
 *
 * Attend  : l'adresse expéditeur et la clé API Resend.
 * Retourne: un pointeur vers Mailer prêt à l'emploi.
 */

func New(from, apiKey string) *Mailer {
	return &Mailer{from: from, apiKey: apiKey}
}

/*
 * Send envoie un email via l'API Resend.
 *
 * Attend  : un struct Email avec le destinataire, le sujet et le corps HTML.
 * Retourne: une erreur si l'envoi échoue.
 */

func (m *Mailer) Send(email Email) error {
	payload := resendPayload{
		From:    m.from,
		To:      []string{email.To},
		Subject: email.Subject,
		HTML:    email.HTML,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, resendApiUrl, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("%s", ErrSendFailed)
	}

	return nil
}
