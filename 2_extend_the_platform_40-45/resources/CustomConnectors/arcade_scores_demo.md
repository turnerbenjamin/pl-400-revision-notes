# Arcade Scores Demo

## Intro

Connectors are just wrappers around APIs that describe the API in terms
that Power Apps can understand.

A connector contains:

- Triggers: Kick off a process in power apps
- Actions: Kick off a process in the API

Setting up a simple connector with a few actions is a simple job that takes a
few minutes. However, there are features of connectors, which require the API to
have certain capabilities. This project focussed on two of these:

- Triggers
- Expanding the OpenAI Definition with Dynamic values

## Authentication

This connector uses Entra ID to authenticate with the API. The first step here
is to create an identity provider for the function app. To do this select
authentication and Add an identity provider. This will create an App
registration for the function.

Next, we need to create a new app registration with permission to act on behalf
of the function app. Create t
