# CHANGELOG

## v0.1.3

Allow updating namespace self-properties `dict_uri`, `dict_lookup`, `dict_ntt_list`.

## v0.1.2

Allow namespaces of type `referential` to not have `dict_uri` specified to
support linking information into systems without API.

## v0.1.1

The cli implementation of some commands was missing. This has been updated.

1. asset::server:property.remove
2. asset::server:property.set
3. asset::server:property.update

## v0.1.0

First implementation of the asset model for the following entities:

1. server
2. runtime environment
3. orchestration environment
4. container

## v0.0.1

First tag of the development version, created to support embedding the
version into the `tom` cli.
