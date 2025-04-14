package app

import (
	"errors"
	"fmt"

	"github.com/turnerbenjamin/go_odata/constants/confirmoption"
	"github.com/turnerbenjamin/go_odata/constants/table_menu_option"
	"github.com/turnerbenjamin/go_odata/service"
	"github.com/turnerbenjamin/go_odata/view"
)

type entityMenu[T view.Entity] struct {
	ui               view.UI
	service          service.EntityService[T]
	listColumns      []view.ListColumn[T]
	getNewEntity     func() (T, error)
	getUpdatedEntity func(T) (T, error)
	entityLabel      string
	searchTerm       string
}

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

		switch table_menu_option.TableMenuOption(menuOutput.UserInput()) {
		case table_menu_option.Back:
			return nil
		case table_menu_option.Search:
			err = em.setSearchTerm()
		case table_menu_option.Create:
			err = em.createEntity()
		case table_menu_option.Update:
			err = em.updateEntity(menuOutput.Target())
		case table_menu_option.Delete:
			err = em.deleteEntity(menuOutput.Target())
		default:
			err = errors.New("invalid menu option: " + menuOutput.UserInput())
		}

		if err != nil {
			return em.displayErrorScreen(err)
		}
	}
}

func (em *entityMenu[T]) notifyNoErrorsFound() error {
	is, err := newInfoScreen("No rows found")
	if err != nil {
		return err
	}

	_, err = em.ui.NavigateTo(is)
	if err != nil {
		return err
	}

	return nil
}

func (em *entityMenu[T]) displayEntityMenu(entityList view.EntityList[T]) (view.ScreenOutput, error) {
	entityListScreen, err := NewEntityListScreen(
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

func (em *entityMenu[T]) createEntity() error {
	entityToCreate, err := em.getNewEntity()
	if err != nil {
		return em.displayErrorScreen(err)
	}

	newEntity, err := em.service.Create(entityToCreate)
	if err != nil {
		return em.displayErrorScreen(err)
	}
	successMessage := fmt.Sprintf("%s: %s", em.entityLabel, newEntity.Label())
	return em.displaySuccessScreen(successMessage)
}

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

func (em *entityMenu[T]) deleteEntity(guid string) error {
	entityToDelete, err := em.service.Get(guid)
	if err != nil {
		return em.displayErrorScreen(err)
	}
	msg := fmt.Sprintf("Are you sure you want to delete %s", entityToDelete.Label())
	confirmationScreen, err := NewConfirmationScreen(msg)

	response, err := em.ui.NavigateTo(confirmationScreen)
	if err != nil {
		return em.displayErrorScreen(err)
	}

	if confirmoption.ConfirmOption(response.UserInput()) != confirmoption.Yes {
		return nil
	}

	err = em.service.Delete(guid)
	if err != nil {
		return err
	}
	successMsg := fmt.Sprintf("%s deleted", em.entityLabel)
	return em.displaySuccessScreen(successMsg)
}

func (em *entityMenu[T]) setSearchTerm() error {
	title := fmt.Sprintf("Set search term: %ss", em.entityLabel)
	msg := "Enter a search term (or leave blank to unset)"

	inputScreen, err := newStringInputScreen(title, msg, "SearchTerm", em.searchTerm, false)
	if err != nil {
		return em.displayErrorScreen(err)
	}

	inputScreenOutputs, err := em.ui.NavigateTo(inputScreen)
	if err != nil {
		return err
	}

	em.searchTerm = inputScreenOutputs.UserInput()
	return nil
}

func (em *entityMenu[T]) displaySuccessScreen(message string) error {
	ss, err := newSuccessScreen(message)
	if err != nil {
		return em.displayErrorScreen(err)
	}
	_, err = em.ui.NavigateTo(ss)
	return err
}

func (em *entityMenu[T]) displayErrorScreen(originalError error) error {
	es, err := getErrorScreen(originalError.Error())
	if err != nil {
		return originalError
	}
	_, err = em.ui.NavigateTo(es)
	if err != nil {
		return originalError
	}
	return nil
}
