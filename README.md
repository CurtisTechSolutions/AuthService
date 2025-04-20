# AuthService

Auth Service is an independent service for authentication, authorization, and session management.

I wanted to build this for my own services, and also to try out GORM.

The auth service is used to separate authentication logic from other service-specific backend logic. This means you can use this auth service across multiple different applications at the same time (multi-tenancy or single-tenant). Using this service allows you to focus on service specific logic. I don't want to use services like Supabase, and decided to build this over the weekend as a fun project.

# Roadmap

- Redis/In-memory KV-store for session ids and query caching
- Logging
- Test cases
- Authorization system (Array of strings? Something like "admin.read" or "admin.write"?). I'm leaning towards JWT for this, so that it's service specific. This auth service should return the access/authorization claims only, and not decide whether a request is authorized for an unrelated service.
- JWT?