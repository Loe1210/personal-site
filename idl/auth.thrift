namespace go auth

struct AdminUser {
    1: i64 id
    2: string username
    3: string nickname
    4: string created_at
    5: string updated_at
}

struct AdminLoginRequest {
    1: string username (api.body="username")
    2: string password (api.body="password")
}

struct AdminLoginResponse {
    1: string token
    2: string expires_at
    3: AdminUser user
}

struct GetCurrentAdminRequest {
}

struct GetCurrentAdminResponse {
    1: AdminUser user
}

service AuthService {
    AdminLoginResponse AdminLogin(1: AdminLoginRequest req) (api.post="/api/admin/login")
    GetCurrentAdminResponse GetCurrentAdmin(1: GetCurrentAdminRequest req) (api.get="/api/admin/me")
}