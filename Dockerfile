FROM postgres:14

# Environment variables for PostgreSQL
ENV POSTGRES_USER=bloguser
ENV POSTGRES_PASSWORD=blogpassword
ENV POSTGRES_DB=blogdb

# Expose the PostgreSQL port
EXPOSE 5432

# Add a volume for data persistence
VOLUME ["/var/lib/postgresql/data"]

# Health check to verify the database is ready
HEALTHCHECK --interval=5s --timeout=5s --retries=5 CMD pg_isready -U $POSTGRES_USER -d $POSTGRES_DB || exit 1

# This Dockerfile creates a PostgreSQL database for a blog API with:
# - Default database name: blogdb
# - Default username: bloguser
# - Default password: blogpassword
# 
# To build and run:
# docker build -t blog-postgres .
# docker run -d -p 5432:5432 --name blog-db blog-postgres
#
# Connection string for Go application:
# "postgres://bloguser:blogpassword@localhost:5432/blogdb?sslmode=disable"