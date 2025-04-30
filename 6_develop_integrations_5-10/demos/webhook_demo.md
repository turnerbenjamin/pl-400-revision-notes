# Webhook demo

This demo simply follows the activity described
[here](https://learn.microsoft.com/en-us/training/modules/integrate-dataverse-azure-solutions/azure-function),
with a few minor modifications.

## Create an Azure Function

A simple Azure function has been produced. This just logs some information from
the request.

## Register Webhook

In PRT register a webhook. We need to provide the function url and set
authentication. Azure functions use WebhookKeys which are sent in the query
string parameters as code={webhookKey}.

![register webhook](./screens/webhook_demo/1_register_webhook.png)

Next we register a step, here a synchronous step has been chosen:

![register step](./screens/webhook_demo/2_register_webhook_step.png)

## Testing

Finally, we can trigger the step and view the logs from the Azure Function:

![testing webhook](./screens/webhook_demo/3_azure_func_log.png)
