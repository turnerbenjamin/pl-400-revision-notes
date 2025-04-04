"use strict";

// Functionality stored in object with publisher prefix to avoid global
// namespace pollution
this.cr950 = this?.window?.cr950 || {};

// IIFE used to avoid pollution of the cr950.xrmNavigationAndPanelDemo
// namespace
this.cr950.xrmNavigationAndPanelDemo = (function () {
  //Dictionary of logical names used in this resource
  const logicalNames = {
    contact: "contact",
    account: "account",
  };

  //Dictionary of form types returned by getFormType
  const pageTypes = {
    entityList: "entitylist",
    entityRecord: "entityrecord",
    dashboard: "dashboard",
    htmlWebResource: "webresource",
    customPage: "custom",
  };

  const targets = {
    inline: 1,
    dialog: 2,
  };

  /**
   * This is a command bar function the isList parameter is configured there.
   *
   * Demonstrates the use of navigateTo  a list and form. We can also navigate
   * to relationship objects, dashboards, html web resources and custom pages
   */

  function navigateToDemo(isList) {
    const pageInput = {
      pageType: isList ? pageTypes.entityList : pageTypes.entityRecord,
      entityName: logicalNames.contact,
    };

    const navigationOptions = {
      target: targets.dialog,
    };

    Xrm.Navigation.navigateTo(pageInput, navigationOptions);
  }

  /**
   * Demonstrates use of the three dialog options. Displays a confirmation
   * dialog. If confirmed an alert dialog is shown. If canceled an error dialog
   * is shown
   */
  async function dialogDemo() {
    const didConfirm = await openAConfirmDialog();
    if (didConfirm) {
      openAnAlertDialog();
    } else {
      openAnErrorDialog();
    }
  }

  async function openAConfirmDialog() {
    const alertStrings = {
      title: "THIS IS A CONFIRMATION DIALOG",
      subtitle: "This is a subtitle",
      text: "This is the text of an confirmation dialog",
      confirmButtonLabel: "Confirm (alert dialog)",
      cancelButtonLabel: "Cancel (Error dialog)",
    };
    const result = await Xrm.Navigation.openConfirmDialog(alertStrings);
    return result.confirmed;
  }

  // Called by dialog demo if confirm button is selected
  function openAnAlertDialog() {
    const alertStrings = {
      title: "THIS IS AN ALERT",
      text: "This is the text of an alert",
      confirmButtonLabel: "",
    };
    Xrm.Navigation.openAlertDialog(alertStrings);
  }

  // Called by dialog demo if confirmation dialog is cancelled. Note that
  // all options are optional but we must specify at least message or errorCode
  function openAnErrorDialog() {
    const errorOptions = {
      errorCode: 418,
      details: "These are details of an error dialog displayed from open log",
      message: "This is the message of an error dialog",
    };
    Xrm.Navigation.openErrorDialog(errorOptions);
  }

  /**
   * Prompts user to open file and then opens it. This will simply add it to
   * downloads.
   */
  async function openFileDemo() {
    const pickFileOptions = {
      accept: "image",
      allowMultipleFiles: false,
      maximumAllowedFileSize: 1024 * 1024 * 1024 * 5, //5GB
    };
    const files = await Xrm.Device.pickFile(pickFileOptions);

    if (files.length) {
      const fileToOpen = files[0];
      Xrm.Navigation.openFile(fileToOpen, {
        openMode: targets.dialog,
      });
      Xrm.Navigation.openAlertDialog({
        title: "File opened",
        text: `Filesize: ${fileToOpen.fileSize}kb`,
        confirmButtonLabel: "Excellent",
      });
    }
  }

  /** Note, if we with to open a form as a dialog we should use navigate to
   * instead
   */
  function openFormDemo() {
    Xrm.Navigation.openForm({
      entityName: logicalNames.account,
      openInNewWindow: true,
    });
  }

  /**
   * Opens a url in a new window. We can pass an options object as a second arg
   * with width and height properties
   */
  function openUrlDemo() {
    const url = "https://google.com";
    Xrm.Navigation.openUrl(url);
  }

  /**
   * Note Xrm.Panel is in preview. It should not be used in production
   * environments. LoadPanel is the only method
   *
   * The API has now been replaced with Xrm.App.sidepanes.createPane. See below
   */
  async function loadPanelDemo() {
    /*
    DEPRECATED:

    const panelTitle = "GOOGLE";
    const url =
      "https://learn.microsoft.com/en-us/power-apps/developer/" +
      "model-driven-apps/clientapi/reference/xrm-panel/loadpanel";
    Xrm.Panel.loadPanel(url, panelTitle);
    */
    const pane = await Xrm.App.sidePanes.createPane({
      title: "ACCOUNTS",
      canClose: true,
    });
    pane.navigate({
      pageType: pageTypes.entityList,
      entityName: logicalNames.account,
    });
  }

  // Expose onLoad handler
  return {
    navigateToDemo,
    dialogDemo,
    openFileDemo,
    openFormDemo,
    openUrlDemo,
    loadPanelDemo,
  };
})();
