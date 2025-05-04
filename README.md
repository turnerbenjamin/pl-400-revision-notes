# PL-400 Revision Notes

## Overview

The folder structure of this repo reflects the skills measured for the PL-400
as of 19th December 2024; although the skills have been reordered. The first
two sections are about extending Power Platform on the frontend and backend
respectively.

The third section looks at Power Platform and the OOTB features it includes. For
the exam it is important to understand these features to determine when we
should extend the platform. When developing with Power Platform the correct
approach is to use OOTB and no/low code features where possible and extend with
code only where necessary. This keeps development simple and open to a broader
range of developers/makers.

The forth section is focussed on environment configuration and application
lifecycle management (ALM). It looks at the security model in Power Platform
and the use of solutions as a mechanism to implement ALM.

The fifth section looks at more advanced features of canvas and model-driven
apps, in particular debugging applications and making performance improvements.
The expectation is that we already understand how to build basic applications.

The final section is focussed on integration, particularly with Azure; although
this is covered in other sections too.

The exam is aimed at individuals with an applied knowledge of these skills. It
is helpful to create a developer environment to explore all of the concepts
covered in the exam. These notes include a wide variety of demos which have cost
a total of £0.05 to produce.

## Resource Notes

A few of the sections contain a resource directory containing example code.
Required elements are often in a .gitignore:

### PCF

In the root directory of the PCF run:

```console
npm i
npm run refreshTypes
```

### Plug-ins

Early bound classes and strong name keys are in .gitignore. In the root
directory of the plugin package run:

```console
sn -k {PACKAGENAME}.snk
pac modelbuilder build -o ./Model -stf ./Model/builderSettings.json
```
