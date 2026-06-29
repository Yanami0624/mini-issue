# CLI examples

This folder contains command-line examples for the mini-issue API.

## Start the server

Create the local MySQL schema first:

```powershell
Get-Content -Raw examples/cli/schema.sql | mysql -uroot
```

Then start the server from the repository root:

```powershell
go run ./cmd
```

The server listens on:

```text
http://localhost:8080
```

## Run the full PowerShell demo

Open another terminal from the repository root:

```powershell
powershell -ExecutionPolicy Bypass -File examples/cli/powershell_issue_flow.ps1
```

You can also pass your own username, password, or issue id:

```powershell
powershell -ExecutionPolicy Bypass -File examples/cli/powershell_issue_flow.ps1 `
  -Username alice `
  -Password 123456 `
  -IssueId 1
```

## What the demo does

1. Registers a user.
2. Logs in and stores the JWT token.
3. Calls `GET /me`.
4. Calls issue APIs with `Authorization: Bearer <token>`:
   - `POST /issues`
   - `GET /issues`
   - `GET /issues/:id`
   - `PATCH /issues/:id`
   - `DELETE /issues/:id`

The script creates a new issue and then reads the issue id back from `GET /issues` unless you pass `-IssueId`.
