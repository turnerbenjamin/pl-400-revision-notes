"use strict";

// Functionality stored in object with publisher prefix to avoid global
// namespace pollution
this.cr950 = this?.window?.cr950 || {};

// IIFE used to avoid pollution of the cr950.formContextDataDemoFunctions
// namespace
this.cr950.formContextDataDemoFunctions = (function () {
  //Dictionary of logical names used in this resource
  const logicalNames = {
    jobTitle: "jobtitle",
    fax: "fax",
    firstName: "firstname",
    lastName: "lastname",
  };

  //Dictionary of form types returned by getFormType
  const formTypes = {
    notSet: 0,
    create: 1,
    update: 2,
    readOnly: 3,
    disabled: 4,
    bulkEdit: 5,
  };

  /**
   *  onLoad handler attached to form properties in make.powerapps.
   *
   *  This flow will set job title to show the handler from which it was set.
   *
   *  It will set the first and last name fields if the form is not valid
   *
   *  A post save handler is added to demonstrate adding handlers
   *  programmatically and the form is saved
   */
  function contactFormOnLoadHandler(executionContext) {
    const formContext = executionContext.getFormContext();
    if (formContext.ui.getFormType() !== formTypes.create) return;

    setFormAttribute(formContext, logicalNames.jobTitle, "ON LOAD");
    setRequiredFieldsIfFormInvalid(formContext);

    formContext.data.entity.addOnPostSave(contactFormPostSaveHandler);
    saveContactFormIfDirty(formContext);
  }

  /**
   *  onSave handler attached programmatically.
   *
   *  This handler sets job title to indicate the handler from which it was set
   *
   *  The form is saved which will trigger a second call to this handler as the
   *  job title field will be dirty. On the second call job title is unchanged
   *  so there will be no infinite loop. This demonstrates is dirty
   *
   *  Finally, job title is set to a new value and the form is refreshed without
   *  being saved first. The job title attribute will revert to the saved value
   */
  async function contactFormPostSaveHandler(executionContext) {
    const formContext = executionContext.getFormContext();
    const saveCount = getSaveCount(formContext);
    const isSecondCall = saveCount === 2;

    setFormAttribute(formContext, logicalNames.jobTitle, "ON SAVE");
    await saveContactFormIfDirty(formContext);

    if (isSecondCall) {
      setFormAttribute(formContext, logicalNames.jobTitle, "DIRTY DATA");
      refreshForm(formContext);
    }
  }

  // Demonstrates the use of is valid on the data object.
  function setRequiredFieldsIfFormInvalid(formContext) {
    if (formContext.data.isValid()) return;
    setFormAttribute(formContext, logicalNames.firstName, "FIRST");
    setFormAttribute(formContext, logicalNames.lastName, "LAST");
  }

  // Helper function to set a form attribute using the data object. This is the
  // long way to access an attribute
  function setFormAttribute(formContext, attributeLogicalName, value) {
    const attribute =
      formContext.data.entity.attributes.get(attributeLogicalName);
    attribute.setValue(value);
  }

  // Uses the getIsDirtyMethod to determine if the form is dirty. Saves the form
  // only if it is dirty. The save count is incremented before a save to
  // demonstrate this.
  async function saveContactFormIfDirty(formContext) {
    const isDirty = formContext.data.getIsDirty();
    if (!isDirty) return;

    incrementSaveCount(formContext);
    try {
      await formContext.data.save();
    } catch (error) {
      console.error(error);
    }
  }

  // Demonstrate refreshing the form. This will remove any dirty data. Note, the
  // bool arg is optional, it is included to be explicit
  async function refreshForm(formContext) {
    const doSave = false;
    try {
      await formContext.data.refresh(doSave);
    } catch (error) {
      console.error(error);
    }
  }

  // Helper function, increments the save count in the fax field
  function incrementSaveCount(formContext) {
    const saveCount = getSaveCount(formContext);
    const faxField = formContext.getAttribute(logicalNames.fax);
    faxField.setValue(`${saveCount + 1}`);
  }

  // Helper function, gets save count from the fax field
  function getSaveCount(formContext) {
    const faxField = formContext.getAttribute(logicalNames.fax);
    if (!faxField.getValue()) return 0;
    return Number.parseInt(faxField.getValue());
  }

  // Expose onLoad handler
  return {
    contactFormOnLoadHandler,
  };
})();
