package nebula

func validateClassFilter(asset *AssetRef, class *AssetClass, subject *SubjectRef, op string) error {
	if class == nil {
		return nil
	}
	if asset != nil || subject != nil {
		return &ValidationError{Op: op, Field: "asset_class", Msg: "standalone class filter cannot be combined with asset or subject selector"}
	}
	return nil // The operation's generated query rules validate enum values.
}

// AssetRef selects one market asset: the pair (asset_id, asset_class) that
// /assets and /search return together. Carrying the pair in one value object
// makes it impossible to send a bare asset_id without its class.
type AssetRef struct {
	ID    string
	Class AssetClass
}

// SubjectRef selects one non-asset person, account, or topic by subject_id.
type SubjectRef struct {
	ID string
}

func (a *AssetRef) optionalPair(op string) error {
	if a == nil {
		return nil
	}
	if a.ID == "" {
		return &ValidationError{Op: op, Field: "asset_id", Msg: "asset_id must be non-empty when the asset selector is used"}
	}
	if a.Class == "" {
		return &ValidationError{Op: op, Field: "asset_class", Msg: "asset_class is required with asset_id"}
	}
	return nil // The operation's generated query rules validate enum values.
}
