namespace go rbac

struct Role {
    1: i64 id
    2: string name
    3: string code
    4: string description
    5: string created_at
    6: string updated_at
}

struct Permission {
    1: i64 id
    2: string name
    3: string code
    4: string resource
    5: string action
    6: string description
    7: string created_at
    8: string updated_at
}

struct CurrentAdminRBACRequest {
}

struct CurrentAdminRBACResponse {
    1: list<Role> roles
    2: list<Permission> permissions
}

struct ListRolesRequest {
}

struct ListRolesResponse {
    1: list<Role> list
}

struct CreateRoleRequest {
    1: string name (api.body="name")
    2: string code (api.body="code")
    3: string description (api.body="description")
}

struct CreateRoleResponse {
    1: Role role
    2: string message
}

struct ListPermissionsRequest {
}

struct ListPermissionsResponse {
    1: list<Permission> list
}

struct CreatePermissionRequest {
    1: string name (api.body="name")
    2: string code (api.body="code")
    3: string resource (api.body="resource")
    4: string action (api.body="action")
    5: string description (api.body="description")
}

struct CreatePermissionResponse {
    1: Permission permission
    2: string message
}

struct BindUserRolesRequest {
    1: i64 user_id (api.path="user_id")
    2: list<i64> role_ids (api.body="role_ids")
}

struct BindUserRolesResponse {
    1: bool success
    2: string message
}

struct BindRolePermissionsRequest {
    1: i64 role_id (api.path="role_id")
    2: list<i64> permission_ids (api.body="permission_ids")
}

struct BindRolePermissionsResponse {
    1: bool success
    2: string message
}

service RBACService {
    CurrentAdminRBACResponse GetCurrentAdminRBAC(1: CurrentAdminRBACRequest req) (api.get="/api/admin/rbac/me")
    ListRolesResponse ListRoles(1: ListRolesRequest req) (api.get="/api/admin/roles")
    CreateRoleResponse CreateRole(1: CreateRoleRequest req) (api.post="/api/admin/roles")
    ListPermissionsResponse ListPermissions(1: ListPermissionsRequest req) (api.get="/api/admin/permissions")
    CreatePermissionResponse CreatePermission(1: CreatePermissionRequest req) (api.post="/api/admin/permissions")
    BindUserRolesResponse BindUserRoles(1: BindUserRolesRequest req) (api.put="/api/admin/users/:user_id/roles")
    BindRolePermissionsResponse BindRolePermissions(1: BindRolePermissionsRequest req) (api.put="/api/admin/roles/:role_id/permissions")
}
