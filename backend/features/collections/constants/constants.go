package constants

const (
	// Types de champs supportés.
	FieldTypeText        = "text"
	FieldTypeTextarea    = "textarea"
	FieldTypeRichtext    = "richtext"
	FieldTypeNumber      = "number"
	FieldTypeBoolean     = "boolean"
	FieldTypeDate        = "date"
	FieldTypeDatetime    = "datetime"
	FieldTypeEmail       = "email"
	FieldTypeURL         = "url"
	FieldTypePhone       = "phone"
	FieldTypeColor       = "color"
	FieldTypeSelect      = "select"
	FieldTypeMultiselect = "multiselect"
	FieldTypeSlug        = "slug"
	FieldTypeJSON        = "json"
	FieldTypePassword    = "password"
	FieldTypeFile        = "file"
	FieldTypeImage       = "image"

	ErrCollectionNotFound = "collection introuvable"
	ErrFieldNotFound      = "champ introuvable"
	ErrNotAuthorized      = "accès non autorisé"
	ErrTableNameTaken     = "ce nom de table est déjà utilisé sur cette connexion"
	ErrColumnNameTaken    = "ce nom de colonne est déjà utilisé dans cette collection"
	ErrSchemaFailed       = "erreur lors de l'application du schéma en base"

	MsgCollectionCreated = "Collection créée."
	MsgCollectionUpdated = "Collection mise à jour."
	MsgCollectionDeleted = "Collection supprimée."
	MsgFieldAdded        = "Champ ajouté."
	MsgFieldRemoved      = "Champ supprimé."
)

// ValidFieldTypes liste les types de champs acceptés.
var ValidFieldTypes = map[string]bool{
	FieldTypeText: true, FieldTypeTextarea: true, FieldTypeRichtext: true,
	FieldTypeNumber: true, FieldTypeBoolean: true, FieldTypeDate: true,
	FieldTypeDatetime: true, FieldTypeEmail: true, FieldTypeURL: true,
	FieldTypePhone: true, FieldTypeColor: true, FieldTypeSelect: true,
	FieldTypeMultiselect: true, FieldTypeSlug: true, FieldTypeJSON: true,
	FieldTypePassword: true, FieldTypeFile: true, FieldTypeImage: true,
}
