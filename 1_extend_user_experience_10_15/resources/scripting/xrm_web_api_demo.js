"use strict";

// Functionality stored in object with publisher prefix to avoid global
// namespace pollution
this.cr950 = this?.window?.cr950 || {};

// IIFE used to avoid pollution of the cr950.xrmWebAPIDemo
// namespace
this.cr950.xrmWebAPIDemo = (function () {
  //Dictionary of logical names used in this resource
  const logicalNames = {
    account: "account",
    contact: "contact",
    accountName: "name",
    accountId: "accountid",
    accountContact: "primarycontactid",
    accountContactLookup: "_primarycontactid_value",
    accountPhoneNumber: "telephone1",
    contactFirstName: "firstname",
    contactLastName: "lastname",
    contactEmail: "emailaddress1",
    contactFullName: "fullname",
    contactId: "contactid",
  };

  const dummyData = {
    accountName: "SAMPLE ACCOUNT",
    contactFirstName: "SAMPLE",
    contactLastName: "CONTACT",
    contactEmail: "sample@email.com",
  };

  /**
   * Demonstrate account creation with the WebApi. Creates both an account and
   * a contact using the deep insert pattern
   */
  function createDemo(primaryControl) {
    executeAsyncWithErrorHandling(async () => {
      const demoAccount = buildDemoAccount();
      await Xrm.WebApi.createRecord(logicalNames.account, demoAccount);
      await primaryControl.refresh();
    }, "Creating sample account and contact");
  }

  /**
   * Demonstrates record deletion by removing the sample account and contact
   * created in the create demo. If the records are not present an error is
   * thrown.
   */
  function deleteDemo(primaryControl) {
    executeAsyncWithErrorHandling(async () => {
      const dummyAccount = await getDummyAccount();
      await Xrm.WebApi.deleteRecord(
        logicalNames.contact,
        dummyAccount[logicalNames.accountContactLookup]
      );
      await Xrm.WebApi.deleteRecord(
        logicalNames.account,
        dummyAccount[logicalNames.accountId]
      );
      await primaryControl.refresh();
    }, "Deleting Sample Account and contact");
  }

  /** Demonstrates record update by adding a phone number to the dummy account
   * An error is thrown if the sample account has not yet been created
   */
  function updateDemo(primaryControl) {
    executeAsyncWithErrorHandling(async () => {
      const dummyAccount = await getDummyAccount();
      const update = {};
      update[logicalNames.accountPhoneNumber] = "123456789";
      await Xrm.WebApi.updateRecord(
        logicalNames.account,
        dummyAccount[logicalNames.accountId],
        update
      );
      await primaryControl.refresh();
    }, "Adding phone number to sample account");
  }

  // Helper function to retrieve the dummy account. Ideally, a retrieve request
  // would be used but I have been unable to construct a request using alternate
  // keys.
  async function getDummyAccount() {
    const select = `$select=${logicalNames.accountName},${logicalNames.accountContactLookup}`;
    const filter = `$filter=${logicalNames.accountName} eq '${dummyData.accountName}'`;
    const limit = "$top=1";
    const records = await Xrm.WebApi.retrieveMultipleRecords(
      logicalNames.account,
      `?${select}&${filter}&${limit}`
    );
    if (!records?.entities[0]) {
      throw new Error("Sample account not found");
    }
    return records?.entities[0];
  }

  // Wrapper function, adds error handling to async functions and displays a
  // progress indicator while the promise is pending
  async function executeAsyncWithErrorHandling(
    asyncFunction,
    messageWhilePending
  ) {
    try {
      Xrm.Utility.showProgressIndicator(messageWhilePending);
      await asyncFunction();
    } catch (error) {
      displayError(error);
    } finally {
      Xrm.Utility.closeProgressIndicator();
    }
  }

  // Shows error as an alert dialog
  function displayError(error) {
    Xrm.Navigation.openErrorDialog({
      message: error.message,
      details: error.stack,
    });
  }

  // Build a demo account. Deep insert pattern used to create both an account
  // and associated contact at the same time. Note that this pattern may only be
  // used in online mode.
  function buildDemoAccount() {
    const demoAccount = {};
    demoAccount[logicalNames.accountName] = dummyData.accountName;
    demoAccount[logicalNames.accountContact] = buildDemoContact();
    return demoAccount;
  }

  // Build a demo contact.
  function buildDemoContact() {
    const demoContact = {};
    demoContact[logicalNames.contactFirstName] = dummyData.contactFirstName;
    demoContact[logicalNames.contactLastName] = dummyData.contactLastName;
    demoContact[logicalNames.contactEmail] = dummyData.contactEmail;
    return demoContact;
  }

  // Expose demo functions
  return {
    createDemo,
    deleteDemo,
    updateDemo,
  };
})();
