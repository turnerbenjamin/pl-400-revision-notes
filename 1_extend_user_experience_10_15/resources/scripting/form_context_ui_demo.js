"use strict";

// Functionality stored in object with publisher prefix to avoid global
// namespace pollution
this.cr950 = this?.window?.cr950 || {};

// IIFE used to avoid pollution of the cr950.formContextUiDemoFunctions
// namespace
this.cr950.formContextUiDemoFunctions = (function () {
  //Dictionary of logical names used in this resource
  const logicalNames = {
    accountName: "name",
    primaryContact: "primarycontactid",
    primaryContactQuickView: "contactquickform",
    primaryContactPhone: "telephone1",
    accountMainPhone: "telephone1",
    accountFax: "fax",
    detailsTab: "DETAILS_TAB",
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

  const notificationTypes = {
    info: "INFO",
    warning: "WARNING",
    error: "ERROR",
  };

  const controlNotificationLevels = {
    recommendation: "RECOMMENDATION",
    error: "ERROR",
  };

  /**
   * The functionality of this demo is split across a number of different
   * events. This handler is added using make.powerapps. It is responsibly for
   * adding the other handlers programmatically.
   */
  function accountForInteractiveExperienceFormOnLoadHandler(executionContext) {
    const formContext = executionContext.getFormContext();
    if (formContext.ui.getFormType() !== formTypes.create) return;

    addOnChangeHandler(
      formContext,
      logicalNames.accountName,
      accountNameOnChangeHandler
    );

    addOnChangeHandler(
      formContext,
      logicalNames.accountFax,
      accountFaxOnChangeHandler
    );

    addOnChangeHandler(
      formContext,
      logicalNames.primaryContact,
      primaryContactOnChangeHandler
    );

    const notificationMessage = "All handlers added";
    formContext.ui.setFormNotification(
      notificationMessage,
      notificationTypes.info,
      notificationMessage
    );
  }

  /**
   * Helper method used to add handlers to the onChange event of attributes
   */
  function addOnChangeHandler(formContext, attributeName, handler) {
    const attribute = formContext.getAttribute(attributeName);
    attribute.addOnChange(handler);
  }

  /**
   * This demo listens for changes on the account name attribute and sets
   * the form entity name to the value in account name
   */
  function accountNameOnChangeHandler(executionContext) {
    const formContext = executionContext.getFormContext();
    const accountName = formContext.data.entity.getPrimaryAttributeValue();
    formContext.ui.setFormEntityName(accountName);
  }

  /**
   * This demo listens for changes on the fax attribute. It gets ui metadata
   * from the form context and displays these in a Confirm Dialogue box. If the
   * box is selected it will navigate to the details tab
   */
  async function accountFaxOnChangeHandler(executionContext) {
    const formContext = executionContext.getFormContext();
    const formMetaDataString = getFormMetaDataString(formContext);
    const detailsTab = formContext.ui.tabs.get(logicalNames.detailsTab);
    const result = await Xrm.Navigation.openConfirmDialog({
      title: "Navigate to Details Tab?",
      text: formMetaDataString,
      confirmButtonLabel: "Let's Go",
      cancelButtonLabel: "No thanks",
    });
    if (result.confirmed) {
      detailsTab.setFocus();
    }
  }

  /** Collects metadata from the ui object and returns a formatted string */
  function getFormMetaDataString(formContext) {
    const availableForms = formContext.ui.formSelector.items
      .get()
      .map((f) => f.getLabel())
      .join(", ");
    const navItems = formContext.ui.navigation.items
      .get()
      .map((f) => f.getLabel())
      .join(", ");
    const viewPortWidth = formContext.ui.getViewPortWidth();
    const viewPortHeight = formContext.ui.getViewPortHeight();

    return [
      `Available Forms ${availableForms}`,
      `Nav Items: ${navItems}`,
      `ViewPort: ${viewPortWidth}x${viewPortHeight}`,
    ].join("\n");
  }

  /**
   * This demo listens for changes on the primary contact attribute. When this
   * attribute changes, a notification is set on the account phone number
   * attribute with an action to set the phone number.
   *
   * If there is a primary contact with a main phone number, the action will use
   * this number else it will use a default number
   */
  function primaryContactOnChangeHandler(executionContext) {
    const formContext = executionContext.getFormContext();
    const contactPhoneNumber = getPrimaryContactPhoneNumber(formContext);
    const defaultPhoneNumber = "0123456789";

    const [accountPhoneAttribute, accountPhoneControl] = getAttributeAndControl(
      formContext,
      logicalNames.accountMainPhone
    );

    recommendSettingMainPhone(
      contactPhoneNumber,
      defaultPhoneNumber,
      accountPhoneAttribute,
      accountPhoneControl
    );
  }

  /**
   * Helper method to set a notification on a given control
   */
  function setControlNotification(control, notificationId, message, actions) {
    control.addNotification({
      messages: [message],
      notificationLevel: controlNotificationLevels.recommendation,
      uniqueId: notificationId,
      actions,
    });
  }

  /**
   * Helper method to get an attribute and control for a given field
   */
  function getAttributeAndControl(formContext, logicalName) {
    const attribute = formContext.getAttribute(logicalName);
    const control = formContext.getControl(logicalName);
    return [attribute, control];
  }

  /**
   * Set a notification on the main phone control of the account with an action
   * to set phone to the contact phone number if not nullish else the default
   * number
   */
  function recommendSettingMainPhone(
    contactPhoneNumber,
    defaultPhoneNumber,
    accountPhoneAttribute,
    accountPhoneControl
  ) {
    const notificationMessage = "Populate the main phone value";
    const notificationId = notificationMessage;

    const populateMainPhoneActions = getPopulateMainPhoneActions(
      notificationId,
      contactPhoneNumber,
      defaultPhoneNumber,
      accountPhoneAttribute,
      accountPhoneControl
    );

    setControlNotification(
      accountPhoneControl,
      notificationId,
      notificationMessage,
      populateMainPhoneActions
    );
  }

  /**
   *  Extract the phone number from the quick view form for contact. Returns
   *  undefined if the contact is not populated or the phone number is not
   *  populated
   */
  function getPrimaryContactPhoneNumber(formContext) {
    const primaryContactQuickViewControl = formContext.ui.quickForms.get(
      logicalNames.primaryContactQuickView
    );
    const phoneNumberAttribute = primaryContactQuickViewControl.getAttribute(
      logicalNames.primaryContactPhone
    );
    return phoneNumberAttribute?.getValue();
  }

  /**
   * Build notification action to set phone number to contact number if defined
   * else the default phone number.
   *
   * If the action is selected the value is set, the notification is removed and
   * the form is refreshed.
   */
  function getPopulateMainPhoneActions(
    notificationId,
    contactPhoneNumber,
    defaultPhoneNumber,
    phoneNumberAttribute,
    phoneNumberControl
  ) {
    const phoneNumberSource = contactPhoneNumber ? "contact" : "default";
    const newValue = contactPhoneNumber || defaultPhoneNumber;

    return [
      {
        message: `Set to ${phoneNumberSource} phone number`,
        actions: [
          function () {
            phoneNumberAttribute.setValue(newValue);
            phoneNumberControl.clearNotification(notificationId);
            phoneNumberControl.data.refresh();
          },
        ],
      },
    ];
  }

  // Expose onLoad handler
  return {
    accountForInteractiveExperienceFormOnLoadHandler,
  };
})();
