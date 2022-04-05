package crud

type ID uint64

type Savable[T any] interface {
	// Save a given entity. Use the returned instance for further operations.
	Save(entity T) (T, error)

	// SaveAll saves all given entities.
	SaveAll(entities []T) ([]T, error)
}

type Readable[T any, ID comparable] interface {
	// GetByID retrieves an entity by its id.
	GetByID(id ID) (T, error)

	// ExistsById returns whether an entity with the given id exists.
	ExistsById(id ID) (bool, error)

	// GetAllById returns all instances of the type T with the given IDs.
	GetAllById(ids []ID) ([]T, error)

	// GetAll returns all instances of the type T
	GetAll() ([]T, error)

	// Count returns the number of entities available.
	Count() (int64, error)
}

type Deletable[T any, ID comparable] interface {
	// Delete deletes a given entity.
	Delete(entity T) error

	// DeleteById deletes the entity with the given id.
	DeleteById(id ID) error

	// DeleteAllById deletes all instances of the type T with the given IDs.
	DeleteAllById(ids []ID) error

	// DeleteAll deletes the given entities.
	DeleteAll(entities []T) error

	// Truncate deletes all entities by truncate.
	Truncate() error
}

type Repository[T any, ID comparable] interface {
	Savable[T, ID]
	Readable[T, ID]
	Deletable[T, ID]
}
