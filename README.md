# go-manila-backend

A modern Go backend for restoring online services on legacy HTC Sense / Manila devices.

The project aims to provide compatibility endpoints for HTC Sense services while using modern APIs and infrastructure behind the scenes.

The initial focus is restoring the **Weather** functionality on devices such as the **HTC Touch Diamond2 running Windows Mobile 6.5 and HTC Sense / Manila**.

## Motivation

HTC Sense relied on remote services that are no longer available or compatible with modern infrastructure.

Instead of modifying Manila itself, `go-manila-backend` provides a compatibility layer:

```text
HTC Sense / Manila
        |
        | Legacy HTTP + XML
        v
go-manila-backend
        |
        | Modern HTTPS APIs
        v
External services