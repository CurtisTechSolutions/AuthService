# AuthService

Auth Service is an independent service for authentication, authorization, and session management.

I wanted to build this for my own services, and also to try out GORM.

The auth service is used to seperate authetication logic from other service-specific backend logic. This means you can use this auth service across multiple different applications at the same time (multi-tenancy or single-tenant). Using this service allows you to focus on service specific logic. I don't want to use services like Supabase, and decided to build this over the weekend as a fun project.
