// Package app implements the application-level functionality for the OData
// client, including screens, navigation, and business logic.
package app

import (
	"fmt"

	confirmOption "github.com/turnerbenjamin/go_odata/constants/confirm_option"
	tableMenuOption "github.com/turnerbenjamin/go_odata/constants/tablemenuoption"
	"github.com/turnerbenjamin/go_odata/service"
	"github.com/turnerbenjamin/go_odata/view"
)

// entityMenu provides CRUD (Create, Read, Update, Delete) functionality for a
// generic entity type T.
// It manages the UI navigation flow and interactions with the entity service.
type entityMenu[T view.Entity] struct {
	// User interface instance for screen navigation
	ui view.UI
	// Service for entity operations
	service service.EntityService[T]
	// Column definitions for entity list display
	listColumns []view.ListColumn[T]
	// Function to get data for a new entity
	getNewEntity func() (T, error)
	// Function to get updated data for an existing entity
	getUpdatedEntity func(T) (T, error)
	// Human-readable label for this entity type
	entityLabel string
	// Current search term for filtering entities
	searchTerm string
}

// run starts the entity menu's main loop, handling user interactions until
// exit. It fetches and displays entities, processes user commands, and manages
// errors.
// Returns an error if any operation fails.
func (em *entityMenu[T]) run() error {
	for {
		entityList, err := em.service.List(em.searchTerm)
		if err != nil {
			return em.displayErrorScreen(err)
		}

		if len(entityList.Data()) == 0 {
			err := em.notifyNoErrorsFound()
			if err != nil || em.searchTerm == "" {
				return err
			}
			em.searchTerm = ""
			continue
		}

		menuOutput, err := em.displayEntityMenu(entityList)
		if err != nil {
			return em.displayErrorScreen(err)
		}

		switch tableMenuOption.TableMenuOption(menuOutput.UserInput()) {
		case tableMenuOption.Back:
			return nil
		case tableMenuOption.Search:
			err = em.setSearchTerm()
		case tableMenuOption.Create:
			err = em.createEntity()
		case tableMenuOption.Update:
			err = em.updateEntity(menuOutput.Target())
		case tableMenuOption.Delete:
			err = em.deleteEntity(menuOutput.Target())
		default:
			err = fmt.Errorf("invalid menu option %s", menuOutput.UserInput())
		}

		if err != nil {
			return em.displayErrorScreen(err)
		}
	}
}

// notifyNoErrorsFound displays a message when no entities match the current
// search criteria.
// Returns an error if the notification screen cannot be displayed.
func (em *entityMenu[T]) notifyNoErrorsFound() error {
	is, err := newInfoScreen("No rows found")
	if err != nil {
		return err
	}

	_, err = em.ui.NavigateTo(is)
	return err
}

// displayEntityMenu creates and shows the entity list screen with the provided
// entity data.
// Returns the user's selection and any error encountered.
func (em *entityMenu[T]) displayEntityMenu(entityList view.EntityList[T]) (view.ScreenOutput, error) {
	entityListScreen, err := newEntityListScreen(
		em.entityLabel,
		listScreenOptions[T]{
			entityList: entityList,
			columns:    em.listColumns,
		},
	)
	if err != nil {
		return nil, err
	}

	outputs, err := em.ui.NavigateTo(entityListScreen)
	if err != nil {
		return nil, err
	}
	return outputs, nil
}

// createEntity handles the workflow for creating a new entity.
// It prompts for entity data, calls the service to create it, and displays a
// success message.
// Returns an error if any step in the process fails.
func (em *entityMenu[T]) createEntity() error {
	entityToCreate, err := em.getNewEntity()
	if err != nil {
		return err
	}

	newEntity, err := em.service.Create(entityToCreate)
	if err != nil {
		return err
	}
	successMessage := fmt.Sprintf("%s: %s", em.entityLabel, newEntity.Label())
	return em.displaySuccessScreen(successMessage)
}

// updateEntity handles the workflow for updating an existing entity.
// It fetches the current entity data, prompts for updates, calls the service to
// update it, and displays a success message.
// The guid parameter identifies the entity to update.
// Returns an error if any step in the process fails.
func (em *entityMenu[T]) updateEntity(guid string) error {
	currentEntity, err := em.service.Get(guid)
	if err != nil {
		return err
	}
	entityToUpdate, err := em.getUpdatedEntity(currentEntity)
	if err != nil {
		return err
	}

	err = em.service.Update(guid, entityToUpdate)
	if err != nil {
		return err
	}
	successMsg := fmt.Sprintf("%s updated", em.entityLabel)
	return em.displaySuccessScreen(successMsg)
}

// deleteEntity handles the workflow for deleting an existing entity.
// It fetches the entity data, prompts for confirmation, calls the service to
// delete it, and displays a success message if confirmed.
// The guid parameter identifies the entity to delete.
// Returns an error if any step in the process fails.
func (em *entityMenu[T]) deleteEntity(guid string) error {
	entityToDelete, err := em.service.Get(guid)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("Are you sure you want to delete %s", entityToDelete.Label())
	confirmationScreen, err := NewConfirmationScreen(msg)
	if err != nil {
		return err
	}

	response, err := em.ui.NavigateTo(confirmationScreen)
	if err != nil {
		return err
	}

	if confirmOption.ConfirmOption(response.UserInput()) != confirmOption.Yes {
		return nil
	}

	err = em.service.Delete(guid)
	if err != nil {
		return err
	}
	successMsg := fmt.Sprintf("%s deleted", em.entityLabel)
	return em.displaySuccessScreen(successMsg)
}

// setSearchTerm prompts the user to enter a search term for filtering entities.
// An empty search term clears the filter.
// Returns an error if the input screen cannot be displayed.
func (em *entityMenu[T]) setSearchTerm() error {
	title := fmt.Sprintf("Set search term: %ss", em.entityLabel)
	msg := "Enter a search term (or leave blank to unset)"

	inputScreen, err := newStringInputScreen(title, msg, "SearchTerm", em.searchTerm, false)
	if err != nil {
		return err
	}

	inputScreenOutputs, err := em.ui.NavigateTo(inputScreen)
	if err != nil {
		return err
	}

	em.searchTerm = inputScreenOutputs.UserInput()
	return nil
}

// displaySuccessScreen shows a success message to the user.
// Returns an error if the success screen cannot be displayed.
func (em *entityMenu[T]) displaySuccessScreen(message string) error {
	ss, err := newSuccessScreen(message)
	if err != nil {
		return err
	}
	_, err = em.ui.NavigateTo(ss)
	return err
}

// displayErrorScreen shows an error message to the user.
// Returns the original error if displaying the error screen fails,
// otherwise returns nil to indicate the error was displayed successfully.
func (em *entityMenu[T]) displayErrorScreen(originalError error) error {
	es, err := newErrorScreen(originalError.Error())
	if err != nil {
		return originalError
	}
	_, err = em.ui.NavigateTo(es)
	if err != nil {
		return originalError
	}
	return nil
}
