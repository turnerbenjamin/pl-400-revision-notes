using System;
using DemoPlugins.Model;
using Microsoft.Xrm.Sdk;

namespace DemoPlugins
{
    public class PlugInParametersDemo : PluginBase
    {
        private string preImageName = "contactPreImage";

        /*
        There are various parameter collections in context (input, output,
        pre-images, post-images)

        This plugin demonstrates usage of these collections

        It is designed to run in the pre-operation stage of the update message
        for contact filtered by the first name field
        */
        public override void ExecuteBusinessLogic(IServiceProvider serviceProvider)
        {
            var targetEntity = GetTargetEntity();
            var preImageEntity = GetEntityPreImage();
            targetEntity[Contact.Fields.EMailAddress1] = GetUpdatedEmail(
                targetEntity,
                preImageEntity
            );
        }

        // Since this plug-in runs when first name is updated and first name is
        // required, first name will be on target. It may not include additional
        // fields.
        //
        // A preimage is registered with a single field, emailAddress, this is
        // used to access the old email address and update the username
        public string GetUpdatedEmail(Entity target, Entity preimage)
        {
            var castTarget = target.ToEntity<Contact>();
            var castPreImage = preimage.ToEntity<Contact>();

            var oldEmail = castPreImage.EMailAddress1 ?? "username@email.com";
            var domain = oldEmail.Split('@')[1];
            return $"{castTarget.FirstName.ToLower()}@{domain}";
        }

        private Entity GetTargetEntity()
        {
            // Helper method wrapping context.InputParameters["Target"] with
            // error handling
            return GetContextParameterOrDefault<Entity>(context.InputParameters, "Target")
                ?? throw new ArgumentException("Unable to access Target from plug-in context");
        }

        private Entity GetEntityPreImage()
        {
            // Helper method wrapping context.PreEntityImages[preImageName] with
            // error handling
            return GetEntityImage(context.PreEntityImages, preImageName)
                ?? throw new ArgumentException($"Unable to access {preImageName} from context");
        }
    }
}
