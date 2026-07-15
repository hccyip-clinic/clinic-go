<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

# Yes, please provide a decision matrix for my notes

Here’s a copy-ready decision matrix you can paste into your notes. For a simple CRUD dashboard, the choice mainly comes down to how much client-side interactivity you expect and how much frontend complexity you want to carry.[^1][^2][^3]

## Decision matrix

| Criterion | Hugo + HTMX + SQLite | Vite + React |
| :-- | :-- | :-- |
| CRUD simplicity | Excellent for table/form/modal workflows with server-rendered HTML. [^4][^3] | Good, but more setup for data fetching, forms, and state. [^5] |
| Frontend complexity | Low; minimal JavaScript and fewer moving parts. [^4][^3] | Higher; component model, routing, state, and query layers. [^5] |
| Interactivity | Best for modest interactivity, partial page swaps, and live sections. [^4][^1] | Best for rich UX, complex client state, and advanced interactions. [^6][^3] |
| Development speed for small CRUD | Very fast if the app is mostly server-rendered forms and lists. [^4][^7] | Fast once the React stack is set up, but initial scaffolding is heavier. [^5] |
| Validation and errors | Server-first validation is natural; HTMX updates fragments after submit. [^8][^9] | Strong client + server validation pattern, but more code. [^5] |
| State management | Mostly server-side; little client state to maintain. [^4][^3] | Requires client-side state handling, though TanStack Query reduces pain. [^10][^5] |
| Data layer | SQLite fits naturally as the app database behind the server. [^11][^8][^12] | Usually needs a separate API in front of SQLite or another DB. [^13][^10] |
| Deployment | Simple if the backend and SQLite are deployed together. [^11][^14] | Simple for static frontend, but backend is still separate. [^13][^10] |
| Team fit | Great when the team prefers backend/server-rendered development. [^4][^2] | Better when the team is strong in React and component-based UI. [^5] |
| Long-term extensibility | Good for admin tools and CRUD-heavy internal apps. [^4][^7] | Better for scaling into a richer product UI over time. [^6][^5] |

## Recommendation rule

- Choose **Hugo + HTMX + SQLite** if the app is mostly CRUD, internal, and you want the simplest maintainable stack.[^3][^4][^7]
- Choose **Vite + React** if you expect the dashboard to grow into a richer client-driven product with more interactive workflows.[^6][^5]


## Practical shortcut

A good shortcut is: if your pages are mostly “list, filter, edit, save,” use HTMX; if your app needs “complex client state, reusable widgets, and heavy interactivity,” use React.[^1][^6][^3]
<span style="display:none">[^15][^16]</span>

<div align="center">⁂</div>

[^1]: https://dev.to/pockit_tools/htmx-in-2026-when-you-dont-need-react-and-when-you-absolutely-do-2mf4

[^2]: https://www.softwareseni.com/when-to-choose-htmx-over-react-a-strategic-decision-framework/

[^3]: https://toolchew.com/en/htmx-vs-react/

[^4]: https://ggoggam.github.io/blog/dashboard

[^5]: https://dev.to/singhamandeep007/building-a-production-ready-react-vite-typescript-boilerplate-architecture-choices-dx-2i7l

[^6]: https://www.buildmvpfast.com/blog/htmx-vs-react-2026-decision-framework

[^7]: https://hackernoon.com/why-more-developers-are-choosing-htmx-over-react-for-crud-apps

[^8]: https://hsong.me/posts/crud-with-htmx-typescript-and-sqlite/

[^9]: https://www.oreilly.com/library/view/server-driven-web-apps/9798888651193/f_0013.xhtml

[^10]: https://tanstack.com/query/latest

[^11]: https://github.com/emarifer/go-htmx-demo

[^12]: https://dev.to/hungle00/htmx-for-building-a-sqlite-query-dashboard-32ck

[^13]: https://dev.to/truongpx396/reactjs-typescript-vite-with-redux-and-tanstack-react-query-in-practice-7eg

[^14]: https://github.com/emarifer/go-frameworkless-htmx-todoapp

[^15]: https://solid-web.com/htmx-vs-react-2026/

[^16]: https://www.softwareseni.com/htmx-vs-react-performance-and-architecture-deep-dive/

