# 08: Dynamic Configuration with YAML ⚙️

In the early stages of a project, hardcoding values (like the port or backend URL) is fine. In the "Intermediate" phase, we moved to Environment Variables. 

But for a "Pro" tier application, you need a **Configuration File**.

## Why use a Config File?
1. **Complexity:** Environment variables become messy when you have complex nested structures (like multiple proxy routes).
2. **Version Control:** You can keep a `config.example.yaml` in your repo to show users how to set up the app.
3. **Hot Reloading (Pro Goal):** In the future, we can make the proxy watch the config file and update its routes without restarting!

## YAML: The Industry Standard
YAML (Yet Another Markup Language) is the most popular format for infrastructure tools (Kubernetes, Docker Compose, Traefik). It is:
- **Human Readable:** Easy to read and write.
- **Hierarchical:** Perfect for defining routes and their specific settings.

## Our Config Structure
We will define:
- `server`: port, admin key, rate limits.
- `routes`: a list of path prefixes and their target URLs.
- `security`: which middlewares to enable.

Next, we will implement the YAML parser in Go!
