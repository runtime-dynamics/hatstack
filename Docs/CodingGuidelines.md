# Go Web Application Coding Guidelines

This document outlines the coding standards and architectural patterns for this Go web application bootstrap project.

**📖 Main Documentation:** See the [main README](../README.md) for project overview, quick start, and getting started guides.

**🐳 Deployment:** This project includes a production-ready `Dockerfile` that can be deployed directly to Google Cloud Run or any container platform. See the [Deployment section](../README.md#deployment) in the main README.

## Table of Contents
1. [Project Architecture](#project-architecture)
2. [Package Structure](#package-structure)
3. [Data Layer Patterns (Repository Pattern)](#data-layer-patterns-repository-pattern)
4. [Service Layer Patterns](#service-layer-patterns)
5. [Web Framework & Routing](#web-framework--routing)
6. [Handler Patterns](#handler-patterns)
7. [Component-Based UI (templ)](#component-based-ui-templ)
8. [Configuration Management](#configuration-management)
9. [Code Style Guidelines](#code-style-guidelines)
10. [Error Handling](#error-handling)
11. [Logging Standards](#logging-standards)

---

## Project Architecture

This is a Go-based web application built on the **H.A.T. Stack** (HTMX, Alpine.js, Templ) with a dual architecture supporting both JSON API and server-rendered HTML.

### Core Technologies

- **Web Framework**: Gin (github.com/gin-gonic/gin)
- **Database**: Google Cloud Datastore
- **Logging**: Zerolog (github.com/rs/zerolog)
- **UI Rendering**: Templ (github.com/a-h/templ)
- **Frontend Interactivity**: HTMX + Alpine.js
- **Styling**: TailwindCSS

### Key Architectural Layers

```
┌─────────────────────────────────────────────────────────────┐
│                    Client Layer                             │
│  - HTMX (Server Interactions)                              │
│  - Alpine.js (Client-Side State)                           │
│  - TailwindCSS (Styling)                                   │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                    Web Layer (web/)                         │
│  - Route Registration (API + App)                          │
│  - Static File Serving                                     │
└─────────────────────────────────────────────────────────────┘
                            ↓
        ┌───────────────────┴───────────────────┐
        ↓                                       ↓
┌──────────────────────┐          ┌──────────────────────────┐
│  API Handlers        │          │  App Handlers            │
│  (web/api/)          │          │  (web/app/)              │
│  - JSON Responses    │          │  - HTML Responses        │
│  - /api/* routes     │          │  - / and other routes    │
└──────────────────────┘          └──────────────────────────┘
        ↓                                       ↓
        └───────────────────┬───────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                Services Layer (services/)                    │
│  - Business Logic                                           │
│  - Validation                                               │
│  - Orchestration                                            │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│              Data Layer (data/)                             │
│  - Repository Pattern                                       │
│  - Data Models                                              │
│  - Datastore Interactions                                   │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│                Google Cloud Datastore                        │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                Views Layer (views/)                          │
│  - Templ Components                                         │
│  - layouts/ - Base page layouts                             │
│  - components/ - Reusable UI components                     │
│  - pages/ - Full page templates                             │
└─────────────────────────────────────────────────────────────┘
```

### Request Flow

**API Request Flow (JSON):**

```
Client → API Handler → Service → Repository → Datastore
                ↓
         JSON Response
```

**Web Request Flow (HTML):**

```
Client → Web Handler → Service → Repository → Datastore
                ↓
         Templ Component → HTML Response
```

---

## Package Structure

### `/config` - Configuration Package

The `config` package provides centralized access to application configuration loaded from environment variables.

**Pattern:**
- Use singleton pattern with `sync.RWMutex` for thread-safe access
- Load configuration once at startup via `LoadConfig()`
- Access configuration via `config.Get()`

**Example:**
```go
cfg := config.Get()
projectID := cfg.GoogleProjectID
// Use configuration values
```

### `/data` - Data Access Layer (Repository Pattern)

The `data` package encapsulates all database interactions using the **Repository Pattern**.

**Key Components:**

- `data.go`: Base datastore client and utilities
- `*_repository.go`: Domain-specific repositories (e.g., `user_repository.go`)
- Data models: Struct definitions for entities
- `Cli()`: Returns singleton datastore client
- `IsNotFound(err)`: Helper to check for entity not found errors

**Pattern:**

```go
// Create a repository
repo := data.NewMyEntityRepository()

// Use repository methods
entity, err := repo.GetByID(ctx, id)
```

### `/services` - Business Logic Layer

The `services` package contains all business logic, validation, and orchestration.

**Key Components:**

- `service.go`: Base service interface and implementation
- `*_service.go`: Domain-specific services (e.g., `user_service.go`)

**Responsibilities:**

- Business logic and validation
- Orchestrating multiple repository calls
- Authorization checks
- Data transformation

**Pattern:**

```go
// Create a service
myService := services.NewMyEntityService(ctx)

// Use service methods
err := myService.Create(entity, ownerID)
```

### `/web/api` - API HTTP Request Handlers (JSON)

Contains all JSON API endpoint handlers organized by domain. These handlers serve `/api/*` routes.

**Organization:**

- `routes.go`: Central API route registration
- Domain-specific handler files as needed

**Responsibilities:**

- Parse HTTP requests
- Call service layer
- Return JSON responses

### `/web/app` - App HTTP Request Handlers (HTML)

Contains all HTML-serving endpoint handlers for the web UI. These handlers serve `/` and other non-API routes.

**Organization:**

- `routes.go`: Central web route registration
- Domain-specific handler files as needed
- `seo/`: SEO-related handlers (robots.txt, sitemap.xml)

**Responsibilities:**

- Parse HTTP requests
- Call service layer
- Render templ components to HTML
- Handle HTMX requests (full pages or fragments)

### `/views` - Templ UI Components

Contains all UI templates using the Templ library.

**Organization:**

- `layouts/`: Base page layouts (e.g., `base.templ`)
- `components/`: Reusable UI components (e.g., `button.templ`, `card.templ`)
- `pages/`: Full page templates (e.g., `home.templ`)

### `/web` - Web Server Entry Point

Main web server initialization and route registration.

**Key Functions:**

- `Start(r *gin.Engine)`: Initializes the web server
- `HandleRoutes(r *gin.Engine, staticDir string)`: Registers all routes (API + App)
- `GetStaticFiles(staticDir string)`: Discovers static files for production serving

---

## Data Layer Patterns (Repository Pattern)

The data layer uses the **Repository Pattern** to decouple business logic from data access. All database operations are performed through repository structs, not directly on models.

### Creating a New Data Model and Repository

#### 1. Define the Data Model Struct

```go
// In data/myentity.go
package data

type MyEntity struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	OwnerID   string    `json:"owner_id"`
    CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

**Conventions:**
- Use `ID` field for the datastore key name
- Use JSON tags for API serialization
- Use datastore tags when needed (e.g., `datastore:",noindex"`)

#### 2. Create a Repository

```go
// In data/myentity_repository.go
package data

import (

"context"

"cloud.google.com/go/datastore"
"github.com/rs/zerolog/log"
)

type MyEntityRepository struct {
	*BaseRepository
}

func NewMyEntityRepository() *MyEntityRepository {
	return &MyEntityRepository{
		BaseRepository: NewBaseRepository(),
    }
}

// GetByID retrieves an entity by its ID
func (r *MyEntityRepository) GetByID(ctx context.Context, id string) (*MyEntity, error) {
	if id == "" {
		return nil, datastore.ErrNoSuchEntity
	}
    key := datastore.NameKey("MyEntity", id, nil)
    entity := &MyEntity{}
	if err := r.Client().Get(ctx, key, entity); err != nil {
		log.Error().Err(err).Msgf("failed to get MyEntity by id: %s", id)
        return nil, err
    }
	entity.ID = id
    return entity, nil
}

// Create creates a new entity
func (r *MyEntityRepository) Create(ctx context.Context, entity *MyEntity) error {
	key := datastore.NameKey("MyEntity", entity.ID, nil)
	log.Debug().Msgf("Creating MyEntity with id: %s", entity.ID)
	if _, err := r.Client().Put(ctx, key, entity); err != nil {
		log.Error().Err(err).Msg("failed to create MyEntity")
		return err
    }
	return nil
}

// Update updates an existing entity
func (r *MyEntityRepository) Update(ctx context.Context, entity *MyEntity) error {
	key := datastore.NameKey("MyEntity", entity.ID, nil)
	log.Debug().Msgf("Updating MyEntity with id: %s", entity.ID)
	if _, err := r.Client().Put(ctx, key, entity); err != nil {
        log.Error().Err(err).Msg("failed to update MyEntity")
        return err
    }
    return nil
}

// Delete removes an entity
func (r *MyEntityRepository) Delete(ctx context.Context, id string) error {
	key := datastore.NameKey("MyEntity", id, nil)
	log.Debug().Msgf("Deleting MyEntity with id: %s", id)
	if err := r.Client().Delete(ctx, key); err != nil {
        log.Error().Err(err).Msg("failed to delete MyEntity")
        return err
    }
    return nil
}

// ListByOwner retrieves all entities for an owner
func (r *MyEntityRepository) ListByOwner(ctx context.Context, ownerID string) ([]MyEntity, error) {
    q := datastore.NewQuery("MyEntity").FilterField("OwnerID", "=", ownerID)
    var results []MyEntity
	_, err := r.Client().GetAll(ctx, q, &results)
    if err != nil {
        log.Error().Err(err).Msg("failed to list MyEntity by OwnerID")
        return nil, err
    }
    return results, nil
}
```

#### 3. Query Operators and FilterField

All datastore queries MUST use `FilterField` (not the deprecated `Filter` method).

**FilterField Signature:**
```go
FilterField(fieldName, operator string, value interface{}) *Query
```

**Supported Operators:**
- `"="` - Equal to
- `"!="` - Not equal to
- `">"` - Greater than
- `"<"` - Less than
- `">="` - Greater than or equal to
- `"<="` - Less than or equal to
- `"in"` - Value in array (requires `[]interface{}` as value)
- `"not-in"` - Value not in array (requires `[]interface{}` as value)

**Examples:**
```go
// Simple equality
q := datastore.NewQuery("MyEntity").FilterField("Status", "=", "active")

// Comparison operators
q := datastore.NewQuery("MyEntity").FilterField("Score", ">=", 100)

// Multiple filters (AND'ed together)
q := datastore.NewQuery("MyEntity").
    FilterField("OwnerID", "=", ownerID).
    FilterField("Status", "=", "active").
    FilterField("CreatedAt", ">=", startTime)

// IN operator (note the []interface{} type)
q := datastore.NewQuery("MyEntity").
    FilterField("Category", "in", []interface{}{"cat1", "cat2", "cat3"})

// NOT-IN operator
q := datastore.NewQuery("MyEntity").
    FilterField("Status", "not-in", []interface{}{"deleted", "archived"})

// Field names with special characters (use strconv.Quote or %q)
import "strconv"
fieldName := strconv.Quote("field with spaces")
q := datastore.NewQuery("MyEntity").FilterField(fieldName, "=", value)
```

**Important Notes:**
- Multiple `FilterField` calls are AND'ed together
- For OR operations, you need multiple queries
- Field names with spaces, quotes, or operators should be quoted using `strconv.Quote` or `fmt.Sprintf("%q", fieldName)`
- The `in` and `not-in` operators require `[]interface{}` type for the value parameter

#### 4. Handling Map Fields (Advanced Pattern)

When you need to store map data in Datastore (which doesn't natively support maps), use the dual-representation pattern from `CharacterData`:

```go
type MyEntity struct {
    ID string
    // JSON-facing map (ignored by datastore)
    MetadataMap map[string]string `json:"metadata" datastore:"-"`
    // Datastore-facing slice (hidden from JSON)
    Metadata []KeyValuePair `json:"-"`
}

type KeyValuePair struct {
    Key   string `datastore:",noindex"`
    Value string `datastore:",noindex"`
}

// Convert maps to slices before save
func (e *MyEntity) populateSlicesFromMaps() {
    e.Metadata = e.Metadata[:0]
    for k, v := range e.MetadataMap {
        e.Metadata = append(e.Metadata, KeyValuePair{Key: k, Value: v})
    }
}

// Convert slices to maps after load
func (e *MyEntity) populateMapsFromSlices() {
    if e.MetadataMap == nil {
        e.MetadataMap = make(map[string]string)
    }
    for k := range e.MetadataMap {
        delete(e.MetadataMap, k)
    }
    for _, kv := range e.Metadata {
        e.MetadataMap[kv.Key] = kv.Value
    }
}

// Call in Upsert before saving
func (e *MyEntity) Upsert() error {
    e.populateSlicesFromMaps()
    // ... rest of upsert logic
}

// Call in GetByID after loading
func (e *MyEntity) GetByID() error {
    // ... get logic
    e.populateMapsFromSlices()
    return nil
}
```

---

## Service Layer Patterns

The service layer contains all business logic, validation, and orchestration. Services sit between handlers and repositories.

### Creating a New Service

```go
// In services/myentity_service.go
package services

import (
	"context"
	"errors"
	"time"

	"runtime-dynamics/data"
)

type MyEntityService struct {
	*BaseService
	repo *data.MyEntityRepository
}

func NewMyEntityService(ctx context.Context) *MyEntityService {
	return &MyEntityService{
		BaseService: NewBaseService(ctx),
		repo:        data.NewMyEntityRepository(),
	}
}

// GetByID retrieves an entity with validation
func (s *MyEntityService) GetByID(id string) (*data.MyEntity, error) {
	if len(id) < 1 {
		return nil, errors.New("entity id cannot be empty")
	}
	return s.repo.GetByID(s.ctx, id)
}

// Create creates a new entity with business logic
func (s *MyEntityService) Create(entity *data.MyEntity, ownerID string) error {
	// Validation
	if entity == nil {
		return errors.New("entity cannot be nil")
	}
	if len(entity.Name) < 3 {
		return errors.New("entity name must be at least 3 characters")
	}

	// Business logic
	entity.OwnerID = ownerID
	entity.CreatedAt = time.Now().UTC()
	entity.UpdatedAt = time.Now().UTC()

	// Persist
	return s.repo.Create(s.ctx, entity)
}

// Update updates an entity with ownership validation
func (s *MyEntityService) Update(entity *data.MyEntity, ownerID string) error {
	// Fetch existing entity
	existing, err := s.repo.GetByID(s.ctx, entity.ID)
	if err != nil {
		return err
	}

	// Verify ownership
	if existing.OwnerID != ownerID {
		return errors.New("not authorized to update this entity")
	}

	// Update timestamp
	entity.UpdatedAt = time.Now().UTC()

	// Persist
	return s.repo.Update(s.ctx, entity)
}

// Delete deletes an entity with ownership validation
func (s *MyEntityService) Delete(id string, ownerID string) error {
	// Fetch entity to verify ownership
	entity, err := s.repo.GetByID(s.ctx, id)
	if err != nil {
		return err
	}

	// Verify ownership
	if entity.OwnerID != ownerID {
		return errors.New("not authorized to delete this entity")
	}

	// Delete
	return s.repo.Delete(s.ctx, id)
}

// ListByOwner retrieves all entities for an owner
func (s *MyEntityService) ListByOwner(ownerID string) ([]data.MyEntity, error) {
	if len(ownerID) < 1 {
		return nil, errors.New("owner id cannot be empty")
	}
	return s.repo.ListByOwner(s.ctx, ownerID)
}
```

### Service Layer Best Practices

1. **No Direct Database Access**: Services must never call `data.Cli()` directly. Always use repositories.
2. **Business Logic Only**: All validation, authorization, and business rules go in services.
3. **Context Propagation**: Always pass context from handlers to services to repositories.
4. **Error Handling**: Return descriptive errors that can be shown to users.
5. **Dependency Injection**: Services receive repositories through their constructors.

---

## Web Framework & Routing

### Entry Point: `web/web.go`

The `web` package provides the main entry point for the web server.

**Key Functions:**

#### `Start(r *gin.Engine) *gin.Engine`
Main initialization function that:
1. Calls `HandleRoutes()` to register all routes
2. Returns the configured engine

#### `HandleRoutes(r *gin.Engine, staticDir string) *gin.Engine`
Registers all application routes:

1. Calls `api.RegisterRoutes(r)` for API routes (JSON, `/api/*`)
2. Calls `app.RegisterWebRoutes(r)` for web routes (HTML, `/` and others)
3. Serves static files (development: directories, production: individual files)

### Route Registration

#### API Routes: `web/api/routes.go`

All JSON API routes are registered in `web/api/routes.go` via the `RegisterRoutes()` function.

```go
func RegisterRoutes(r *gin.Engine) {
	// All API routes are prefixed with /api
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/health", HealthCheckHandler)
		// Add more API routes as needed
	}
}
```

#### App Routes: `web/app/routes.go`

All HTML web routes are registered in `web/app/routes.go` via the `RegisterWebRoutes()` function.

```go
func RegisterWebRoutes(r *gin.Engine) {
	// SEO and crawler files
	r.GET("/robots.txt", seo.RobotsTxtHandler)
	r.GET("/sitemap.xml", seo.SitemapXMLHandler)

	// Homepage (public)
	r.GET("/", HomePageHandler)
	
	// Add more web routes as needed
}
```

### Utility Functions for API Handlers

API handlers in `web/api/routes.go` include helper functions for JSON responses:

```go
// Render error with status code
renderError(c, err, http.StatusBadRequest, "error message")

// Render success response
renderSuccess(c)

// Render content (data)
renderFinalContent(c, data, "key", err)
```

**Note:** App handlers do NOT use these functions. They render templ components directly.

---

## Handler Patterns

The application has **two types of handlers**: API handlers (JSON) and Web handlers (HTML).

### API Handler Pattern (JSON)

API handlers serve JSON responses for `/api/*` routes. They are located in the `web/api/` package.

**Standard API Handler Pattern:**

```go
package api

import (
	"net/http"

	"runtime-dynamics/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func MyAPIHandler(c *gin.Context) {
    // 1. Extract and validate parameters
    id := c.Param("id")
    if len(id) < 1 {
        renderError(c, nil, http.StatusBadRequest, "id cannot be empty")
        return
    }
    
    // 2. Bind JSON body if needed
	var requestData MyRequestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		renderError(c, err, http.StatusBadRequest, "bad json format")
		return
	}

	// 3. Call service layer (NO direct data access!)
	myService := services.NewMyEntityService(c.Request.Context())
	result, err := myService.GetByID(id)
	if err != nil {
		log.Error().Err(err).Msg("failed to get entity")
		renderError(c, err, http.StatusInternalServerError, "failed to get entity")
        return
    }

	// 4. Return JSON response
	renderFinalContent(c, result, "entity", nil)
}
```

**API Handler Rules:**

- ✅ Call service layer methods
- ✅ Return JSON using helper functions in `web/api/routes.go`
- ❌ Never call repositories or `data.Cli()` directly
- ❌ Never render HTML

### App Handler Pattern (HTML/Templ)

App handlers serve HTML responses for `/` and other non-API routes. They are located in the `web/app/` package.

**Standard App Handler Pattern:**

```go
package app

import (
	"net/http"

	"runtime-dynamics/services"
	"runtime-dynamics/views/pages"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func MyPageHandler(c *gin.Context) {
	// 1. Extract parameters
	id := c.Param("id")

	// 2. Call service layer (same as API handlers)
	myService := services.NewMyEntityService(c.Request.Context())
	entity, err := myService.GetByID(id)
	if err != nil {
		log.Error().Err(err).Msg("failed to get entity")
		c.String(http.StatusNotFound, "entity not found")
		return
	}

	// 3. Render templ component
	component := pages.MyEntityPage(entity)
	component.Render(c.Request.Context(), c.Writer)
}
```

**HTMX Fragment Handler Pattern:**

For HTMX requests that return HTML fragments (not full pages):

```go
func CreateEntityFragmentHandler(c *gin.Context) {
	// Parse form data
	name := c.PostForm("name")

	// Call service
	myService := services.NewMyEntityService(c.Request.Context())
	entity, err := myService.Create(name)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// Return HTML fragment (not full page)
	component := components.EntityCard(entity)
	component.Render(c.Request.Context(), c.Writer)
}
```

**App Handler Rules:**

- ✅ Call service layer methods
- ✅ Render templ components to HTML
- ✅ Handle HTMX requests (full pages or fragments)
- ❌ Never call repositories or `data.Cli()` directly
- ❌ Never use API helper functions (renderError, etc.)
- ❌ Never call `/api/*` endpoints from templ components

---

## Component-Based UI (templ)

The application uses **Templ** for server-rendered HTML components. Templ files (`.templ`) are compiled to Go code.

### Directory Structure

```
views/
├── layouts/          # Base page layouts
│   └── base.templ    # Main layout with <html>, <head>, <body>
├── components/       # Reusable UI components
│   ├── button.templ
│   ├── card.templ
│   └── modal.templ
└── pages/           # Full page templates
    └── home.templ
```

### Creating a New Templ Component

#### 1. Layout Component

```templ
// In views/layouts/base.templ
package layouts

templ Base(title string) {
    <!DOCTYPE html>
    <html lang="en">
        <head>
            <meta charset="UTF-8"/>
            <title>{ title }</title>
            <script src="https://cdn.tailwindcss.com"></script>
            <script src="https://unpkg.com/htmx.org@2.0.4"></script>
            <script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"></script>
        </head>
        <body class="bg-gray-900 text-gray-100">
            { children... }
        </body>
    </html>
}
```

#### 2. Reusable Component

```templ
// In views/components/button.templ
package components

templ Button(text string, variant string, attrs templ.Attributes) {
    <button
        { attrs... }
        class={ "px-4 py-2 rounded-lg", variant == "primary" ? "bg-blue-600" : "bg-gray-600" }
    >
        { text }
    </button>
}
```

#### 3. Page Component

```templ
// In views/pages/mypage.templ
package pages

import (
    "runtime-dynamics/views/layouts"
    "runtime-dynamics/views/components"
    "runtime-dynamics/data"
)

templ MyPage(entity *data.MyEntity) {
    @layouts.Base("My Page") {
        <div class="container mx-auto px-4 py-8">
            <h1 class="text-3xl font-bold mb-4">{ entity.Name }</h1>
            @components.Button("Click Me", "primary", templ.Attributes{
                "hx-post": "/action",
                "hx-target": "#result",
            })
        </div>
    }
}
```

### Using HTMX in Templ Components

HTMX attributes are added directly to HTML elements:

```templ
templ EntityCard(entity data.MyEntity) {
    <div class="card">
        <h3>{ entity.Name }</h3>
        <button
            hx-delete={ fmt.Sprintf("/entities/%s", entity.ID) }
            hx-confirm="Are you sure?"
            hx-target="closest div.card"
            hx-swap="outerHTML"
        >
            Delete
        </button>
    </div>
}
```

### Using Alpine.js in Templ Components

Alpine.js directives for client-side state:

```templ
templ Modal(id string, title string) {
    <div
        x-data="{ open: false }"
        x-show="open"
        @modal-open.window="if ($event.detail === '{ id }') open = true"
        @modal-close.window="if ($event.detail === '{ id }') open = false"
    >
        <div class="modal-content">
            <h3>{ title }</h3>
            { children... }
        </div>
    </div>
}
```

### Generating Go Code from Templ

After creating or modifying `.templ` files, generate Go code:

```bash
# Generate once
templ generate

# Watch for changes (development)
templ generate --watch
```

### Templ Best Practices

1. **Component Composition**: Build complex UIs from small, reusable components
2. **Props**: Pass data as function parameters
3. **Children**: Use `{ children... }` for nested content
4. **Attributes**: Use `templ.Attributes` for dynamic HTML attributes
5. **HTMX Integration**: Use `hx-*` attributes for server interactions
6. **Alpine.js Integration**: Use `x-*` attributes for client-side state
7. **Type Safety**: Leverage Go's type system for component props

---

## Configuration Management

### Accessing Configuration

**Always use `config.Get()` to access configuration:**

```go
cfg := config.Get()
if cfg == nil {
    // Handle error
}

// Access fields
projectID := cfg.GoogleProjectID
datastoreName := cfg.DataStoreName
```

### Configuration Fields

The `AppConfig` struct contains environment-based configuration values. Add fields as needed for your application:

- `DataStoreName` - Datastore database name
- `GoogleProjectID` - GCP project ID
- Add additional configuration fields as your application requires

### Thread Safety

The config package uses `sync.RWMutex` for thread-safe access:
- `LoadConfig()` uses write lock
- `Get()` uses read lock

---

## Code Style Guidelines

### General Go Conventions

1. **Follow standard Go formatting**: Use `gofmt` or `goimports`
2. **Package naming**: Short, lowercase, single-word names
3. **Exported vs unexported**: Use PascalCase for exported, camelCase for unexported
4. **Error handling**: Always check errors, never ignore them

### Naming Conventions

#### Variables
```go
// Good
userID := "123"
accountData := &data.Account{}
cfg := config.Get()

// Avoid
user_id := "123"  // No snake_case
ud := &data.Account{}  // Too abbreviated in broad scope
```

#### Functions
```go
// Good
func GetAccountByID(id string) (*Account, error)
func (a *Account) Update() error
func CreateCharacter(char *CharacterData) error

// Avoid
func get_account(id string) (*Account, error)  // No snake_case
func UpdateAcc(a *Account) error  // Avoid abbreviations
```

#### Constants
```go
// Good
const MaxRetries = 3
const DefaultTimeout = 30 * time.Second

// Avoid
const max_retries = 3  // No snake_case
```

### Import Organization

Group imports in this order:
1. Standard library
2. External dependencies
3. Internal packages

```go
import (
    "context"
    "errors"
    "time"
    
    "cloud.google.com/go/datastore"
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
    
    "github.com/Nitecon/starxapi/config"
    "github.com/Nitecon/starxapi/data"
    "github.com/Nitecon/starxapi/helpers"
)
```

### Function Length

- Keep functions focused and concise
- Extract complex logic into helper functions
- Aim for functions under 50 lines when possible

### Comments

```go
// Package-level comment
// Describes the package purpose

// Exported function comment
// Describes what the function does, parameters, and return values
func GetAccountByID(id string) (*Account, error) {
    // Implementation comments for complex logic
}

// Unexported functions should have comments if non-obvious
func extractBearerToken(gc *gin.Context) (string, bool) {
    // ...
}
```

---

## Error Handling

### Service Layer Errors

```go
// Services return descriptive errors
myService := services.NewMyEntityService(ctx)
entity, err := myService.GetByID(id)
if err != nil {
log.Error().Err(err).Msg("failed to get entity")
    return err
}
```

### Repository Layer Errors

```go
// Repositories log and return errors
repo := data.NewMyEntityRepository()
entity, err := repo.GetByID(ctx, id)
if err != nil {
    if data.IsNotFound(err) {
        // Handle not found case
		return nil, errors.New("entity not found")
    }
    return nil, err
}
```

### API Handler Errors

```go
// Use helpers for JSON error responses
myService := services.NewMyEntityService(ctx.Request.Context())
if err := myService.Create(entity, oid); err != nil {
log.Error().Err(err).Msg("failed to create entity")
helpers.RenderError(ctx, http.StatusInternalServerError, err.Error())
    return
}
```

### Web Handler Errors

```go
// Return appropriate HTTP responses or redirects
myService := services.NewMyEntityService(c.Request.Context())
entity, err := myService.GetByID(id)
if err != nil {
log.Error().Err(err).Msg("failed to get entity")
c.String(http.StatusNotFound, "entity not found")
    return
}
```

### Error Messages

- **User-facing**: Generic, don't expose internals
    - Good: `"failed to save entity"`
    - Bad: `"datastore put failed: connection timeout"`

- **Logs**: Detailed, include context
    - Good: `log.Error().Err(err).Msgf("failed to upsert character %s", id)`
    - Bad: `log.Error().Msg("error")`

---

## Logging Standards

### Use Zerolog

The project uses `github.com/rs/zerolog/log` for structured logging.

**Log Level Configuration:**
- Set via `DEBUG` environment variable
- `DEBUG=true` → Debug level (development)
- `DEBUG=""` or unset → Info level (production)

### Log Levels - When to Use Each

#### DEBUG Level (`log.Debug()`)
**Use for**: Development-only detailed information that helps trace execution flow.

**Examples:**
```go
// Repository layer - detailed operations
log.Debug().Msgf("ListByPlayerID: Got %d results and %d keys", len(results), len(keys))
log.Debug().Msgf("GetByID: Using numeric key ID=%d", intID)

// Service layer - method entry/parameters
log.Debug().Msgf("UpsertCharacter called, CharacterId: %s, Name: %s, ownerID: %s", char.CharacterId, char.Name, ownerID)

// Handler layer - request details
log.Debug().Msgf("Request body: %+v", requestData)
```

**DON'T use for:**
- ❌ Production-critical information
- ❌ Error conditions
- ❌ User actions or state changes

#### INFO Level (`log.Info()`)
**Use for**: Important application events and successful operations.

**Examples:**
```go
// Startup/shutdown
log.Info().Msg("Starting StarXAPI (Version: source)")
log.Info().Msg("Listening on :8080")

// Successful operations
log.Info().Msgf("Successfully created character: ID=%s, Name=%s", char.ID, char.Name)
log.Info().Msgf("User authenticated: %s", userID)

// State changes
log.Info().Msgf("Character deleted: ID=%s", id)
log.Info().Msgf("Server registered: %s", serverID)

// Service layer - operation results
log.Info().Msgf("ListCharactersByPlayerID: found %d characters for playerID: %s", len(chars), playerID)
```

**DON'T use for:**
- ❌ Every request (too verbose)
- ❌ Internal implementation details
- ❌ Temporary debugging

#### WARN Level (`log.Warn()`)
**Use for**: Unexpected situations that are handled but may indicate issues.

**Examples:**
```go
// Missing optional data
log.Warn().Msg("No characters found for player, returning empty list")

// Recoverable errors
log.Warn().Msgf("Failed to parse optional field: %s", fieldName)

// Security concerns
log.Warn().Msgf("Multiple failed login attempts from: %s", ipAddress)

// Deprecated usage
log.Warn().Msg("Using deprecated API endpoint")
```

**DON'T use for:**
- ❌ Expected conditions (use Debug)
- ❌ Fatal errors (use Error)
- ❌ Normal validation failures

#### ERROR Level (`log.Error()`)
**Use for**: Failures that prevent operations from completing.

**ALWAYS include `.Err(err)` when you have an error object.**

**Examples:**
```go
// Database failures
log.Error().Err(err).Msgf("Failed to get character by ID: %s", id)
log.Error().Err(err).Msg("Failed to upsert character")

// Service failures
log.Error().Err(err).Msgf("Could not authenticate user: %s", userID)
log.Error().Err(err).Msg("Failed to process payment")

// External service failures
log.Error().Err(err).Msg("Failed to connect to game server")
```

**DON'T use for:**
- ❌ Validation errors (use Warn or Info)
- ❌ "Not found" errors (use Warn or Debug)
- ❌ Expected business logic failures

#### FATAL Level (`log.Fatal()`)
**Use for**: Critical errors that require application shutdown.

**Examples:**
```go
// Startup failures
log.Fatal().Err(err).Msg("Failed to create datastore client")
log.Fatal().Err(err).Msg("Failed to load configuration")
log.Fatal().Err(err).Msg("Failed to bind to port")
```

**DON'T use for:**
- ❌ Request handling errors
- ❌ Recoverable failures
- ❌ Anything during normal operation

### Logging Patterns

#### Include Error Objects
```go
// ALWAYS use .Err() when you have an error
log.Error().Err(err).Msg("failed to update entity")

// NOT this:
log.Error().Msgf("failed to update entity: %v", err)  // ❌ Wrong
```

#### Include Context Fields
```go
// Use structured fields for searchable data
log.Info().
    Str("user_id", userID).
    Str("character_id", charID).
    Int("level", level).
    Msg("character created")
```

#### Use Msgf for Dynamic Messages
```go
// Good
log.Info().Msgf("Processing request for user: %s", userID)

// Also good for structured data
log.Error().Err(err).Msgf("Failed to delete character ID: %s for owner: %s", id, ownerID)
```

### Layer-Specific Logging Guidelines

#### Repository Layer
```go
// DEBUG: Detailed operations
log.Debug().Msgf("Querying characters with filter: %+v", filter)

// INFO: Successful operations with counts
log.Info().Msgf("Successfully upserted character with CharacterId: %s", char.CharacterId)

// ERROR: Database failures
log.Error().Err(err).Msg("failed to list characters by PlayerId")
```

#### Service Layer
```go
// DEBUG: Method entry and parameters
log.Debug().Msgf("DeleteCharacterByID called with datastoreID: %s, ownerID: %s", datastoreID, ownerID)

// INFO: Successful business operations
log.Info().Msgf("Successfully deleted character: ID=%s", char.ID)

// WARN: Business rule violations
log.Warn().Msgf("Ownership verification failed: character PlayerId=%s, requested ownerID=%s", char.PlayerId, ownerID)

// ERROR: Operation failures
log.Error().Err(err).Msgf("Failed to upsert character, CharacterId: %s", char.CharacterId)
```

#### Handler Layer (API & Web)
```go
// INFO: Request received (important endpoints only)
log.Info().Msgf("Delete character request received for ID: %s", id)

// WARN: Authentication/authorization failures
log.Warn().Msg("Delete character failed: not authenticated")

// ERROR: Request processing failures
log.Error().Err(err).Msgf("Failed to delete character ID: %s for owner: %s", id, oid)
```

### What to Log

**DO log:**
- ✅ All errors with `.Err(err)` and context
- ✅ Authentication events (login, logout, failures)
- ✅ Authorization failures (access denied)
- ✅ Important state changes (create, update, delete)
- ✅ WebSocket connections/disconnections
- ✅ External service calls (start and result)
- ✅ Startup/shutdown events
- ✅ Configuration loading

**DON'T log:**
- ❌ Sensitive data (passwords, tokens, API keys, session IDs)
- ❌ Personal information (emails, names) without anonymization
- ❌ Full request/response bodies in production
- ❌ Every single request (use DEBUG level)
- ❌ Redundant information already in error objects

### Production vs Development

**Development (DEBUG=true):**
- All log levels visible
- Detailed parameter logging
- Request/response details
- Execution flow tracing

**Production (DEBUG not set):**
- INFO and above only
- Minimal parameter logging
- No request/response bodies
- Focus on errors and important events

### Examples by Scenario

#### Successful Operation
```go
// Service layer
log.Info().Msgf("Successfully created character: ID=%s, CharacterId=%s, Name=%s", char.ID, char.CharacterId, char.Name)

// Handler layer
log.Info().Msgf("Character created for user: %s", ownerID)
```

#### Failed Operation
```go
// Repository layer
log.Error().Err(err).Msgf("Failed to get character by ID: %s", id)

// Service layer
log.Error().Err(err).Msgf("Could not delete character: ID=%s", id)

// Handler layer
log.Error().Err(err).Msgf("Failed to process delete request for character: %s", id)
```

#### Validation Failure
```go
// Service layer
log.Warn().Msgf("Validation failed: empty owner ID")

// Handler layer
log.Warn().Msgf("Invalid request: missing required field 'name'")
```

#### Not Found
```go
// Service layer
log.Debug().Msgf("Character not found by CharacterId: %s", characterID)

// Handler layer (if important)
log.Info().Msgf("Character not found: %s", id)
```

---

## Best Practices Summary

### Repository Layer (Data Access)

- ✅ Use Repository Pattern for all data access
- ✅ Create `*Repository` structs that embed `BaseRepository`
- ✅ All repository methods accept `context.Context` as first parameter
- ✅ Use `datastore.NameKey()` for entity keys
- ✅ Log all database operations at Debug level
- ✅ Return `datastore.ErrNoSuchEntity` for not-found cases
- ❌ Never call repositories directly from handlers

### Service Layer (Business Logic)

- ✅ All business logic goes in services
- ✅ Services orchestrate repository calls
- ✅ Perform validation and authorization in services
- ✅ Services accept context and return descriptive errors
- ✅ Use dependency injection for repositories
- ❌ Never call `data.Cli()` directly from services
- ❌ Never put business logic in handlers

### API Handlers (JSON)

- ✅ Register all API routes in `web/api/routes.go` under `/api/*`
- ✅ Call service layer methods only
- ✅ Use helper functions in `web/api/routes.go` for JSON responses
- ✅ Validate input before processing
- ✅ Use `c.ShouldBindJSON()` for request body parsing
- ❌ Never call repositories or `data.Cli()` directly
- ❌ Never render HTML from API handlers

### App Handlers (HTML)

- ✅ Register all web routes in `web/app/routes.go` under `/` and other routes
- ✅ Call service layer methods only
- ✅ Render templ components to HTML
- ✅ Handle HTMX requests (full pages or fragments)
- ❌ Never call repositories or `data.Cli()` directly
- ❌ Never use API helper functions (renderError, etc.)
- ❌ Never call `/api/*` endpoints from templ components

### Templ Components

- ✅ Organize into `layouts/`, `components/`, and `pages/`
- ✅ Use component composition for complex UIs
- ✅ Use HTMX attributes for server interactions
- ✅ Use Alpine.js for client-side state
- ✅ Run `templ generate` after modifying `.templ` files
- ❌ Never call `/api/*` JSON endpoints from templ

### Configuration
- ✅ Access config via `config.Get()`
- ✅ Never hardcode configuration values
- ✅ Use environment variables for all config

### General
- ✅ Use structured logging with zerolog
- ✅ Handle all errors explicitly
- ✅ Follow Go naming conventions
- ✅ Keep functions focused and concise
- ✅ Write comments for exported functions
- ✅ Use meaningful variable names

---

## Quick Reference

### Creating a New API Endpoint (JSON)

1. **Create repository** in `data/myentity_repository.go`
2. **Create service** in `services/myentity_service.go`
3. **Create handler** in `web/api/routes.go` or separate file
4. **Register route** in `web/api/routes.go` under `/api/*`
5. **Handler calls service**, service calls repository
6. **Return JSON** using helper functions in `web/api/routes.go`

### Creating a New App Endpoint (HTML)

1. **Create repository** in `data/myentity_repository.go` (if not exists)
2. **Create service** in `services/myentity_service.go` (if not exists)
3. **Create templ component** in `views/pages/mypage.templ`
4. **Create handler** in `web/app/` directory
5. **Register route** in `web/app/routes.go`
6. **Handler calls service**, renders templ component
7. **Run** `templ generate` to compile templates

### Creating a New Data Model with Repository

1. **Define struct** in `data/myentity.go`
2. **Create repository** in `data/myentity_repository.go`
3. **Implement repository methods**: `GetByID()`, `Create()`, `Update()`, `Delete()`, `ListByOwner()`
4. **All methods accept** `context.Context` as first parameter
5. **Add logging** to all database operations

### Creating a New Service

1. **Create service file** in `services/myentity_service.go`
2. **Define service struct** that embeds `BaseService`
3. **Add repository** as a field
4. **Implement business logic methods**
5. **Perform validation and authorization** in service methods
6. **Return descriptive errors**

### Creating a New Templ Component

1. **Create `.templ` file** in appropriate directory (`layouts/`, `components/`, or `pages/`)
2. **Define component** with `templ ComponentName(params) { ... }`
3. **Add HTMX attributes** for server interactions
4. **Add Alpine.js directives** for client-side state
5. **Run** `templ generate` to compile to Go code
6. **Import and use** in handlers or other components

### Adding a Configuration Value

1. **Add field** to `AppConfig` struct in `config/config.go`
2. **Load from environment** in `LoadConfig()` function
3. **Add validation** if required
4. **Access via** `config.Get().YourField`

### Development Workflow

```bash
# Install tools
task tools:install

# Run development server with live reload
task dev

# Generate templ files
task templ:generate

# Run tests
task test

# Build application
task build
```

---

**Last Updated**: 2025-11-26  
**Version**: 1.0 (H.A.T. Stack Bootstrap)

#### Code Examples

- Use fenced code blocks with language specification:

```go
func Example() {
    // Code here
}
```

- For shell commands, use `bash` or `powershell`:

```bash
go build -o ./tmp/main.exe ./cmd
```

#### Links

- Use relative links between docs: `[Architecture](Architecture.md)`
- Use descriptive link text (not "click here")
- Verify links after renaming files

#### Lists

- Use `-` for unordered lists
- Use `1.` for ordered lists
- Indent nested lists with 2 spaces

### 13.5. Suggested Organization

**Core Documentation**:
- `Architecture.md` — System and component architecture, design decisions
- `Setup.md` — Installation, configuration, environment setup
- `CodingGuidelines.md` — Development standards, patterns, best practices
- `API.md` — API endpoints, request/response formats
- `WebSockets.md` — WebSocket protocol, events, message formats

**Operational Documentation**:
- `CloudRunDeployment.md` — Deployment procedures and configuration
- `ServerAuthentication.md` — Authentication setup and integration
- `Quickstart.md` — Quick start guide for new developers

**Process Documentation**:
- `DocRequirements.md` — Documentation requirements and standards
- `UpdateDocs.md` — Guide for updating documentation

**Feature Documentation**:
- `AuthIntegration.md` — Authentication feature implementation
- `DesignSystem.md` — UI/UX design system
- `BrandingUpdate.md` — Branding and visual identity changes

### 13.6. Versioning

If documentation differs by version:

1. **Note the applicable version** at the top of each file:

```markdown
# Document Title

[← Back to README](../README.md)

**Applicable Version**: v2.0+  
**Last Updated**: 2025-11-03
```

2. **For breaking changes**, create version-specific docs:
    - `MigrationV1ToV2.md`
    - `BreakingChangesV2.md`

### 13.7. Ownership and Updates

**Feature Development Rule**: Each new feature **must** add or update relevant documentation in `Docs/` as part of the same PR.

**Required Updates**:
- Architecture changes → Update `Architecture.md`
- New API endpoints → Update `API.md`
- New configuration → Update `Setup.md`
- New patterns → Update `CodingGuidelines.md`
- Deployment changes → Update `CloudRunDeployment.md`

**Review Checklist**:
- [ ] Documentation added/updated for new feature
- [ ] `README.md` table of contents updated
- [ ] Backlinks added to new files
- [ ] All links verified and working
- [ ] Code examples tested
- [ ] Version noted if applicable

### 13.8. Documentation Quality Standards

**Clarity**:
- Write for your audience (developers, operators, users)
- Define technical terms on first use
- Use examples to illustrate concepts

**Completeness**:
- Cover all aspects of the feature/system
- Include prerequisites and dependencies
- Document error cases and troubleshooting

**Accuracy**:
- Test all code examples
- Verify all commands work
- Update docs when code changes

**Maintainability**:
- Use consistent formatting
- Follow the style guide
- Keep docs close to code (update together)

### 13.8.1. What NOT to Document

**PROHIBITED Content in `Docs/`:**

Documentation in the `Docs/` directory must NEVER contain:

❌ **Task Completion Information**
- No "Changes Made" sections
- No "Files Modified" lists
- No "Testing Checklist" with checkboxes
- No "Build Status" sections
- No implementation step-by-step details

❌ **Implementation Details**
- No "How We Fixed It" narratives
- No "Root Cause" debugging stories
- No "Before/After" code comparisons for fixes
- No "Migration Steps Completed" summaries
- No PR/commit-style change descriptions

❌ **TODO Lists and Future Work**
- No "Next Steps" sections with action items
- No "Future Improvements" wish lists
- No "Known Limitations" that are really bugs to fix
- No "Pending Tasks" or work tracking

❌ **Historical Change Logs**
- No "What Changed in This Update" sections
- No dated change summaries
- No "Issues Fixed" lists
- No "Improvements Made" sections

**Where This Content BELONGS:**

- **Git Commit Messages** - Implementation details, what changed, why
- **Pull Request Descriptions** - Changes made, testing done, checklist
- **Issue Tracker** - Bugs, feature requests, TODO items
- **CHANGELOG.md** (if needed) - Version history, breaking changes
- **Internal Wiki/Notes** - Team-specific implementation notes

**Documentation Purpose:**

Documentation in `Docs/` exists to answer:
- ✅ **"How do I use this?"** - Usage guides, examples
- ✅ **"How does this work?"** - Architecture, design patterns
- ✅ **"How do I set this up?"** - Installation, configuration
- ✅ **"What does this do?"** - Feature descriptions, API reference
- ✅ **"How do I troubleshoot?"** - Common issues, solutions

NOT to answer:
- ❌ "What did we change?" - Use git history
- ❌ "What's left to do?" - Use issue tracker
- ❌ "How did we fix that bug?" - Use commit messages
- ❌ "What was the problem?" - Use PR descriptions

### 13.8.2. Documentation Content Guidelines

**Focus on USER NEEDS:**

```markdown
✅ GOOD - User-Focused:
## Authentication

The application uses Firebase Authentication with cookie-based sessions.

### Login
Send credentials to `/api/auth/login`:
```json
{
  "email": "user@example.com",
  "password": "password"
}
```

The server sets a `starx_session` cookie valid for 7 days.

### Accessing Protected Routes
Include the session cookie in requests. The server validates the token automatically.
```

```markdown
❌ BAD - Implementation-Focused:
## Cookie Authentication Fix

### Problem
Users were redirected back to login after successful authentication.

### Root Cause
1. Secure flag was set to true in development
2. GetOwner() wasn't checking cookies

### Changes Made
- Updated handlers/auth.go line 45
- Modified helpers/auth.go GetOwner() function
- Added production detection logic

### Testing Checklist
- [x] Login works in dev
- [x] Cookies set correctly
- [x] Dashboard accessible

### Files Modified
1. handlers/auth.go
2. helpers/auth.go
3. webhandlers/dashboard.go
```

**Write for ONBOARDING:**

Imagine a new developer joining the team. They need to:
- Understand what the system does
- Learn how to build and run it
- Know how to use the APIs
- Understand the architecture and patterns

They do NOT need to:
- Know the history of every bug fix
- See a list of completed tasks
- Read about implementation struggles
- Review change logs from development

### 13.8.3. Review Checklist for Documentation

Before committing documentation, verify:

**Content Check:**
- [ ] No "Changes Made" or "Files Modified" sections
- [ ] No "Testing Checklist" with completion status
- [ ] No "Next Steps" or "TODO" action items
- [ ] No "Issues Fixed" or "Problems Solved" narratives
- [ ] No "Before/After" implementation comparisons
- [ ] No "Build Status" or "Migration Complete" statements

**Purpose Check:**
- [ ] Explains HOW TO USE the feature
- [ ] Describes WHAT the feature does
- [ ] Documents HOW IT WORKS (architecture)
- [ ] Provides EXAMPLES and usage patterns
- [ ] Includes TROUBLESHOOTING for users
- [ ] Answers questions a NEW DEVELOPER would have

**Quality Check:**
- [ ] Written in present tense (not past tense change log)
- [ ] Focuses on current state (not historical changes)
- [ ] Provides actionable information (not history)
- [ ] Includes working code examples
- [ ] Links to related documentation

**If Your Documentation Contains:**
- "Changes Made" → Move to PR description
- "Testing Checklist" → Move to PR description
- "Next Steps" → Move to issue tracker
- "Issues Fixed" → Move to commit message
- "Files Modified" → Already in git history
- "Build Status" → Not needed in docs

### 13.8.4. Enforcement

**Code Review Requirements:**

Reviewers MUST reject PRs that add documentation containing:
- Task completion information
- Implementation change details
- Testing checklists
- "Next Steps" or TODO lists
- Historical change narratives

**Corrective Action:**

If such documentation is found:
1. Remove the file or section
2. Extract useful information (if any)
3. Rewrite in user-focused format
4. Move implementation details to appropriate location (PR, commit, issue)

**Examples of Acceptable Documentation:**

✅ `Architecture.md` - System design, component interaction  
✅ `API.md` - Endpoint reference, request/response formats  
✅ `Setup.md` - Installation, configuration steps  
✅ `CodingGuidelines.md` - Development patterns, standards  
✅ `Troubleshooting.md` - Common issues and solutions

**Examples of UNACCEPTABLE Documentation:**

❌ `CookieFix.md` - Bug fix implementation details  
❌ `BrandingUpdate.md` - Task completion checklist  
❌ `MigrationSummary.md` - Migration steps completed  
❌ `AuthUxImprovements.md` - Issues fixed narrative  
❌ `UpdateDocs.md` - Internal prompts/instructions

### 13.9. Example Documentation File

```markdown
# Feature Name

[← Back to README](../README.md)

**Applicable Version**: v2.0+  
**Last Updated**: 2025-11-03

## Overview

Brief description of the feature and its purpose.

## Architecture

How the feature fits into the system.

## Usage

### Prerequisites

- Requirement 1
- Requirement 2

### Configuration

```yaml
config:
  setting: value
```

### Examples

```go
// Code example
func Example() {
    // Implementation
}
```

## API Reference

### Endpoint Name

**Method**: `POST`  
**Path**: `/api/endpoint`

**Request**:
```json
{
  "field": "value"
}
```

**Response**:
```json
{
  "result": "success"
}
```

## Troubleshooting

### Common Issues

**Issue**: Description  
**Solution**: How to fix

## Related Documentation

- [Architecture](Architecture.md)
- [API Documentation](API.md)

---

[← Back to README](../README.md)
```

---

**Last Updated**: 2025-11-03  
**Version**: 2.0 (H.A.T. Stack Architecture)
