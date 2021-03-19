# Updating Chremoas
Due to the nature of the dependency tree (and the new go module sum business) updating the projects in the correct order is paramount to your sanity.

### Update Order

#### Group 1
1) services-common

#### Group 2 (depends on services-common)
1) chremoas
1) discord-gateway
1) esi-srv
1) perms-srv

#### Group 3
1) role-srv (depends on: discord-gateway perms-srv)

#### Group 4
1) auth-srv (depends on: role-srv)

#### Group 5
1) auth-cmd (depends on: auth-srv chremoas)
1) auth-esi-poller (depends on: auth-srv esi-srv)
1) auth-web (depends on: auth-srv)
1) filter-cmd (depends on: chremoas perms-srv role-srv)
1) lookup-cmd (depends on: chremoas esi-srv)
1) perms-cmd (depends on: chremoas perms-srv role-srv)
1) role-cmd (depends on: chremoas perms-srv role-srv)
1) sig-cmd (depends on: chremoas perms-srv role-srv)

