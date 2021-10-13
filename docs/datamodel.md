# Datamodel

## asset::server

 Attribute              | Value Type | Cardinality | Permanence | Possible Values
------------------------|------------|-------------|------------|-------------------
 name                   | string     | unique      | ephemeral  |
 type                   | enum       | not unique  | perpetual  | physical, virtual
 serial                 | string     | unique      | perpetual  |
 dc_id                  | string     | unique      | perpetual  |
 interface_list         | list       | not unique  | ephemeral  | acc0, mgmt0
   if_{}_link_addr_list | list       | not unique  | ephemeral  |

## asset::runtime

 Attribute              | Value Type | Cardinality | Permanence | Values
------------------------|------------|-------------|------------|-------------
 name                   | string     | unique      | ephemeral  |
 type                   | string     | not unique  | perpetual  |
 bios_uuid              | string     | unique      | ephemeral  |
 accountable            | string     | not unique  | ephemeral  |
 responsible_team       | string     | not unique  | ephemeral  |
 lifecycle_phase        | enum       | not unique  | ephemeral  | acquire, use, renewal
 fqdn                   | string     | not unique  | ephemeral  |
 interface_list         | list       | not unique  | ephemeral  |
   if_{}_link_addr_list | list       | not unique  | ephemeral  |
   if_{}_ip_addr_list   | list       | not unique  | ephemeral  |
 listen_socket_list     | list       | not unique  | ephemeral  |
 version                | string     | not unique  | ephemeral  |
 user_mgmt

## asset::orchestration

 Attribute              | Value Type | Cardinality | Permanence | Values
------------------------|------------|-------------|------------|-------------
 name                   | string     | unique      | ephemeral  |
 type                   | string     | not unique  | perpetual  |

## asset::container

 Attribute              | Value Type | Cardinality | Permanence | Values
------------------------|------------|-------------|------------|-------------
 name
 type
 registry_link
 
