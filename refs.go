package nebula

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

func (a *AssetRef) requirePair(op string) error {
	if a == nil || a.ID == "" {
		return &ValidationError{Op: op, Field: "asset_id", Msg: "asset selector is required"}
	}
	if a.Class == "" {
		return &ValidationError{Op: op, Field: "asset_class", Msg: "asset_class is required with asset_id"}
	}
	return a.Class.valid(op, "asset_class")
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
	return a.Class.valid(op, "asset_class")
}

func (a AssetClass) valid(op, field string) error {
	switch a {
	case AssetClassEquity, AssetClassCommodity, AssetClassIndex, AssetClassCrypto:
		return nil
	}
	return &ValidationError{Op: op, Field: field, Msg: "must be one of equity, commodity, index, crypto"}
}
