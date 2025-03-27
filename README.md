# PL-400 Revision Notes

## Overview

Revision notes and code resources for the MS PL-400 exam. These
are a work in progress and have lots of gaps and inaccuracies.

## Resource Notes

### PCF

node_modules ignored. In the root directory of the PCF run:

```console
npm i
npm run refreshTypes
```

### Plug-ins

strong name keys and early bound classes ignored. In the root directory of the
plugin package run:

```console
sn -k PACKAGENAME.snk
pac modelbuilder build -o ./Model -stf ./Model/builderSettings.json
```
