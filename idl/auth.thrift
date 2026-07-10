namespace go auth

struct User {
    1: i64 id
    2: string username
    3: string nickname
    4: string created_at
    5: string updated_at
}

struct UserLoginRequest {
    1: string username (api.body="username")
    2: string password (api.body="password")
}

struct UserLoginResponse {
    1: User user
}

struct GetCurrentUserRequest {
}

struct GetCurrentUserResponse {
    1: User user
}

struct LogoutResponse {
    1: string message
}

service AuthService {
    UserLoginResponse UserLogin(1: UserLoginRequest req) (api.post="/api/admin/login")
    GetCurrentUserResponse GetCurrentUser(1: GetCurrentUserRequest req) (api.get="/api/admin/me")
}

