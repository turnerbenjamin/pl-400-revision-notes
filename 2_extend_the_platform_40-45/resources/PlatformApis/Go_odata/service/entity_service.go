// Package service provides interfaces and implementations for data access
// operations against OData endpoints. It offers generic entity services
// that handle CRUD operations with support for pagination and filtering.
package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/turnerbenjamin/go_odata/model"
	requestbuilder "github.com/turnerbenjamin/go_odata/request_builder"
	"github.com/turnerbenjamin/go_odata/view"
)

// Query parameter keys used for API requests

const (
	queryParamKeySelect = "$select"
	queryParmKeyFilter  = "$filter"
)

// HTTP header name constants used for API requests
const (
	headerContentType = "Content-Type"
	headerPrefer      = "Prefer"
	authHeader        = "Authorization"
	acceptHeader      = "Accept"
)

// HTTP header value constants used for API requests
const (
	contentTypeJSON            = "application/json"
	preferReturnRepresentation = "return=representation"
	preferMaxPageSizeFormat    = "odata.maxpagesize=%d"
	bearerTokenPrefix          = "Bearer "
)

// EntityService provides a generic interface for CRUD operations on entities
// that implement the view.Entity interface. The type parameter T represents the
// entity type.
type EntityService[T view.Entity] interface {
	// List retrieves all entities, optionally filtered by searchTerm
	List(searchTerm string) (view.EntityList[T], error)

	// Get retrieves a specific entity by its GUID
	Get(guid string) (T, error)

	// Create adds a new entity and returns the created entity with
	// server-generated fields
	Create(entityToCreate T) (newEntity T, err error)

	// Update modifies an existing entity identified by GUID
	Update(guid string, entityToUpdate T) error

	// Delete removes an entity identified by GUID
	Delete(guid string) error
}

// EntityServiceOptions contains configuration parameters for creating an
// EntityService instance
type EntityServiceOptions struct {
	// DataverseService handles the actual HTTP communication with the API
	DataverseService DataverseService

	// BaseUrl is the root URL of the API
	BaseUrl *url.URL

	// ResourcePath is the path segment for the specific entity type (e.g.,
	// "accounts")
	ResourcePath string

	// PageLimit sets the maximum number of entities to retrieve per page
	PageLimit int

	// SelectsFields specifies which entity fields to retrieve from the API
	SelectsFields []string

	// SearchFields defines which fields are included in search operations
	SearchFields []string
}

// entityService implements EntityService for a specific entity type T
type entityService[T view.Entity] struct {
	dataverseService DataverseService
	pageLimit        int
	resourceUrl      *url.URL
	selects          string
	zeroValue        T
	searchFields     []string
}

// NewEntityService creates a new EntityService implementation for the specified
// entity type with the provided options.
func NewEntityService[T view.Entity](options EntityServiceOptions) EntityService[T] {
	selectsString := strings.Join(options.SelectsFields, ",")

	resourceUrl := *options.BaseUrl
	resourceUrl.Path = path.Join(resourceUrl.Path, options.ResourcePath)

	return &entityService[T]{
		dataverseService: options.DataverseService,
		pageLimit:        options.PageLimit,
		resourceUrl:      &resourceUrl,
		selects:          selectsString,
		searchFields:     options.SearchFields,
	}
}

// Create adds a new entity to the system and returns the created entity with
// server-generated fields populated.
func (s *entityService[T]) Create(entity T) (T, error) {

	path := s.resourceUrl.String()
	payload, err := json.Marshal(entity)
	if err != nil {
		return s.zeroValue, fmt.Errorf("failed to serialise entity %w", err)
	}

	req, err := requestbuilder.NewRequestBuilder(http.MethodPost, path, bytes.NewReader(payload)).
		Build()
	if err != nil {
		return s.zeroValue, err
	}

	req.Header.Set(headerContentType, contentTypeJSON)
	req.Header.Set(headerPrefer, preferReturnRepresentation)

	res, err := s.dataverseService.Execute(req)
	if err != nil {
		return s.zeroValue, fmt.Errorf("failed to create entity: %w", err)
	}

	if !res.IsSuccessful {
		errMsg := model.ParseErrorMessage(res.Body)
		return s.zeroValue, errors.New(errMsg)
	}

	var newEntity T
	if err := json.Unmarshal(res.Body, &newEntity); err != nil {
		return s.zeroValue, fmt.Errorf("failed to unmarshal created entity: %w", err)
	}
	return newEntity, nil
}

