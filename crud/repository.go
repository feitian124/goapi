package crud

type ID uint64

type Entity interface {
}

type Savable[E Entity] interface {
	// Save a given entity. Use the returned instance for further operations.
	Save(entity E) (E, error)

	// SaveAll saves all given entities.
	SaveAll(entities []E) ([]E, error)
}

type Readable[E Entity, ID comparable] interface {
	// GetByID retrieves an entity by its id.
	GetByID(id ID) (E, error)

	// ExistsById returns whether an entity with the given id exists.
	ExistsById(id ID) (bool, error)

	// GetAllById returns all instances of the type E with the given IDs.
	GetAllById(ids []ID) ([]E, error)

	// GetAll returns all instances of the type E
	GetAll() ([]E, error)

	// Count returns the number of entities available.
	Count() (int64, error)
}

type Deletable[E Entity, ID comparable] interface {
	// Delete deletes a given entity.
	Delete(entity E) error

	// DeleteById deletes the entity with the given id.
	DeleteById(id ID) error

	// DeleteAllById deletes all instances of the type E with the given IDs.
	DeleteAllById(ids []ID) error

	// DeleteAll deletes the given entities.
	DeleteAll(entities []E) error

	// Truncate deletes all entities by truncate.
	Truncate() error
}

type Repository[E Entity, ID comparable] interface {
	Savable[E, ID]
	Readable[E, ID]
	Deletable[E, ID]
}
