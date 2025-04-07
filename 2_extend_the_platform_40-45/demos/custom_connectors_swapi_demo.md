# Custom Connectors SWAPI Demo

In this demo, a custom connector is created for the pre-existing SWAPI API.

## Definition

The general tab contains some display information and high-level connection
details:

![Create-General](./screens/cc_swapi_create_general.png)

The SWAPI API does not have any authentication:

![Create-Security](./screens/cc_swapi_create_security.png)

Three actions have been set-up for different endpoints.

![Create-Definition](./screens/cc_swapi_create_definition_actions_list.png)

## Custom Code

This example makes heavy use of custom code to transform the response.
Specifically, it is used to remove pagination and to change the search param
name.

However, I found that custom code in a connector is not well designed for
extensive transformations:

- We can only provide a single code block for multiple actions
- There are limited namespaces available at runtime
- The code cannot exceed 1mb
- The execution time of the code cannot exceed 5 seconds

A better approach would have been to:

- use a proxy to transform the response
- use a custom connector to connect to the proxy
- use policies to transform the query string parameters

If you are interested, the code can be found
[here](../resources/CustomConnectors/SwapiConnector/SwapiCustomLogic.cs).

![Create-Code](./screens/cc_swapi_create_custom_code.png)

## Usage

A simple flow has been created to test this custom connector:

![Access custom connector](./screens/cc_swapi_usage_find_connector.png)

![Get people notification](./screens/cc_swapi_usage_get_all_characters.png)

![Get person notification](./screens/cc_swapi_usage_get_character.png)

![Get person notification](./screens/cc_swapi_usage_get_planet.png)
