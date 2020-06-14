package domain

// Beer is a definition of a beer.
type Beer struct {
	ID      string
	Name    string
	Type    Type
	Brewer  string
	Country string
}

// Validate validates a beer.
func (b *Beer) Validate() error {
	if b.ID == "" {
		return NewValidationError("beer ID is empty")
	}
	return nil
}

// CreateBeerParams describes parameters for creating a beer.
type CreateBeerParams struct {
	Name    string
	Type    Type
	Brewer  string
	Country string
}

// Validate validates the CreateBeerParams.
func (b *CreateBeerParams) Validate() error {
	if b.Name == "" {
		return NewValidationError("beer name is empty")
	}
	return nil
}

// GetBeerParams describes parameters for getting a beer.
type GetBeerParams struct {
	ID string
}

// Validate validates the GetBeerParams.
func (b *GetBeerParams) Validate() error {
	if b.ID == "" {
		return NewValidationError("beer ID is empty")
	}
	return nil
}

// UpdateBeerParams describes parameters for updating a beer.
type UpdateBeerParams struct {
	ID      string
	Name    *string
	Type    *Type
	Brewer  *string
	Country *string
}

// Validate validates the UpdateBeerParams.
func (b *UpdateBeerParams) Validate() error {
	if b.ID == "" {
		return NewValidationError("beer ID is empty")
	}
	return nil
}

// DeleteBeerParams describes parameters for deleting a beer.
type DeleteBeerParams struct {
	ID string
}

// Validate validates the DeleteBeerParams.
func (b *DeleteBeerParams) Validate() error {
	if b.ID == "" {
		return NewValidationError("beer ID is empty")
	}
	return nil
}

// ListBeersParams describes parameters for listing beers.
type ListBeersParams struct {
	// Page is the page number of the beers.
	Page int
}

// Validate validates the ListBeersParams.
func (b *ListBeersParams) Validate() error {
	if b.Page < 1 {
		return NewValidationError("page number less than one")
	}
	return nil
}
