## Kind: design

PR #1208 is a Dependabot bump (react 18.3.1 to 19.3.1, @types/react to 19.3.0). The `Dashboard UI Tests` check is red: `react-dom` (18.3.1) and `@types/react-dom` (18.3.7) were not bumped with `react`, so 7 Jest suites fail with `ReactCurrentDispatcher` undefined.

Only a commit fixes it (move `react-dom` and `@types/react-dom` to 19.x, refresh `yarn.lock`, fix any React 19 API breakage). The fleet adds no commit to a Dependabot branch it has not taken over. I posted the fix mark for head f33893f68309c20085591819dc61af4eed51fcd8; the dispatcher takes the branch over on that mark and starts a fix round, or a person fixes it, closes it, or asks Dependabot (`@dependabot recreate`).