// List retrieves entities, optionally filtered by searchTerm.
// It returns a paginated collection that handles fetching additional pages as
// needed.
func (s *entityService[T]) List(searchTerm string) (view.EntityList[T], error) {

	path := s.resourceUrl.String()
	rb := requestbuilder.NewRequestBuilder(http.MethodGet, path, nil).
		AddQueryParam(queryParamKeySelect, s.selects)

	if searchTerm != "" && len(s.searchFields) > 0 {
		filter := s.buildFilterQuery(searchTerm)
		rb.AddQueryParam(queryParmKeyFilter, filter)
	}

	req, err := rb.Build()
	if err != nil {
		return nil, err
	}

	req.Header.Set(headerPrefer, fmt.Sprintf(preferMaxPageSizeFormat, s.pageLimit))

	res, err := s.dataverseService.Execute(req)
	if err != nil {

		return nil, fmt.Errorf("failed to retrieve entities: %w", err)
	}

	if !res.IsSuccessful {
		errMsg := model.ParseErrorMessage(res.Body)
		return nil, errors.New(errMsg)
	}

	gmr := model.NewGetManyResponseFromJson[T](res.Body)

	return model.CreateEntityList(*gmr, s.getNextResult), nil
}

// Get retrieves a single entity by its unique identifier (GUID).
func (s *entityService[T]) Get(guid string) (T, error) {

	path := s.buildUrlWithGuid(guid)
	req, err := requestbuilder.NewRequestBuilder(http.MethodGet, path, nil).
		AddQueryParam(queryParamKeySelect, s.selects).
		Build()
	if err != nil {
		return s.zeroValue, err
	}

	res, err := s.dataverseService.Execute(req)
	if err != nil {
		return s.zeroValue, fmt.Errorf("failed to retrieve entity (%s): %w", guid, err)
	}

	if !res.IsSuccessful {
		errMsg := model.ParseErrorMessage(res.Body)
		return s.zeroValue, errors.New(errMsg)
	}

	var entity T
	if err := json.Unmarshal(res.Body, &entity); err != nil {
		return s.zeroValue, fmt.Errorf("failed to unmarshal retrieved entity: %w", err)
	}

	if err != nil {
		return s.zeroValue, err
	}
	return entity, nil
}

// Update modifies an existing entity identified by GUID with the properties
// from entityToUpdate.
func (s *entityService[T]) Update(guid string, entityToUpdate T) error {

	path := s.buildUrlWithGuid(guid)
	payload, err := json.Marshal(entityToUpdate)

	if err != nil {
		return fmt.Errorf("failed to serialise entity %w", err)
	}
	req, err := requestbuilder.NewRequestBuilder(http.MethodPatch, path, bytes.NewReader(payload)).
		Build()
	if err != nil {
		return err
	}

	req.Header.Set(headerContentType, contentTypeJSON)

	res, err := s.dataverseService.Execute(req)
	if err != nil {
		return err
	}

	if !res.IsSuccessful {
		errMsg := model.ParseErrorMessage(res.Body)
		return errors.New(errMsg)
	}
	return nil
}

// Delete removes an entity identified by its GUID.
func (s *entityService[T]) Delete(guid string) error {

	path := s.buildUrlWithGuid(guid)
	req, err := requestbuilder.NewRequestBuilder(http.MethodDelete, path, nil).
		Build()
	if err != nil {
		return err
	}

	res, err := s.dataverseService.Execute(req)
	if err != nil {
		return fmt.Errorf("failed to delete entity (%s): %w", guid, err)
	}

	if !res.IsSuccessful {
		errMsg := model.ParseErrorMessage(res.Body)
		return errors.New(errMsg)
	}
	return nil
}

// buildUrlWithGuid constructs a URL targeting a specific entity by appending
// its GUID to the resource URL.
func (s *entityService[T]) buildUrlWithGuid(guid string) string {
	return fmt.Sprintf("%s(%s)", s.resourceUrl.String(), guid)
}

// buildFilterQuery constructs an OData filter query string from a search term
// using the configured search fields.
func (s *entityService[T]) buildFilterQuery(searchTerm string) string {
	if len(s.searchFields) == 0 {
		return ""
	}
	filters := make([]string, len(s.searchFields))
	for i, sf := range s.searchFields {
		filters[i] = fmt.Sprintf("contains(%s,'%s')", sf, searchTerm)
	}

	return url.QueryEscape(strings.Join(filters, " or "))
}

// getNextResult fetches the next page of results during list pagination
// using the provided URL from the previous response's "@odata.nextLink".
func (s *entityService[T]) getNextResult(url string) (*model.GetManyResponse[T], error) {
	req, err := requestbuilder.NewRequestBuilder(http.MethodGet, url, nil).Build()
	if err != nil {
		return nil, err
	}

	req.Header.Set(headerPrefer, fmt.Sprintf(preferMaxPageSizeFormat, s.pageLimit))

	dr, err := s.dataverseService.Execute(req)
	if err != nil {
		return nil, err
	}
	return model.NewGetManyResponseFromJson[T](dr.Body), nil
}
