# 15: Environment Variable Overrides 🌍

A "Pro" Tier application follows the **Twelve-Factor App** methodology. One of its core principles is: **Store config in the environment.**

## The Hybrid Approach
We already have a `config.yaml` file, which is great for complex structures like routes. But what if you want to quickly change the `ADMIN_KEY` in production without modifying the file?

We will implement a system that:
1.  Loads the `config.yaml`.
2.  Checks the environment for specific variables (e.g. `SENTINEL_ADMIN_KEY`).
3.  If found, **overwrites** the value from the YAML.

## Naming Convention
We will use the prefix `SENTINEL_` for all our overrides to avoid clashing with other system variables.
- `SENTINEL_PORT` -> Overrides `server.port`
- `SENTINEL_ADMIN_KEY` -> Overrides `server.admin_key`

This makes API Sentinel perfectly suited for Docker and CI/CD pipelines.

Next, we will update the Config loader!
