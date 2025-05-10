# Ent Mixins

This package provides reusable mixins for [Ent](https://entgo.io/) schemas.

## Available Mixins

### ID Mixin

The `IDMixin` provides a robust ID system for your entities with both UUID and optional human-readable identifiers:

```go
import "github.com/kopexa-grc/common/ent/mixins"

type User struct {
    ent.Schema
}

func (User) Mixin() []ent.Mixin {
    return []ent.Mixin{
        mixins.IDMixin{
            HumanIdentifierPrefix: "USR",  // Optional: Adds human-readable IDs like "USR-ABC123"
            SingleFieldIndex: true,        // Optional: Makes display_id unique
            DisplayIDLength: 6,            // Optional: Length of the display ID (default: 6)
        },
    }
}
```

#### Features

1. **Primary ID**
   - UUID-based primary identifier
   - Automatically generated
   - Globally unique
   - Immutable after creation

2. **Human-Readable Display ID** (Optional)
   - Configurable prefix (e.g., "USR" for users)
   - Fixed-length alphanumeric suffix
   - Automatically generated from UUID
   - Optional unique constraint
   - Collision-resistant based on length:
     - 6 chars: ~0.005% for 10,000 IDs
     - 8 chars: ~0.0001% for 1,000,000 IDs

3. **Performance Optimizations**
   - Built-in indexing for both ID types
   - Efficient ID generation using SHA256 and Base32
   - Optimized for high-volume operations

#### Example Output

```go
// With HumanIdentifierPrefix: "USR"
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "display_id": "USR-ABC123"
}

// Without HumanIdentifierPrefix
{
    "id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Audit Mixin

The `AuditMixin` provides automatic audit logging capabilities for your entities:

```go
import "github.com/kopexa-grc/common/ent/mixins"

type Document struct {
    ent.Schema
}

func (Document) Mixin() []ent.Mixin {
    return []ent.Mixin{
        mixins.AuditMixin{},
    }
}
```

#### Features

1. **Automatic Timestamps**
   - `created_at`: Immutable timestamp of creation
   - `updated_at`: Timestamp of the last update, automatically updated

2. **Actor Tracking**
   - `created_by`: Optional immutable ID of the creator
   - `updated_by`: Optional ID of the last updater

3. **Automatic Population**
   - Fields are automatically populated during mutations
   - Uses the actor from the context to set creator/updater IDs

#### Example Output

```go
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "created_at": "2023-01-01T12:00:00Z",
    "created_by": "user-123",
    "updated_at": "2023-01-02T15:30:00Z",
    "updated_by": "user-456"
}
```

## Usage

To use a mixin in your Ent schema:

1. Import the mixins package
2. Add the mixin to your schema's `Mixin()` method
3. Configure the mixin options as needed
4. The mixin's fields and behaviors will be automatically added to your schema

## Best Practices

- Use mixins to share common fields and behaviors across multiple schemas
- Keep mixins focused on a single responsibility
- Document the purpose and behavior of each mixin
- Consider the impact on database performance when adding mixins
- Choose appropriate display ID lengths based on your expected volume
- Use meaningful prefixes for different entity types 