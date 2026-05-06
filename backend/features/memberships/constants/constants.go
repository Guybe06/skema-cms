package constants

import "time"

const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleMember = "member"

	StatusActive  = "active"
	StatusPending = "pending"
	StatusRevoked = "revoked"

	InviteTokenExpiry = 48 * time.Hour

	ErrAlreadyMember      = "cet utilisateur est déjà membre de l'organisation"
	ErrMemberNotFound     = "membre introuvable"
	ErrCannotRemoveOwner  = "le propriétaire ne peut pas être retiré de son organisation"
	ErrCannotChangeOwner  = "le rôle du propriétaire ne peut pas être modifié ici"
	ErrInvalidRole        = "rôle invalide, valeurs acceptées : admin, member"
	ErrInviteTokenInvalid = "invitation invalide ou expirée"
	ErrNotAuthorized      = "action réservée au propriétaire ou aux administrateurs"

	MsgInviteSent    = "Invitation envoyée."
	MsgInviteAccepted = "Invitation acceptée, vous êtes maintenant membre."
	MsgRoleUpdated   = "Rôle mis à jour."
	MsgMemberRemoved = "Membre retiré."
)
