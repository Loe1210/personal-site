namespace go auth

struct AuthContext { 1: i64 user_id, 2: string username, 3: list<string> roles }
struct ValidateSessionRequest { 1: string session_id }
struct CheckPermissionRequest { 1: i64 user_id, 2: string code }
service AuthService {
  AuthContext ValidateSession(1: ValidateSessionRequest request)
  bool CheckPermission(1: CheckPermissionRequest request)
}