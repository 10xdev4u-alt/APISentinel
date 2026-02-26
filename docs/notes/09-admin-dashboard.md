# 09: Visualizing Security - The Admin Dashboard 📊

Logs and metrics are great for machines, but humans need **visuals**. 

A security professional shouldn't have to `tail -f audit.log` all day. They need a central place to see the health of the system.

## Our Tech Choice: `html/template`
Instead of building a separate React/Vue frontend (which would add complexity), we will use Go's standard `html/template` library. 
- **Fast:** Server-side rendered (SSR).
- **Simple:** No build steps required.
- **Secure:** Go's templates are context-aware and automatically escape data to prevent XSS (the irony is not lost on us!).

## Dashboard Features
1. **Live Stats:** Real-time request and block counters.
2. **Audit Feed:** The most recent security violations in a readable table.
3. **Blocklist View:** See which IPs are currently banned.

## Brutalist Design
In line with our project's aesthetic, we will use a **Brutalist** UI:
- Heavy borders.
- Monospace fonts.
- High contrast (Black/White/Red).
- No unnecessary animations.

Next, we build the Dashboard Handler!
